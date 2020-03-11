package main

import (
    tex "github.com/yellia1989/tex-go"
    "github.com/yellia1989/tex-go/tools/log"
    "github.com/yellia1989/tex-go/example/server/echo"
)

func main() {
    comm := tex.NewCommunicator()

    log.SetFrameworkLevel(log.DEBUG)

    defer func() {
        comm.Close()
        log.FlushLogger()
    }()

    echoPrx := new(echo.EchoService)
    comm.StringToProxy("test.EchoServer.EchoServiceObj@tcp -h 127.0.0.1 -p 8080 -t 3600000", echoPrx)

    var resp string
    err := echoPrx.Hello("yellia", &resp)
    if err != nil {
        log.Errorf("err:%s", err)
    } else {
        log.Debugf("resp:%s\n", resp)
    }
}
