package net

import (
    "net"
)

func isTimeoutErr(err error) bool {
    if err, ok := err.(net.Error); ok {
        return err.Timeout()
    }

    return false
}
