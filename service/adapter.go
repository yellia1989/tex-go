package tex

import (
    "sync"
    "sync/atomic"
    "strconv"
    "fmt"
    "time"
    "errors"
    "encoding/binary"
    "github.com/yellia1989/tex-go/sdp/protocol"
    "github.com/yellia1989/tex-go/tools/sdp/codec"
    "github.com/yellia1989/tex-go/tools/net"
    "github.com/yellia1989/tex-go/tools/log"
    "github.com/yellia1989/tex-go/tools/rtimer"
)

const (
    adapterActiveInterval = 10 * time.Second // 每10秒检查一下活跃连接
    adapterConsfail = 5 // 最大持续接受消息失败次数
    adapterMinfail = 2 // 最小接受小时失败次数
    adapterFailpation = 50/100 // 消息接受失败率
    adapterTrytime = 30 * time.Second // 对于不活跃的连接30秒重连一次
)

type adapterProxy struct {
    ep *Endpoint
    cli *net.Cli
    done chan bool // 关闭通道

    mu sync.Mutex
    reqQueueLen int // 请求队列长度

    sendTotal uint32   // 请求总数
    failTotal uint32   // 请求失败总数
    consfailTotal uint32 // 持续失败总数
    active int32 // 是否可用标志
    nextTryTime int64 // active=0时下一次重连时间
    connfailed int32 // 是否是连接失败

    req sync.Map // 请求队列
}

func (adapter *adapterProxy) invoke(req *protocol.RequestPacket, resp **protocol.ResponsePacket) error {
    // 请求队列已经达到最大值,直接报错
    mu := &adapter.mu
    mu.Lock()
    if adapter.reqQueueLen > cliCfg.adapterSendQueueCap {
        mu.Unlock()
        return fmt.Errorf("adapter req queue full")
    }

    adapter.reqQueueLen += 1
    mu.Unlock()

    ch := make(chan *protocol.ResponsePacket)
    adapter.req.Store(req.IRequestId, ch)
    defer func() {
        mu.Lock()
        adapter.reqQueueLen -= 1
        mu.Unlock()

        adapter.req.Delete(req.IRequestId)
        close(ch)
    }()

    begintime := time.Now()
    if err := adapter.send(req); err != nil {
        // 这种情况一般是连接失败,因为net.Cli.Send方法自带重连功能,
        // 只有连接失败的时候才返回错误
        // 注意:
        // send发送是阻塞的，所以当服务不可用时会有大量的消息阻塞在这,
        // 直到达到cliCfg配置的消息上限
        return fmt.Errorf("connect err:%s", err.Error())
    }

    select {
    case <-rtimer.After(time.Duration(req.ITimeout) * time.Millisecond):
        // 请求超时，取消请求
        atomic.AddUint32(&adapter.failTotal, 1)
        atomic.AddUint32(&adapter.consfailTotal, 1)
        log.FErrorf("wait for response timeout, reqid:%d, adapter:%s", req.IRequestId, adapter.ep)
        return errors.New("req timeout")
    case resp2 := <-ch:
        // 收到回复
        atomic.StoreUint32(&adapter.consfailTotal, 0)
        log.FDebugf("got response, ret:%d, cost:%d ms, reqid:%d, adapter:%s", resp2.IRet, time.Since(begintime).Milliseconds(), req.IRequestId, adapter.ep)
        if resp2.IRet == protocol.SDPSERVERSUCCESS {
            *resp = resp2
        } else {
            return fmt.Errorf("remote server err, ret:%d", resp2.IRet)
        }
    }

    return nil
}

func (adapter *adapterProxy) send(req *protocol.RequestPacket) error {
    b1 := codec.SdpToString(req)

    total := len(b1)+4
    b2 := make([]byte, total)
    binary.BigEndian.PutUint32(b2, uint32(total))
    copy(b2[4:], b1)

    if err := adapter.cli.Send(b2); err != nil {
        // 这种情况一般是服务器连接不上，直接设为不可用
        adapter.setInactive(1)
        return err
    }

    // 只统计发送成功的
    atomic.AddUint32(&adapter.sendTotal, 1)

    log.FDebugf("push request(send), reqid:%d, adapter:%s", req.IRequestId, adapter.ep)

    return nil
}

func (adapter *adapterProxy) isActive(reconnect bool) bool {
    active := atomic.LoadInt32(&adapter.active)
    if active == 1 {
        return true
    }

    nextTryTime := atomic.LoadInt64(&adapter.nextTryTime)
    if reconnect || time.Now().Before(time.Unix(nextTryTime, 0)) {
        // 重连时间到了
        atomic.StoreInt32(&adapter.active, 1)
        return true
    }

    return false
}

func (adapter *adapterProxy) Parse(bytes []byte) (int,int) {
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

func (adapter *adapterProxy) Recv(pkg []byte) {
    defer func() {
        err := recover()
        if err != nil {
            log.FErrorf("parse ResponsePacket err:%v, adapter:%s", err, adapter.ep)
        }
    }()
    resp := &protocol.ResponsePacket{}
    codec.StringToSdp(pkg[4:], resp)

    req, ok := adapter.req.Load(resp.IRequestId)
    if !ok {
        // 请求超时了，直接丢弃
        log.FErrorf("got response but timeout, reqid:%d, adapter:%s", resp.IRequestId, adapter.ep)
        return
    }

    ch := req.(chan *protocol.ResponsePacket)
    ch <- resp
}

func (adapter *adapterProxy) checkActive() {
    loop := time.NewTicker(adapterActiveInterval)
    for {
        select {
        case <-adapter.done:
            break
        case <-loop.C:
            // 持续失败达到一定次数后强制关闭连接
            consfailTotal := atomic.LoadUint32(&adapter.consfailTotal)
            if consfailTotal >= adapterConsfail {
                adapter.setInactive(0)
                log.FDebugf("disable connection(cont fail), continuous fail:%d, adapter:%s", consfailTotal, adapter.ep)
                continue
            }
            // 调用次数总量达到一定失败比例后强制关闭连接
            failTotal := atomic.LoadUint32(&adapter.failTotal)
            sendTotal := atomic.LoadUint32(&adapter.sendTotal)
            if failTotal >= adapterMinfail && (failTotal / sendTotal >= adapterFailpation) {
                adapter.setInactive(0)
                log.FDebugf("disable connection(statbility), fail:%d, total:%d, adapter:%s", failTotal, sendTotal, adapter.ep)
                continue
            }
        }
    }
    loop.Stop()
}

func (adapter *adapterProxy) close() {
    adapter.setInactive(0)
    adapter.done <- true
}

func (adapter *adapterProxy) setInactive(connfailed int32) {
    atomic.StoreUint32(&adapter.sendTotal, 0)
    atomic.StoreUint32(&adapter.failTotal, 0)
    atomic.StoreUint32(&adapter.consfailTotal, 0)
    atomic.StoreInt32(&adapter.active, 0)
    atomic.StoreInt64(&adapter.nextTryTime, time.Now().Add(adapterTrytime).Unix())
    atomic.StoreInt32(&adapter.connfailed, connfailed)
    adapter.cli.Close()
    log.FDebugf("close adapter:%s", adapter.ep)
}

func newAdapter(ep *Endpoint) (*adapterProxy, error) {
    adapter := &adapterProxy{active: 1, ep: ep, done: make(chan bool)}

    address := ep.IP + ":" + strconv.Itoa(ep.Port)
    adapter.cli = net.NewCli(address,
            &net.CliCfg{
                Proto: ep.Proto,
                WriteQueueCap: cliCfg.adapterSendQueueCap,
                IdleTimeout: cliCfg.adapterIdleTimeout,
            }, adapter)

    go adapter.checkActive()

    return adapter, nil
}
