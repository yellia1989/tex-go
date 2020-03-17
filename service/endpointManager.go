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
    comm *Communicator
    refreshInterval time.Duration
    direct bool
    query *rpc.Query

    mu sync.Mutex
    ready bool
    mAdapter map[Endpoint]*adapterProxy
    vEndpoint []*Endpoint
    index int
    depth int
}

func newEpMgr(objName string, comm *Communicator) (*endpointManager, error) {
    epmgr := &endpointManager{sObjName: objName, comm: comm, refreshInterval: cliCfg.endpointRefreshInterval}
    epmgr.mAdapter = make(map[Endpoint]*adapterProxy)

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
        epmgr.query = new(rpc.Query)
        comm.StringToProxy(comm.sLocator, epmgr.query)
        epmgr.refreshEndpoint()
        // 定时更新endpoint列表
        // go epmgr.refreshEndpoint()
    }

    return epmgr, nil
}

func (epmgr *endpointManager) refreshEndpoint() {
    // TODO
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
    vEndpoint2 := make([]*Endpoint, 1)
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
    vEndpoint2 := make([]*Endpoint, 1)
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
        <-v.done
        log.FDebugf("adapter:%s has been closed", &k)
    }
}
