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
    mAdapter map[Endpoint]*adapterProxy
    vEndpoint []*Endpoint
    index int
    ready bool
}

func newEpMgr(objName string, comm *Communicator) (*endpointManager) {
    epmgr := &endpointManager{sObjName: objName, comm: comm, refreshInterval: cliCfg.endpointRefreshInterval}
    epmgr.mAdapter = make(map[Endpoint]*adapterProxy)
    epmgr.sDivision = cliCfg.division

    p := strings.Index(objName, "@")
    if p != -1 {
        epmgr.sObjName = objName[:p]
        epmgr.direct = true
        vEndpoint := strings.Split(objName[p+1:], ":")
        for _, ep := range vEndpoint {
            if ep, err := NewEndpoint(ep); err == nil {
                epmgr.vEndpoint = append(epmgr.vEndpoint, ep)
            }
        }
        // 根据endpoint创建adapter
        for _, ep := range epmgr.vEndpoint {
            if adapter, err := newAdapter(ep); err == nil {
                epmgr.mAdapter[*ep] = adapter
            }
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
        comm.StringToProxy(comm.sLocator, epmgr.query)

        go func() {
            ticker := time.NewTicker(epmgr.refreshInterval)
            for {
                select {
                case <-ticker.C:
                    epmgr.refreshEndpoint()
                case <-epmgr.done:
                    ticker.Stop()
                    return
                }
            }
        }()
    }

    return epmgr
}

func (epmgr *endpointManager) refreshEndpoint() {
    var vActiveEps []string
    var vInactiveEps []string
    log.FDebugf("registry query endpoint start, obj:%s", epmgr.sObjName)
    ret, err := epmgr.query.GetEndpoints(epmgr.sObjName, epmgr.sDivision, &vActiveEps, &vInactiveEps)
    if ret != 0 || err != nil {
        serr := ""
        if err != nil {
            serr = err.Error()
        }
        log.FErrorf("registry query endpoint failed, obj:%s, division:%s, err:%s, ret:%d", epmgr.sObjName, epmgr.sDivision, serr, ret)
        return
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
        log.FErrorf("no active endpoint, obj:%s", epmgr.sObjName)
        return
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
            go adapter.checkActive()
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
}

func (epmgr *endpointManager) selectAdapter(bHash bool, hashCode uint64) (*adapterProxy, error) {
    mu := &epmgr.mu
    mu.Lock()
    if !epmgr.ready {
        epmgr.refreshEndpoint()
    }
    epmgr.mu.Unlock()

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
    vEndpoint3 := make([]*Endpoint, 0)
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
        } else {
            vEndpoint3 = append(vEndpoint3, vEndpoint[p])
        }
        vEndpoint = append(vEndpoint[:p],vEndpoint[p+1:]...)
        if len(vEndpoint) == 0 {
            break
        }
    }

    // 到此为止，我们已经没有了活跃连接可用
    // 这时可以尝试不稳定连接
    if len(vEndpoint2) != 0 {
        p := int(hashCode % uint64(len(vEndpoint2)))
        ep := vEndpoint2[p]
        adapter := epmgr.mAdapter[*ep]
        if adapter.isActive(true) {
            return adapter, nil
        }
    }

    // 不使用直接连接失败(connfailed==1)的连接
    if len(vEndpoint3) != 0 {
        p := int(hashCode % uint64(len(vEndpoint3)))
        ep := vEndpoint3[p]
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
    vEndpoint3 := make([]*Endpoint, 0)
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
        } else {
            vEndpoint3 = append(vEndpoint3, ep)
        }
        count -= 1
    }

    // 到此为止，我们已经没有了活跃连接可用
    // 这时可以先尝试不稳定连接
    if len(vEndpoint2) != 0 {
        p := rand.Int() % len(vEndpoint2)
        ep := vEndpoint2[p]
        adapter := epmgr.mAdapter[*ep]
        if adapter.isActive(true) {
            return adapter, nil
        }
    }

    // 再使用直接连接失败(connfailed==1)的连接
    if len(vEndpoint3) != 0 {
        p := rand.Int() % len(vEndpoint3)
        ep := vEndpoint3[p]
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
