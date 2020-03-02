package net

import (
   "net"
   "time"
   "sync"
   "fmt"
   "io"
   "github.com/yellia1989/tex-go/tools/log"
)

const (
    cli_idle_timeout = time.Minute * 5
    cli_write_queuecap = 100
)

type CliCfg struct {
    Proto string // tcp,udp
    WriteQueueCap  int // 每个连接的待发送队列的长度
    IdleTimeout time.Duration // 连接最长空闲时间
}

type Cli struct {
    address string
    cfg *CliCfg
    pkgHandle CliPkgHandle

    mu sync.Mutex
    close bool
    conn net.Conn
    idleTime time.Time

    writech chan []byte
}

func NewCli(address string, cfg *CliCfg , pkgHandle CliPkgHandle) *Cli {
    if cfg.IdleTimeout < cli_idle_timeout {
        cfg.IdleTimeout = cli_idle_timeout
    }
    if cfg.WriteQueueCap < cli_write_queuecap {
        cfg.WriteQueueCap = cli_write_queuecap
    }
    cli := &Cli{address: address, cfg: cfg, pkgHandle: pkgHandle, close: true}
    
    return cli
}

func (cli *Cli) Send(pkg []byte) error {
    if len(pkg) == 0 {
        return fmt.Errorf("empty or nil pkg")
    }

    cli.mu.Lock()
    if cli.close {
        if err := cli.connect(); err != nil {
            cli.mu.Unlock()
            return err
        }
    }
    cli.idleTime = time.Now()
    cli.mu.Unlock()

    cli.writech <- pkg
    return nil
}

func (cli *Cli) connect() error {
    if !cli.close {
        panic("conn is connected")
    }

    // 默认1秒连接超时
    conn, err := net.DialTimeout(cli.cfg.Proto, cli.address, time.Second * 1)
    if err != nil {
        return err
    }
    cli.conn = conn
    cli.close = false
    cli.writech = make(chan []byte, cli.cfg.WriteQueueCap)

    // 开启一个独立的携程写
    go cli.doWrite()
    // 开启一个独立的携程读
    go cli.doRead()

    return nil
}

func (cli *Cli) Close() {
    cli.mu.Lock()
    defer cli.mu.Unlock()

    if cli.close {
        return
    }

    cli.close = true
    cli.conn.Close()
    close(cli.writech)
}

func (cli *Cli) SafeClose() {
    cli.mu.Lock()
    if cli.close {
        cli.mu.Unlock()
        return
    }
    cli.mu.Unlock()

    pkg := make([]byte, 0)
    cli.writech <- pkg
}

func (cli *Cli) doWrite() {
    defer func() {
        cli.Close()
    }()

    for {
        select {
        case pkg := <-cli.writech :
            total := len(pkg)
            if total == 0 {
                // 优雅关闭
                return
            }
            for {
                n, err := cli.conn.Write(pkg)
                if err != nil {
                    log.FErrorf("write err:%s", err.Error())
                    return
                }
                if n > 0 {
                    total -= n
                }
                if total == 0 {
                    break
                }
                pkg = pkg[n:]
            }
        }
    }
}

func (cli *Cli) doRead() {
    defer func() {
        cli.Close()
    }()

    tmpbuf := make([]byte, 1024*4)
    var pkgbuf []byte
    for {
        if err := cli.conn.SetReadDeadline(time.Now().Add(time.Millisecond*500)); err != nil {
            log.FError("set conn read deadline err:%s", err.Error())
            return
        }
        n, err := cli.conn.Read(tmpbuf)
        if err != nil {
            if isTimeoutErr(err) {
                cli.mu.Lock()
                if cli.idleTime.Add(cli.cfg.IdleTimeout).Before(time.Now()) {
                    cli.mu.Unlock()
                    log.FDebugf("conn is unactive, will be closed")
                    return
                }
                cli.mu.Unlock()
                continue
            }
            if (err == io.EOF) {
                log.FDebugf("conn has been closed by server:%s", err.Error())
            } else {
                log.FErrorf("read err:%s", err.Error())
            }
            return
        }
        // 解析包
        pkgbuf = append(pkgbuf, tmpbuf[:n]...)
        for {
            pkglen, status := cli.pkgHandle.Parse(pkgbuf)
            if status == PACKAGE_LESS {
                break
            }
            if status == PACKAGE_FULL {
                pkg := make([]byte, pkglen)
                copy(pkg, pkgbuf[:pkglen])
                cli.pkgHandle.Recv(pkg)
                pkgbuf = pkgbuf[pkglen:]
                if len(pkgbuf) > 0 {
                    continue
                }
                pkgbuf = nil
                break
            }
            log.FErrorf("parse package error")
            return
        }
    }
}
