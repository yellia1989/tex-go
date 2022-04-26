package tex

import (
    "flag"
    "fmt"
    "os"
    "strconv"
    "time"
    "github.com/yellia1989/tex-go/tools/log"
    "github.com/yellia1989/tex-go/tools/util"
    "github.com/yellia1989/tex-go/sdp/rpc"
)

var (
    servicesCfg map[string]*serviceCfg
    cliCfg clientCfg
    configFile string

    App     string
    Server  string
    Zone    string
    loop_interval int
    Node    string
    Admin   string
)

func init() {
    servicesCfg = make(map[string]*serviceCfg)
    flag.StringVar(&configFile, "config", "", "config file")
    flag.Usage = func () {
        fmt.Fprintf(os.Stderr, "Usage: %s --config file\n", os.Args[0])
        os.Exit(1)
    }

    cliCfg.invokeTimeout = 5 * time.Second
    cliCfg.endpointRefreshInterval = 60 * time.Second
    cliCfg.adapterSendQueueCap = 10000
    cliCfg.adapterIdleTimeout = 10 * time.Minute
}

type serviceCfg struct {
    service string // test.EchoServer.EchoServiceObj
    endpoint *Endpoint // tcp -h 127.0.0.1 -p 8080 -t 3600000
    isTex bool
    threads  int  // 2
    maxconns int  // 1024
    queuecap int  // 10240
    queuetimeout time.Duration // 5000
}

type clientCfg struct {
    locator string
    division string
    invokeTimeout time.Duration
    endpointRefreshInterval time.Duration
    adapterSendQueueCap int
    adapterIdleTimeout time.Duration
}

func parseCfg() error {
    flag.Parse()
    if len(configFile) == 0 {
        flag.Usage()
    }

    cfg := util.NewConfig()
    cfg.ParseFile(configFile)

    // 解析服务器配置
    svrCfg := cfg.GetSubCfg("mfw/application/server")
    App = svrCfg.GetCfg("app", "app")
    Server = svrCfg.GetCfg("server", "server")
    Zone = cfg.GetCfg("mfw/application/setdivision", "")
    Node = svrCfg.GetCfg("node", "")
    Admin = svrCfg.GetCfg("admin", "")

    // 日志相关
    defLogger := log.GetDefaultLogger()
    logpath := svrCfg.GetCfg("logpath", "")
    if len(logpath) != 0 {
        lognum := svrCfg.GetInt("lognum", 10)
        logsize := util.AtoMB(svrCfg.GetCfg("logsize", "100M"))
        if logsize < 100 {
            logsize = 100 
        }
        if logpath[len(logpath)-1] != '/' {
            logpath += "/"
        }
        logpath += App + "/" + Server + "/" + Zone
        defLogger.SetLogName(App+"."+Server)
        defLogger.SetFileRoller(logpath, lognum, int(logsize))
    }
    log.SetLevel(log.StringToLevel(svrCfg.GetCfg("loglevel", "INFO")))
    log.SetFrameworkLevel(log.StringToLevel(svrCfg.GetCfg("framework-loglevel", "INFO")))

    // service相关
    i := 1
    for {
        cfg := svrCfg.GetSubCfg("Service_"+strconv.Itoa(i)) 
        if cfg == nil {
            break
        }
        objCfg := &serviceCfg{}
        objCfg.service = cfg.GetCfg("service","")

        var err error
        objCfg.endpoint, err = NewEndpoint(cfg.GetCfg("endpoint", ""))
        if err != nil {
            return err
        }

        objCfg.isTex = cfg.GetCfg("protocol", "mfw") == "mfw"
        objCfg.threads = cfg.GetInt("threads", 1)
        objCfg.maxconns = cfg.GetInt("maxconns", 1024)
        objCfg.queuecap = cfg.GetInt("queuecap", 10240)
        queuetimeout := cfg.GetCfg("queuetimeout", "5000ms")
        if queuetimeout != "5000ms" {
            queuetimeout += "ms"
        }
        objCfg.queuetimeout = util.AtoDuration(queuetimeout)
        if _, ok := servicesCfg[objCfg.service]; ok {
            return fmt.Errorf("duplicate service obj:%s", objCfg.service)
        }
        servicesCfg[objCfg.service] = objCfg
        i++
    }
    if Admin != "" {
        objCfg := &serviceCfg{}
        objCfg.service = "AdminObj"

        var err error
        objCfg.endpoint, err = NewEndpoint(Admin)
        if err != nil {
            return err
        }

        objCfg.isTex = true
        objCfg.threads = 1
        objCfg.maxconns = 1024
        objCfg.queuecap = 10240
        queuetimeout :=  "5000ms"
        objCfg.queuetimeout = util.AtoDuration(queuetimeout)
        servicesCfg[objCfg.service] = objCfg

        admin := new(rpc.Admin)
        adminImpl := new(AdminServiceImpl)
        AddService("AdminObj", admin, adminImpl)
    }

    // 解析客户端配置
    cliconfig := cfg.GetSubCfg("mfw/application/client")
    cliCfg.locator = cliconfig.GetCfg("locator", "")
    invokeTimeout := cliconfig.GetCfg("async-invoke-timeout", "5s")
    if invokeTimeout != "5s" {
        invokeTimeout += "ms"
    }
    cliCfg.invokeTimeout = util.AtoDuration(invokeTimeout)
    endpointRefreshInterval := cliconfig.GetCfg("refresh-endpoint-interval", "10s")
    if endpointRefreshInterval != "10s" {
        endpointRefreshInterval += "ms"
    }
    cliCfg.endpointRefreshInterval = util.AtoDuration(endpointRefreshInterval)
    cliCfg.adapterSendQueueCap = cliconfig.GetInt("send-queue-cap", 10000)
    adapterIdleTimeout := cliconfig.GetCfg("idle-time", "10m")
    if adapterIdleTimeout != "10m" {
        adapterIdleTimeout += "ms"
    }
    cliCfg.adapterIdleTimeout = util.AtoDuration(adapterIdleTimeout)

    return nil
}
