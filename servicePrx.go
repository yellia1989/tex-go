package tex

import (
    "time"
    "strings"
    "sync/atomic"
    "github.com/yellia1989/tex-go/protocol/protocol"
    "github.com/yellia1989/tex-go/tools/log"
)

type ServicePrx interface {
    SetPrxImpl(impl ServicePrxImpl)
}

type ServicePrxImpl interface {
    // 只有当error == nil时，resp才有效
    Invoke(sFuncName string, params []byte, resp *protocol.ResponsePacket) error
    SetTimeout(timeout time.Duration)
}

type servicePrxImpl struct {
    name string
    comm *Communicator
    epmgr *endpointManager

    invokeTimeout time.Duration
    reqid uint32
}

func (impl *servicePrxImpl) Invoke(sFuncName string, params []byte, resp *protocol.ResponsePacket) error {
    // 构造请求消息
    req := &protocol.RequestPacket{
        IRequestId: atomic.AddUint32(&impl.reqid, 1),
        SServiceName: impl.name,
        SFuncName: sFuncName,
        SReqPayload: string(params),
        ITimeout: uint32(impl.invokeTimeout.Milliseconds()),
    }

    // 选择一个adapterProxy发送消息
    adapter, err := impl.epmgr.selectAdapter(false, 0)
    if err != nil {
        return err
    }

    // 等待消息返回
    err = adapter.invoke(req, resp)
    if err != nil {
        log.FErrorf("invoke obj:%s, func:%s, err:%s", impl.name, sFuncName, err.Error())
        return err
    }

    return nil
}

func (impl *servicePrxImpl) SetTimeout(timeout time.Duration) {
    impl.invokeTimeout = timeout
}

func (impl *servicePrxImpl) close() {
    impl.epmgr.close()
}

func newPrxImpl(name string, comm *Communicator) (*servicePrxImpl, error) {
    serviceName := name
    if p := strings.Index(name, "@"); p != -1 {
        serviceName = name[:p]
    }
    impl := &servicePrxImpl{name: serviceName, comm: comm, invokeTimeout: cliCfg.invokeTimeout}

    var err error
    impl.epmgr, err = newEpMgr(name, comm)
    if err != nil {
        return nil, err
    }

    return impl,nil
}
