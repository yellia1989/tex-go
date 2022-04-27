package net

import (
    "sync"
    "sync/atomic"
    "time"
    "net"
    "io"
    "fmt"
    "context"
    "github.com/yellia1989/tex-go/tools/gpool"
    "github.com/yellia1989/tex-go/tools/log"
)

const (
    svr_write_queuecap = 10
    svr_work_queuetimeout = time.Millisecond*5
)

// 传输协议接口
type netHandle interface {
    Run()
}

type heartbeatHandle func()

// 服务器配置
type SvrCfg struct {
    Name string
    Proto string // tcp,udp
    Address string // listen address

    WorkThread int // 包处理协程个数
    WorkQueueCap int // 包处理队列长度
    WorkQueueTimeout   time.Duration

    MaxConn        int // 最大连接数
    WriteQueueCap  int // 每个连接的待发送队列的长度

    IdleTimeout    time.Duration // 每个连接的最长空闲时间

    TCPReadBuffer  int
    TCPWriteBuffer int
    TCPNoDelay     bool

    Heartbeat heartbeatHandle // 心跳函数
}

// 连接
type Conn struct {
    ID uint32 // 连接id
    IsTcp bool // 是否是tcp连接
    conn net.Conn // 连接fd
    svr *Svr // 服务器

    close int32 // 连接关闭
    writech chan []byte // 写通道
    done sync.WaitGroup // 等待连接读写关闭

    idleTime time.Time // 最后一次活跃时间点
}

func (c *Conn) Send(pkg []byte) (err error) {
    if len(pkg) == 0 {
        return fmt.Errorf("empty or nil pkg")
    }

    if c.isClose() {
        return fmt.Errorf("conn: %u has been closed", c.ID)
    }

    // writech的大小应该根据下面情况综合考虑来设置
    // 1. 发包频率
    // 2. 收包频率
    c.writech <- pkg
    return nil
}

func (c *Conn) SafeClose() {
    if c.isClose() {
        return
    }

    pkg := make([]byte, 0)
    c.writech <- pkg

    log.FDebugf("conn:%d safe close", c.ID)
}

func (c *Conn) isClose() bool {
    return atomic.LoadInt32(&c.close) == 1 
}

func (c *Conn) Close() {
    if !atomic.CompareAndSwapInt32(&c.close, 0, 1) {
        // 已经关闭了
        return
    }
    c.conn.Close()
    close(c.writech)
}

func (c *Conn) doRead() {
    defer func() {
        log.FDebugf("conn:%d close read", c.ID)
        c.Close()
        c.done.Done()
    }()

    c.idleTime = time.Now()
    tmpbuf := make([]byte, 1024*4)
    var pkgbuf []byte
    for !c.isClose() {
        if c.svr.cfg.IdleTimeout != 0 {
            if err := c.conn.SetReadDeadline(time.Now().Add(time.Millisecond*50)); err != nil {
                log.FErrorf("conn:%d set read timeout err:%s", err.Error())
                return
            }
        }
        n, err := c.conn.Read(tmpbuf)
        if err != nil {
            if c.isClose() {
                return
            }

            if isTimeoutErr(err) && c.svr.cfg.IdleTimeout != 0 {
                // 不活跃直接关闭
                if c.idleTime.Add(c.svr.cfg.IdleTimeout).Before(time.Now()) {
                    log.FDebugf("conn:%d is unactive, will be closed", c.ID)
                    return
                }
                continue
            }
            if (err == io.EOF) {
                log.FDebugf("conn:%d client closed connection:%s", c.ID, err.Error())
            } else {
                log.FErrorf("conn:%d read err:%s", c.ID, err.Error())
            }
            return
        }
        c.idleTime = time.Now()
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
                c.svr.recvPkg(c, pkg)
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

func (c *Conn) doWrite() {
    defer func() {
        log.FDebugf("conn:%d close write", c.ID)
        c.Close()
        c.done.Done()
    }()

    for {
        select {
        case pkg := <-c.writech :
            total := len(pkg)
            if total == 0 {
                // 优雅关闭
                return
            }
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
    pkgHandle SvrPkgHandle // 包处理
    isclose int32 // 服务器是否关闭
    close chan struct{}

    netHandle netHandle // 网络字节流处理
    workPool *gpool.Pool // 工作线程

    queueSize int32 // 工作队列长度

    mu sync.Mutex
    id uint32 // conn auto incr id
    connNum uint32 // 当前连接数
    conns map[uint32]*Conn //网络连接
}

func NewSvr(cfg *SvrCfg, pkgHandle SvrPkgHandle) *Svr {
    if cfg.WriteQueueCap <= svr_write_queuecap {
        cfg.WriteQueueCap = svr_write_queuecap
    }
    if cfg.WorkQueueTimeout <= svr_work_queuetimeout {
        cfg.WorkQueueTimeout = svr_work_queuetimeout
    }

    s := &Svr{cfg: cfg, pkgHandle: pkgHandle}
    s.conns = make(map[uint32]*Conn)
    s.close = make(chan struct{})

    if s.cfg.Proto == "tcp" {
        s.netHandle = &tcpHandle{svr: s}
    } else if s.cfg.Proto == "udp" {
        s.netHandle = &udpHandle{svr: s}
    } else {
        panic("unsupport proto:" + s.cfg.Proto)
    }

    return s
}

func (s *Svr) isClose() bool {
    return atomic.LoadInt32(&s.isclose) == 1
}

func (s *Svr) Run() {
    defer func() {
        log.FDebugf("service:%s stopped", s.cfg.Name)
    }()

    log.FDebugf("service:%s starting", s.cfg.Name)

    // 开启工作协程
    s.workPool = gpool.NewPool(s.cfg.WorkThread, s.cfg.WorkQueueCap)
    // 停止工作协程
    defer func() {
        s.workPool.Release()
        log.FDebugf("service:%s work threads stop", s.cfg.Name)
    }()
    log.FDebugf("service:%s work threads=%d cap=%d start", s.cfg.Name, s.cfg.WorkThread, s.cfg.WorkQueueCap)

    var heartbeat chan struct{}
    if s.cfg.Heartbeat != nil {
        heartbeat = make(chan struct{})
        go func() {
            ticker := time.NewTicker(time.Second * 3)
            defer func() {
                ticker.Stop()
                log.FDebugf("service:%s heartbeat thread stop", s.cfg.Name)
            }()
            log.FDebugf("service:%s heartbeat thread start", s.cfg.Name)
            for {
                select {
                case <-ticker.C:
                    s.workPool.JobQueue <- gpool.Job(s.cfg.Heartbeat)
                case <-s.close:
                    heartbeat <- struct{}{}
                    return
                }
            }
        }()
    }

    network := make(chan struct{})
    go func () {
        defer func() {
            log.FDebugf("service:%s net thread stop", s.cfg.Name)
        }()
        log.FDebugf("service:%s net thread start", s.cfg.Name)
        // 开启网络监听
        s.netHandle.Run()
        network <- struct{}{}
    }()
    <-network

    if heartbeat != nil {
        <-heartbeat
    }
}

func (s *Svr) Stop() {
    log.FDebugf("service:%s stop...", s.cfg.Name)

    if !atomic.CompareAndSwapInt32(&s.isclose, 0, 1) {
        return
    }

    // 等待所有连接关闭
    log.FDebugf("service:%s, begin to stop all connection", s.cfg.Name)
    stop := make(chan struct{})
    go func() {
        ticker := time.NewTicker(time.Millisecond * 100)
        defer ticker.Stop()
        for {
            select {
            case <-ticker.C:
                s.mu.Lock()
                defer s.mu.Unlock()
                if len(s.conns) == 0 {
                    stop <- struct{}{}
                    return
                }
                for _,conn := range s.conns {
                    conn.SafeClose()
                }
            }
        }
    }()
    <-stop
    log.FDebugf("service:%s, all connection has been stopped", s.cfg.Name)

    close(s.close)
}

func (s *Svr) delConnection(id uint32) {
    s.mu.Lock()
    delete(s.conns, id)
    s.mu.Unlock()
    log.FDebugf("delete conn:%d", id)
}

func (s *Svr) closeConnection(id uint32) {
    s.mu.Lock()
    defer s.mu.Unlock()
    conn, ok := s.conns[id]
    if !ok {
        return
    }
    conn.SafeClose()
}

func (s *Svr) addConnection(c net.Conn) {
    if s.isClose() {
        c.Close()
        return
    }

    s.mu.Lock()
    connNum := len(s.conns)
    if connNum >= s.cfg.MaxConn {
        s.mu.Unlock()
        // 超过了最大连接数,直接关闭连接
        c.Close()
        log.FErrorf("exceed max conn:%d, cur:%d", s.cfg.MaxConn, connNum)
        return
    }

    s.id++
    conn := &Conn{ID: s.id, conn: c, svr: s}

    _, conn.IsTcp = c.(*net.TCPConn)

    // writech的大小决定了conn调用write时是否阻塞
    conn.writech = make(chan []byte, s.cfg.WriteQueueCap)

    s.conns[conn.ID] = conn
    s.mu.Unlock()

    // 等待连接关闭
    conn.done.Add(2)
    go func() {
        conn.done.Wait()
        s.delConnection(conn.ID)
    }()

    // 开启读写协程
    go conn.doRead()
    go conn.doWrite()

    log.FDebugf("accept conn:%d remote addr:%s", conn.ID, c.RemoteAddr())
}

func (s *Svr) Send(id uint32, pkg []byte) {
    s.mu.Lock()
    conn, ok := s.conns[id]
    if !ok {
        s.mu.Unlock()
        return 
    }
    s.mu.Unlock()

    conn.Send(pkg)
}

func (s *Svr) SendToAll(pkg []byte) {
    s.mu.Lock()
    conns := make([]*Conn, len(s.conns))
    for _, v := range s.conns {
        conns = append(conns, v)
    }
    s.mu.Unlock()

    for _, v := range conns {
        v.Send(pkg)
    }
}

func (s *Svr) recvPkg(c *Conn, pkg []byte) {
    overload := false
    queueSize := atomic.LoadInt32(&s.queueSize)
    if queueSize > int32(s.cfg.WorkQueueCap) {
        // 超过服务器负载直接丢弃
        return
    }
    atomic.AddInt32(&s.queueSize, 1)
    if queueSize > int32(s.cfg.WorkQueueCap)/2 {
        overload = true
    }

    ctx := contextWithCurrent(context.Background(), c)
    recvTime := time.Now()
    handler := func() {
        defer func() {
            atomic.AddInt32(&s.queueSize, -1)
        }()
        timeout := recvTime.Add(s.cfg.WorkQueueTimeout).Before(time.Now())
        s.pkgHandle.HandleRecv(ctx, pkg, overload, timeout)
    }

    c.svr.workPool.JobQueue <- handler
}
