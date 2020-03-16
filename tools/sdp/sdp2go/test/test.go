// 此文件为sdp2go工具自动生成,请不要手动编辑

package test

import (
	"fmt"
	"github.com/yellia1989/tex-go/tools/sdp/codec"
)

type SimpleStruct struct {
	B  bool            `json:"b"`
	By byte            `json:"by"`
	S  int16           `json:"s"`
	Us uint16          `json:"us"`
	I  int32           `json:"i"`
	Ui uint32          `json:"ui"`
	L  int64           `json:"l"`
	Ul uint64          `json:"ul"`
	F  float32         `json:"f"`
	D  float64         `json:"d"`
	Ss string          `json:"ss"`
	Vi []int32         `json:"vi"`
	Mi map[int32]int32 `json:"mi"`
}

func (st *SimpleStruct) ResetDefault() {
}
func (st *SimpleStruct) ReadStruct(up *codec.UnPacker) error {
	var err error
	var length uint32
	var has bool
	var ty uint32
	st.ResetDefault()
	err = up.ReadBool(&st.B, 0, false)
	if err != nil {
		return err
	}
	err = up.ReadByte(&st.By, 1, false)
	if err != nil {
		return err
	}
	err = up.ReadInt16(&st.S, 2, false)
	if err != nil {
		return err
	}
	err = up.ReadUint16(&st.Us, 3, false)
	if err != nil {
		return err
	}
	err = up.ReadInt32(&st.I, 4, false)
	if err != nil {
		return err
	}
	err = up.ReadUint32(&st.Ui, 5, false)
	if err != nil {
		return err
	}
	err = up.ReadInt64(&st.L, 6, false)
	if err != nil {
		return err
	}
	err = up.ReadUint64(&st.Ul, 7, false)
	if err != nil {
		return err
	}
	err = up.ReadFloat32(&st.F, 8, false)
	if err != nil {
		return err
	}
	err = up.ReadFloat64(&st.D, 9, false)
	if err != nil {
		return err
	}
	err = up.ReadString(&st.Ss, 10, false)
	if err != nil {
		return err
	}

	has, ty, err = up.SkipToTag(11, false)
	if !has || err != nil {
		return err
	}
	if ty != codec.SdpType_Vector {
		return fmt.Errorf("tag:%d got wrong type %d", 11, ty)
	}

	_, length, err = up.ReadNumber32()
	if err != nil {
		return err
	}
	st.Vi = make([]int32, length, length)
	for i := uint32(0); i < length; i++ {
		err = up.ReadInt32(&st.Vi[i], 0, true)
		if err != nil {
			return err
		}
	}

	has, ty, err = up.SkipToTag(12, false)
	if !has || err != nil {
		return err
	}
	if ty != codec.SdpType_Map {
		return fmt.Errorf("tag:%d got wrong type %d", 12, ty)
	}

	_, length, err = up.ReadNumber32()
	if err != nil {
		return err
	}
	st.Mi = make(map[int32]int32)
	for i := uint32(0); i < length; i++ {
		var k int32
		err = up.ReadInt32(&k, 0, true)
		if err != nil {
			return err
		}
		var v int32
		err = up.ReadInt32(&v, 0, true)
		if err != nil {
			return err
		}
		st.Mi[k] = v
	}

	_ = length
	_ = has
	_ = ty

	return err
}
func (st *SimpleStruct) ReadStructFromTag(up *codec.UnPacker, tag uint32, require bool) error {
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
func (st *SimpleStruct) WriteStruct(p *codec.Packer) error {
	var err error
	var length int
	if false || st.B != false {
		err = p.WriteBool(0, st.B)
		if err != nil {
			return err
		}
	}
	if false || st.By != 0 {
		err = p.WriteByte(1, st.By)
		if err != nil {
			return err
		}
	}
	if false || st.S != 0 {
		err = p.WriteInt16(2, st.S)
		if err != nil {
			return err
		}
	}
	if false || st.Us != 0 {
		err = p.WriteUint16(3, st.Us)
		if err != nil {
			return err
		}
	}
	if false || st.I != 0 {
		err = p.WriteInt32(4, st.I)
		if err != nil {
			return err
		}
	}
	if false || st.Ui != 0 {
		err = p.WriteUint32(5, st.Ui)
		if err != nil {
			return err
		}
	}
	if false || st.L != 0 {
		err = p.WriteInt64(6, st.L)
		if err != nil {
			return err
		}
	}
	if false || st.Ul != 0 {
		err = p.WriteUint64(7, st.Ul)
		if err != nil {
			return err
		}
	}
	if false || st.F != 0 {
		err = p.WriteFloat32(8, st.F)
		if err != nil {
			return err
		}
	}
	if false || st.D != 0 {
		err = p.WriteFloat64(9, st.D)
		if err != nil {
			return err
		}
	}
	if false || st.Ss != "" {
		err = p.WriteString(10, st.Ss)
		if err != nil {
			return err
		}
	}

	length = len(st.Vi)
	if false || length != 0 {
		err = p.WriteHeader(11, codec.SdpType_Vector)
		if err != nil {
			return err
		}
		err = p.WriteNumber32(uint32(length))
		if err != nil {
			return err
		}
		for _, v := range st.Vi {
			if true || v != 0 {
				err = p.WriteInt32(0, v)
				if err != nil {
					return err
				}
			}
		}
	}

	length = len(st.Mi)
	if false || length != 0 {
		err = p.WriteHeader(12, codec.SdpType_Map)
		if err != nil {
			return err
		}
		err = p.WriteNumber32(uint32(length))
		if err != nil {
			return err
		}
		for _k, _v := range st.Mi {
			if true || _k != 0 {
				err = p.WriteInt32(0, _k)
				if err != nil {
					return err
				}
			}
			if true || _v != 0 {
				err = p.WriteInt32(0, _v)
				if err != nil {
					return err
				}
			}
		}
	}

	_ = length
	return err
}
func (st *SimpleStruct) WriteStructFromTag(p *codec.Packer, tag uint32, require bool) error {
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

type RequireStruct struct {
	Ss SimpleStruct `json:"ss"`
}

func (st *RequireStruct) ResetDefault() {
	st.Ss.ResetDefault()
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
	err = st.Ss.WriteStructFromTag(p, 0, true)
	if err != nil {
		return err
	}

	_ = length
	return err
}
func (st *RequireStruct) WriteStructFromTag(p *codec.Packer, tag uint32, require bool) error {
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

type DefaultStruct struct {
	B  bool   `json:"b"`
	By byte   `json:"by"`
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
	err = up.ReadByte(&st.By, 1, false)
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
	if false || st.B != true {
		err = p.WriteBool(0, st.B)
		if err != nil {
			return err
		}
	}
	if false || st.By != 1 {
		err = p.WriteByte(1, st.By)
		if err != nil {
			return err
		}
	}
	if false || st.S != 10 {
		err = p.WriteInt16(2, st.S)
		if err != nil {
			return err
		}
	}
	if false || st.I != 1 {
		err = p.WriteInt32(3, st.I)
		if err != nil {
			return err
		}
	}
	if false || st.L != 0x0FFFFFFFFFFFFFFF {
		err = p.WriteInt64(4, st.L)
		if err != nil {
			return err
		}
	}
	if false || st.Ss != "yellia" {
		err = p.WriteString(5, st.Ss)
		if err != nil {
			return err
		}
	}

	_ = length
	return err
}
func (st *DefaultStruct) WriteStructFromTag(p *codec.Packer, tag uint32, require bool) error {
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
