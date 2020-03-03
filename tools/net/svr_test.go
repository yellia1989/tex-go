package net

import (
    "testing"
    "net"
    "sync"
    "time"
    "fmt"
    "github.com/yellia1989/tex-go/tools/log"
)

type EchoHandle struct {
}

func (s *EchoHandle) Parse(bytes []byte) (int,int) {
    return len(bytes),PACKAGE_FULL
}

func (s *EchoHandle) HandleRecv(pkg []byte) []byte {
    log.FDebugf("svr recv:%s", string(pkg))
    return pkg
}

func (s *EchoHandle) HandleTimeout(pkg []byte) []byte {
    return pkg
}

type EchoCli struct {
}

func (cli *EchoCli) Recv(pkg []byte) {
    log.Debugf("client recv:%s", pkg)
}

func (cli *EchoCli) Parse(pkg []byte) (int,int) {
    return len(pkg),PACKAGE_FULL
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
    stopSvr.Add(102)

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

    go func() {
        cfg := CliCfg{Proto:"tcp"}
        cli := NewCli(":8888", &cfg, &EchoCli{})

        defer func() {
            cli.Close()
            stopSvr.Done()
        }()

        // 每隔2秒钟发送一个hello
        // 发送5次断开连接,再次发送5次然后退出
        ticker := time.NewTicker(time.Second * 2)
        cnt := 0
        total := 11
        for {
            select {
            case <-ticker.C:
                msg := fmt.Sprintf("%d:hello", total)
                log.FDebugf("client send:%s", msg)
                cli.Send([]byte(msg))
                cnt += 1
                total--
                if total <= 0 {
                    return
                }
                if cnt > 5 {
                    cnt = 0
                    cli.SafeClose()
                }
            }
        }
    }()

    stopSvr.Wait()

    time.Sleep(time.Second * 10)

    // 等待服务器结束
    svr.Stop()
    <-stop
    log.FlushLogger()
}
