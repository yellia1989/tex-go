// 此文件为sdp2go工具自动生成,请不要手动编辑

package rpc

import (
	"context"
	"github.com/yellia1989/tex-go/sdp/protocol"
	"github.com/yellia1989/tex-go/service/model"
	"github.com/yellia1989/tex-go/tools/log"
	"github.com/yellia1989/tex-go/tools/net"
	"github.com/yellia1989/tex-go/tools/sdp/codec"
	"time"
)

type Admin struct {
	proxy model.ServicePrxImpl
}

func (s *Admin) SetPrxImpl(impl model.ServicePrxImpl) {
	s.proxy = impl
}
func (s *Admin) SetTimeout(timeout time.Duration) {
	s.proxy.SetTimeout(timeout)
}
func (s *Admin) Shutdown() error {
	p := codec.NewPacker()
	var err error
	var has bool
	var ty uint32
	var length uint32
	var rsp *protocol.ResponsePacket
	err = s.proxy.Invoke("shutdown", p.ToBytes(), &rsp)
	if err != nil {
		return err
	}
	_ = has
	_ = ty
	_ = length
	return nil
}
func (s *Admin) Notify(sCmd string, sResult *string) (int32, error) {
	p := codec.NewPacker()
	var ret int32
	var err error
	var has bool
	var ty uint32
	var length uint32
	if true || sCmd != "" {
		err = p.WriteString(1, sCmd)
		if err != nil {
			return ret, err
		}
	}
	var rsp *protocol.ResponsePacket
	err = s.proxy.Invoke("notify", p.ToBytes(), &rsp)
	if err != nil {
		return ret, err
	}
	up := codec.NewUnPacker([]byte(rsp.SRspPayload))
	err = up.ReadInt32(&ret, 0, true)
	if err != nil {
		return ret, err
	}
	err = up.ReadString(&(*sResult), 2, true)
	if err != nil {
		return ret, err
	}
	_ = has
	_ = ty
	_ = length
	return ret, nil
}

type _AdminImpl interface {
	Shutdown(ctx context.Context) error
	Notify(ctx context.Context, sCmd string, sResult *string) (int32, error)
}

func _AdminShutdownImpl(ctx context.Context, serviceImpl interface{}, up *codec.UnPacker, p *codec.Packer) error {
	var err error
	var length uint32
	var ty uint32
	var has bool
	impl := serviceImpl.(_AdminImpl)
	err = impl.Shutdown(ctx)
	if err != nil {
		return err
	}
	_ = length
	_ = ty
	_ = has
	return nil
}
func _AdminNotifyImpl(ctx context.Context, serviceImpl interface{}, up *codec.UnPacker, p *codec.Packer) error {
	var err error
	var length uint32
	var ty uint32
	var has bool
	impl := serviceImpl.(_AdminImpl)
	var p1 string
	err = up.ReadString(&p1, 1, true)
	if err != nil {
		return err
	}
	var p2 string
	var ret int32
	ret, err = impl.Notify(ctx, p1, &p2)
	if err != nil {
		return err
	}
	if true || ret != 0 {
		err = p.WriteInt32(0, ret)
		if err != nil {
			return err
		}
	}
	if true || p2 != "" {
		err = p.WriteString(2, p2)
		if err != nil {
			return err
		}
	}
	_ = length
	_ = ty
	_ = has
	return nil
}

func (s *Admin) Dispatch(ctx context.Context, serviceImpl interface{}, req *protocol.RequestPacket) {
	current := net.ContextGetCurrent(ctx)

	log.FDebugf("handle tex request, peer: %s:%d, obj: %s, func: %s, reqid: %d", current.IP, current.Port, req.SServiceName, req.SFuncName, req.IRequestId)

	texret := protocol.SDPSERVERUNKNOWNERR
	up := codec.NewUnPacker([]byte(req.SReqPayload))
	p := codec.NewPacker()

	var err error
	switch req.SFuncName {
	case "shutdown":
		err = _AdminShutdownImpl(ctx, serviceImpl, up, p)
		if err != nil {
			break
		}
		texret = protocol.SDPSERVERSUCCESS
	case "notify":
		err = _AdminNotifyImpl(ctx, serviceImpl, up, p)
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
