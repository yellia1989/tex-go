package echo

import (
    "context"
    "github.com/yellia1989/tex-go/protocol/protocol"
    "github.com/yellia1989/tex-go/tools/sdp/codec"
    "github.com/yellia1989/tex-go/tools/net"
    "github.com/yellia1989/tex-go/tools/log"
)

type EchoService struct {
    name string
}

type echoServiceImpl interface {
    Hello(ctx context.Context, req string, resp *string) error
}

func (s *EchoService) Dispatch(ctx context.Context, serviceImpl interface{}, req *protocol.RequestPacket) {
    current := net.ContextGetCurrent(ctx)

    log.FDebugf("handle tex request, peer: %s:%d, obj: %s, func: %s, reqid: %d", current.IP, current.Port, req.SServiceName, req.SFuncName, req.IRequestId) 

    // 服务名称不匹配
    if s.name != req.SServiceName {
        current.SendTexResponse(protocol.SDPSERVERNOSERVICEERR, nil)
        return
    }

    ret := protocol.SDPSERVERUNKNOWNERR
    up := codec.NewUnPacker([]byte(req.SReqPayload))
    p := codec.NewPacker()

    switch req.SFuncName {
    case "Hello":
        impl := serviceImpl.(echoServiceImpl)

        var p1 string
        if err := up.ReadString(&p1, 0, true); err != nil {
            break
        }

        var p2 string
        if err := impl.Hello(ctx, p1, &p2); err != nil {
            break
        }

        if err := p.WriteString(0, p2); err != nil {
            break
        }

        ret = 0
    default:
        ret = protocol.SDPSERVERNOFUNCERR
    }

    if current.Rsp() {
        current.SendTexResponse(int32(ret), p.ToBytes())
    }
}
