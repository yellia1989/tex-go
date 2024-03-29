package main

import (
    tex "github.com/yellia1989/tex-go/service"
    "github.com/yellia1989/tex-go/sdp/rpc"
    "github.com/yellia1989/tex-go/tools/log"
    "time"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    log.SetFrameworkLevel(log.DEBUG)
    log.SetLevel(log.DEBUG)

    defer func() {
        log.FlushLogger()
    }()

    c := make(chan os.Signal)
    signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

    comm := tex.NewCommunicator("tex.mfwregistry.QueryObj@tcp -h 192.168.0.16 -p 2000 -t 3600000")

    var vActiveEps []string
    var vInactiveEps []string

    exit:
    for {
        loop := time.NewTicker(time.Second * 10)
        select {
        case <-c:
            loop.Stop()
            break exit 
        case <-loop.C:
            query := new(rpc.Query)
            if err := comm.StringToProxy("tex.mfwregistry.QueryObj", query); err != nil {
                log.Errorf("failed to alloc proxy, err:%s", err.Error())
                return
            }
            ret, err := query.GetEndpoints("aqua.GameServer.GameServiceObj", "aqua.zone.2", &vActiveEps, &vInactiveEps)
            if err != nil {
                log.Debugf("query err:%s", err.Error())
            } else if (ret != 0) {
                log.Debugf("query err, ret:%d", ret)
            } else {
                log.Debugf("query success active:%v, inactive:%v", vActiveEps, vInactiveEps)
            }
        }
    }
}
