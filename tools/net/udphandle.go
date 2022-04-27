package net

import (
    "net"
    _ "github.com/yellia1989/tex-go/tools/log"
)

type udpHandle struct {
    svr *Svr
    conn *net.UDPConn
}

func (h *udpHandle) Run() {
    /*
    defer func() {
        if h.conn != nil {
            h.conn.Close()
        }
    }()

    addr, err := net.ResolveUDPAddr("udp4", cfg.Address)
    if err != nil {
        log.FErrorf("listen on %s failed", cfg.Address)
        return
    }
    h.conn, err = net.ListenUDP("udp4", addr)
    if err != nil {
        log.FErrorf("listen on %s failed", cfg.Address)
        return
    }
    log.FDebugf("start listen on:%s", addr)

    tmpbuf := make([]byte, 65535)
    for atomic.LoadInt32(&h.svr.isclose) == 0 {
        n, udpAddr, err := h.conn.ReadFromUDP(tmpbuf)
        if err != nil {
            if isTimeoutErr(err) {
                continue
            } else {
                log.FErrorf("udp read err:%s", err.Error())
                return
            }
        }
        pkg := make([]byte, n)
        copy(pkg, tmpbuf[0:n])
    }*/
}
