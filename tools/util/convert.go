package util

import (
    "time"
    "strings"
    "strconv"
)

// 字符串格式如下
// ns 纳秒
// us/µs 微秒
// ms 毫秒
// s 秒
// m 分钟
// h 小时
// 可以组合使用
func AtoDuration(v string) (d time.Duration) {
    var err error
    if d, err = time.ParseDuration(v); err != nil {
        panic("invalid time duration:"+v)
    }
    return
}

const (
	// B byte
	B uint64 = 1
	// K kilobyte
	K uint64 = 1 << (10 * iota)
	// M megabyte
	M
	// G gigabyte
	G
	// T TeraByte
	T
	// P PetaByte
	P
	// E ExaByte
	E
)

var unitMap = map[string]uint64{
	"B":  B,
	"K":  K,
	"KB": K,
	"M":  M,
	"MB": M,
	"G":  G,
	"GB": G,
	"T":  T,
	"TB": T,
	"P":  P,
	"PB": P,
	"E":  E,
	"EB": E,
}

func AtoMB(ssize string) uint64 {
    ssize = strings.ToUpper(strings.TrimSpace(ssize))
	if ssize == "" {
		return 0
	}

    size := uint64(0)
    unit := "B"
    p := strings.IndexAny(ssize, "BKMGTPE")
    if p != -1 {
        unit = ssize[p:]
        ssize = ssize[:p]
    }

    var err error
    if size,err = strconv.ParseUint(ssize, 10, 64); err != nil {
        return 0
    }

	iunit,ok := unitMap[unit]
	if ok {
	    size = size * iunit / 1024 / 1024
    }

	return size
}
