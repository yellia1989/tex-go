package net

import (
    "net"
    "context"
)

func isTimeoutErr(err error) bool {
    if err, ok := err.(net.Error); ok {
        return err.Timeout()
    }

    return false
}

const (
    PACKAGE_LESS = iota
    PACKAGE_FULL
    PACKAGE_ERROR
)

// 服务器接收到数据包的处理接口
type SvrPkgHandle interface {
    // 将二进制流按照特定的协议解析成单个的包
    Parse(bytes []byte) (int,int)
    HandleRecv(ctx context.Context, pkg []byte, overload bool, queuetimeout bool)
}

// 客户端接收到数据包的处理接口
type CliPkgHandle interface {
    // 将二进制流按照特定的协议解析成单个的包
    Parse(bytes []byte)(int,int)
    // 单个数据包正常处理
    Recv(pkg []byte)
}
