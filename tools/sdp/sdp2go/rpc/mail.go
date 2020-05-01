// 此文件为sdp2go工具自动生成,请不要手动编辑

package rpc

import (
	"bytes"
	"fmt"
	"github.com/yellia1989/tex-go/tools/sdp/codec"
	"github.com/yellia1989/tex-go/tools/sdp/util"
	"strconv"
)

type CmdIDNum struct {
	IId  uint32 `json:"iId"`
	INum uint32 `json:"iNum"`
}

func (st *CmdIDNum) resetDefault() {
}
func (st *CmdIDNum) Copy() *CmdIDNum {
	ret := NewCmdIDNum()
	ret.IId = st.IId
	ret.INum = st.INum
	return ret
}
func NewCmdIDNum() *CmdIDNum {
	ret := &CmdIDNum{}
	ret.resetDefault()
	return ret
}
func (st *CmdIDNum) Visit(buff *bytes.Buffer, t int) {
	util.Tab(buff, t+1, util.Fieldname("iId")+fmt.Sprintf("%v\n", st.IId))
	util.Tab(buff, t+1, util.Fieldname("iNum")+fmt.Sprintf("%v\n", st.INum))
}
func (st *CmdIDNum) ReadStruct(up *codec.UnPacker) error {
	var err error
	var length uint32
	var has bool
	var ty uint32
	st.resetDefault()
	err = up.ReadUint32(&st.IId, 0, false)
	if err != nil {
		return err
	}
	err = up.ReadUint32(&st.INum, 1, false)
	if err != nil {
		return err
	}

	_ = length
	_ = has
	_ = ty

	return err
}
func (st *CmdIDNum) ReadStructFromTag(up *codec.UnPacker, tag uint32, require bool) error {
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
func (st *CmdIDNum) WriteStruct(p *codec.Packer) error {
	var err error
	var length uint32
	if false || st.IId != 0 {
		err = p.WriteUint32(0, st.IId)
		if err != nil {
			return err
		}
	}
	if false || st.INum != 0 {
		err = p.WriteUint32(1, st.INum)
		if err != nil {
			return err
		}
	}

	_ = length
	return err
}
func (st *CmdIDNum) WriteStructFromTag(p *codec.Packer, tag uint32, require bool) error {
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

type MailDataInfo struct {
	IMailId                uint32     `json:"iMailId"`
	SFrom                  string     `json:"sFrom"`
	VToUser                []uint64   `json:"vToUser"`
	STime                  string     `json:"sTime"`
	STitle                 string     `json:"sTitle"`
	SContent               string     `json:"sContent"`
	IDiamond               uint32     `json:"iDiamond"`
	ICoin                  uint32     `json:"iCoin"`
	VItems                 []CmdIDNum `json:"vItems"`
	VSendZoneIds           []uint32   `json:"vSendZoneIds"`
	IFlag                  uint32     `json:"iFlag"`
	VRcvZoneIds            []uint32   `json:"vRcvZoneIds"`
	IArenaCoin             uint32     `json:"iArenaCoin"`
	IDelTimeAfterOpen      uint32     `json:"iDelTimeAfterOpen"`
	SUserFileName          string     `json:"sUserFileName"`
	IKingCoin              uint32     `json:"iKingCoin"`
	VCustomItem            []string   `json:"vCustomItem"`
	IDelTimeAfterRcvAttach uint32     `json:"iDelTimeAfterRcvAttach"`
}

func (st *MailDataInfo) resetDefault() {
}
func (st *MailDataInfo) Copy() *MailDataInfo {
	ret := NewMailDataInfo()
	ret.IMailId = st.IMailId
	ret.SFrom = st.SFrom
	ret.VToUser = make([]uint64, len(st.VToUser))
	for i, v := range st.VToUser {
		ret.VToUser[i] = v
	}
	ret.STime = st.STime
	ret.STitle = st.STitle
	ret.SContent = st.SContent
	ret.IDiamond = st.IDiamond
	ret.ICoin = st.ICoin
	ret.VItems = make([]CmdIDNum, len(st.VItems))
	for i, v := range st.VItems {
		ret.VItems[i] = *(v.Copy())
	}
	ret.VSendZoneIds = make([]uint32, len(st.VSendZoneIds))
	for i, v := range st.VSendZoneIds {
		ret.VSendZoneIds[i] = v
	}
	ret.IFlag = st.IFlag
	ret.VRcvZoneIds = make([]uint32, len(st.VRcvZoneIds))
	for i, v := range st.VRcvZoneIds {
		ret.VRcvZoneIds[i] = v
	}
	ret.IArenaCoin = st.IArenaCoin
	ret.IDelTimeAfterOpen = st.IDelTimeAfterOpen
	ret.SUserFileName = st.SUserFileName
	ret.IKingCoin = st.IKingCoin
	ret.VCustomItem = make([]string, len(st.VCustomItem))
	for i, v := range st.VCustomItem {
		ret.VCustomItem[i] = v
	}
	ret.IDelTimeAfterRcvAttach = st.IDelTimeAfterRcvAttach
	return ret
}
func NewMailDataInfo() *MailDataInfo {
	ret := &MailDataInfo{}
	ret.resetDefault()
	return ret
}
func (st *MailDataInfo) Visit(buff *bytes.Buffer, t int) {
	util.Tab(buff, t+1, util.Fieldname("iMailId")+fmt.Sprintf("%v\n", st.IMailId))
	util.Tab(buff, t+1, util.Fieldname("sFrom")+fmt.Sprintf("%v\n", st.SFrom))
	util.Tab(buff, t+1, util.Fieldname("vToUser")+strconv.Itoa(len(st.VToUser)))
	if len(st.VToUser) == 0 {
		buff.WriteString(", []\n")
	} else {
		buff.WriteString(", [\n")
	}
	for _, v := range st.VToUser {
		util.Tab(buff, t+1+1, util.Fieldname("")+fmt.Sprintf("%v\n", v))
	}
	if len(st.VToUser) != 0 {
		util.Tab(buff, t+1, "]\n")
	}
	util.Tab(buff, t+1, util.Fieldname("sTime")+fmt.Sprintf("%v\n", st.STime))
	util.Tab(buff, t+1, util.Fieldname("sTitle")+fmt.Sprintf("%v\n", st.STitle))
	util.Tab(buff, t+1, util.Fieldname("sContent")+fmt.Sprintf("%v\n", st.SContent))
	util.Tab(buff, t+1, util.Fieldname("iDiamond")+fmt.Sprintf("%v\n", st.IDiamond))
	util.Tab(buff, t+1, util.Fieldname("iCoin")+fmt.Sprintf("%v\n", st.ICoin))
	util.Tab(buff, t+1, util.Fieldname("vItems")+strconv.Itoa(len(st.VItems)))
	if len(st.VItems) == 0 {
		buff.WriteString(", []\n")
	} else {
		buff.WriteString(", [\n")
	}
	for _, v := range st.VItems {
		util.Tab(buff, t+1+1, util.Fieldname("")+"{\n")
		v.Visit(buff, t+1+1+1)
		util.Tab(buff, t+1+1, "}\n")
	}
	if len(st.VItems) != 0 {
		util.Tab(buff, t+1, "]\n")
	}
	util.Tab(buff, t+1, util.Fieldname("vSendZoneIds")+strconv.Itoa(len(st.VSendZoneIds)))
	if len(st.VSendZoneIds) == 0 {
		buff.WriteString(", []\n")
	} else {
		buff.WriteString(", [\n")
	}
	for _, v := range st.VSendZoneIds {
		util.Tab(buff, t+1+1, util.Fieldname("")+fmt.Sprintf("%v\n", v))
	}
	if len(st.VSendZoneIds) != 0 {
		util.Tab(buff, t+1, "]\n")
	}
	util.Tab(buff, t+1, util.Fieldname("iFlag")+fmt.Sprintf("%v\n", st.IFlag))
	util.Tab(buff, t+1, util.Fieldname("vRcvZoneIds")+strconv.Itoa(len(st.VRcvZoneIds)))
	if len(st.VRcvZoneIds) == 0 {
		buff.WriteString(", []\n")
	} else {
		buff.WriteString(", [\n")
	}
	for _, v := range st.VRcvZoneIds {
		util.Tab(buff, t+1+1, util.Fieldname("")+fmt.Sprintf("%v\n", v))
	}
	if len(st.VRcvZoneIds) != 0 {
		util.Tab(buff, t+1, "]\n")
	}
	util.Tab(buff, t+1, util.Fieldname("iArenaCoin")+fmt.Sprintf("%v\n", st.IArenaCoin))
	util.Tab(buff, t+1, util.Fieldname("iDelTimeAfterOpen")+fmt.Sprintf("%v\n", st.IDelTimeAfterOpen))
	util.Tab(buff, t+1, util.Fieldname("sUserFileName")+fmt.Sprintf("%v\n", st.SUserFileName))
	util.Tab(buff, t+1, util.Fieldname("iKingCoin")+fmt.Sprintf("%v\n", st.IKingCoin))
	util.Tab(buff, t+1, util.Fieldname("vCustomItem")+strconv.Itoa(len(st.VCustomItem)))
	if len(st.VCustomItem) == 0 {
		buff.WriteString(", []\n")
	} else {
		buff.WriteString(", [\n")
	}
	for _, v := range st.VCustomItem {
		util.Tab(buff, t+1+1, util.Fieldname("")+fmt.Sprintf("%v\n", v))
	}
	if len(st.VCustomItem) != 0 {
		util.Tab(buff, t+1, "]\n")
	}
	util.Tab(buff, t+1, util.Fieldname("iDelTimeAfterRcvAttach")+fmt.Sprintf("%v\n", st.IDelTimeAfterRcvAttach))
}
func (st *MailDataInfo) ReadStruct(up *codec.UnPacker) error {
	var err error
	var length uint32
	var has bool
	var ty uint32
	st.resetDefault()
	err = up.ReadUint32(&st.IMailId, 0, false)
	if err != nil {
		return err
	}
	err = up.ReadString(&st.SFrom, 1, false)
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
		st.VToUser = make([]uint64, length, length)
		for i := uint32(0); i < length; i++ {
			err = up.ReadUint64(&st.VToUser[i], 0, true)
			if err != nil {
				return err
			}
		}
	}
	err = up.ReadString(&st.STime, 3, false)
	if err != nil {
		return err
	}
	err = up.ReadString(&st.STitle, 4, false)
	if err != nil {
		return err
	}
	err = up.ReadString(&st.SContent, 5, false)
	if err != nil {
		return err
	}
	err = up.ReadUint32(&st.IDiamond, 7, false)
	if err != nil {
		return err
	}
	err = up.ReadUint32(&st.ICoin, 8, false)
	if err != nil {
		return err
	}

	has, ty, err = up.SkipToTag(9, false)
	if err != nil {
		return err
	}
	if has {
		if ty != codec.SdpType_Vector {
			return fmt.Errorf("tag:%d got wrong type %d", 9, ty)
		}

		_, length, err = up.ReadNumber32()
		if err != nil {
			return err
		}
		st.VItems = make([]CmdIDNum, length, length)
		for i := uint32(0); i < length; i++ {
			err = st.VItems[i].ReadStructFromTag(up, 0, true)
			if err != nil {
				return err
			}
		}
	}

	has, ty, err = up.SkipToTag(10, false)
	if err != nil {
		return err
	}
	if has {
		if ty != codec.SdpType_Vector {
			return fmt.Errorf("tag:%d got wrong type %d", 10, ty)
		}

		_, length, err = up.ReadNumber32()
		if err != nil {
			return err
		}
		st.VSendZoneIds = make([]uint32, length, length)
		for i := uint32(0); i < length; i++ {
			err = up.ReadUint32(&st.VSendZoneIds[i], 0, true)
			if err != nil {
				return err
			}
		}
	}
	err = up.ReadUint32(&st.IFlag, 11, false)
	if err != nil {
		return err
	}

	has, ty, err = up.SkipToTag(12, false)
	if err != nil {
		return err
	}
	if has {
		if ty != codec.SdpType_Vector {
			return fmt.Errorf("tag:%d got wrong type %d", 12, ty)
		}

		_, length, err = up.ReadNumber32()
		if err != nil {
			return err
		}
		st.VRcvZoneIds = make([]uint32, length, length)
		for i := uint32(0); i < length; i++ {
			err = up.ReadUint32(&st.VRcvZoneIds[i], 0, true)
			if err != nil {
				return err
			}
		}
	}
	err = up.ReadUint32(&st.IArenaCoin, 14, false)
	if err != nil {
		return err
	}
	err = up.ReadUint32(&st.IDelTimeAfterOpen, 15, false)
	if err != nil {
		return err
	}
	err = up.ReadString(&st.SUserFileName, 16, false)
	if err != nil {
		return err
	}
	err = up.ReadUint32(&st.IKingCoin, 17, false)
	if err != nil {
		return err
	}

	has, ty, err = up.SkipToTag(19, false)
	if err != nil {
		return err
	}
	if has {
		if ty != codec.SdpType_Vector {
			return fmt.Errorf("tag:%d got wrong type %d", 19, ty)
		}

		_, length, err = up.ReadNumber32()
		if err != nil {
			return err
		}
		st.VCustomItem = make([]string, length, length)
		for i := uint32(0); i < length; i++ {
			err = up.ReadString(&st.VCustomItem[i], 0, true)
			if err != nil {
				return err
			}
		}
	}
	err = up.ReadUint32(&st.IDelTimeAfterRcvAttach, 20, false)
	if err != nil {
		return err
	}

	_ = length
	_ = has
	_ = ty

	return err
}
func (st *MailDataInfo) ReadStructFromTag(up *codec.UnPacker, tag uint32, require bool) error {
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
func (st *MailDataInfo) WriteStruct(p *codec.Packer) error {
	var err error
	var length uint32
	if false || st.IMailId != 0 {
		err = p.WriteUint32(0, st.IMailId)
		if err != nil {
			return err
		}
	}
	if false || st.SFrom != "" {
		err = p.WriteString(1, st.SFrom)
		if err != nil {
			return err
		}
	}

	length = uint32(len(st.VToUser))
	if false || length != 0 {
		err = p.WriteHeader(2, codec.SdpType_Vector)
		if err != nil {
			return err
		}
		err = p.WriteNumber32(length)
		if err != nil {
			return err
		}
		for _, v := range st.VToUser {
			if true || v != 0 {
				err = p.WriteUint64(0, v)
				if err != nil {
					return err
				}
			}
		}
	}
	if false || st.STime != "" {
		err = p.WriteString(3, st.STime)
		if err != nil {
			return err
		}
	}
	if false || st.STitle != "" {
		err = p.WriteString(4, st.STitle)
		if err != nil {
			return err
		}
	}
	if false || st.SContent != "" {
		err = p.WriteString(5, st.SContent)
		if err != nil {
			return err
		}
	}
	if false || st.IDiamond != 0 {
		err = p.WriteUint32(7, st.IDiamond)
		if err != nil {
			return err
		}
	}
	if false || st.ICoin != 0 {
		err = p.WriteUint32(8, st.ICoin)
		if err != nil {
			return err
		}
	}

	length = uint32(len(st.VItems))
	if false || length != 0 {
		err = p.WriteHeader(9, codec.SdpType_Vector)
		if err != nil {
			return err
		}
		err = p.WriteNumber32(length)
		if err != nil {
			return err
		}
		for _, v := range st.VItems {
			err = v.WriteStructFromTag(p, 0, true)
			if err != nil {
				return err
			}
		}
	}

	length = uint32(len(st.VSendZoneIds))
	if false || length != 0 {
		err = p.WriteHeader(10, codec.SdpType_Vector)
		if err != nil {
			return err
		}
		err = p.WriteNumber32(length)
		if err != nil {
			return err
		}
		for _, v := range st.VSendZoneIds {
			if true || v != 0 {
				err = p.WriteUint32(0, v)
				if err != nil {
					return err
				}
			}
		}
	}
	if false || st.IFlag != 0 {
		err = p.WriteUint32(11, st.IFlag)
		if err != nil {
			return err
		}
	}

	length = uint32(len(st.VRcvZoneIds))
	if false || length != 0 {
		err = p.WriteHeader(12, codec.SdpType_Vector)
		if err != nil {
			return err
		}
		err = p.WriteNumber32(length)
		if err != nil {
			return err
		}
		for _, v := range st.VRcvZoneIds {
			if true || v != 0 {
				err = p.WriteUint32(0, v)
				if err != nil {
					return err
				}
			}
		}
	}
	if false || st.IArenaCoin != 0 {
		err = p.WriteUint32(14, st.IArenaCoin)
		if err != nil {
			return err
		}
	}
	if false || st.IDelTimeAfterOpen != 0 {
		err = p.WriteUint32(15, st.IDelTimeAfterOpen)
		if err != nil {
			return err
		}
	}
	if false || st.SUserFileName != "" {
		err = p.WriteString(16, st.SUserFileName)
		if err != nil {
			return err
		}
	}
	if false || st.IKingCoin != 0 {
		err = p.WriteUint32(17, st.IKingCoin)
		if err != nil {
			return err
		}
	}

	length = uint32(len(st.VCustomItem))
	if false || length != 0 {
		err = p.WriteHeader(19, codec.SdpType_Vector)
		if err != nil {
			return err
		}
		err = p.WriteNumber32(length)
		if err != nil {
			return err
		}
		for _, v := range st.VCustomItem {
			if true || v != "" {
				err = p.WriteString(0, v)
				if err != nil {
					return err
				}
			}
		}
	}
	if false || st.IDelTimeAfterRcvAttach != 0 {
		err = p.WriteUint32(20, st.IDelTimeAfterRcvAttach)
		if err != nil {
			return err
		}
	}

	_ = length
	return err
}
func (st *MailDataInfo) WriteStructFromTag(p *codec.Packer, tag uint32, require bool) error {
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
