// 此文件为sdp2go工具自动生成,请不要手动编辑

package test

import (
	"fmt"
	"github.com/yellia1989/tex-go/tools/sdp/codec"
	"github.com/yellia1989/tex-go/tools/sdp/sdp2go/test2"
)

type RequireStruct struct {
	Ss test2.SimpleStruct `json:"ss"`
}

func (st *RequireStruct) ResetDefault() {
}
func (st *RequireStruct) ReadStruct(up *codec.UnPacker) error {
	var err error
	var length uint32
	var has bool
	var ty uint32
	st.ResetDefault()
	err = st.Ss.ReadStructFromTag(up, 0, true)
	if err != nil {
		return err
	}

	_ = length
	_ = has
	_ = ty

	return err
}
func (st *RequireStruct) ReadStructFromTag(up *codec.UnPacker, tag uint32, require bool) error {
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
func (st *RequireStruct) WriteStruct(p *codec.Packer) error {
	var err error
	var length int
	err = st.Ss.WriteStructFromTag(p, 0)
	if err != nil {
		return err
	}

	_ = length
	return err
}
func (st *RequireStruct) WriteStructFromTag(p *codec.Packer, tag uint32) error {
	var err error

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

	return nil
}

type DefaultStruct struct {
	B  bool   `json:"b"`
	By int8   `json:"by"`
	S  int16  `json:"s"`
	I  int32  `json:"i"`
	L  int64  `json:"l"`
	Ss string `json:"ss"`
}

func (st *DefaultStruct) ResetDefault() {
	st.B = true
	st.By = 1
	st.S = 10
	st.I = 1
	st.L = 0x0FFFFFFFFFFFFFFF
	st.Ss = "yellia"
}
func (st *DefaultStruct) ReadStruct(up *codec.UnPacker) error {
	var err error
	var length uint32
	var has bool
	var ty uint32
	st.ResetDefault()
	err = up.ReadBool(&st.B, 0, false)
	if err != nil {
		return err
	}
	err = up.ReadInt8(&st.By, 1, false)
	if err != nil {
		return err
	}
	err = up.ReadInt16(&st.S, 2, false)
	if err != nil {
		return err
	}
	err = up.ReadInt32(&st.I, 3, false)
	if err != nil {
		return err
	}
	err = up.ReadInt64(&st.L, 4, false)
	if err != nil {
		return err
	}
	err = up.ReadString(&st.Ss, 5, false)
	if err != nil {
		return err
	}

	_ = length
	_ = has
	_ = ty

	return err
}
func (st *DefaultStruct) ReadStructFromTag(up *codec.UnPacker, tag uint32, require bool) error {
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
func (st *DefaultStruct) WriteStruct(p *codec.Packer) error {
	var err error
	var length int
	if st.B != true {
		err = p.WriteBool(0, st.B)
		if err != nil {
			return err
		}
	}
	if st.By != 1 {
		err = p.WriteInt8(1, st.By)
		if err != nil {
			return err
		}
	}
	if st.S != 10 {
		err = p.WriteInt16(2, st.S)
		if err != nil {
			return err
		}
	}
	if st.I != 1 {
		err = p.WriteInt32(3, st.I)
		if err != nil {
			return err
		}
	}
	if st.L != 0x0FFFFFFFFFFFFFFF {
		err = p.WriteInt64(4, st.L)
		if err != nil {
			return err
		}
	}
	if st.Ss != "yellia" {
		err = p.WriteString(5, st.Ss)
		if err != nil {
			return err
		}
	}

	_ = length
	return err
}
func (st *DefaultStruct) WriteStructFromTag(p *codec.Packer, tag uint32) error {
	var err error

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

	return nil
}
