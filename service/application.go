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
    Init() error
    Loop()
    Terminate()
}

var (
    comm *Communicator
    nodeHelper NodeHelper
    quit chan struct{}
)

func Run(svr app) {
    defer func() {
        log.FlushLogger()
    }()

    log.FDebug("server start...")

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
        log.FErrorf("init server err:%s", err.Error())
        return }

    // 开启服务器
    if err := startService(); err != nil {
        log.FErrorf("start service err:%s", err.Error())
        return
    }

    // 开启协程定时通知node正在启动中防止被判断为心跳超时
    initdone := make(chan struct{})
    go func () {
        ticker := time.NewTicker(time.Second * 2)
        defer ticker.Stop()
        for {
            select {
                case <-initdone :
                    go nodeHelper.keepAlive("", false)
                    return
                case <-ticker.C :
                    go nodeHelper.keepAlive("", true)
            }
        }
    }()
    // 初始化应用程序
    if err := svr.Init(); err != nil {
        initdone <- struct{}{}
        log.FErrorf("server init err:%s", err.Error())
        return
    }
    initdone <- struct{}{}

    // 监听信号
    c := make(chan os.Signal)
    signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

    log.FDebug("server started")

    // 启动主循环等待服务器结束
    ticker := time.NewTicker(time.Second)
    run := true
    for run {
        select {
        case <-c :
            run = false
        case <-quit :
            run = false
        case <-ticker.C :
            svr.Loop()
        }
    }
    // 结束服务器
    stopService()

    // 结束应用程序
    svr.Terminate()

    // 结束客户端
    stopClient()

    log.FDebug("server stopped")
}

func StringToProxy(name string, proxy ServicePrx) error {
    return comm.StringToProxy(name, proxy)
}

func initClient() error {
    comm = NewCommunicator(cliCfg.locator)

    return nil
}

func initServer() error {
    if err := nodeHelper.init(Node); err != nil {
        return err
    }

    return nil
}

func stopClient() {
    comm.Close()
}
