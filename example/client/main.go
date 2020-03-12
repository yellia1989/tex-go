package main

import (
    "sync"
    "time"
    "fmt"
    "os"
    "syscall"
    "os/signal"
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
    comm.StringToProxy("test.EchoServer.EchoServiceObj@tcp -h 127.0.0.1 -p 9000 -t 3600000", echoPrx)

    closecli := make(chan bool)

    c := make(chan os.Signal)
    signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-c
        close(closecli)
    }()

    var done sync.WaitGroup
    num := 1

    for i := 0; i < num; i++ {
        done.Add(1)
        go func(id int) {
            log.Debugf("client:%d start", id)
            loop := time.NewTicker(time.Second)
            run := true
            for run {
                select {
                case <-closecli:
                    run = false
                    break
                case <-loop.C:
                    var resp string
                    err := echoPrx.Hello(fmt.Sprintf("client:%d, yellia", id), &resp)
                    if err != nil {
                        log.Errorf("err:%s", err)
                        loop.Stop()
                    } else {
                        log.Debugf("resp:%s", resp)
                    }
                }
            }
            loop.Stop()
            done.Done()
            log.Debugf("client:%d stop", id)
        }(i+1)
    }

    log.Debug("wait client to exit")
    done.Wait()
}
