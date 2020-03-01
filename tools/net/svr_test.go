package net

import (
    "testing"
    "net"
    "sync"
    "time"
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
        WorkQueueCap: 1000,
        MaxConn: 1000,
        TCPReadBuffer: 128*1204,
        TCPWriteBuffer: 128*1024,
        TCPNoDelay: true,
    }

    log.SetFrameworkLevel(log.DEBUG)

    svr, err := NewSvr(cfg, &EchoHandle{})
    if err != nil {
        t.Fatalf("create svr err:%s", err)
    }

    stop := make(chan bool)
    go func() {
        svr.Run()
        stop <- true
    }()
    // 等待服务器启动成功
    time.Sleep(time.Second*2)

    var stopSvr sync.WaitGroup
    stopSvr.Add(101)

    for i := 0; i < 100; i++ {
        // 之所以开启携程是模拟10个客户端并发连接
        go func() {
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
                conn.Close()
            })
            stopSvr.Done()
        } ()
    }

    go func() {
        t.Run("closeconnection", func (t *testing.T) {
            addr, err := net.ResolveTCPAddr("tcp", ":8888")
            if err != nil {
                t.Fatalf("dial error:%s", err)
            }
            conn, err := net.DialTCP("tcp", nil, addr)
            if err != nil {
                t.Fatalf("dial error:%s", err)
            }

            conn.Close()
        })
        stopSvr.Done()
    }()

    stopSvr.Wait()
    // 等待服务器结束
    svr.Stop()
    <-stop
    log.FlushLogger()
}
