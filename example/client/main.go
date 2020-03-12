package main

import (
    "sync"
    "time"
    "os"
    "syscall"
    "os/signal"
    tex "github.com/yellia1989/tex-go"
    "github.com/yellia1989/tex-go/tools/log"
    "github.com/yellia1989/tex-go/example/server/echo"
)

var done sync.WaitGroup

func main() {
    log.SetFrameworkLevel(log.INFO)

    defer func() {
        log.FlushLogger()
    }()

    closecli := make(chan bool)

    c := make(chan os.Signal)
    signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        select {
        case <-c:
            close(closecli)
        case <-time.After(time.Second * 15):
            // 测试15秒, 统计每秒hello次数
            close(closecli)
        }
    }()

    q := newQps(time.Second, closecli)
    num := 1000

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
            run := true
            for run {
                select {
                case <-closecli:
                    run = false
                    break
                default:
                    var resp string
                    start := time.Now()
                    err := echoPrx.Hello("hello", &resp)
                    q.add(time.Since(start))
                    if err != nil {
                        log.Errorf("err:%s", err)
                    }
                }
            }
            done.Done()
        }(i+1)
    }

    done.Wait()
    log.Debug("wait client to exit")
}

type Qps struct {
    mu sync.Mutex

    total uint32
    qps uint32
    times uint32

    qpstime time.Duration
}

func (p *Qps) add(t time.Duration) {
    p.mu.Lock()
    defer p.mu.Unlock()

    p.qps++
    p.qpstime += t
}

func (p *Qps) round() {
    p.mu.Lock()

    qps := p.qps
    t := p.qpstime

    p.total += p.qps
    p.qps = 0
    p.times++
    
    p.qpstime = 0
    p.mu.Unlock()

    log.Debugf("QPS:%d,time:%d ms", qps, t.Milliseconds()/int64(qps))
}

func (p *Qps) avg() {
    p.mu.Lock()
    defer p.mu.Unlock()
   
    log.Debugf("AVG QPS:%d", p.total/p.times)
}

func newQps(roundtime time.Duration, stop chan bool) *Qps {
    loop := time.NewTicker(roundtime)

    q := &Qps{}

    go func(){
        done.Add(1)
        defer func() {
            done.Done()
        }()
        for {
            select {
            case <-stop:
                loop.Stop()
                q.avg() 
                return
            case <-loop.C:
                q.round()
            }
        }
    }()

    return q
}
