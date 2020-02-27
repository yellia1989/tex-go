package net

import (
    "sync"
    "sync/atomic"
    "time"
    "net"
    "io"
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
    MaxConn        int

    HandleTimeout  time.Duration

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
}

func (c *Conn) Send(pkg []byte) {
    if c.close || len(pkg) == 0 {
        return
    }
    c.writech <- pkg
}

func (c *Conn) SafeClose() {
    pkg := make([]byte,0)
    c.writech <- pkg
}

func (c *Conn) Close() {
    if c.close {
        return
    }
    if err := c.conn.Close(); err != nil {
        return
    }

    c.svr.delConnection(c.ID)
    c.close = true
}

func (c *Conn) doRead() {
    defer c.Close()

    tmpbuf := make([]byte, 1024*4)
    var pkgbuf []byte
    for !c.svr.close {
        n, err := c.conn.Read(tmpbuf)
        if err != nil {
            if (err == io.EOF) {
                log.Debugf("conn:%d client closed connection:%s", c.ID, err.Error())
            } else {
                log.Errorf("conn:%d read err:%s", c.ID, err.Error())
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
                pkg := pkgbuf[:pkglen]
                c.recvPkg(pkg)
                pkgbuf = pkgbuf[pkglen:]
                if len(pkgbuf) > 0 {
                    continue
                }
                pkgbuf = nil
                break
            }
            log.Errorf("conn:%d parse package error", c.ID)
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
    defer c.Close()

    for {
        select {
        case pkg := <-c.writech :
            if c.svr.close {
                // 未发完的包丢弃
                return
            }
            total := len(pkg)
            if total == 0 {
                // 优雅关闭
                return
            }
            for {
                n, err := c.conn.Write(pkg)
                if err != nil {
                    log.Errorf("conn:%d write err:%s", c.ID, err.Error())
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
}

func NewSvr(cfg *SvrCfg, pkgHandle PackageHandle) (*Svr, error) {
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
    log.Debug("start server")

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

    log.Debug("svr stop")
}

func (s *Svr) Stop() {
    s.close = true
}

func (s *Svr) delConnection(id uint32) {
    s.conns.Delete(id)
    log.Debugf("conn:%d is closed", id)
}

func (s *Svr) CloseConnection(id uint32) {
    conn, ok := s.conns.Load(id)
    if !ok {
        return
    }
    conn.(*Conn).SafeClose()
}

func (s *Svr) addConnection(c net.Conn) {
    id := atomic.AddUint32(&s.id, 1)
    conn := &Conn{ID: id, conn: c, close: false, svr: s}

    _, conn.IsTcp = c.(*net.TCPConn)

    // writech的大小决定了conn调用write时是否阻塞
    conn.writech = make(chan []byte, 10)

    _, ok := s.conns.LoadOrStore(id, conn)
    if ok {
        panic("add new conn failed, id:" + strconv.Itoa(int(id)))
    }

    // 开启读写协程
    go conn.doRead()
    go conn.doWrite()

    log.Debugf("accept conn:%d remote addr:%s", conn.ID, c.RemoteAddr())
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

