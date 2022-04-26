package tex

import (
    "fmt"
	"context"
    "strings"
    "github.com/yellia1989/tex-go/sdp/protocol"
    "github.com/yellia1989/tex-go/tools/log"
)

type AdminServiceImpl struct {
}

func (admin *AdminServiceImpl) Shutdown(ctx context.Context) error {
    quit <- struct{}{}
    return nil
}

func (admin *AdminServiceImpl) Notify(ctx context.Context, sCmd string, sResult *string) (int32, error) {
    log.FDebug("notify cmd: %s", sCmd)

    if sCmd == "" {
        return protocol.MFW_INVALID_PARAM, nil
    }

    vCmd := strings.Split(sCmd, " ")
    cmd := vCmd[0]
    mParam := make(map[string]string)
    if len(vCmd) > 1 {
        for _,v := range strings.Split(vCmd[1], "&") {
            tmp := strings.Split(v, "=")
            if len(tmp) == 2 {
                mParam[tmp[0]] = tmp[1]
            }
        }
    }

    switch cmd {
    case "set_log_level":
        return admin.setLogLevel(mParam)
    default:
        return protocol.MFW_INVALID_PARAM, nil
    }

    return 0, nil
}

func (admin *AdminServiceImpl) setLogLevel(mParam map[string]string) (ret int32, err error) {
    defer func() {
        perr := recover()
        if perr != nil {
            err = fmt.Errorf("invalid log level, err: %s", perr)
        }
    }()

    ret = protocol.MFW_UNKNOWN

    if level, ok := mParam["loglevel"]; ok {
        log.SetLevel(log.StringToLevel(level))
    }

    if level, ok := mParam["framework-loglevel"]; ok {
        log.SetFrameworkLevel(log.StringToLevel(level))
    }

    ret = 0

    return
}
