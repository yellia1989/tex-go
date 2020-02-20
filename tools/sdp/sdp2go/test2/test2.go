// 此文件为sdp2go工具自动生成,请不要手动编辑

package test2

import (
	"fmt"
	"github.com/yellia1989/tex-go/tools/sdp/codec"
)

type Age int32

const (
	Age_10 = 1
	Age_20 = 2
	Age_30 = 2
)
const (
	NAME string = "yellia"
)

type SimpleStruct struct {
	B   bool            `json:"b"`
	By  int8            `json:"by"`
	S   int16           `json:"s"`
	Us  uint16          `json:"us"`
	I   int32           `json:"i"`
	Ui  uint32          `json:"ui"`
	L   int64           `json:"l"`
	Ul  uint64          `json:"ul"`
	F   float32         `json:"f"`
	D   float64         `json:"d"`
	Ss  string          `json:"ss"`
	Vi  []int32         `json:"vi"`
	Mi  map[int32]int32 `json:"mi"`
	Age Age             `json:"age"`
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
	err = up.ReadInt8(&st.By, 1, false)
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
		err = up.ReadInt32(&v, 1, true)
		if err != nil {
			return err
		}
		st.Mi[k] = v
	}
	err = up.ReadInt32((*int32)(&st.Age), 13, false)
	if err != nil {
		return err
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
	err = p.WriteBool(0, st.B)
	if err != nil {
		return err
	}
	err = p.WriteInt8(1, st.By)
	if err != nil {
		return err
	}
	err = p.WriteInt16(2, st.S)
	if err != nil {
		return err
	}
	err = p.WriteUint16(3, st.Us)
	if err != nil {
		return err
	}
	err = p.WriteInt32(4, st.I)
	if err != nil {
		return err
	}
	err = p.WriteUint32(5, st.Ui)
	if err != nil {
		return err
	}
	err = p.WriteInt64(6, st.L)
	if err != nil {
		return err
	}
	err = p.WriteUint64(7, st.Ul)
	if err != nil {
		return err
	}
	err = p.WriteFloat32(8, st.F)
	if err != nil {
		return err
	}
	err = p.WriteFloat64(9, st.D)
	if err != nil {
		return err
	}
	err = p.WriteString(10, st.Ss)
	if err != nil {
		return err
	}

	err = p.WriteHeader(11, codec.SdpType_Vector)
	if err != nil {
		return err
	}
	length = len(st.Vi)
	err = p.WriteNumber32(uint32(length))
	if err != nil {
		return err
	}
	for _, v := range st.Vi {
		err = p.WriteInt32(0, v)
		if err != nil {
			return err
		}
	}

	err = p.WriteHeader(12, codec.SdpType_Map)
	if err != nil {
		return err
	}
	length = len(st.Mi)
	err = p.WriteNumber32(uint32(length))
	if err != nil {
		return err
	}
	for _k, _v := range st.Mi {
		err = p.WriteInt32(0, _k)
		if err != nil {
			return err
		}
		err = p.WriteInt32(1, _v)
		if err != nil {
			return err
		}
	}
	err = p.WriteInt32(13, int32(st.Age))
	if err != nil {
		return err
	}

	_ = length
	return err
}
func (st *SimpleStruct) WriteStructFromTag(p *codec.Packer, tag uint32) error {
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
