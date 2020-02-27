package net

import (
    "testing"
    "time"
    _ "github.com/yellia1989/tex-go/tools/log"
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

func TestSvrRun(t *testing.T) {
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

    svr, err := NewSvr(cfg, &EchoHandle{})
    if err != nil {
        t.Fatalf("create svr err:%s", err)
    }
    stop := make(chan struct{})
    go func() {
        svr.Run()
        stop <- struct{}{}
    }()
    select {
        case <-time.After(time.Second * 10):
            svr.Stop()
    }

    <-stop
    //log.FlushLogger()
}
