// 此文件为sdp2go工具自动生成,请不要手动编辑

package rpc

import (
	"bytes"
	"context"
	"fmt"
	"github.com/yellia1989/tex-go/sdp/protocol"
	"github.com/yellia1989/tex-go/service/model"
	"github.com/yellia1989/tex-go/tools/log"
	"github.com/yellia1989/tex-go/tools/net"
	"github.com/yellia1989/tex-go/tools/sdp/codec"
	"github.com/yellia1989/tex-go/tools/sdp/util"
	"time"
)

type ObjEndpoint struct {
	SObj      string `json:"sObj"`
	SDivision string `json:"sDivision"`
	SEp       string `json:"sEp"`
}

func (st *ObjEndpoint) ResetDefault() {
}
func (st *ObjEndpoint) Visit(buff *bytes.Buffer, t int) {
	util.Tab(buff, t+1, util.Fieldname("sObj")+fmt.Sprintf("%v\n", st.SObj))
	util.Tab(buff, t+1, util.Fieldname("sDivision")+fmt.Sprintf("%v\n", st.SDivision))
	util.Tab(buff, t+1, util.Fieldname("sEp")+fmt.Sprintf("%v\n", st.SEp))
}
func (st *ObjEndpoint) ReadStruct(up *codec.UnPacker) error {
	var err error
	var length uint32
	var has bool
	var ty uint32
	st.ResetDefault()
	err = up.ReadString(&st.SObj, 0, false)
	if err != nil {
		return err
	}
	err = up.ReadString(&st.SDivision, 1, false)
	if err != nil {
		return err
	}
	err = up.ReadString(&st.SEp, 2, false)
	if err != nil {
		return err
	}

	_ = length
	_ = has
	_ = ty

	return err
}
func (st *ObjEndpoint) ReadStructFromTag(up *codec.UnPacker, tag uint32, require bool) error {
	var err error
	var has bool
	var ty uint32
	st.ResetDefault()

	has, ty, err = up.SkipToTag(tag, require)
	if !has || err != nil {
		return err
	}

	if ty != codec.SdpType_StructBegin {
		return fmt.Errorf("tag:%d got wrong type %d", tag, ty)
	}

	err = st.ReadStruct(up)
	if err != nil {
		return err
	}
	err = up.SkipStruct()
	if err != nil {
		return err
	}

	_ = has
	_ = ty
	return nil
}
func (st *ObjEndpoint) WriteStruct(p *codec.Packer) error {
	var err error
	var length int
	if false || st.SObj != "" {
		err = p.WriteString(0, st.SObj)
		if err != nil {
			return err
		}
	}
	if false || st.SDivision != "" {
		err = p.WriteString(1, st.SDivision)
		if err != nil {
			return err
		}
	}
	if false || st.SEp != "" {
		err = p.WriteString(2, st.SEp)
		if err != nil {
			return err
		}
	}

	_ = length
	return err
}
func (st *ObjEndpoint) WriteStructFromTag(p *codec.Packer, tag uint32, require bool) error {
	var err error

	if require {
		err = p.WriteHeader(tag, codec.SdpType_StructBegin)
		if err != nil {
			return err
		}
		err = st.WriteStruct(p)
		if err != nil {
			return err
		}
		err = p.WriteHeader(0, codec.SdpType_StructEnd)
		if err != nil {
			return err
		}
	} else {
		p2 := codec.NewPacker()
		err = st.WriteStruct(p2)
		if err != nil {
			return err
		}
		if p2.Len() != 0 {
			err = p.WriteHeader(tag, codec.SdpType_StructBegin)
			if err != nil {
				return err
			}
			err = p.WriteData(p2.ToBytes())
			if err != nil {
				return err
			}
			err = p.WriteHeader(0, codec.SdpType_StructEnd)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type Query struct {
	proxy model.ServicePrxImpl
}

func (s *Query) SetPrxImpl(impl model.ServicePrxImpl) {
	s.proxy = impl
}
func (s *Query) SetTimeout(timeout time.Duration) {
	s.proxy.SetTimeout(timeout)
}
func (s *Query) GetEndpoints(sObj string, sDivision string, vActiveEps *[]string, vInactiveEps *[]string) (int32, error) {
	p := codec.NewPacker()
	var ret int32
	var err error
	var has bool
	var ty uint32
	var length uint32
	if true || sObj != "" {
		err = p.WriteString(1, sObj)
		if err != nil {
			return ret, err
		}
	}
	if true || sDivision != "" {
		err = p.WriteString(2, sDivision)
		if err != nil {
			return ret, err
		}
	}
	var rsp *protocol.ResponsePacket
	err = s.proxy.Invoke("getEndpoints", p.ToBytes(), &rsp)
	if err != nil {
		return ret, err
	}
	up := codec.NewUnPacker([]byte(rsp.SRspPayload))
	err = up.ReadInt32(&ret, 0, true)
	if err != nil {
		return ret, err
	}

	has, ty, err = up.SkipToTag(3, true)
	if !has || err != nil {
		return ret, err
	}
	if ty != codec.SdpType_Vector {
		return ret, fmt.Errorf("tag:%d got wrong type %d", 3, ty)
	}

	_, length, err = up.ReadNumber32()
	if err != nil {
		return ret, err
	}
	(*vActiveEps) = make([]string, length, length)
	for i := uint32(0); i < length; i++ {
		err = up.ReadString(&(*vActiveEps)[i], 0, true)
		if err != nil {
			return ret, err
		}
	}

	has, ty, err = up.SkipToTag(4, true)
	if !has || err != nil {
		return ret, err
	}
	if ty != codec.SdpType_Vector {
		return ret, fmt.Errorf("tag:%d got wrong type %d", 4, ty)
	}

	_, length, err = up.ReadNumber32()
	if err != nil {
		return ret, err
	}
	(*vInactiveEps) = make([]string, length, length)
	for i := uint32(0); i < length; i++ {
		err = up.ReadString(&(*vInactiveEps)[i], 0, true)
		if err != nil {
			return ret, err
		}
	}
	_ = has
	_ = ty
	_ = length
	return ret, nil
}
func (s *Query) AddEndpoint(sObj string, sDivision string, ep string) (int32, error) {
	p := codec.NewPacker()
	var ret int32
	var err error
	var has bool
	var ty uint32
	var length uint32
	if true || sObj != "" {
		err = p.WriteString(1, sObj)
		if err != nil {
			return ret, err
		}
	}
	if true || sDivision != "" {
		err = p.WriteString(2, sDivision)
		if err != nil {
			return ret, err
		}
	}
	if true || ep != "" {
		err = p.WriteString(3, ep)
		if err != nil {
			return ret, err
		}
	}
	var rsp *protocol.ResponsePacket
	err = s.proxy.Invoke("addEndpoint", p.ToBytes(), &rsp)
	if err != nil {
		return ret, err
	}
	up := codec.NewUnPacker([]byte(rsp.SRspPayload))
	err = up.ReadInt32(&ret, 0, true)
	if err != nil {
		return ret, err
	}
	_ = has
	_ = ty
	_ = length
	return ret, nil
}
func (s *Query) RemoveEndpoint(sObj string, sDivision string, ep string) (int32, error) {
	p := codec.NewPacker()
	var ret int32
	var err error
	var has bool
	var ty uint32
	var length uint32
	if true || sObj != "" {
		err = p.WriteString(1, sObj)
		if err != nil {
			return ret, err
		}
	}
	if true || sDivision != "" {
		err = p.WriteString(2, sDivision)
		if err != nil {
			return ret, err
		}
	}
	if true || ep != "" {
		err = p.WriteString(3, ep)
		if err != nil {
			return ret, err
		}
	}
	var rsp *protocol.ResponsePacket
	err = s.proxy.Invoke("removeEndpoint", p.ToBytes(), &rsp)
	if err != nil {
		return ret, err
	}
	up := codec.NewUnPacker([]byte(rsp.SRspPayload))
	err = up.ReadInt32(&ret, 0, true)
	if err != nil {
		return ret, err
	}
	_ = has
	_ = ty
	_ = length
	return ret, nil
}
func (s *Query) GetAllEndpoints(vObj *[]ObjEndpoint) (int32, error) {
	p := codec.NewPacker()
	var ret int32
	var err error
	var has bool
	var ty uint32
	var length uint32
	var rsp *protocol.ResponsePacket
	err = s.proxy.Invoke("getAllEndpoints", p.ToBytes(), &rsp)
	if err != nil {
		return ret, err
	}
	up := codec.NewUnPacker([]byte(rsp.SRspPayload))
	err = up.ReadInt32(&ret, 0, true)
	if err != nil {
		return ret, err
	}

	has, ty, err = up.SkipToTag(1, true)
	if !has || err != nil {
		return ret, err
	}
	if ty != codec.SdpType_Vector {
		return ret, fmt.Errorf("tag:%d got wrong type %d", 1, ty)
	}

	_, length, err = up.ReadNumber32()
	if err != nil {
		return ret, err
	}
	(*vObj) = make([]ObjEndpoint, length, length)
	for i := uint32(0); i < length; i++ {
		err = (*vObj)[i].ReadStructFromTag(up, 0, true)
		if err != nil {
			return ret, err
		}
	}
	_ = has
	_ = ty
	_ = length
	return ret, nil
}

type _QueryImpl interface {
	GetEndpoints(ctx context.Context, sObj string, sDivision string, vActiveEps *[]string, vInactiveEps *[]string) (int32, error)
	AddEndpoint(ctx context.Context, sObj string, sDivision string, ep string) (int32, error)
	RemoveEndpoint(ctx context.Context, sObj string, sDivision string, ep string) (int32, error)
	GetAllEndpoints(ctx context.Context, vObj *[]ObjEndpoint) (int32, error)
}

func _QueryGetEndpointsImpl(ctx context.Context, serviceImpl interface{}, up *codec.UnPacker, p *codec.Packer) error {
	var err error
	var length int
	impl := serviceImpl.(_QueryImpl)
	var p1 string
	err = up.ReadString(&p1, 1, true)
	if err != nil {
		return err
	}
	var p2 string
	err = up.ReadString(&p2, 2, true)
	if err != nil {
		return err
	}
	var p3 []string
	var p4 []string
	var ret int32
	ret, err = impl.GetEndpoints(ctx, p1, p2, &p3, &p4)
	if err != nil {
		return err
	}
	if true || ret != 0 {
		err = p.WriteInt32(0, ret)
		if err != nil {
			return err
		}
	}

	length = len(p3)
	if true || length != 0 {
		err = p.WriteHeader(3, codec.SdpType_Vector)
		if err != nil {
			return err
		}
		err = p.WriteNumber32(uint32(length))
		if err != nil {
			return err
		}
		for _, v := range p3 {
			if true || v != "" {
				err = p.WriteString(0, v)
				if err != nil {
					return err
				}
			}
		}
	}

	length = len(p4)
	if true || length != 0 {
		err = p.WriteHeader(4, codec.SdpType_Vector)
		if err != nil {
			return err
		}
		err = p.WriteNumber32(uint32(length))
		if err != nil {
			return err
		}
		for _, v := range p4 {
			if true || v != "" {
				err = p.WriteString(0, v)
				if err != nil {
					return err
				}
			}
		}
	}
	_ = length
	return nil
}
func _QueryAddEndpointImpl(ctx context.Context, serviceImpl interface{}, up *codec.UnPacker, p *codec.Packer) error {
	var err error
	var length int
	impl := serviceImpl.(_QueryImpl)
	var p1 string
	err = up.ReadString(&p1, 1, true)
	if err != nil {
		return err
	}
	var p2 string
	err = up.ReadString(&p2, 2, true)
	if err != nil {
		return err
	}
	var p3 string
	err = up.ReadString(&p3, 3, true)
	if err != nil {
		return err
	}
	var ret int32
	ret, err = impl.AddEndpoint(ctx, p1, p2, p3)
	if err != nil {
		return err
	}
	if true || ret != 0 {
		err = p.WriteInt32(0, ret)
		if err != nil {
			return err
		}
	}
	_ = length
	return nil
}
func _QueryRemoveEndpointImpl(ctx context.Context, serviceImpl interface{}, up *codec.UnPacker, p *codec.Packer) error {
	var err error
	var length int
	impl := serviceImpl.(_QueryImpl)
	var p1 string
	err = up.ReadString(&p1, 1, true)
	if err != nil {
		return err
	}
	var p2 string
	err = up.ReadString(&p2, 2, true)
	if err != nil {
		return err
	}
	var p3 string
	err = up.ReadString(&p3, 3, true)
	if err != nil {
		return err
	}
	var ret int32
	ret, err = impl.RemoveEndpoint(ctx, p1, p2, p3)
	if err != nil {
		return err
	}
	if true || ret != 0 {
		err = p.WriteInt32(0, ret)
		if err != nil {
			return err
		}
	}
	_ = length
	return nil
}
func _QueryGetAllEndpointsImpl(ctx context.Context, serviceImpl interface{}, up *codec.UnPacker, p *codec.Packer) error {
	var err error
	var length int
	impl := serviceImpl.(_QueryImpl)
	var p1 []ObjEndpoint
	var ret int32
	ret, err = impl.GetAllEndpoints(ctx, &p1)
	if err != nil {
		return err
	}
	if true || ret != 0 {
		err = p.WriteInt32(0, ret)
		if err != nil {
			return err
		}
	}

	length = len(p1)
	if true || length != 0 {
		err = p.WriteHeader(1, codec.SdpType_Vector)
		if err != nil {
			return err
		}
		err = p.WriteNumber32(uint32(length))
		if err != nil {
			return err
		}
		for _, v := range p1 {
			err = v.WriteStructFromTag(p, 0, true)
			if err != nil {
				return err
			}
		}
	}
	_ = length
	return nil
}

func (s *Query) Dispatch(ctx context.Context, serviceImpl interface{}, req *protocol.RequestPacket) {
	current := net.ContextGetCurrent(ctx)

	log.FDebugf("handle tex request, peer: %s:%d, obj: %s, func: %s, reqid: %d", current.IP, current.Port, req.SServiceName, req.SFuncName, req.IRequestId)

	texret := protocol.SDPSERVERUNKNOWNERR
	up := codec.NewUnPacker([]byte(req.SReqPayload))
	p := codec.NewPacker()

	var err error
	switch req.SFuncName {
	case "getEndpoints":
		err = _QueryGetEndpointsImpl(ctx, serviceImpl, up, p)
		if err != nil {
			break
		}
		texret = protocol.SDPSERVERSUCCESS
	case "addEndpoint":
		err = _QueryAddEndpointImpl(ctx, serviceImpl, up, p)
		if err != nil {
			break
		}
		texret = protocol.SDPSERVERSUCCESS
	case "removeEndpoint":
		err = _QueryRemoveEndpointImpl(ctx, serviceImpl, up, p)
		if err != nil {
			break
		}
		texret = protocol.SDPSERVERSUCCESS
	case "getAllEndpoints":
		err = _QueryGetAllEndpointsImpl(ctx, serviceImpl, up, p)
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
