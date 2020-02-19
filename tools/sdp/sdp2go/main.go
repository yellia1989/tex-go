package main

import (
    "flag"
    "os"
    "strings"
    "fmt"
)

// 生成文件保存路径
var dir = flag.String("dir", "./", "dir to save generated code")

func usage() {
    bin := os.Args[0]
    if p := strings.LastIndex(bin, "/"); p != -1 {
        bin = bin[p+1:]
    }
    fmt.Printf("Usage: %s --dir= file[,file]\n", bin)
    flag.PrintDefaults()
}

func main() {
    flag.Usage = usage
    flag.Parse()

    if flag.NArg() == 0 {
        usage()
        os.Exit(0)
    }

    for _, file := range flag.Args() {
        sdp2Go := newSdp2Go(file, *dir)
        sdp2Go.generate()
    }
}
