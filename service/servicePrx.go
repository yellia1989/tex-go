package tex

import (
    "time"
    "strings"
    "sync/atomic"
    "github.com/yellia1989/tex-go/sdp/protocol"
    "github.com/yellia1989/tex-go/tools/log"
    "github.com/yellia1989/tex-go/service/model"
)

type ServicePrx interface {
    SetPrxImpl(impl model.ServicePrxImpl)
}

type servicePrxImpl struct {
    name string
    comm *Communicator
    epmgr *endpointManager

    invokeTimeout uint32

    reqid uint32
}

func (impl *servicePrxImpl) Invoke(sFuncName string, params []byte, resp **protocol.ResponsePacket) error {
    // 构造请求消息
    req := &protocol.RequestPacket{
        IRequestId: atomic.AddUint32(&impl.reqid, 1),
        SServiceName: impl.name,
        SFuncName: sFuncName,
        SReqPayload: string(params),
        ITimeout: impl.invokeTimeout,
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
    impl.invokeTimeout = uint32(timeout.Milliseconds())
}

func (impl *servicePrxImpl) Close() {
    impl.epmgr.Close()
}

func newPrxImpl(name string, comm *Communicator) (*servicePrxImpl) {
    serviceName := name
    if p := strings.Index(name, "%"); p != -1 {
        // 去掉分区
        serviceName = name[:p]
    }
    if p := strings.Index(serviceName, "@"); p != -1 {
        // 去掉endpoint
        serviceName = serviceName[:p]
    }

    return &servicePrxImpl{name: serviceName, comm: comm, invokeTimeout: uint32(cliCfg.invokeTimeout.Milliseconds()), epmgr: newEpMgr(name, comm)}
}
