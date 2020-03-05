package tex

import (
    "time"
    "github.com/yellia1989/tex-go/tools/log"
)

// 应用程序必须实现的接口
type App interface {
    Init()
    Loop()
    Terminate()
}

var (
    shutdown chan bool
)

func init() {
    shutdown = make(chan bool)
}

func Run(app App) {
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

    // 初始化应用程序
    app.Init()

    // 开启服务器
    if err := startServer(); err != nil {
        log.FErrorf("start server err:%s", err.Error())
        return
    }

    // 启动主循环等待服务器结束
    ticker := time.NewTicker(time.Second)
    for {
        select {
        case <-shutdown :
            goto stop
        case <-ticker.C :
            app.Loop()
        }
    }
    stop:
    // 结束服务器
    stopServer()

    // 结束应用程序
    app.Terminate()
}

func stop() {
    shutdown <- true
}

func parseCfg() error {
    // TODO
    return nil
}

func initClient() error {
    // TODO
    return nil
}

func initServer() error {
    // TODO
    return nil
}
