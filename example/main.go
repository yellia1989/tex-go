package main

import (
    tex "github.com/yellia1989/tex-go/service"
    "github.com/yellia1989/tex-go/sdp/rpc"
    "github.com/yellia1989/tex-go/tools/log"
)

func main() {
    log.SetFrameworkLevel(log.DEBUG)
    log.SetLevel(log.DEBUG)

    comm := tex.NewCommunicator("")

    query := new(rpc.Query)
    comm.StringToProxy("tex.mfwregistry.QueryObj@tcp -h 192.168.0.16 -p 2000 -t 3600000", query)

    var vActiveEps []string
    var vInactiveEps []string
    ret, err := query.GetEndpoints("aqua.GameServer.GameServiceObj", "aqua.zone.2", &vActiveEps, &vInactiveEps)
    if err != nil {
        log.Debugf("query err:%s", err.Error())
    } else if (ret != 0) {
        log.Debugf("query err, ret:%d", ret)
    } else {
        log.Debugf("query success active:%v, inactive:%v", vActiveEps, vInactiveEps)
    }

    log.FlushLogger()
}
