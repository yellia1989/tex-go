// 此文件为sdp2go工具自动生成,请不要手动编辑

package echo

import (
	"context"
	tex "github.com/yellia1989/tex-go/service"
	"github.com/yellia1989/tex-go/service/protocol/protocol"
	"github.com/yellia1989/tex-go/tools/log"
	"github.com/yellia1989/tex-go/tools/net"
	"github.com/yellia1989/tex-go/tools/sdp/codec"
	"time"
)

type EchoService struct {
	proxy tex.ServicePrxImpl
}

func (s *EchoService) SetPrxImpl(impl tex.ServicePrxImpl) {
	s.proxy = impl
}
func (s *EchoService) SetTimeout(timeout time.Duration) {
	s.proxy.SetTimeout(timeout)
}
func (s *EchoService) Hello(req string, resp *string) error {
	p := codec.NewPacker()
	var err error
	var has bool
	var ty uint32
	var length uint32
	if true || req != "" {
		err = p.WriteString(1, req)
		if err != nil {
			return err
		}
	}
	var rsp protocol.ResponsePacket
	err = s.proxy.Invoke("hello", p.ToBytes(), &rsp)
	if err != nil {
		return err
	}
	up := codec.NewUnPacker([]byte(rsp.SRspPayload))
	err = up.ReadString(&(*resp), 2, true)
	if err != nil {
		return err
	}
	_ = has
	_ = ty
	_ = length
	return nil
}

type _EchoServiceImpl interface {
	Hello(ctx context.Context, req string, resp *string) error
}

func _HelloImpl(ctx context.Context, serviceImpl interface{}, up *codec.UnPacker, p *codec.Packer) error {
	var err error
	var length int
	impl := serviceImpl.(_EchoServiceImpl)
	var p1 string
	err = up.ReadString(&p1, 1, true)
	if err != nil {
		return err
	}
	var p2 string
	err = impl.Hello(ctx, p1, &p2)
	if err != nil {
		return err
	}
	if true || p2 != "" {
		err = p.WriteString(2, p2)
		if err != nil {
			return err
		}
	}
	_ = length
	return nil
}

func (s *EchoService) Dispatch(ctx context.Context, serviceImpl interface{}, req *protocol.RequestPacket) {
	current := net.ContextGetCurrent(ctx)

	log.FDebugf("handle tex request, peer: %s:%d, obj: %s, func: %s, reqid: %d", current.IP, current.Port, req.SServiceName, req.SFuncName, req.IRequestId)

	texret := protocol.SDPSERVERUNKNOWNERR
	up := codec.NewUnPacker([]byte(req.SReqPayload))
	p := codec.NewPacker()

	var err error
	switch req.SFuncName {
	case "hello":
		err = _HelloImpl(ctx, serviceImpl, up, p)
		if err != nil {
			break
		}
		texret = protocol.SDPSERVERSUCCESS
	default:
		texret = protocol.SDPSERVERNOFUNCERR
	}

	if err != nil {
		log.FErrorf("handle tex request, peer: %s:%d, obj: %s, func: %s, reqid: %d, err: %s", current.IP, current.Port, req.SServiceName, req.SFuncName, req.IRequestId, err.Error())
	}

	if current.Rsp() {
		current.SendTexResponse(int32(texret), p.ToBytes())
	}
}
