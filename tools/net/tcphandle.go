package net

import (
    "net"
    "time"
    "github.com/yellia1989/tex-go/tools/log"
)

type tcpHandle struct {
    svr *Svr
    lis *net.TCPListener
}

func (h *tcpHandle) Run() {
    defer func() {
        if h.lis != nil {
            h.lis.Close()
        }
    }()

    cfg := h.svr.cfg
    addr, err := net.ResolveTCPAddr("tcp4", cfg.Address)
    if err != nil {
        log.Errorf("listen on %s failed", cfg.Address)
        return
    }
    h.lis, err = net.ListenTCP("tcp4", addr)
    if err != nil {
        log.Errorf("listen on %s failed", cfg.Address)
        return
    }
    log.Debugf("start listen on:%s", addr)

    for !h.svr.close {
        if err := h.lis.SetDeadline(time.Now().Add(time.Millisecond*500)); err != nil {
            log.Errorf("set accept timeout failed:%s", err.Error())
            return
        }
        conn, err := h.lis.AcceptTCP()
        if err != nil {
            if isTimeoutErr(err) {
            } else {
                log.Errorf("accept error:%s", err.Error())    
            }
            continue
        }

        conn.SetReadBuffer(cfg.TCPReadBuffer)
        conn.SetWriteBuffer(cfg.TCPWriteBuffer)
        conn.SetNoDelay(cfg.TCPNoDelay)
        h.svr.addConnection(conn)
    }

    log.Debugf("stop listen on:%s", addr)
}
