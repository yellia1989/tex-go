package net

import (
    "sync"
    "sync/atomic"
    "time"
    "net"
    "io"
    "fmt"
    "strconv"
    "github.com/yellia1989/tex-go/tools/rtimer"
    "github.com/yellia1989/tex-go/tools/gpool"
    "github.com/yellia1989/tex-go/tools/log"
)

const (
    PACKAGE_LESS = iota
    PACKAGE_FULL
    PACKAGE_ERROR
)

const (
    write_timeout = time.Millisecond * 10
    write_queuecap = 10
)

// 服务器接收到数据包的处理接口
type PackageHandle interface {
    // 将二进制流按照特定的协议解析成单个的包
    Parse(bytes []byte) (int,int)
    // 单个数据包正常处理
    HandleRecv(pkg []byte) []byte
    // 数据包超时处理
    HandleTimeout(pkg []byte) []byte
}

// 传输协议接口
type netHandle interface {
    Run()
}

// 服务器配置
type SvrCfg struct {
    Proto string // tcp,udp
    Address string // listen address

    WorkThread int // 包处理协程个数
    WorkQueueCap int // 包处理队列长度
    MaxConn        int // 最大连接数
    WriteQueueCap  int // 每个连接的待发送队列的长度
    WriteTimeout   time.Duration

    HandleTimeout  time.Duration
    IdleTimeout    time.Duration

    TCPReadBuffer  int
    TCPWriteBuffer int
    TCPNoDelay     bool
}

// 连接
type Conn struct {
    ID uint32 // 连接id
    IsTcp bool // 是否是tcp连接
    conn net.Conn // 连接fd
    svr *Svr // 服务器

    close bool // 连接关闭
    writech chan []byte // 写通道
    done sync.WaitGroup // 等待连接读写关闭

    idleTime time.Time // 最后一次活跃时间点
}

func (c *Conn) Send(pkg []byte) error {
    if len(pkg) == 0 {
        return fmt.Errorf("empty packet")
    }

    if c.close {
        return fmt.Errorf("conn has been closed")
    }

    // writech的大小应该根据下面情况综合考虑来设置
    // 1. 发包频率
    // 2. 收包频率
    done := make(chan struct{})
    go func() {
        defer func() {
            if err := recover(); err != nil {
                // 防止往已经关闭的ch写数据导致的crash
                done <- struct{}{}
            }
        }()
        c.writech <- pkg
        done <- struct{}{}
    }()

    select {
        case <-rtimer.After(c.svr.cfg.WriteTimeout):
            return fmt.Errorf("send packet timeout")
        case <-done:
    }

    return nil
}

func (c *Conn) SafeClose() {
    close(c.writech)
}

func (c *Conn) Close() {
    if c.close {
        return
    }
    c.close = true
    c.conn.Close()
    close(c.writech)
}

func (c *Conn) doRead() {
    defer func() {
        c.Close()
        c.done.Done()
    }()

    c.idleTime = time.Now()
    tmpbuf := make([]byte, 1024*4)
    var pkgbuf []byte
    for !c.svr.close {
        if c.svr.cfg.IdleTimeout != 0 {
            if err := c.conn.SetReadDeadline(time.Now().Add(time.Millisecond*500)); err != nil {
                log.FErrorf("conn:%d set read timeout err:%s", err.Error())
                return
            }
        }
        n, err := c.conn.Read(tmpbuf)
        if err != nil {
            if isTimeoutErr(err) {
                // 不活跃直接关闭
                if c.idleTime.Add(c.svr.cfg.IdleTimeout).Before(time.Now()) {
                    log.FDebugf("conn:%d is unactive, will be closed", c.ID)
                    return
                }
                c.idleTime = time.Now()
                continue
            }
            if (err == io.EOF) {
                log.FDebugf("conn:%d client closed connection:%s", c.ID, err.Error())
            } else {
                log.FErrorf("conn:%d read err:%s", c.ID, err.Error())
            }
            return
        }
        // 解析包
        pkgbuf = append(pkgbuf, tmpbuf[:n]...)
        for {
            pkglen, status := c.svr.pkgHandle.Parse(pkgbuf)
            if status == PACKAGE_LESS {
                break
            }
            if status == PACKAGE_FULL {
                pkg := make([]byte, pkglen)
                copy(pkg, pkgbuf[:pkglen])
                c.recvPkg(pkg)
                pkgbuf = pkgbuf[pkglen:]
                if len(pkgbuf) > 0 {
                    continue
                }
                pkgbuf = nil
                break
            }
            log.FErrorf("conn:%d parse package error", c.ID)
            return
        }
    }
}

func (c *Conn) recvPkg(pkg []byte) {
    now := time.Now()
    handler := func() {
        rsp := c.svr.recvPkg(now, pkg)
        c.Send(rsp)
    }

    c.svr.workPool.JobQueue <- handler
}

func (c *Conn) doWrite() {
    defer func() {
        c.Close()
        c.done.Done()
    }()

    for {
        select {
        case pkg, ok := <-c.writech :
            if !ok {
                // 优雅关闭
                return
            }
            total := len(pkg)
            for {
                n, err := c.conn.Write(pkg)
                if err != nil {
                    log.FErrorf("conn:%d write err:%s", c.ID, err.Error())
                    return
                }
                if n > 0 {
                    total -= n
                }
                if total == 0 {
                    break
                }
                pkg = pkg[n:]
            }
        }
    }
}

// 服务器
type Svr struct {
    cfg *SvrCfg // 配置
    pkgHandle PackageHandle // 包处理
    close bool // 服务器是否关闭

    netHandle netHandle // 网络字节流处理
    conns sync.Map //[uint32]*Conn 网络连接
    id uint32 // conn auto incr id

    workPool *gpool.Pool // 工作线程

    connNum int32 // 当前连接数
}

func NewSvr(cfg *SvrCfg, pkgHandle PackageHandle) (*Svr, error) {
    if cfg.WriteQueueCap <= write_queuecap {
        cfg.WriteQueueCap = write_queuecap
    }
    if cfg.WriteTimeout <= write_timeout {
        cfg.WriteTimeout = write_timeout
    }

    s := &Svr{cfg: cfg, pkgHandle: pkgHandle, close: false}

    if s.cfg.Proto == "tcp" {
        s.netHandle = &tcpHandle{svr: s}
    } else if s.cfg.Proto == "udp" {
        s.netHandle = &udpHandle{svr: s}
    } else {
        panic("unsupport proto:" + s.cfg.Proto)
    }

    return s, nil
}

func (s *Svr) Run() {
    log.FDebug("start server")

    // 开启工作协程
    s.workPool = gpool.NewPool(s.cfg.WorkThread, s.cfg.WorkQueueCap)

    network := make(chan struct{})
    go func () {
        // 开启网络监听
        s.netHandle.Run()
        network <- struct{}{}
    }()
    <-network

    // 停止工作协程
    s.workPool.Release()

    log.FDebug("svr stop")
}

func (s *Svr) Stop() {
    s.close = true
}

func (s *Svr) delConnection(id uint32) {
    s.conns.Delete(id)
    atomic.AddInt32(&s.connNum, -1)
    log.FDebugf("conn:%d is deleted", id)
}

func (s *Svr) CloseConnection(id uint32) {
    conn, ok := s.conns.Load(id)
    if !ok {
        return
    }
    conn.(*Conn).SafeClose()
}

func (s *Svr) addConnection(c net.Conn) {

    if s.connNum >= int32(s.cfg.MaxConn) {
        // 超过了最大连接数,直接关闭连接
        c.Close()
        log.FErrorf("exceed max conn:%d, cur:%d", s.cfg.MaxConn, s.connNum)
        return
    }

    id := atomic.AddUint32(&s.id, 1)
    conn := &Conn{ID: id, conn: c, close: false, svr: s}

    _, conn.IsTcp = c.(*net.TCPConn)

    // writech的大小决定了conn调用write时是否阻塞
    conn.writech = make(chan []byte, s.cfg.WriteQueueCap)

    _, ok := s.conns.LoadOrStore(id, conn)
    if ok {
        panic("add new conn failed, id:" + strconv.Itoa(int(id)))
    }
    atomic.AddInt32(&s.connNum, 1)

    // 等待连接关闭
    conn.done.Add(2)
    go func() {
        conn.done.Wait()
        s.delConnection(id)
    }()

    // 开启读写协程
    go conn.doRead()
    go conn.doWrite()

    log.FDebugf("accept conn:%d remote addr:%s", conn.ID, c.RemoteAddr())
}

func (s *Svr) Send(id uint32, pkg []byte) {
    v, ok := s.conns.Load(id)
    if !ok {
        return 
    }
    conn := v.(*Conn)
    conn.Send(pkg)
}

func (s *Svr) SendToAll(pkg []byte) {
    s.conns.Range(func (k, v interface{}) bool {
        v.(*Conn).Send(pkg)
        return true
    })
}

func (s *Svr) recvPkg(recvTime time.Time, pkg []byte) []byte {
    cfg := s.cfg
    var rsp []byte
    if cfg.HandleTimeout == 0 {
        rsp = s.pkgHandle.HandleRecv(pkg)
    } else {
        done := make(chan struct{})
        go func() {
            rsp = s.pkgHandle.HandleRecv(pkg)
            done <- struct{}{}
        }()
        endtime := recvTime.Add(cfg.HandleTimeout)
        select {
        case <-rtimer.After(endtime.Sub(time.Now())):
            rsp = s.pkgHandle.HandleTimeout(pkg)
        case <-done:
        }
    }
    return rsp
}

