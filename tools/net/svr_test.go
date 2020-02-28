package net

import (
    "testing"
    "time"
    "net"
    "github.com/yellia1989/tex-go/tools/log"
)

type EchoHandle struct {
}

func (s *EchoHandle) Parse(bytes []byte) (int,int) {
    return len(bytes),PACKAGE_FULL
}

func (s *EchoHandle) HandleRecv(pkg []byte) []byte {
    return pkg
}

func (s *EchoHandle) HandleTimeout(pkg []byte) []byte {
    return pkg
}

func TestSvr(t *testing.T) {
    cfg := &SvrCfg{
        Proto: "tcp",
        Address: ":8888",
        WorkThread: 1,
        WorkQueueCap: 10000,
        MaxConn: 10000,
        HandleTimeout: 0 * time.Second,
        TCPReadBuffer: 128*1024*1204,
        TCPWriteBuffer: 128*1024*1024,
        TCPNoDelay: true,
    }

    log.SetFrameworkLevel(log.DEBUG)

    svr, err := NewSvr(cfg, &EchoHandle{})
    if err != nil {
        t.Fatalf("create svr err:%s", err)
    }

    start := make(chan struct{})
    stop := make(chan struct{})
    go func() {
        start <- struct{}{}
        svr.Run()
        stop <- struct{}{}
    }()
    // 等待服务器启动成功
    <-start

    for i := 0; i < 100; i++ {
        t.Run("accept new connection", func (t *testing.T) {
            addr, err := net.ResolveTCPAddr("tcp", ":8888")
            if err != nil {
                t.Fatalf("dial error:%s", err)
            }
            conn, err := net.DialTCP("tcp", nil, addr)
            if err != nil {
                t.Fatalf("dial error:%s", err)
            }

            hello := []byte("hello")
            n, err := conn.Write(hello)
            if err != nil {
                t.Fatalf("write error:%s", err)
            }

            buf := make([]byte, n)
            n2, err := conn.Read(buf)
            if err != nil {
                t.Fatalf("read error:%s", err)
            }
            if n != n2 || string(buf) != string(hello) {
                t.Fatalf("write:%s vs read:%s", string(hello), string(buf))
            }
        })
    }

    // 等待服务器结束
    svr.Stop()
    <-stop
    log.FlushLogger()
}
