package tex

import (
    "sync"
    "fmt"
    "context"
    "github.com/yellia1989/tex-go/tools/net"
    "github.com/yellia1989/tex-go/protocol/protocol"
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

func startServer() (err error) {
    defer func() {
        err := recover()
        if err != nil {
            err = fmt.Errorf("%s", err)
        }
    }()
    for k, v := range services {
        _, ok := servicesCfg[k]
        if !ok {
            panic(fmt.Sprintf("service:%s can't find cfg", k))
        }
        
        svr := net.NewSvr(&net.SvrCfg{
           // TODO 
        },&texSvrPkgHandle{
            service: v.service,
            serviceImpl: v.serviceImpl,
        })
        svrRun[k] = svr
        svrDone.Add(1)

        go svr.Run()
    }
    return nil
}

func stopServer() {
    for _, svr := range svrRun {
        svr.Stop()
    }
    svrDone.Wait()
}
