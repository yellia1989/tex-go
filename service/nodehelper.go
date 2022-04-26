package tex

import (
    "os"
    "github.com/yellia1989/tex-go/sdp/rpc"
    "github.com/yellia1989/tex-go/tools/log"
)

type NodeHelper struct {
    proxy *rpc.Node
    pid uint32
};

func (nh *NodeHelper) init(node string) error {
    if node == "" {
        return nil
    }

    nh.proxy = new(rpc.Node)
    err := StringToProxy(node, nh.proxy)
    if err != nil {
        return err
    }
    nh.pid = uint32(os.Getpid())

    return nil
}

func (nh *NodeHelper) keepAlive(adapterName string, initing bool) {
    if nh.proxy == nil {
        return
    }

    if ret, err := nh.proxy.KeepAlive(App, Server, Zone, nh.pid, adapterName, initing); err != nil {
        log.FError("keepAlive failed, ret: %d, err: %s", ret, err.Error())
    }
}
