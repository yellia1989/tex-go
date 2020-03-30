// sdp生成文件中用到的公共函数放在这里
package util

import (
    "bytes"
    "strings"
)

func Tab(buff *bytes.Buffer, tab int, code string) {
    buff.WriteString(strings.Repeat(" ", tab*4) + code)
}

func Fieldname(name string) string {
    if name != "" {
        return name + ": "
    }
    return ""
}
