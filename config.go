package tex

import (
    "flag"
    "fmt"
    "os"
    "strconv"
    "time"
    "github.com/yellia1989/tex-go/tools/log"
    "github.com/yellia1989/tex-go/tools/util"
)

var (
    servicesCfg map[string]*serviceCfg
    configFile string

    App     string
    Server  string
    Zone    string
    loop_interval int
)

func init() {
    servicesCfg = make(map[string]*serviceCfg)
    flag.StringVar(&configFile, "config", "", "config file")
    flag.Usage = func () {
        fmt.Fprintf(os.Stderr, "Usage: %s --config file\n", os.Args[0])
        os.Exit(1)
    }
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

func parseCfg() error {
    flag.Parse()
    if len(configFile) == 0 {
        flag.Usage()
    }

    cfg := util.NewConfig()
    cfg.ParseFile(configFile)

    // 解析服务器配置
    svrCfg := cfg.GetSubCfg("tex/application/server")
    App = svrCfg.GetCfg("app", "app")
    Server = svrCfg.GetCfg("server", "server")
    Zone = cfg.GetCfg("tex/application/setdivision", "")

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
    log.SetFrameworkLevel(log.StringToLevel(svrCfg.GetCfg("framework-loglevel", "DEBUG")))

    // service相关
    i := 1
    for {
        cfg := svrCfg.GetSubCfg("Service_"+strconv.Itoa(i)) 
        if cfg == nil {
            break
        }
        objCfg := &serviceCfg{}
        objCfg.service = cfg.GetCfg("service","")
        objCfg.endpoint = NewEndpoint(cfg.GetCfg("endpoint", ""))
        objCfg.isTex = cfg.GetCfg("protocol", "tex") == "tex"
        objCfg.threads = cfg.GetInt("threads", 1)
        objCfg.maxconns = cfg.GetInt("maxconns", 1024)
        objCfg.queuecap = cfg.GetInt("queuecap", 10240)
        queuetimeout := cfg.GetCfg("queuetimeout", "500ms")
        if queuetimeout != "500ms" {
            queuetimeout += "ms"
        }
        objCfg.queuetimeout = util.AtoDuration(queuetimeout)
        if _, ok := servicesCfg[objCfg.service]; ok {
            return fmt.Errorf("duplicate service obj:%s", objCfg.service)
        }
        servicesCfg[objCfg.service] = objCfg
        i++
    }

    return nil
}
