package main

import (
    "bytes"
    "fmt"
    "os"
    "io/ioutil"
    "go/format"
    "strconv"
    "strings"
)

func upperFirstLetter(name string) string {
    if len(name) == 0 {
        return ""
    }

    if len(name) == 1 {
        return strings.ToUpper(string(name[0]))
    }
    return strings.ToUpper(string(name[0])) + name[1:]
}

func (en *enumInfo) rename() {
    en.name = upperFirstLetter(en.name)
    for i := range en.members {
        en.members[i].name = upperFirstLetter(en.members[i].name)
    }
}

func (cst *constInfo) rename() {
    cst.name = upperFirstLetter(cst.name)
}

func (st *structInfo) rename() {
    st.name = upperFirstLetter(st.name)
    for i := range st.members {
        st.members[i].oldname = st.members[i].name
        st.members[i].name = upperFirstLetter(st.members[i].name)
    }
}

func (itf *interfaceInfo) rename() {
    itf.name = upperFirstLetter(itf.name)
    for i := range itf.funcs {
        itf.funcs[i].name = upperFirstLetter(itf.funcs[i].name)
    }
}

type sdp2Go struct {
    code bytes.Buffer // 用来生成go文件的buf
    dir string // 保存文件的目录
    p *parse // 解析器
}

func (s2g *sdp2Go) Write(code string) {
    s2g.code.WriteString(code+"\n")
}

func (s2g *sdp2Go) generate() {
    for _, v := range s2g.p.includes {
        _s2g := &sdp2Go{dir: s2g.dir, p: v}
        _s2g.generate()
    }

    // 生成include的go文件
    hasInterface := len(s2g.p.interfaces) != 0
    s2g.Write("// 此文件为sdp2go工具自动生成,请不要手动编辑\n")
    s2g.Write("package " + s2g.p.module)
    s2g.Write("import (")
    if hasInterface {
        s2g.Write(`"context"`)
        s2g.Write(`"time"`)
    }
    s2g.Write(`"fmt"`)
    s2g.Write(`"github.com/yellia1989/tex-go/tools/sdp/codec"`)
    if hasInterface {
        s2g.Write(`"github.com/yellia1989/tex-go/service/protocol/protocol"`)
        s2g.Write(`"github.com/yellia1989/tex-go/tools/net"`)
        s2g.Write(`"github.com/yellia1989/tex-go/tools/log"`)
        s2g.Write(`tex "github.com/yellia1989/tex-go/service"`)
    }
    for m, _ := range s2g.p.dependModule {
        s2g.Write(`"`+ m +`"`)
    }
    s2g.Write(")")

    // 生成枚举
    for _, v := range s2g.p.enums {
        s2g.genEnum(&v)
    }

    // 生成常量
    s2g.genConst(s2g.p.consts)

    // 生成结构体
    for _, v := range s2g.p.structs {
        s2g.genStruct(&v)
    }

    // 生成接口
    for _, v := range s2g.p.interfaces {
        s2g.genInterface(&v)
    }

    s2g.saveToFile()
}

func (s2g *sdp2Go) saveToFile() {
    // 格式化生成的代码
    beauty, err := format.Source(s2g.code.Bytes())
    if err != nil {
        s2g.err("format source err:" + err.Error())
    }
    //beauty := s2g.code.Bytes()
    //var err error

    err = os.MkdirAll(s2g.dir+s2g.p.module, 0766)
    if err != nil {
        s2g.err("create dir err:" + err.Error())
    }

    file := s2g.p.file
    p := strings.LastIndex(file, ".")
    if p == -1 {
        file += ".go"
    } else {
        file = file[0:p] + ".go"
    }
    err = ioutil.WriteFile(s2g.dir+s2g.p.module+"/"+file, beauty, 0666)
    if err != nil {
        s2g.err("create file err:" + err.Error())
    }
}

func (s2g *sdp2Go) err(err string) {
    panic(err)
}

func (s2g *sdp2Go) genEnum(en *enumInfo) {
    en.rename()
    s2g.Write("type " + en.name + " int32")
    s2g.Write("const (")
    for _, v := range en.members {
        s2g.Write(v.name + " = " + strconv.Itoa(int(v.value)))
    }
    s2g.Write(")")
}

func (s2g *sdp2Go) genConst(csts []constInfo) {
    if len(csts) == 0 {
        return
    }

    s2g.Write("const (")
    for _, v := range csts {
        v.rename()
        s2g.Write(v.name + " " + s2g.genType(v.ty) + " = " + v.value)
    }
    s2g.Write(")")
}

func (s2g *sdp2Go) genType(ty *varType) string {
    ret := ""
    switch ty.ty {
    case tkTBool:
        ret = "bool"
    case tkTInt:
        if ty.unsigned {
            ret = "uint32"
        } else {
            ret = "int32"
        }
    case tkTShort:
        if ty.unsigned {
            ret = "uint16"
        } else {
            ret = "int16"
        }
    case tkTByte:
        if ty.unsigned {
            ret = "uint8"
        } else {
            ret = "int8"
        }
    case tkTLong:
        if ty.unsigned {
            ret = "uint64"
        } else {
            ret = "int64"
        }
    case tkTFloat:
        ret = "float32"
    case tkTDouble:
        ret = "float64"
    case tkTString:
        ret = "string"
    case tkTVector:
        ret = "[]" + s2g.genType(ty.typeK)
    case tkTMap:
        ret = "map[" + s2g.genType(ty.typeK) + "]" + s2g.genType(ty.typeV)
    case tkName:
        names := strings.Split(ty.customTyName, "::")
        for i := range names {
            if i == len(names)-1 {
                names[i] = upperFirstLetter(names[i])
            }
        }
        ret = strings.Join(names, ".")
    default:
        s2g.err("unsupported type " + TokenMap[ty.ty])
    }

    return ret
}

func (s2g *sdp2Go) genStruct(st *structInfo) {
    st.rename()

    // 结构体定义
    s2g.Write("type " + st.name + " struct {")
    for _, v := range st.members {
        s2g.Write(v.name + " " + s2g.genType(v.ty) + " `json:\"" + v.oldname + "\"`")
    }
    s2g.Write("}")

    // 默认值
    s2g.Write("func (st *" + st.name + ") ResetDefault(){")
    for _, v := range st.members {
        if v.defVal == "" {
            continue
        }
        s2g.Write("st." + v.name + " = " + v.defVal)
    }
    s2g.Write("}")

    // 读函数
    s2g.Write(`func (st *` + st.name + `) ReadStruct(up *codec.UnPacker) error {
        var err error
        var length uint32
        var has bool
        var ty uint32
        st.ResetDefault()`)
    for _, v := range st.members {
        s2g.genReadVar(&v, "st.", false, true)
    }
    s2g.Write(`
        _ = length
        _ = has
        _ = ty

        return err
    }`)
    s2g.Write(`func (st *` + st.name + `) ReadStructFromTag(up *codec.UnPacker, tag uint32, require bool) error {
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
    _= ty
    return nil
    }`)

    // 写函数
    s2g.Write(`func (st *` + st.name + `) WriteStruct(p *codec.Packer) error {
        var err error
        var length int`)
    for _, v := range st.members {
        s2g.genWriteVar(&v, "st.", false, true)
    }
    s2g.Write(`
    _ = length 
    return err
    }`)
    s2g.Write(`func (st *` + st.name + `) WriteStructFromTag(p *codec.Packer, tag uint32) error {
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
    }`)
}

func (s2g *sdp2Go) genInterface(v *interfaceInfo) {
    v.rename()

    s2g.Write(`
type ` + v.name + ` struct {
    proxy tex.ServicePrxImpl
}
func (s *` + v.name + `) SetPrxImpl(impl tex.ServicePrxImpl) {
    s.proxy = impl
}
func (s *` + v.name + `) SetTimeout(timeout time.Duration) {
    s.proxy.SetTimeout(timeout)
}`)

    for _, f := range v.funcs {
        s2g.genInterfaceProxyFunc(v, &f)
    }

    s2g.Write("type _" + v.name + "Impl interface {")
    for _, f := range v.funcs {
        buff := bytes.Buffer{}
        buff.WriteString(f.name + "(")
        // 参数
        if len(f.args) != 0 {
            buff.WriteString("ctx context.Context")
        }
        for _, arg := range f.args {
            buff.WriteString(", " + arg.name + " ")
            out := ""
            if arg.out {
                out = "*" 
            }
            buff.WriteString(out + s2g.genType(arg.ty))
        }
        buff.WriteString(") ")
        // 返回值
        ret := "("
        if f.hasRet {
            ret += s2g.genType(f.retTy)
            ret += ","
        }
        buff.WriteString(ret + "error)")
        s2g.Write(buff.String())
    }
    s2g.Write(`}`)

    // Dispatch
    s2g.Write(`
func (s *` + v.name + `) Dispatch(ctx context.Context, serviceImpl interface{}, req *protocol.RequestPacket) {
current := net.ContextGetCurrent(ctx)

log.FDebugf("handle tex request, peer: %s:%d, obj: %s, func: %s, reqid: %d", current.IP, current.    Port, req.SServiceName, req.SFuncName, req.IRequestId)

texret := protocol.SDPSERVERUNKNOWNERR
up := codec.NewUnPacker([]byte(req.SReqPayload))
p := codec.NewPacker()

var err error
switch req.SFuncName {`)
    for _, f := range v.funcs {
        s2g.genInterfaceFunc(v, &f)
    }
    
    s2g.Write(`
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
`)
}

func (s2g *sdp2Go) genInterfaceFunc(itf *interfaceInfo, f *funcInfo) {
    s2g.Write(`case "` + f.name + `":`)
    s2g.Write("impl := serviceImpl.(_" + itf.name + "Impl)")

    // 读参数
    for i, arg := range f.args {
        if arg.out {
            continue
        }
        name := "p"+ strconv.Itoa(i+1)
        s2g.Write("var " + name + " " + s2g.genType(arg.ty))
        dummy := &structMember{}
        dummy.ty = arg.ty
        dummy.name = name
        dummy.tag = uint32(i+1)
        dummy.require = true
        s2g.genReadVar(dummy, "", false, false)
    }
 
    // 调用实现
    for i, arg := range f.args {
        if !arg.out {
            continue
        }
        name := "p"+ strconv.Itoa(i+1)
        s2g.Write("var " + name + " " + s2g.genType(arg.ty))
    }
    buff := bytes.Buffer{}
    if f.hasRet {
        buff.WriteString("var ret " + s2g.genType(f.retTy) + "\n")
        buff.WriteString("ret, ")
    }
    buff.WriteString("err = impl." + f.name + "(ctx")
    for i, arg := range f.args {
        name := "p"+ strconv.Itoa(i+1)
        buff.WriteString(", ")
        if !arg.out {
            buff.WriteString(name)
        } else {
            buff.WriteString("&" + name)
        }
    }
    s2g.Write(buff.String() + ")")
    s2g.Write(genCheckErr(false, false))

    // 返回值
    if f.hasRet {
        dummy := &structMember{}
        dummy.name = "ret"
        dummy.ty = f.retTy
        dummy.tag = 0
        s2g.genWriteVar(dummy, "", false, false)
    }

    // 返回参数
    for i, arg := range f.args {
        name := "p"+ strconv.Itoa(i+1)
        if !arg.out {
            continue
        } 
        dummy := &structMember{}
        dummy.name = name
        dummy.ty = arg.ty
        dummy.tag = uint32(i+1)
        s2g.genWriteVar(dummy, "", false, false)
    }

    s2g.Write("texret = protocol.SDPSERVERSUCCESS")
}

func (s2g *sdp2Go) genInterfaceProxyFunc(itf *interfaceInfo, f *funcInfo) {
    buff := bytes.Buffer{}
    buff.WriteString("func (s *" + itf.name + ") " + f.name + "(")
    // 参数
    for i, arg := range f.args {
        if i != 0 {
            buff.WriteString(", ")
        }
        buff.WriteString(arg.name + " ")
        out := ""
        if arg.out {
            out = "*" 
        }
        buff.WriteString(out + s2g.genType(arg.ty))
    }
    buff.WriteString(") ")
    // 返回值
    ret := "("
    if f.hasRet {
        ret += s2g.genType(f.retTy)
        ret += ", "
    }
    buff.WriteString(ret + "error) {")
    s2g.Write(buff.String())
    // 函数体
    s2g.Write("p := codec.NewPacker()")
    if f.hasRet {
        s2g.Write("var ret " + s2g.genType(f.retTy))
    }
    s2g.Write("var err error")
    
    // 请求参数
    hasOut := false
    for i, arg := range f.args {
        if arg.out {
            hasOut = true
        } else {
            dummy := &structMember{}
            dummy.name = arg.name
            dummy.ty = arg.ty
            dummy.tag = uint32(i+1)
            s2g.genWriteVar(dummy, "", f.hasRet, true)
        }
    }

    s2g.Write("var rsp protocol.ResponsePacket")
    s2g.Write(`err = s.proxy.Invoke("` + f.name + `", p.ToBytes(), &rsp)`)
    s2g.Write(genCheckErr(f.hasRet, true))

    if hasOut || f.hasRet {
        s2g.Write("up := codec.NewUnPacker([]byte(rsp.SRspPayload))")
    }

    // 返回值
    if f.hasRet {
        dummy := &structMember{}
        dummy.ty = f.retTy
        dummy.name = "ret"
        dummy.tag = 0
        dummy.require = true
        s2g.genReadVar(dummy, "", f.hasRet, true)
    }

    // 返回参数
    for i, arg := range f.args {
        if arg.out {
            dummy := &structMember{}
            dummy.name = "(*"+arg.name+")"
            dummy.ty = arg.ty
            dummy.tag = uint32(i+1)
            dummy.require = true
            s2g.genReadVar(dummy, "", f.hasRet, true)
        }
    }

    if f.hasRet {
        s2g.Write("return ret, nil")
    } else {
        s2g.Write("return nil")
    }

    s2g.Write("}")
}

func genCheckErr(checkRet bool, isReturn bool) string {
    var errStr string
    if checkRet {
        errStr = "return ret, err"
    } else {
        errStr = "return err"
    }
    if !isReturn {
        errStr = "break"
    }
    return `if err != nil {
    ` + errStr + `
    }`
}

func (s2g *sdp2Go) genWriteVar(v *structMember, prefix string, checkRet bool, isReturn bool) {
    switch v.ty.ty {
    case tkTVector:
        s2g.genWriteVector(v, prefix, checkRet)
    case tkTMap:
        s2g.genWriteMap(v, prefix, checkRet)
    case tkName:
        if v.ty.customTy == tkEnum {
            tag := strconv.Itoa(int(v.tag))
            if v.defVal != "" {
                s2g.Write("if " + prefix+v.name + " != " + v.defVal + " {")
            }
            s2g.Write("err = p.WriteInt32(" + tag + ", int32(" + prefix + v.name + "))")
            s2g.Write(genCheckErr(checkRet, isReturn))
            if v.defVal != "" {
                s2g.Write("}")
            }
        } else {
            s2g.genWriteStruct(v, prefix, checkRet, isReturn)
        }
    default:
        tag := strconv.Itoa(int(v.tag))
        if v.defVal != "" {
            s2g.Write("if " + prefix+v.name + " != " + v.defVal + " {")
        }
        s2g.Write("err = p.Write" + upperFirstLetter(s2g.genType(v.ty)) + "(" + tag + ", " + prefix + v.name + ")")
        s2g.Write(genCheckErr(checkRet, isReturn))
        if v.defVal != "" {
            s2g.Write("}")
        }
    }
}

func (s2g *sdp2Go) genReadVar(v *structMember, prefix string, checkRet bool, isReturn bool) {
    switch v.ty.ty {
    case tkTVector:
        s2g.genReadVector(v, prefix, checkRet)
    case tkTMap:
        s2g.genReadMap(v, prefix, checkRet)
    case tkName:
        if v.ty.customTy == tkEnum {
            tag := strconv.Itoa(int(v.tag))
            require := "false"
            if v.require {
                require = "true"
            }
            s2g.Write("err = up.ReadInt32((*int32)(&" + prefix + v.name + "), " + tag + ", " + require + ")")
            s2g.Write(genCheckErr(checkRet, isReturn))
        } else {
            s2g.genReadStruct(v, prefix, checkRet, isReturn)
        }
    default:
        tag := strconv.Itoa(int(v.tag))
        require := "false"
        if v.require {
            require = "true"
        }
        s2g.Write("err = up.Read" + upperFirstLetter(s2g.genType(v.ty)) + "(&" + prefix + v.name + ", " + tag + ", " + require + ")")
        s2g.Write(genCheckErr(checkRet, isReturn))
    }
}

func (s2g *sdp2Go) genReadVector(v *structMember, prefix string, checkRet bool) {
    tag := strconv.Itoa(int(v.tag))
    require := "false"
    if v.require {
        require = "true"
    }

    s2g.Write(`
has, ty, err = up.SkipToTag(` + tag + `, ` + require + `)
if !has || err != nil {
    return err
}
if ty != codec.SdpType_Vector {
    return fmt.Errorf("tag:%d got wrong type %d", ` + tag + `, ty)
}

_, length, err = up.ReadNumber32()
if err != nil {
    return err
}
` + prefix + v.name + ` = make(` + s2g.genType(v.ty) + `, length, length)
for i := uint32(0); i < length; i++ {`)

    vmember := &structMember{}
    vmember.require = true
    vmember.ty = v.ty.typeK
    vmember.name = v.name + "[i]"
    s2g.genReadVar(vmember, prefix, checkRet, true)

    s2g.Write("}")
}

func (s2g *sdp2Go) genWriteVector(v *structMember, prefix string, checkRet bool) {
    tag := strconv.Itoa(int(v.tag))

    s2g.Write(`
err = p.WriteHeader(` + tag + `, codec.SdpType_Vector)
if err != nil {
    return err
}
length = len(`+ prefix + v.name +`)
err = p.WriteNumber32(uint32(length))
if err != nil {
    return err
}
for _,v := range `+ prefix + v.name+` {`)
    vmember := &structMember{}
    vmember.name = "v"
    vmember.ty = v.ty.typeK
    s2g.genWriteVar(vmember, "", checkRet, true)

    s2g.Write("}")
}

func (s2g *sdp2Go) genReadMap(v *structMember, prefix string, checkRet bool) {
    tag := strconv.Itoa(int(v.tag))
    require := "false"
    if v.require {
        require = "true"
    }

    s2g.Write(`
has, ty, err = up.SkipToTag(` + tag + `, ` + require + `)
if !has || err != nil {
    return err
}
if ty != codec.SdpType_Map {
    return fmt.Errorf("tag:%d got wrong type %d", ` + tag + `, ty)
}

_, length, err = up.ReadNumber32()
if err != nil {
    return err
}
` + prefix + v.name + ` = make(` + s2g.genType(v.ty) + `)
for i := uint32(0); i < length; i++ {`)

    s2g.Write("var k " + s2g.genType(v.ty.typeK))
    kmember := &structMember{}
    kmember.require = true
    kmember.ty = v.ty.typeK
    kmember.name = "k"
    s2g.genReadVar(kmember, "", checkRet, true)

    s2g.Write("var v " + s2g.genType(v.ty.typeV))
    vmember := &structMember{}
    vmember.require = true
    vmember.ty = v.ty.typeV
    vmember.name = "v"
    vmember.tag = 1
    s2g.genReadVar(vmember, "", checkRet, true)
    s2g.Write(prefix + v.name + "[k] = v")

    s2g.Write("}")
}

func (s2g *sdp2Go) genWriteMap(v *structMember, prefix string, checkRet bool) {
    tag := strconv.Itoa(int(v.tag))

    s2g.Write(`
err = p.WriteHeader(` + tag + `, codec.SdpType_Map)
if err != nil {
    return err
}
length = len(`+ prefix + v.name +`)
err = p.WriteNumber32(uint32(length))
if err != nil {
    return err
}
for _k,_v := range `+ prefix + v.name+` {`)
    kmember := &structMember{}
    kmember.name = "_k"
    kmember.ty = v.ty.typeK
    s2g.genWriteVar(kmember, "", checkRet, true)

    vmember := &structMember{}
    vmember.name = "_v"
    vmember.ty = v.ty.typeV
    vmember.tag = 1
    s2g.genWriteVar(vmember, "", checkRet, true)

    s2g.Write("}")
}

func (s2g *sdp2Go) genReadStruct(v *structMember, prefix string, checkRet bool, isReturn bool) {
    tag := strconv.Itoa(int(v.tag))
    require := "false"
    if v.require {
        require = "true"
    }

    s2g.Write("err = " + prefix + v.name + ".ReadStructFromTag(up, " + tag + ", " + require + ")")
    s2g.Write(genCheckErr(checkRet, isReturn))
}

func (s2g *sdp2Go) genWriteStruct(v *structMember, prefix string, checkRet bool, isReturn bool) {
    tag := strconv.Itoa(int(v.tag))
    s2g.Write("err = " + prefix + v.name + ".WriteStructFromTag(p, " + tag + ")")
    s2g.Write(genCheckErr(checkRet, isReturn))
}

func newSdp2Go(file string, dir string) *sdp2Go {
    defer func() {
        if err := recover(); err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
    }()

    if dir != "" {
        tmp := []byte(dir)
        if tmp[len(tmp)-1] != byte('/') {
            dir += "/"
        }
    }

    return &sdp2Go{dir: dir, p: newParse(file)}
}
