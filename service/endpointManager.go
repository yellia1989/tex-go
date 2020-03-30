/*
* endpointManager负责维护obj对应的可用的adapterProxy列表
* direct=true时 adapterProxy列表不会定时更新状态
* direct=false时 adaterProxy列表会定时通过registry更新
*/

package tex

import (
    "sync"
    "time"
    "fmt"
    "strings"
    "sync/atomic"
    "math/rand"
    "github.com/yellia1989/tex-go/tools/log"
    "github.com/yellia1989/tex-go/sdp/rpc"
)

type endpointManager struct {
    sObjName string
    sDivision string
    comm *Communicator
    refreshInterval time.Duration
    direct bool
    query *rpc.Query
    done chan bool

    mu sync.Mutex
    ready bool
    mAdapter map[Endpoint]*adapterProxy
    vEndpoint []*Endpoint
    index int
}

func newEpMgr(objName string, comm *Communicator) (*endpointManager, error) {
    epmgr := &endpointManager{sObjName: objName, comm: comm, refreshInterval: cliCfg.endpointRefreshInterval}
    epmgr.mAdapter = make(map[Endpoint]*adapterProxy)
    epmgr.sDivision = cliCfg.division

    p := strings.Index(objName, "@")
    if p != -1 {
        epmgr.sObjName = objName[:p]
        epmgr.direct = true
        epmgr.ready = true
        vEndpoint := strings.Split(objName[p+1:], ":")
        for _, ep := range vEndpoint {
            ep, err := NewEndpoint(ep)
            if err != nil {
                return nil, err
            }
            epmgr.vEndpoint = append(epmgr.vEndpoint, ep)
        }
        if len(epmgr.vEndpoint) == 0 {
            return nil, fmt.Errorf("empty endpoints for obj:%s", epmgr.sObjName)
        }
        // 根据endpoint创建adapter
        for _, ep := range epmgr.vEndpoint {
            adapter, err := newAdapter(ep)
            if err != nil {
                return nil,fmt.Errorf("create adapter for ep:%s, obj:%s, err:%s", ep, epmgr.sObjName, err.Error())
            }
            epmgr.mAdapter[*ep] = adapter
        }
    } else {
        epmgr.sObjName = objName
        p := strings.Index(objName, "%")
        if p != -1 {
            epmgr.sObjName = objName[:p]
            epmgr.sDivision = objName[p+1:]
        }
        epmgr.query = new(rpc.Query)
        epmgr.done = make(chan bool)
        if strings.Index(comm.sLocator, "@") == -1 {
            return nil, fmt.Errorf("invalid locator")
        }
        comm.StringToProxy(comm.sLocator, epmgr.query)
        if err := epmgr.refreshEndpoint(); err != nil {
            return nil, err
        }

        go func() {
            ticker := time.NewTicker(epmgr.refreshInterval)
            select {
            case <-ticker.C:
                epmgr.refreshEndpoint()
            case <-epmgr.done:
                ticker.Stop()
                return
            }
        }()
    }

    return epmgr, nil
}

func (epmgr *endpointManager) refreshEndpoint() error {
    var vActiveEps []string
    var vInactiveEps []string
    log.FDebugf("registry query endpoint start, obj:%s", epmgr.sObjName)
    ret, err := epmgr.query.GetEndpoints(epmgr.sObjName, epmgr.sDivision, &vActiveEps, &vInactiveEps)
    if ret != 0 || err != nil {
        log.FErrorf("registry query endpoint failed, obj:%s, err:%s, ret:%d", epmgr.sObjName, err.Error(), ret)
        return err
    }
    log.FDebugf("registry query endpoint success, obj:%s", epmgr.sObjName)

    // 根据endpoint创建adapter
    vEndpoint := make([]*Endpoint, 0)
    for _, addr := range vActiveEps {
        ep, err := NewEndpoint(addr)
        if err != nil {
            log.FErrorf("parse endpoint failed, ep:%s, obj:%s, err:%s", addr, epmgr.sObjName, err.Error())
            continue
        }
        vEndpoint = append(vEndpoint, ep)
    }
    mAdapter := make(map[Endpoint]*adapterProxy)
    for _, ep := range vEndpoint {
        adapter, err := newAdapter(ep)
        if err != nil {
            log.FErrorf("create adapter for ep:%s, obj:%s, err:%s", ep, epmgr.sObjName, err.Error())
            continue
        }
        mAdapter[*ep] = adapter
    }
    if len(mAdapter) == 0 {
        return fmt.Errorf("no active endpoint")
    }

    // 待关闭的adapter
    var closeAdapter []*adapterProxy
    changed := false
    epmgr.mu.Lock()
    // 删除现有无用的adapter
    for ep, adapter := range epmgr.mAdapter {
        _, ok := mAdapter[ep]
        if !ok {
            closeAdapter = append(closeAdapter, adapter)
            changed = true
        }
    }
    // 添加新的adapter
    for ep, adapter := range mAdapter {
        _, ok := epmgr.mAdapter[ep]
        if !ok {
            epmgr.mAdapter[ep] = adapter
            changed = true
        }
    }
    if changed {
        epmgr.vEndpoint = vEndpoint
    }

    if !epmgr.ready {
        epmgr.ready = true
    }

    epmgr.mu.Unlock()

    for _, adapter := range closeAdapter {
        adapter.close()
    }

    return nil
}

func (epmgr *endpointManager) selectAdapter(bHash bool, hashCode uint64) (*adapterProxy, error) {
    if bHash {
        return epmgr.selectHashAdapter(hashCode)
    }
    return epmgr.selectNextAdapter()
}

func (epmgr *endpointManager) selectHashAdapter(hashCode uint64) (*adapterProxy, error) {
    mu := &epmgr.mu
    mu.Lock()

    defer mu.Unlock()

    l := len(epmgr.vEndpoint)
    if l == 0 {
        return nil, fmt.Errorf("empty endpoint for obj:%s", epmgr.sObjName)
    }

    vEndpoint := epmgr.vEndpoint
    vEndpoint2 := make([]*Endpoint, 0)
    for {
        p := int(hashCode % uint64(len(vEndpoint)))
        ep := vEndpoint[p]
        adapter, ok := epmgr.mAdapter[*ep]
        if !ok {
            // endpoint和proxy应该是一一对应的
            panic(fmt.Errorf("no adapter for endpoint:%s, obj:%s", ep, epmgr.sObjName))
        }
        if adapter.isActive(false) {
            return adapter, nil
        }
        if atomic.LoadInt32(&adapter.connfailed) != 1 {
            vEndpoint2 = append(vEndpoint2, vEndpoint[p])
        }
        vEndpoint = append(vEndpoint[:p],vEndpoint[p+1:]...)
        if len(vEndpoint) == 0 {
            break
        }
    }

    // 到此为止，我们已经没有了活跃连接可用
    // 这时可以尝试不稳定连接,但不使用直接连接失败(connfailed==1)的连接
    if len(vEndpoint2) != 0 {
        p := int(hashCode % uint64(len(vEndpoint2)))
        ep := vEndpoint2[p]
        adapter := epmgr.mAdapter[*ep]
        if adapter.isActive(true) {
            return adapter, nil
        }
    }

    return nil, fmt.Errorf("no active adapter for obj:%s", epmgr.sObjName)
}

func (epmgr *endpointManager) selectNextAdapter() (*adapterProxy, error) {
    mu := &epmgr.mu
    mu.Lock()

    defer mu.Unlock()

    l := len(epmgr.vEndpoint)
    if l == 0 {
        return nil, fmt.Errorf("empty endpoint for obj:%s", epmgr.sObjName)
    }

    // 总共尝试的次数,超过这个次数说明不可能找到可用的adapter
    count := l
    vEndpoint2 := make([]*Endpoint, 0)
    for count > 0 {
        epmgr.index += 1
        if epmgr.index >= l {
            epmgr.index = 0
        }
        ep := epmgr.vEndpoint[epmgr.index % l]
        adapter, ok := epmgr.mAdapter[*ep]
        if !ok {
            // endpoint和proxy应该是一一对应的
            panic(fmt.Errorf("no adapter for endpoint:%s, obj:%s", ep, epmgr.sObjName))
        }
        if adapter.isActive(false) {
            return adapter, nil
        }
        if atomic.LoadInt32(&adapter.connfailed) != 1 {
            vEndpoint2 = append(vEndpoint2, ep)
        }
        count -= 1
    }

    // 到此为止，我们已经没有了活跃连接可用
    // 这时可以尝试不稳定连接,但不使用直接连接失败(connfailed==1)的连接
    if len(vEndpoint2) != 0 {
        p := rand.Int() % len(vEndpoint2)
        ep := vEndpoint2[p]
        adapter := epmgr.mAdapter[*ep]
        if adapter.isActive(true) {
            return adapter, nil
        }
    }

    return nil, fmt.Errorf("no active adapter for obj:%s", epmgr.sObjName)
}

func (epmgr *endpointManager) close() {
    for k, v := range epmgr.mAdapter {
        v.close() 
        log.FDebugf("adapter:%s has been closed", &k)
    }
    epmgr.done <- true
}