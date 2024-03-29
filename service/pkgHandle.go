package tex

import (
    "context"
    "encoding/binary"
    "github.com/yellia1989/tex-go/tools/net"
    "github.com/yellia1989/tex-go/tools/log"
    "github.com/yellia1989/tex-go/tools/sdp/codec"
    "github.com/yellia1989/tex-go/sdp/protocol"
)

type texSvrPkgHandle struct {
    name string
    service Service
    serviceImpl interface{}
}

const (
    // 最大包长度
    MAX_PACKET_SIZE = 100 * 1024 * 1024
)

func (h *texSvrPkgHandle) Parse(bytes []byte) (int, int) {
    if len(bytes) < 4 {
        return 0, net.PACKAGE_LESS
    }
    iHeaderLen := int(binary.BigEndian.Uint32(bytes[0:4]))
    if iHeaderLen < 4 || iHeaderLen > MAX_PACKET_SIZE {
        return 0, net.PACKAGE_ERROR
    }
    if len(bytes) < iHeaderLen {
        return 0, net.PACKAGE_LESS
    }
    return iHeaderLen, net.PACKAGE_FULL
}

func (h *texSvrPkgHandle) HandleRecv(ctx context.Context, pkg []byte, overload bool, queuetimeout bool) {
    current := net.ContextGetCurrent(ctx)
    req := &current.Request

    defer func() {
        err := recover()
        if err != nil {
            log.FErrorf("peer: %s:%d parse RequestPacket err:%s", current.IP, current.Port, err)
            current.Close()
        }
    }()
    codec.StringToSdp(pkg[4:], req)

    // 服务名称不匹配
    if h.name != req.SServiceName {
        log.FErrorf("handle dismatch servicename, peer: %s:%d, obj: %s,func: %s,reqid: %d", current.IP, current.Port, req.SServiceName, req.SFuncName, req.IRequestId)
        current.SendTexResponse(protocol.SDPSERVERNOSERVICEERR, nil)
        return
    }

    if queuetimeout {
        log.FErrorf("handle queuetimeout, peer: %s:%d, obj: %s,func: %s,reqid: %d", current.IP, current.Port, req.SServiceName, req.SFuncName, req.IRequestId)
        current.SendTexResponse(protocol.SDPSERVERQUEUETIMEOUT, nil)
        return
    }

    if overload {
        log.FErrorf("handle overload, peer: %s:%d, obj: %s,func: %s,reqid: %d", current.IP, current.Port, req.SServiceName, req.SFuncName, req.IRequestId)
        current.SendTexResponse(protocol.SDPSERVEROVERLOAD, nil)
        return
    }

    h.service.Dispatch(ctx, h.serviceImpl, req)
}
