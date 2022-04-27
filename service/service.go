package tex

import (
    "sync"
    "fmt"
    "context"
    "time"
    "github.com/yellia1989/tex-go/tools/net"
    "github.com/yellia1989/tex-go/tools/log"
    "github.com/yellia1989/tex-go/sdp/protocol"
)

type Service interface {
    Dispatch(ctx context.Context, serviceImpl interface{}, req *protocol.RequestPacket)    
}

type serviceDetail struct {
    service Service
    serviceImpl interface{}
}

var (
    // 本application提供的所有service,不包括http service
    services map[string]serviceDetail
    svrRun map[string]*net.Svr
    svrDone sync.WaitGroup
)

func init() {
    services = make(map[string]serviceDetail)
    svrRun = make(map[string]*net.Svr)
}

func AddService(obj string, service Service, serviceImpl interface{}) {
    _, ok := services[obj]
    if ok {
        panic("duplicate service obj:" + obj)
    }
    services[obj] = serviceDetail{service: service, serviceImpl: serviceImpl}
}

func startService() (err error) {
    defer func() {
        perr := recover()
        if perr != nil {
            err = fmt.Errorf("%s", perr)
        }
    }()
    
    ch := make(chan string)
    for k, v := range services {
        cfg, ok := servicesCfg[k]
        if !ok {
            panic(fmt.Sprintf("service:%s can't find cfg", k))
        }
        
        adapterName := k
        svr := net.NewSvr(&net.SvrCfg{
            Name: adapterName,
            Proto: cfg.endpoint.Proto,
            Address: cfg.endpoint.Address(),
            WorkThread: cfg.threads,
            WorkQueueCap: cfg.queuecap,
            WorkQueueTimeout: cfg.queuetimeout,
            MaxConn: cfg.maxconns,
            IdleTimeout: cfg.endpoint.Idletimeout,
            TCPNoDelay: true,
            Heartbeat: func() {
                if adapterName != "AdminObj" {
                    go nodeHelper.keepAlive(adapterName, false)
                }
            },
        },&texSvrPkgHandle{
            name: k,
            service: v.service,
            serviceImpl: v.serviceImpl,
        })
        svrRun[k] = svr
        svrDone.Add(1)

        go func(service string) {
            log.FDebugf("service:%s start", service)
            svr.Run()
            log.FDebugf("service:%s stop", service)
            svrDone.Done()
            ch <- service
        }(k)
    }

    // 等待2秒所有服务器监听成功
    select {
    case name := <-ch:
        err = fmt.Errorf("start service:%s failed", name)
    case <-time.After(time.Second * 2):
    }

    return
}

func stopService() {
    for _, svr := range svrRun {
        svr.Stop()
    }
    svrDone.Wait()
}
