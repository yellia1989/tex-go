// 此文件为sdp2go工具自动生成,请不要手动编辑

package protocol

import (
	"fmt"
	"github.com/yellia1989/tex-go/tools/sdp/codec"
)

type RequestPacket struct {
	BIsOneWay    bool              `json:"bIsOneWay"`
	IRequestId   uint32            `json:"iRequestId"`
	SServiceName string            `json:"sServiceName"`
	SFuncName    string            `json:"sFuncName"`
	SReqPayload  string            `json:"sReqPayload"`
	ITimeout     uint32            `json:"iTimeout"`
	Context      map[string]string `json:"context"`
}

func (st *RequestPacket) ResetDefault() {
}
func (st *RequestPacket) ReadStruct(up *codec.UnPacker) error {
	var err error
	var length uint32
	var has bool
	var ty uint32
	st.ResetDefault()
	err = up.ReadBool(&st.BIsOneWay, 0, false)
	if err != nil {
		return err
	}
	err = up.ReadUint32(&st.IRequestId, 1, false)
	if err != nil {
		return err
	}
	err = up.ReadString(&st.SServiceName, 2, false)
	if err != nil {
		return err
	}
	err = up.ReadString(&st.SFuncName, 3, false)
	if err != nil {
		return err
	}
	err = up.ReadString(&st.SReqPayload, 4, false)
	if err != nil {
		return err
	}
	err = up.ReadUint32(&st.ITimeout, 5, false)
	if err != nil {
		return err
	}

	has, ty, err = up.SkipToTag(6, false)
	if !has || err != nil {
		return err
	}
	if ty != codec.SdpType_Map {
		return fmt.Errorf("tag:%d got wrong type %d", 6, ty)
	}

	_, length, err = up.ReadNumber32()
	if err != nil {
		return err
	}
	st.Context = make(map[string]string)
	for i := uint32(0); i < length; i++ {
		var k string
		err = up.ReadString(&k, 0, true)
		if err != nil {
			return err
		}
		var v string
		err = up.ReadString(&v, 0, true)
		if err != nil {
			return err
		}
		st.Context[k] = v
	}

	_ = length
	_ = has
	_ = ty

	return err
}
func (st *RequestPacket) ReadStructFromTag(up *codec.UnPacker, tag uint32, require bool) error {
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
func (st *RequestPacket) WriteStruct(p *codec.Packer) error {
	var err error
	var length int
	if false || st.BIsOneWay != false {
		err = p.WriteBool(0, st.BIsOneWay)
		if err != nil {
			return err
		}
	}
	if false || st.IRequestId != 0 {
		err = p.WriteUint32(1, st.IRequestId)
		if err != nil {
			return err
		}
	}
	if false || st.SServiceName != "" {
		err = p.WriteString(2, st.SServiceName)
		if err != nil {
			return err
		}
	}
	if false || st.SFuncName != "" {
		err = p.WriteString(3, st.SFuncName)
		if err != nil {
			return err
		}
	}
	if false || st.SReqPayload != "" {
		err = p.WriteString(4, st.SReqPayload)
		if err != nil {
			return err
		}
	}
	if false || st.ITimeout != 0 {
		err = p.WriteUint32(5, st.ITimeout)
		if err != nil {
			return err
		}
	}

	length = len(st.Context)
	if false || length != 0 {
		err = p.WriteHeader(6, codec.SdpType_Map)
		if err != nil {
			return err
		}
		err = p.WriteNumber32(uint32(length))
		if err != nil {
			return err
		}
		for _k, _v := range st.Context {
			if true || _k != "" {
				err = p.WriteString(0, _k)
				if err != nil {
					return err
				}
			}
			if true || _v != "" {
				err = p.WriteString(0, _v)
				if err != nil {
					return err
				}
			}
		}
	}

	_ = length
	return err
}
func (st *RequestPacket) WriteStructFromTag(p *codec.Packer, tag uint32, require bool) error {
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

type ResponsePacket struct {
	IRet        int32             `json:"iRet"`
	IRequestId  uint32            `json:"iRequestId"`
	SRspPayload string            `json:"sRspPayload"`
	Context     map[string]string `json:"context"`
}

func (st *ResponsePacket) ResetDefault() {
}
func (st *ResponsePacket) ReadStruct(up *codec.UnPacker) error {
	var err error
	var length uint32
	var has bool
	var ty uint32
	st.ResetDefault()
	err = up.ReadInt32(&st.IRet, 0, false)
	if err != nil {
		return err
	}
	err = up.ReadUint32(&st.IRequestId, 1, false)
	if err != nil {
		return err
	}
	err = up.ReadString(&st.SRspPayload, 2, false)
	if err != nil {
		return err
	}

	has, ty, err = up.SkipToTag(3, false)
	if !has || err != nil {
		return err
	}
	if ty != codec.SdpType_Map {
		return fmt.Errorf("tag:%d got wrong type %d", 3, ty)
	}

	_, length, err = up.ReadNumber32()
	if err != nil {
		return err
	}
	st.Context = make(map[string]string)
	for i := uint32(0); i < length; i++ {
		var k string
		err = up.ReadString(&k, 0, true)
		if err != nil {
			return err
		}
		var v string
		err = up.ReadString(&v, 0, true)
		if err != nil {
			return err
		}
		st.Context[k] = v
	}

	_ = length
	_ = has
	_ = ty

	return err
}
func (st *ResponsePacket) ReadStructFromTag(up *codec.UnPacker, tag uint32, require bool) error {
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
func (st *ResponsePacket) WriteStruct(p *codec.Packer) error {
	var err error
	var length int
	if false || st.IRet != 0 {
		err = p.WriteInt32(0, st.IRet)
		if err != nil {
			return err
		}
	}
	if false || st.IRequestId != 0 {
		err = p.WriteUint32(1, st.IRequestId)
		if err != nil {
			return err
		}
	}
	if false || st.SRspPayload != "" {
		err = p.WriteString(2, st.SRspPayload)
		if err != nil {
			return err
		}
	}

	length = len(st.Context)
	if false || length != 0 {
		err = p.WriteHeader(3, codec.SdpType_Map)
		if err != nil {
			return err
		}
		err = p.WriteNumber32(uint32(length))
		if err != nil {
			return err
		}
		for _k, _v := range st.Context {
			if true || _k != "" {
				err = p.WriteString(0, _k)
				if err != nil {
					return err
				}
			}
			if true || _v != "" {
				err = p.WriteString(0, _v)
				if err != nil {
					return err
				}
			}
		}
	}

	_ = length
	return err
}
func (st *ResponsePacket) WriteStructFromTag(p *codec.Packer, tag uint32, require bool) error {
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
