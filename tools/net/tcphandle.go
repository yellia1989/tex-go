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
        log.FError(err.Error())
        return
    }
    h.lis, err = net.ListenTCP("tcp4", addr)
    if err != nil {
        log.FError(err.Error())
        return
    }
    log.FDebugf("start listen on:%s", addr)

    for !h.svr.isClose() {
        if err := h.lis.SetDeadline(time.Now().Add(time.Millisecond*10)); err != nil {
            log.FErrorf("set accept timeout failed:%s", err.Error())
            return
        }
        conn, err := h.lis.AcceptTCP()
        if err != nil {
            if !isTimeoutErr(err) {
                log.FErrorf("accept error:%s", err.Error())    
            }
            continue
        }

        if h.svr.isClose() {
            conn.Close()
            break
        }

        if err := conn.SetReadBuffer(cfg.TCPReadBuffer); err != nil {
            log.FErrorf("set tcp conn read buffer err:%s", err.Error())
            conn.Close()
            continue
        }
        if err := conn.SetWriteBuffer(cfg.TCPWriteBuffer); err != nil {
            log.FErrorf("set tcp conn write buffer err:%s", err.Error())
            conn.Close()
            continue
        }
        if err := conn.SetNoDelay(cfg.TCPNoDelay); err != nil {
            log.FErrorf("set tcp no delay err:%s", err.Error())
            conn.Close()
            continue
        }
        h.svr.addConnection(conn)
    }

    log.FDebugf("stop listen on:%s", addr)
}
