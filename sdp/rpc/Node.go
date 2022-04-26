// 此文件为sdp2go工具自动生成,请不要手动编辑

package rpc

import (
	"bytes"
	"context"
	"github.com/yellia1989/tex-go/sdp/protocol"
	"github.com/yellia1989/tex-go/service/model"
	"github.com/yellia1989/tex-go/tools/log"
	"github.com/yellia1989/tex-go/tools/net"
	"github.com/yellia1989/tex-go/tools/sdp/codec"
	"github.com/yellia1989/tex-go/tools/sdp/util"
	"time"
	"fmt"
)

type PatchRequest struct {
	SFileName string `json:"sFileName" form:"sFileName"`
	SMd5      string `json:"sMd5" form:"sMd5"`
	IFileSize uint32 `json:"iFileSize" form:"iFileSize"`
}

func (st *PatchRequest) resetDefault() {
}
func (st *PatchRequest) Copy() *PatchRequest {
	ret := NewPatchRequest()
	ret.SFileName = st.SFileName
	ret.SMd5 = st.SMd5
	ret.IFileSize = st.IFileSize
	return ret
}
func NewPatchRequest() *PatchRequest {
	ret := &PatchRequest{}
	ret.resetDefault()
	return ret
}
func (st *PatchRequest) Visit(buff *bytes.Buffer, t int) {
	util.Tab(buff, t+1, util.Fieldname("sFileName")+fmt.Sprintf("%v\n", st.SFileName))
	util.Tab(buff, t+1, util.Fieldname("sMd5")+fmt.Sprintf("%v\n", st.SMd5))
	util.Tab(buff, t+1, util.Fieldname("iFileSize")+fmt.Sprintf("%v\n", st.IFileSize))
}
func (st *PatchRequest) ReadStruct(up *codec.UnPacker) error {
	var err error
	var length uint32
	var has bool
	var ty uint32
	st.resetDefault()
	err = up.ReadString(&st.SFileName, 0, false)
	if err != nil {
		return err
	}
	err = up.ReadString(&st.SMd5, 1, false)
	if err != nil {
		return err
	}
	err = up.ReadUint32(&st.IFileSize, 2, false)
	if err != nil {
		return err
	}

	_ = length
	_ = has
	_ = ty

	return err
}
func (st *PatchRequest) ReadStructFromTag(up *codec.UnPacker, tag uint32, require bool) error {
	var err error
	var has bool
	var ty uint32

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
func (st *PatchRequest) WriteStruct(p *codec.Packer) error {
	var err error
	var length uint32
	if false || st.SFileName != "" {
		err = p.WriteString(0, st.SFileName)
		if err != nil {
			return err
		}
	}
	if false || st.SMd5 != "" {
		err = p.WriteString(1, st.SMd5)
		if err != nil {
			return err
		}
	}
	if false || st.IFileSize != 0 {
		err = p.WriteUint32(2, st.IFileSize)
		if err != nil {
			return err
		}
	}

	_ = length
	return err
}
func (st *PatchRequest) WriteStructFromTag(p *codec.Packer, tag uint32, require bool) error {
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

type PatchPercent struct {
	IPercent uint32 `json:"iPercent" form:"iPercent"`
	BSuccess bool   `json:"bSuccess" form:"bSuccess"`
	SResult  string `json:"sResult" form:"sResult"`
}

func (st *PatchPercent) resetDefault() {
}
func (st *PatchPercent) Copy() *PatchPercent {
	ret := NewPatchPercent()
	ret.IPercent = st.IPercent
	ret.BSuccess = st.BSuccess
	ret.SResult = st.SResult
	return ret
}
func NewPatchPercent() *PatchPercent {
	ret := &PatchPercent{}
	ret.resetDefault()
	return ret
}
func (st *PatchPercent) Visit(buff *bytes.Buffer, t int) {
	util.Tab(buff, t+1, util.Fieldname("iPercent")+fmt.Sprintf("%v\n", st.IPercent))
	util.Tab(buff, t+1, util.Fieldname("bSuccess")+fmt.Sprintf("%v\n", st.BSuccess))
	util.Tab(buff, t+1, util.Fieldname("sResult")+fmt.Sprintf("%v\n", st.SResult))
}
func (st *PatchPercent) ReadStruct(up *codec.UnPacker) error {
	var err error
	var length uint32
	var has bool
	var ty uint32
	st.resetDefault()
	err = up.ReadUint32(&st.IPercent, 0, false)
	if err != nil {
		return err
	}
	err = up.ReadBool(&st.BSuccess, 1, false)
	if err != nil {
		return err
	}
	err = up.ReadString(&st.SResult, 2, false)
	if err != nil {
		return err
	}

	_ = length
	_ = has
	_ = ty

	return err
}
func (st *PatchPercent) ReadStructFromTag(up *codec.UnPacker, tag uint32, require bool) error {
	var err error
	var has bool
	var ty uint32

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
func (st *PatchPercent) WriteStruct(p *codec.Packer) error {
	var err error
	var length uint32
	if false || st.IPercent != 0 {
		err = p.WriteUint32(0, st.IPercent)
		if err != nil {
			return err
		}
	}
	if false || st.BSuccess != false {
		err = p.WriteBool(1, st.BSuccess)
		if err != nil {
			return err
		}
	}
	if false || st.SResult != "" {
		err = p.WriteString(2, st.SResult)
		if err != nil {
			return err
		}
	}

	_ = length
	return err
}
func (st *PatchPercent) WriteStructFromTag(p *codec.Packer, tag uint32, require bool) error {
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

type Node struct {
	proxy model.ServicePrxImpl
}

func (s *Node) SetPrxImpl(impl model.ServicePrxImpl) {
	s.proxy = impl
}
func (s *Node) SetTimeout(timeout time.Duration) {
	s.proxy.SetTimeout(timeout)
}
func (s *Node) Stop(sApp string, sServer string, sDivision string, sResult *string) (int32, error) {
	p := codec.NewPacker()
	var ret int32
	var err error
	var has bool
	var ty uint32
	var length uint32
	if true || sApp != "" {
		err = p.WriteString(1, sApp)
		if err != nil {
			return ret, err
		}
	}
	if true || sServer != "" {
		err = p.WriteString(2, sServer)
		if err != nil {
			return ret, err
		}
	}
	if true || sDivision != "" {
		err = p.WriteString(3, sDivision)
		if err != nil {
			return ret, err
		}
	}
	var rsp *protocol.ResponsePacket
	err = s.proxy.Invoke("stop", p.ToBytes(), &rsp)
	if err != nil {
		return ret, err
	}
	up := codec.NewUnPacker([]byte(rsp.SRspPayload))
	err = up.ReadInt32(&ret, 0, true)
	if err != nil {
		return ret, err
	}
	err = up.ReadString(&(*sResult), 4, true)
	if err != nil {
		return ret, err
	}
	_ = has
	_ = ty
	_ = length
	return ret, nil
}
func (s *Node) Start(sApp string, sServer string, sDivsioin string, sResult *string) (int32, error) {
	p := codec.NewPacker()
	var ret int32
	var err error
	var has bool
	var ty uint32
	var length uint32
	if true || sApp != "" {
		err = p.WriteString(1, sApp)
		if err != nil {
			return ret, err
		}
	}
	if true || sServer != "" {
		err = p.WriteString(2, sServer)
		if err != nil {
			return ret, err
		}
	}
	if true || sDivsioin != "" {
		err = p.WriteString(3, sDivsioin)
		if err != nil {
			return ret, err
		}
	}
	var rsp *protocol.ResponsePacket
	err = s.proxy.Invoke("start", p.ToBytes(), &rsp)
	if err != nil {
		return ret, err
	}
	up := codec.NewUnPacker([]byte(rsp.SRspPayload))
	err = up.ReadInt32(&ret, 0, true)
	if err != nil {
		return ret, err
	}
	err = up.ReadString(&(*sResult), 4, true)
	if err != nil {
		return ret, err
	}
	_ = has
	_ = ty
	_ = length
	return ret, nil
}
func (s *Node) Restart(sApp string, sServer string, sDivision string, sResult *string) (int32, error) {
	p := codec.NewPacker()
	var ret int32
	var err error
	var has bool
	var ty uint32
	var length uint32
	if true || sApp != "" {
		err = p.WriteString(1, sApp)
		if err != nil {
			return ret, err
		}
	}
	if true || sServer != "" {
		err = p.WriteString(2, sServer)
		if err != nil {
			return ret, err
		}
	}
	if true || sDivision != "" {
		err = p.WriteString(3, sDivision)
		if err != nil {
			return ret, err
		}
	}
	var rsp *protocol.ResponsePacket
	err = s.proxy.Invoke("restart", p.ToBytes(), &rsp)
	if err != nil {
		return ret, err
	}
	up := codec.NewUnPacker([]byte(rsp.SRspPayload))
	err = up.ReadInt32(&ret, 0, true)
	if err != nil {
		return ret, err
	}
	err = up.ReadString(&(*sResult), 4, true)
	if err != nil {
		return ret, err
	}
	_ = has
	_ = ty
	_ = length
	return ret, nil
}
func (s *Node) Patch(sApp string, sServer string, sDivision string, patchReq PatchRequest, sResult *string) (int32, error) {
	p := codec.NewPacker()
	var ret int32
	var err error
	var has bool
	var ty uint32
	var length uint32
	if true || sApp != "" {
		err = p.WriteString(1, sApp)
		if err != nil {
			return ret, err
		}
	}
	if true || sServer != "" {
		err = p.WriteString(2, sServer)
		if err != nil {
			return ret, err
		}
	}
	if true || sDivision != "" {
		err = p.WriteString(3, sDivision)
		if err != nil {
			return ret, err
		}
	}
	err = patchReq.WriteStructFromTag(p, 4, true)
	if err != nil {
		return ret, err
	}
	var rsp *protocol.ResponsePacket
	err = s.proxy.Invoke("patch", p.ToBytes(), &rsp)
	if err != nil {
		return ret, err
	}
	up := codec.NewUnPacker([]byte(rsp.SRspPayload))
	err = up.ReadInt32(&ret, 0, true)
	if err != nil {
		return ret, err
	}
	err = up.ReadString(&(*sResult), 5, true)
	if err != nil {
		return ret, err
	}
	_ = has
	_ = ty
	_ = length
	return ret, nil
}
func (s *Node) GetPatchPercent(sApp string, sServer string, sDivision string, patchPercent *PatchPercent) (int32, error) {
	p := codec.NewPacker()
	var ret int32
	var err error
	var has bool
	var ty uint32
	var length uint32
	if true || sApp != "" {
		err = p.WriteString(1, sApp)
		if err != nil {
			return ret, err
		}
	}
	if true || sServer != "" {
		err = p.WriteString(2, sServer)
		if err != nil {
			return ret, err
		}
	}
	if true || sDivision != "" {
		err = p.WriteString(3, sDivision)
		if err != nil {
			return ret, err
		}
	}
	var rsp *protocol.ResponsePacket
	err = s.proxy.Invoke("getPatchPercent", p.ToBytes(), &rsp)
	if err != nil {
		return ret, err
	}
	up := codec.NewUnPacker([]byte(rsp.SRspPayload))
	err = up.ReadInt32(&ret, 0, true)
	if err != nil {
		return ret, err
	}
	err = (*patchPercent).ReadStructFromTag(up, 4, true)
	if err != nil {
		return ret, err
	}
	_ = has
	_ = ty
	_ = length
	return ret, nil
}
func (s *Node) Notify(sApp string, sServer string, sDivision string, sCmd string, sResult *string) (int32, error) {
	p := codec.NewPacker()
	var ret int32
	var err error
	var has bool
	var ty uint32
	var length uint32
	if true || sApp != "" {
		err = p.WriteString(1, sApp)
		if err != nil {
			return ret, err
		}
	}
	if true || sServer != "" {
		err = p.WriteString(2, sServer)
		if err != nil {
			return ret, err
		}
	}
	if true || sDivision != "" {
		err = p.WriteString(3, sDivision)
		if err != nil {
			return ret, err
		}
	}
	if true || sCmd != "" {
		err = p.WriteString(4, sCmd)
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
	err = up.ReadString(&(*sResult), 5, true)
	if err != nil {
		return ret, err
	}
	_ = has
	_ = ty
	_ = length
	return ret, nil
}
func (s *Node) KeepAlive(sApp string, sServer string, sDivision string, iPid uint32, sAdapterName string, bIniting bool) (int32, error) {
	p := codec.NewPacker()
	var ret int32
	var err error
	var has bool
	var ty uint32
	var length uint32
	if true || sApp != "" {
		err = p.WriteString(1, sApp)
		if err != nil {
			return ret, err
		}
	}
	if true || sServer != "" {
		err = p.WriteString(2, sServer)
		if err != nil {
			return ret, err
		}
	}
	if true || sDivision != "" {
		err = p.WriteString(3, sDivision)
		if err != nil {
			return ret, err
		}
	}
	if true || iPid != 0 {
		err = p.WriteUint32(4, iPid)
		if err != nil {
			return ret, err
		}
	}
	if true || sAdapterName != "" {
		err = p.WriteString(5, sAdapterName)
		if err != nil {
			return ret, err
		}
	}
	if true || bIniting != false {
		err = p.WriteBool(6, bIniting)
		if err != nil {
			return ret, err
		}
	}
	var rsp *protocol.ResponsePacket
	err = s.proxy.Invoke("keepAlive", p.ToBytes(), &rsp)
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

type _NodeImpl interface {
	Stop(ctx context.Context, sApp string, sServer string, sDivision string, sResult *string) (int32, error)
	Start(ctx context.Context, sApp string, sServer string, sDivsioin string, sResult *string) (int32, error)
	Restart(ctx context.Context, sApp string, sServer string, sDivision string, sResult *string) (int32, error)
	Patch(ctx context.Context, sApp string, sServer string, sDivision string, patchReq PatchRequest, sResult *string) (int32, error)
	GetPatchPercent(ctx context.Context, sApp string, sServer string, sDivision string, patchPercent *PatchPercent) (int32, error)
	Notify(ctx context.Context, sApp string, sServer string, sDivision string, sCmd string, sResult *string) (int32, error)
	KeepAlive(ctx context.Context, sApp string, sServer string, sDivision string, iPid uint32, sAdapterName string, bIniting bool) (int32, error)
}

func _NodeStopImpl(ctx context.Context, serviceImpl interface{}, up *codec.UnPacker, p *codec.Packer) error {
	var err error
	var length uint32
	var ty uint32
	var has bool
	impl := serviceImpl.(_NodeImpl)
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
	var p4 string
	var ret int32
	ret, err = impl.Stop(ctx, p1, p2, p3, &p4)
	if err != nil {
		return err
	}
	if true || ret != 0 {
		err = p.WriteInt32(0, ret)
		if err != nil {
			return err
		}
	}
	if true || p4 != "" {
		err = p.WriteString(4, p4)
		if err != nil {
			return err
		}
	}
	_ = length
	_ = ty
	_ = has
	return nil
}
func _NodeStartImpl(ctx context.Context, serviceImpl interface{}, up *codec.UnPacker, p *codec.Packer) error {
	var err error
	var length uint32
	var ty uint32
	var has bool
	impl := serviceImpl.(_NodeImpl)
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
	var p4 string
	var ret int32
	ret, err = impl.Start(ctx, p1, p2, p3, &p4)
	if err != nil {
		return err
	}
	if true || ret != 0 {
		err = p.WriteInt32(0, ret)
		if err != nil {
			return err
		}
	}
	if true || p4 != "" {
		err = p.WriteString(4, p4)
		if err != nil {
			return err
		}
	}
	_ = length
	_ = ty
	_ = has
	return nil
}
func _NodeRestartImpl(ctx context.Context, serviceImpl interface{}, up *codec.UnPacker, p *codec.Packer) error {
	var err error
	var length uint32
	var ty uint32
	var has bool
	impl := serviceImpl.(_NodeImpl)
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
	var p4 string
	var ret int32
	ret, err = impl.Restart(ctx, p1, p2, p3, &p4)
	if err != nil {
		return err
	}
	if true || ret != 0 {
		err = p.WriteInt32(0, ret)
		if err != nil {
			return err
		}
	}
	if true || p4 != "" {
		err = p.WriteString(4, p4)
		if err != nil {
			return err
		}
	}
	_ = length
	_ = ty
	_ = has
	return nil
}
func _NodePatchImpl(ctx context.Context, serviceImpl interface{}, up *codec.UnPacker, p *codec.Packer) error {
	var err error
	var length uint32
	var ty uint32
	var has bool
	impl := serviceImpl.(_NodeImpl)
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
	var p4 PatchRequest
	err = p4.ReadStructFromTag(up, 4, true)
	if err != nil {
		return err
	}
	var p5 string
	var ret int32
	ret, err = impl.Patch(ctx, p1, p2, p3, p4, &p5)
	if err != nil {
		return err
	}
	if true || ret != 0 {
		err = p.WriteInt32(0, ret)
		if err != nil {
			return err
		}
	}
	if true || p5 != "" {
		err = p.WriteString(5, p5)
		if err != nil {
			return err
		}
	}
	_ = length
	_ = ty
	_ = has
	return nil
}
func _NodeGetPatchPercentImpl(ctx context.Context, serviceImpl interface{}, up *codec.UnPacker, p *codec.Packer) error {
	var err error
	var length uint32
	var ty uint32
	var has bool
	impl := serviceImpl.(_NodeImpl)
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
	var p4 PatchPercent
	var ret int32
	ret, err = impl.GetPatchPercent(ctx, p1, p2, p3, &p4)
	if err != nil {
		return err
	}
	if true || ret != 0 {
		err = p.WriteInt32(0, ret)
		if err != nil {
			return err
		}
	}
	err = p4.WriteStructFromTag(p, 4, true)
	if err != nil {
		return err
	}
	_ = length
	_ = ty
	_ = has
	return nil
}
func _NodeNotifyImpl(ctx context.Context, serviceImpl interface{}, up *codec.UnPacker, p *codec.Packer) error {
	var err error
	var length uint32
	var ty uint32
	var has bool
	impl := serviceImpl.(_NodeImpl)
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
	var p4 string
	err = up.ReadString(&p4, 4, true)
	if err != nil {
		return err
	}
	var p5 string
	var ret int32
	ret, err = impl.Notify(ctx, p1, p2, p3, p4, &p5)
	if err != nil {
		return err
	}
	if true || ret != 0 {
		err = p.WriteInt32(0, ret)
		if err != nil {
			return err
		}
	}
	if true || p5 != "" {
		err = p.WriteString(5, p5)
		if err != nil {
			return err
		}
	}
	_ = length
	_ = ty
	_ = has
	return nil
}
func _NodeKeepAliveImpl(ctx context.Context, serviceImpl interface{}, up *codec.UnPacker, p *codec.Packer) error {
	var err error
	var length uint32
	var ty uint32
	var has bool
	impl := serviceImpl.(_NodeImpl)
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
	var p4 uint32
	err = up.ReadUint32(&p4, 4, true)
	if err != nil {
		return err
	}
	var p5 string
	err = up.ReadString(&p5, 5, true)
	if err != nil {
		return err
	}
	var p6 bool
	err = up.ReadBool(&p6, 6, true)
	if err != nil {
		return err
	}
	var ret int32
	ret, err = impl.KeepAlive(ctx, p1, p2, p3, p4, p5, p6)
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
	_ = ty
	_ = has
	return nil
}

func (s *Node) Dispatch(ctx context.Context, serviceImpl interface{}, req *protocol.RequestPacket) {
	current := net.ContextGetCurrent(ctx)

	log.FDebugf("handle tex request, peer: %s:%d, obj: %s, func: %s, reqid: %d", current.IP, current.Port, req.SServiceName, req.SFuncName, req.IRequestId)

	texret := protocol.SDPSERVERUNKNOWNERR
	up := codec.NewUnPacker([]byte(req.SReqPayload))
	p := codec.NewPacker()

	var err error
	switch req.SFuncName {
	case "stop":
		err = _NodeStopImpl(ctx, serviceImpl, up, p)
		if err != nil {
			break
		}
		texret = protocol.SDPSERVERSUCCESS
	case "start":
		err = _NodeStartImpl(ctx, serviceImpl, up, p)
		if err != nil {
			break
		}
		texret = protocol.SDPSERVERSUCCESS
	case "restart":
		err = _NodeRestartImpl(ctx, serviceImpl, up, p)
		if err != nil {
			break
		}
		texret = protocol.SDPSERVERSUCCESS
	case "patch":
		err = _NodePatchImpl(ctx, serviceImpl, up, p)
		if err != nil {
			break
		}
		texret = protocol.SDPSERVERSUCCESS
	case "getPatchPercent":
		err = _NodeGetPatchPercentImpl(ctx, serviceImpl, up, p)
		if err != nil {
			break
		}
		texret = protocol.SDPSERVERSUCCESS
	case "notify":
		err = _NodeNotifyImpl(ctx, serviceImpl, up, p)
		if err != nil {
			break
		}
		texret = protocol.SDPSERVERSUCCESS
	case "keepAlive":
		err = _NodeKeepAliveImpl(ctx, serviceImpl, up, p)
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
