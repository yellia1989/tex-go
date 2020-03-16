package main

import (
    tex "github.com/yellia1989/tex-go/service"
    "github.com/yellia1989/tex-go/tools/log"
    "github.com/yellia1989/tex-go/example/rpc/server/echo"
)

type EchoServer struct {
}

func (s *EchoServer) Init() {
    // 应用初始化
    log.Debug("server init")
}

func (s *EchoServer) Loop() {
    // 应用主循环
    //log.Debug("server loop")
}

func (s *EchoServer) Terminate() {
    // 应用停止
    log.Debug("server terminate")
}

func main() {
    service := &echo.EchoService{}
    serviceImpl := &EchoServiceImpl{}
    tex.AddService("test.EchoServer.EchoServiceObj", service, serviceImpl)

    service2 := &echo.EchoService{}
    serviceImpl2 := &EchoServiceImpl{}
    tex.AddService("test.EchoServer.EchoServiceObj2", service2, serviceImpl2)

    tex.Run(&EchoServer{})
}
