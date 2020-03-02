package net

import (
   _ "net"
   _ "time"
)

type CliCfg struct {
    Proto string // tcp,udp

    WriteQueueCap  int // 每个连接的待发送队列的长度
}
