package tex

import (
    "time"
    "os"
    "os/signal"
    "syscall"
    "github.com/yellia1989/tex-go/tools/log"
)

// 应用程序必须实现的接口
type app interface {
    Init()
    Loop()
    Terminate()
}

var (
    comm *Communicator
)

func init() {
    comm = NewCommunicator()
}

func Run(svr app) {
    defer func() {
        log.FlushLogger()
    }()

    // 初始化配置
    if err := parseCfg(); err != nil {
        log.FErrorf("parse cfg err:%s", err.Error())
        return
    }

    // 初始化客户端
    if err := initClient(); err != nil {
        log.FErrorf("init client err:%s", err.Error())
        return
    }

    // 初始化服务器
    if err := initServer(); err != nil {
        log.FErrorf("init client err:%s", err.Error())
        return
    }

    // 开启服务器
    if err := startServer(); err != nil {
        log.FErrorf("start server err:%s", err.Error())
        return
    }

    // 初始化应用程序
    svr.Init()

    // 监听信号
    c := make(chan os.Signal)
    signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

    // 启动主循环等待服务器结束
    ticker := time.NewTicker(time.Second)
    run := true
    for run {
        select {
        case <-c :
            run = false
        case <-ticker.C :
            svr.Loop()
        }
    }
    // 结束服务器
    stopServer()

    // 结束应用程序
    svr.Terminate()

    // 结束客户端
    stopClient()
}

func StringToProxy(name string, proxy ServicePrx) error {
    return comm.StringToProxy(name, proxy)
}

func initClient() error {
    // TODO
    return nil
}

func initServer() error {
    // TODO
    return nil
}

func stopClient() {
    comm.Close()
}
