package net

import (
    "context"
    "strings"
    "strconv"
    "time"
    "sync/atomic"
    "encoding/binary"
    "github.com/yellia1989/tex-go/sdp/protocol"
    "github.com/yellia1989/tex-go/tools/sdp/codec"
)

type Current struct {
    ID uint32   // 连接id
    IP string   // 连接ip
    Port int    // 连接port
    rsp int32    // 是否立即响应当前请求
    svr *Svr    // 服务器
    Request protocol.RequestPacket
    start time.Time
}

func (c *Current) SendResponse(pkg []byte) {
    c.svr.Send(c.ID, pkg)
}

func (c *Current) SendTexResponse(ret int32, pkg []byte) {
    if c.Request.BIsOneWay {
        return
    }

    resp := protocol.ResponsePacket{}
    resp.ResetDefault()
    resp.IRet = ret
    resp.IRequestId = c.Request.IRequestId
    resp.SRspPayload = string(pkg)

    b1 := codec.SdpToString(&resp)

    total := len(b1)+4
    b2 := make([]byte, total)
    binary.BigEndian.PutUint32(b2, uint32(total))
    copy(b2[4:], b1)

    c.svr.Send(c.ID, b2)
}

func (c *Current) Close() {
    c.svr.closeConnection(c.ID)
}

func (c *Current) AsyncRsp() {
    atomic.CompareAndSwapInt32(&c.rsp, 1, 0)
}

func (c *Current) Rsp() bool {
    return atomic.LoadInt32(&c.rsp) == 1
}

func contextWithCurrent(ctx context.Context, conn *Conn) context.Context {
    current := &Current{ID: conn.ID, rsp: 1}
    addr := conn.conn.RemoteAddr().String()
    tmp := strings.SplitN(addr, ":", 2)
    current.IP = tmp[0]
    current.Port,_ = strconv.Atoi(tmp[1])
    current.svr = conn.svr
    current.start = time.Now()
    return context.WithValue(ctx, 0x484900, current)
}

func ContextGetCurrent(ctx context.Context) *Current {
    current, ok := ctx.Value(0x484900).(*Current)
    if !ok {
        panic("failed to get current from context")
    }
    return current
}
