// 此文件为sdp2go工具自动生成,请不要手动编辑

package Test2

import (
	"fmt"
	"github.com/yellia1989/tex-go/tools/sdp/codec"
)

type NUMBER int32

const (
	NUMBER_1 = 1
	NUMBER_2 = 2
)

type Student struct {
	IUid    uint64            `json:"iUid"`
	SName   string            `json:"sName"`
	IAge    uint32            `json:"iAge"`
	MSecret map[string]string `json:"mSecret"`
}

func (st *Student) ResetDefault() {
	st.IUid = 1
}
func (st *Student) ReadStruct(up *codec.UnPacker) error {
	var err error
	var length uint32
	var has bool
	var ty uint32
	st.ResetDefault()
	err = up.ReadUint64(&st.IUid, 0, false)
	if err != nil {
		return err
	}
	err = up.ReadString(&st.SName, 1, false)
	if err != nil {
		return err
	}
	err = up.ReadUint32(&st.IAge, 2, false)
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
	st.MSecret = make(map[string]string)
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
		st.MSecret[k] = v
	}

	_ = length
	_ = has
	_ = ty

	return err
}
func (st *Student) ReadStructFromTag(up *codec.UnPacker, tag uint32, require bool) error {
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
func (st *Student) WriteStruct(p *codec.Packer) error {
	var err error
	var length int
	if false || st.IUid != 1 {
		err = p.WriteUint64(0, st.IUid)
		if err != nil {
			return err
		}
	}
	if false || st.SName != "" {
		err = p.WriteString(1, st.SName)
		if err != nil {
			return err
		}
	}
	if false || st.IAge != 0 {
		err = p.WriteUint32(2, st.IAge)
		if err != nil {
			return err
		}
	}

	length = len(st.MSecret)
	if false || length != 0 {
		err = p.WriteHeader(3, codec.SdpType_Map)
		if err != nil {
			return err
		}
		err = p.WriteNumber32(uint32(length))
		if err != nil {
			return err
		}
		for _k, _v := range st.MSecret {
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
func (st *Student) WriteStructFromTag(p *codec.Packer, tag uint32, require bool) error {
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

type Teacher struct {
	IId   uint32  `json:"iId"`
	SName string  `json:"sName"`
	S1    Student `json:"s1"`
	S2    Student `json:"s2"`
}

func (st *Teacher) ResetDefault() {
	st.S1.ResetDefault()
	st.S2.ResetDefault()
}
func (st *Teacher) ReadStruct(up *codec.UnPacker) error {
	var err error
	var length uint32
	var has bool
	var ty uint32
	st.ResetDefault()
	err = up.ReadUint32(&st.IId, 0, false)
	if err != nil {
		return err
	}
	err = up.ReadString(&st.SName, 1, false)
	if err != nil {
		return err
	}
	err = st.S1.ReadStructFromTag(up, 2, false)
	if err != nil {
		return err
	}
	err = st.S2.ReadStructFromTag(up, 3, true)
	if err != nil {
		return err
	}

	_ = length
	_ = has
	_ = ty

	return err
}
func (st *Teacher) ReadStructFromTag(up *codec.UnPacker, tag uint32, require bool) error {
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
func (st *Teacher) WriteStruct(p *codec.Packer) error {
	var err error
	var length int
	if false || st.IId != 0 {
		err = p.WriteUint32(0, st.IId)
		if err != nil {
			return err
		}
	}
	if false || st.SName != "" {
		err = p.WriteString(1, st.SName)
		if err != nil {
			return err
		}
	}
	err = st.S1.WriteStructFromTag(p, 2, false)
	if err != nil {
		return err
	}
	err = st.S2.WriteStructFromTag(p, 3, true)
	if err != nil {
		return err
	}

	_ = length
	return err
}
func (st *Teacher) WriteStructFromTag(p *codec.Packer, tag uint32, require bool) error {
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

type Teachers struct {
	VTeacher []Teacher `json:"vTeacher"`
}

func (st *Teachers) ResetDefault() {
}
func (st *Teachers) ReadStruct(up *codec.UnPacker) error {
	var err error
	var length uint32
	var has bool
	var ty uint32
	st.ResetDefault()

	has, ty, err = up.SkipToTag(0, false)
	if !has || err != nil {
		return err
	}
	if ty != codec.SdpType_Vector {
		return fmt.Errorf("tag:%d got wrong type %d", 0, ty)
	}

	_, length, err = up.ReadNumber32()
	if err != nil {
		return err
	}
	st.VTeacher = make([]Teacher, length, length)
	for i := uint32(0); i < length; i++ {
		err = st.VTeacher[i].ReadStructFromTag(up, 0, true)
		if err != nil {
			return err
		}
	}

	_ = length
	_ = has
	_ = ty

	return err
}
func (st *Teachers) ReadStructFromTag(up *codec.UnPacker, tag uint32, require bool) error {
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
func (st *Teachers) WriteStruct(p *codec.Packer) error {
	var err error
	var length int

	length = len(st.VTeacher)
	if false || length != 0 {
		err = p.WriteHeader(0, codec.SdpType_Vector)
		if err != nil {
			return err
		}
		err = p.WriteNumber32(uint32(length))
		if err != nil {
			return err
		}
		for _, v := range st.VTeacher {
			err = v.WriteStructFromTag(p, 0, true)
			if err != nil {
				return err
			}
		}
	}

	_ = length
	return err
}
func (st *Teachers) WriteStructFromTag(p *codec.Packer, tag uint32, require bool) error {
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

type Class struct {
	IId      uint32    `json:"iId"`
	SName    string    `json:"sName"`
	VStudent []Student `json:"vStudent"`
	VData    []byte    `json:"vData"`
	VTeacher []Teacher `json:"vTeacher"`
}

func (st *Class) ResetDefault() {
}
func (st *Class) ReadStruct(up *codec.UnPacker) error {
	var err error
	var length uint32
	var has bool
	var ty uint32
	st.ResetDefault()
	err = up.ReadUint32(&st.IId, 0, true)
	if err != nil {
		return err
	}
	err = up.ReadString(&st.SName, 1, false)
	if err != nil {
		return err
	}

	has, ty, err = up.SkipToTag(2, false)
	if !has || err != nil {
		return err
	}
	if ty != codec.SdpType_Vector {
		return fmt.Errorf("tag:%d got wrong type %d", 2, ty)
	}

	_, length, err = up.ReadNumber32()
	if err != nil {
		return err
	}
	st.VStudent = make([]Student, length, length)
	for i := uint32(0); i < length; i++ {
		err = st.VStudent[i].ReadStructFromTag(up, 0, true)
		if err != nil {
			return err
		}
	}
	var sVData string
	err = up.ReadString(&sVData, 3, false)
	if err != nil {
		return err
	}
	st.VData = []byte(sVData)

	has, ty, err = up.SkipToTag(4, true)
	if !has || err != nil {
		return err
	}
	if ty != codec.SdpType_Vector {
		return fmt.Errorf("tag:%d got wrong type %d", 4, ty)
	}

	_, length, err = up.ReadNumber32()
	if err != nil {
		return err
	}
	st.VTeacher = make([]Teacher, length, length)
	for i := uint32(0); i < length; i++ {
		err = st.VTeacher[i].ReadStructFromTag(up, 0, true)
		if err != nil {
			return err
		}
	}

	_ = length
	_ = has
	_ = ty

	return err
}
func (st *Class) ReadStructFromTag(up *codec.UnPacker, tag uint32, require bool) error {
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
func (st *Class) WriteStruct(p *codec.Packer) error {
	var err error
	var length int
	if true || st.IId != 0 {
		err = p.WriteUint32(0, st.IId)
		if err != nil {
			return err
		}
	}
	if false || st.SName != "" {
		err = p.WriteString(1, st.SName)
		if err != nil {
			return err
		}
	}

	length = len(st.VStudent)
	if false || length != 0 {
		err = p.WriteHeader(2, codec.SdpType_Vector)
		if err != nil {
			return err
		}
		err = p.WriteNumber32(uint32(length))
		if err != nil {
			return err
		}
		for _, v := range st.VStudent {
			err = v.WriteStructFromTag(p, 0, true)
			if err != nil {
				return err
			}
		}
	}
	length = len(st.VData)
	if false || length != 0 {
		stmp := string(st.VData)
		err = p.WriteString(3, stmp)
		if err != nil {
			return err
		}
	}

	length = len(st.VTeacher)
	if true || length != 0 {
		err = p.WriteHeader(4, codec.SdpType_Vector)
		if err != nil {
			return err
		}
		err = p.WriteNumber32(uint32(length))
		if err != nil {
			return err
		}
		for _, v := range st.VTeacher {
			err = v.WriteStructFromTag(p, 0, true)
			if err != nil {
				return err
			}
		}
	}

	_ = length
	return err
}
func (st *Class) WriteStructFromTag(p *codec.Packer, tag uint32, require bool) error {
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
