// 此文件为sdp2go工具自动生成,请不要手动编辑

package Test2

import (
	"bytes"
	"fmt"
	"github.com/yellia1989/tex-go/tools/sdp/codec"
	"github.com/yellia1989/tex-go/tools/sdp/util"
	"strconv"
)

type NUMBER int32

const (
	NUMBER_1 = 1
	NUMBER_2 = 2
)

func (en NUMBER) String() string {
	ret := ""
	switch en {
	case NUMBER_1:
		ret = "NUMBER_1"
	case NUMBER_2:
		ret = "NUMBER_2"
	}
	return ret
}

type Student struct {
	IUid    uint64            `json:"iUid"`
	SName   string            `json:"sName"`
	IAge    uint32            `json:"iAge"`
	MSecret map[string]string `json:"mSecret"`
}

func (st *Student) resetDefault() {
	st.IUid = 1
}
func (st *Student) Copy() *Student {
	ret := NewStudent()
	ret.IUid = st.IUid
	ret.SName = st.SName
	ret.IAge = st.IAge
	ret.MSecret = make(map[string]string)
	for k, v := range st.MSecret {
		ret.MSecret[k] = v
	}
	return ret
}
func NewStudent() *Student {
	ret := &Student{}
	ret.resetDefault()
	return ret
}
func (st *Student) Visit(buff *bytes.Buffer, t int) {
	util.Tab(buff, t+1, util.Fieldname("iUid")+fmt.Sprintf("%v\n", st.IUid))
	util.Tab(buff, t+1, util.Fieldname("sName")+fmt.Sprintf("%v\n", st.SName))
	util.Tab(buff, t+1, util.Fieldname("iAge")+fmt.Sprintf("%v\n", st.IAge))
	util.Tab(buff, t+1, util.Fieldname("mSecret")+strconv.Itoa(len(st.MSecret)))
	if len(st.MSecret) == 0 {
		buff.WriteString(", {}\n")
	} else {
		buff.WriteString(", {\n")
	}
	for k, v := range st.MSecret {
		util.Tab(buff, t+1+1, "(\n")
		util.Tab(buff, t+1+2, util.Fieldname("")+fmt.Sprintf("%v\n", k))
		util.Tab(buff, t+1+2, util.Fieldname("")+fmt.Sprintf("%v\n", v))
		util.Tab(buff, t+1+1, ")\n")
	}
	if len(st.MSecret) != 0 {
		util.Tab(buff, t+1, "}\n")
	}
}
func (st *Student) ReadStruct(up *codec.UnPacker) error {
	var err error
	var length uint32
	var has bool
	var ty uint32
	st.resetDefault()
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
	if err != nil {
		return err
	}
	if has {
		if ty != codec.SdpType_Map {
			return fmt.Errorf("tag:%d got wrong type %d", 3, ty)
		}

		_, length, err = up.ReadNumber32()
		if err != nil {
			return err
		}
		st.MSecret = make(map[string]string)
		for i := uint32(0); i < length; i++ {
			var stMSecretk string
			err = up.ReadString(&stMSecretk, 0, true)
			if err != nil {
				return err
			}
			var stMSecretv string
			err = up.ReadString(&stMSecretv, 0, true)
			if err != nil {
				return err
			}
			st.MSecret[stMSecretk] = stMSecretv
		}
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
	var length uint32
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

	length = uint32(len(st.MSecret))
	if false || length != 0 {
		err = p.WriteHeader(3, codec.SdpType_Map)
		if err != nil {
			return err
		}
		err = p.WriteNumber32(length)
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

func (st *Teacher) resetDefault() {
	st.S1.resetDefault()
	st.S2.resetDefault()
}
func (st *Teacher) Copy() *Teacher {
	ret := NewTeacher()
	ret.IId = st.IId
	ret.SName = st.SName
	ret.S1 = *(st.S1.Copy())
	ret.S2 = *(st.S2.Copy())
	return ret
}
func NewTeacher() *Teacher {
	ret := &Teacher{}
	ret.resetDefault()
	return ret
}
func (st *Teacher) Visit(buff *bytes.Buffer, t int) {
	util.Tab(buff, t+1, util.Fieldname("iId")+fmt.Sprintf("%v\n", st.IId))
	util.Tab(buff, t+1, util.Fieldname("sName")+fmt.Sprintf("%v\n", st.SName))
	util.Tab(buff, t+1, util.Fieldname("s1")+"{\n")
	st.S1.Visit(buff, t+1+1)
	util.Tab(buff, t+1, "}\n")
	util.Tab(buff, t+1, util.Fieldname("s2")+"{\n")
	st.S2.Visit(buff, t+1+1)
	util.Tab(buff, t+1, "}\n")
}
func (st *Teacher) ReadStruct(up *codec.UnPacker) error {
	var err error
	var length uint32
	var has bool
	var ty uint32
	st.resetDefault()
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
	var length uint32
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

func (st *Teachers) resetDefault() {
}
func (st *Teachers) Copy() *Teachers {
	ret := NewTeachers()
	ret.VTeacher = make([]Teacher, len(st.VTeacher))
	for i, v := range st.VTeacher {
		ret.VTeacher[i] = *(v.Copy())
	}
	return ret
}
func NewTeachers() *Teachers {
	ret := &Teachers{}
	ret.resetDefault()
	return ret
}
func (st *Teachers) Visit(buff *bytes.Buffer, t int) {
	util.Tab(buff, t+1, util.Fieldname("vTeacher")+strconv.Itoa(len(st.VTeacher)))
	if len(st.VTeacher) == 0 {
		buff.WriteString(", []\n")
	} else {
		buff.WriteString(", [\n")
	}
	for _, v := range st.VTeacher {
		util.Tab(buff, t+1+1, util.Fieldname("")+"{\n")
		v.Visit(buff, t+1+1+1)
		util.Tab(buff, t+1+1, "}\n")
	}
	if len(st.VTeacher) != 0 {
		util.Tab(buff, t+1, "]\n")
	}
}
func (st *Teachers) ReadStruct(up *codec.UnPacker) error {
	var err error
	var length uint32
	var has bool
	var ty uint32
	st.resetDefault()

	has, ty, err = up.SkipToTag(0, false)
	if err != nil {
		return err
	}
	if has {
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
	var length uint32

	length = uint32(len(st.VTeacher))
	if false || length != 0 {
		err = p.WriteHeader(0, codec.SdpType_Vector)
		if err != nil {
			return err
		}
		err = p.WriteNumber32(length)
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

func (st *Class) resetDefault() {
}
func (st *Class) Copy() *Class {
	ret := NewClass()
	ret.IId = st.IId
	ret.SName = st.SName
	ret.VStudent = make([]Student, len(st.VStudent))
	for i, v := range st.VStudent {
		ret.VStudent[i] = *(v.Copy())
	}
	ret.VData = make([]byte, len(st.VData))
	for i, v := range st.VData {
		ret.VData[i] = v
	}
	ret.VTeacher = make([]Teacher, len(st.VTeacher))
	for i, v := range st.VTeacher {
		ret.VTeacher[i] = *(v.Copy())
	}
	return ret
}
func NewClass() *Class {
	ret := &Class{}
	ret.resetDefault()
	return ret
}
func (st *Class) Visit(buff *bytes.Buffer, t int) {
	util.Tab(buff, t+1, util.Fieldname("iId")+fmt.Sprintf("%v\n", st.IId))
	util.Tab(buff, t+1, util.Fieldname("sName")+fmt.Sprintf("%v\n", st.SName))
	util.Tab(buff, t+1, util.Fieldname("vStudent")+strconv.Itoa(len(st.VStudent)))
	if len(st.VStudent) == 0 {
		buff.WriteString(", []\n")
	} else {
		buff.WriteString(", [\n")
	}
	for _, v := range st.VStudent {
		util.Tab(buff, t+1+1, util.Fieldname("")+"{\n")
		v.Visit(buff, t+1+1+1)
		util.Tab(buff, t+1+1, "}\n")
	}
	if len(st.VStudent) != 0 {
		util.Tab(buff, t+1, "]\n")
	}
	util.Tab(buff, t+1, util.Fieldname("vData")+strconv.Itoa(len(st.VData)))
	if len(st.VData) == 0 {
		buff.WriteString(", []\n")
	} else {
		buff.WriteString(", [\n")
	}
	for _, v := range st.VData {
		util.Tab(buff, t+1+1, util.Fieldname("")+fmt.Sprintf("%v\n", v))
	}
	if len(st.VData) != 0 {
		util.Tab(buff, t+1, "]\n")
	}
	util.Tab(buff, t+1, util.Fieldname("vTeacher")+strconv.Itoa(len(st.VTeacher)))
	if len(st.VTeacher) == 0 {
		buff.WriteString(", []\n")
	} else {
		buff.WriteString(", [\n")
	}
	for _, v := range st.VTeacher {
		util.Tab(buff, t+1+1, util.Fieldname("")+"{\n")
		v.Visit(buff, t+1+1+1)
		util.Tab(buff, t+1+1, "}\n")
	}
	if len(st.VTeacher) != 0 {
		util.Tab(buff, t+1, "]\n")
	}
}
func (st *Class) ReadStruct(up *codec.UnPacker) error {
	var err error
	var length uint32
	var has bool
	var ty uint32
	st.resetDefault()
	err = up.ReadUint32(&st.IId, 0, true)
	if err != nil {
		return err
	}
	err = up.ReadString(&st.SName, 1, false)
	if err != nil {
		return err
	}

	has, ty, err = up.SkipToTag(2, false)
	if err != nil {
		return err
	}
	if has {
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
	}
	var sstVData string
	err = up.ReadString(&sstVData, 3, false)
	if err != nil {
		return err
	}
	st.VData = []byte(sstVData)

	has, ty, err = up.SkipToTag(4, true)
	if err != nil {
		return err
	}
	if has {
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
	var length uint32
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

	length = uint32(len(st.VStudent))
	if false || length != 0 {
		err = p.WriteHeader(2, codec.SdpType_Vector)
		if err != nil {
			return err
		}
		err = p.WriteNumber32(length)
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
	length = uint32(len(st.VData))
	if false || length != 0 {
		stmp := string(st.VData)
		err = p.WriteString(3, stmp)
		if err != nil {
			return err
		}
	}

	length = uint32(len(st.VTeacher))
	if true || length != 0 {
		err = p.WriteHeader(4, codec.SdpType_Vector)
		if err != nil {
			return err
		}
		err = p.WriteNumber32(length)
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

type School struct {
	MClass map[uint32]Class `json:"mClass"`
}

func (st *School) resetDefault() {
}
func (st *School) Copy() *School {
	ret := NewSchool()
	ret.MClass = make(map[uint32]Class)
	for k, v := range st.MClass {
		ret.MClass[k] = *(v.Copy())
	}
	return ret
}
func NewSchool() *School {
	ret := &School{}
	ret.resetDefault()
	return ret
}
func (st *School) Visit(buff *bytes.Buffer, t int) {
	util.Tab(buff, t+1, util.Fieldname("mClass")+strconv.Itoa(len(st.MClass)))
	if len(st.MClass) == 0 {
		buff.WriteString(", {}\n")
	} else {
		buff.WriteString(", {\n")
	}
	for k, v := range st.MClass {
		util.Tab(buff, t+1+1, "(\n")
		util.Tab(buff, t+1+2, util.Fieldname("")+fmt.Sprintf("%v\n", k))
		util.Tab(buff, t+1+2, util.Fieldname("")+"{\n")
		v.Visit(buff, t+1+2+1)
		util.Tab(buff, t+1+2, "}\n")
		util.Tab(buff, t+1+1, ")\n")
	}
	if len(st.MClass) != 0 {
		util.Tab(buff, t+1, "}\n")
	}
}
func (st *School) ReadStruct(up *codec.UnPacker) error {
	var err error
	var length uint32
	var has bool
	var ty uint32
	st.resetDefault()

	has, ty, err = up.SkipToTag(0, false)
	if err != nil {
		return err
	}
	if has {
		if ty != codec.SdpType_Map {
			return fmt.Errorf("tag:%d got wrong type %d", 0, ty)
		}

		_, length, err = up.ReadNumber32()
		if err != nil {
			return err
		}
		st.MClass = make(map[uint32]Class)
		for i := uint32(0); i < length; i++ {
			var stMClassk uint32
			err = up.ReadUint32(&stMClassk, 0, true)
			if err != nil {
				return err
			}
			var stMClassv Class
			err = stMClassv.ReadStructFromTag(up, 0, true)
			if err != nil {
				return err
			}
			st.MClass[stMClassk] = stMClassv
		}
	}

	_ = length
	_ = has
	_ = ty

	return err
}
func (st *School) ReadStructFromTag(up *codec.UnPacker, tag uint32, require bool) error {
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
func (st *School) WriteStruct(p *codec.Packer) error {
	var err error
	var length uint32

	length = uint32(len(st.MClass))
	if false || length != 0 {
		err = p.WriteHeader(0, codec.SdpType_Map)
		if err != nil {
			return err
		}
		err = p.WriteNumber32(length)
		if err != nil {
			return err
		}
		for _k, _v := range st.MClass {
			if true || _k != 0 {
				err = p.WriteUint32(0, _k)
				if err != nil {
					return err
				}
			}
			err = _v.WriteStructFromTag(p, 0, true)
			if err != nil {
				return err
			}
		}
	}

	_ = length
	return err
}
func (st *School) WriteStructFromTag(p *codec.Packer, tag uint32, require bool) error {
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
