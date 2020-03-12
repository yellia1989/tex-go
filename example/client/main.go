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
    log.SetFrameworkLevel(log.INFO)

    defer func() {
        log.FlushLogger()
    }()

    closecli := make(chan bool)

    c := make(chan os.Signal)
    signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-c
        close(closecli)
    }()

    var done sync.WaitGroup
    num := 10000

    for i := 0; i < num; i++ {
        done.Add(1)
        go func(id int) {
            log.Debugf("client:%d start", id)
            comm := tex.NewCommunicator()
            defer func() {
                comm.Close()
            }()

            echoPrx := new(echo.EchoService)
            comm.StringToProxy("test.EchoServer.EchoServiceObj@tcp -h 127.0.0.1 -p 9000 -t 3600000", echoPrx)
            loop := time.NewTicker(time.Millisecond * 1000)
            var total time.Duration
            invokenum := float64(0)
            run := true
            for run {
                select {
                case <-closecli:
                    run = false
                    break
                case <-loop.C:
                    var resp string
                    start := time.Now()
                    err := echoPrx.Hello(fmt.Sprintf("client:%d, yellia", id), &resp)
                    total += time.Since(start)
                    if err != nil {
                        log.Errorf("err:%s", err)
                    } else {
                        invokenum++
                    }
                }
            }
            cost := float64(total.Milliseconds())
            if invokenum != 0 {
                log.Debugf("client:%d stop, cost:%.2f ms, invokenum:%.2f, ops:%.2f ms", id, cost, invokenum, cost/invokenum)
            }
            loop.Stop()
            done.Done()
        }(i+1)
    }

    done.Wait()
    log.Debug("wait client to exit")
}
