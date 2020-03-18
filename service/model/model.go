// 只是为了解决循环依赖

package model

import (
    "time"
    "github.com/yellia1989/tex-go/sdp/protocol"
)

type ServicePrxImpl interface {
    // 只有当error == nil时，resp才有效
    Invoke(sFuncName string, params []byte, resp **protocol.ResponsePacket) error
    SetTimeout(timeout time.Duration)
}
