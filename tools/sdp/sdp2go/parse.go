package main

import (
    "io/ioutil"
    "strconv"
    "sort"
    "strings"
)

// 类型
type varType struct {
    ty TK         // 基础类型
    unsigned bool   // 是否是无符号
    customTyName string // 自定义类型名称
    customTy    TK // 自定义类型类型tkEnum,tkStruct
    typeK *varType // vector的类型或map的key类型
    typeV *varType // map的value类型
}

// 常量类型
type constInfo struct {
    ty *varType // 类型
    name string
    value string
}

// 枚举类型
type enumMember struct {
    name string
    value int32
}
type enumInfo struct {
    name string
    members []enumMember
}

// 结构体类型
type structMember struct {
    tag uint32
    require bool
    ty *varType
    name string
    oldname string
    defVal string
    defTy TK
}
type structMemberSorter []structMember
func (a structMemberSorter) Len() int {return len(a)}
func (a structMemberSorter) Swap(i, j int) {a[i],a[j] = a[j],a[i]}
func (a structMemberSorter) Less(i, j int) bool {return a[i].tag < a[j].tag}
type structInfo struct {
    name string
    members []structMember
}

type parse struct {
    file string // 要解析的文件

    // 模块名称
    module string
    // 依赖的模块名称
    dependModule map[string]bool

    // includes
    includes_ []string
    includes []*parse

    // 常量
    consts []constInfo
    // 枚举
    enums []enumInfo
    // 结构体
    structs []structInfo

    lex *LexState // 词法分析器
    t *Token // 当前token
    lastT *Token  // 上一个token
}

func (p *parse) parse() {
OUT:
    for {
        p.next()
        t := p.t
        switch t.T {
        case tkEos:
            break OUT
        case tkInclude:
            p.parseInclude()
        case tkModule:
            p.parseModule()
        default:
            p.parseErr("expect include or module")
        }
    }
    p.analyzeDepend()
}

func (p *parse) next() {
    p.lastT = p.t
    p.t = p.lex.NextToken()
}

func (p *parse) expect(t TK) {
    p.next()
    if p.t.T != t {
        p.parseErr("expect " + TokenMap[t])
    }
}

func (p *parse) parseErr(err string) {
    line := "0"
    if p.t != nil {
        line = strconv.Itoa(p.t.Line)
    }

    panic(p.file + ": " + line + ". " + err)
}

func (p *parse) parseInclude() {
    p.expect(tkString)
    p.includes_ = append(p.includes_, p.t.S.S)
}

func (p *parse) parseModule() {
    p.expect(tkName)

    if p.module != "" {
        p.parseErr("can't define module more than once")
    }
    p.module = p.t.S.S

    p.expect(tkBracel)
    for {
        p.next()
        t := p.t
        switch t.T {
        case tkBracer:
            p.expect(tkSemi);
            return
        case tkConst:
            p.parseConst()
        case tkEnum:
            p.parseEnum()
        case tkStruct:
            p.parseStruct()
        default:
            p.parseErr("parse module error, not expected " + TokenMap[t.T])
        }
    }
}

func (p *parse) parseConst() {
    cst := constInfo{}

    p.next()
    switch p.t.T {
    case tkTBool,tkTByte,tkTShort,
        tkTInt,tkTLong,tkTFloat,
        tkTDouble,tkTString,tkUnsigned:
        cst.ty = p.parseType()
    default:
        p.parseErr("unsupported const type:" + TokenMap[p.t.T])
    }

    p.expect(tkName)
    cst.name = p.t.S.S

    p.expect(tkEq)
    p.next()
    switch p.t.T {
    case tkInteger,tkFloat:
        if !isNumberType(cst.ty.ty) {
            p.parseErr("expect a non-number")
        }
        cst.value = p.t.S.S
    case tkString:
        if cst.ty.ty != tkTString {
            p.parseErr("expect a non-string")
        }
        cst.value = `"` + p.t.S.S + `"`
    case tkTrue:
        if cst.ty.ty != tkTBool {
            p.parseErr("expect a non-bool")
        } 
        cst.value = "true"
    case tkFalse:
        if cst.ty.ty != tkTBool {
            p.parseErr("expect a non-bool")
        } 
        cst.value = "false"
    default:
        p.parseErr("unsupported const type " + TokenMap[p.t.T])
    }
    p.expect(tkSemi)

    p.consts = append(p.consts, cst)
}

func (p *parse) parseType() *varType {
    vtype := &varType{ty: p.t.T}

    switch vtype.ty {
    case tkName:
        // 自定义类型
        vtype.customTyName = p.t.S.S
    case tkTInt,tkTBool,tkTShort,
        tkTLong,tkTByte,tkTFloat,
        tkTDouble,tkTString:
        // 基础类型不需要处理
    case tkTVector:
        p.expect(tkShl)
        p.next()
        vtype.typeK = p.parseType()
        p.expect(tkShr)
    case tkTMap:
        p.expect(tkShl)
        p.next()
        vtype.typeK = p.parseType()
        p.expect(tkComma)
        p.next()
        vtype.typeV = p.parseType()
        p.expect(tkShr)
    case tkUnsigned:
        p.next()
        vtype2 := p.parseType()
        switch vtype2.ty {
        case tkTInt,tkTShort,tkTByte,tkTLong:
           vtype2.unsigned = true 
        default:
            p.parseErr("unsupported unsigned " + TokenMap[vtype2.ty])
        }
        return vtype2
    default:
        p.parseErr("unsupported type " + TokenMap[vtype.ty])
    }

    return vtype
}

func (p *parse) parseEnum() {
    en := enumInfo{}

    p.expect(tkName)
    en.name = p.t.S.S
    // 枚举名字不能重复
    for _, v := range p.enums {
        if v.name == en.name {
            p.parseErr(en.name + " Redefine")
        }
    }

    p.expect(tkBracel)
    defer func() {
        p.expect(tkSemi)
        p.enums = append(p.enums, en)
    }()

    var i int32
    for {
        p.next()
        switch p.t.T {
        case tkBracer:
            return
        case tkName:
            name := p.t.S.S
            p.next()
            switch p.t.T {
            case tkComma:
                member := enumMember{name: name, value: i}
                en.members = append(en.members, member)
                i++
            case tkBracer:
                member := enumMember{name: name, value: i}
                en.members = append(en.members, member)
                return
            case tkEq:
                p.expect(tkInteger)
                i = int32(p.t.S.I)
                member := enumMember{name: name, value: i}
                en.members = append(en.members, member)
                i++
                p.next()
                if p.t.T == tkBracer {
                    return
                } else if p.t.T == tkComma {
                    // 继续解析下一条
                } else {
                    p.parseErr("expect , or }")
                }
            default:
                p.parseErr("not expected " + TokenMap[p.t.T] + " in enum")
            }
        default:
            p.parseErr("not expected " + TokenMap[p.t.T] + " in enum")
        }
    }
}

func (p *parse) parseStruct() {
    st := structInfo{}
    p.expect(tkName)
    st.name = p.t.S.S

    // 结构体名字不能重复
    for _, v := range p.structs {
        if v.name == st.name {
            p.parseErr(st.name + " Redefine")
        }
    }

    p.expect(tkBracel)
    for {
        member := p.parseStructMember()
        if member == nil {
            break
        }
        st.members = append(st.members, *member)
    }
    p.expect(tkSemi)

    // 检查tag是否重复
    tags := make(map[uint32]bool)
    for _, v := range st.members {
        if tags[v.tag] {
            p.parseErr("tag = " + strconv.Itoa(int(v.tag)) + " duplicate")
        }
        tags[v.tag] = true
    }
    // 对tag排序
    sort.Sort(structMemberSorter(st.members))

    p.structs = append(p.structs, st)
}

func (p *parse) parseStructMember() *structMember {
    p.next()
    if p.t.T == tkBracer {
        // 空结构体
        return nil
    }

    // tag
    if p.t.T != tkInteger {
        p.parseErr("expect tag")
    }
    member := &structMember{}
    member.tag = uint32(p.t.S.I)

    // require/optional
    p.next()
    if p.t.T == tkRequire {
        member.require = true
    } else if p.t.T == tkOptional {
        member.require = false
    } else {
        p.parseErr("expect require or optional")
    }

    // type
    p.next()
    if !isType(p.t.T) && p.t.T != tkName && p.t.T != tkUnsigned {
        p.parseErr("expect type")
    } else {
        member.ty = p.parseType()
    }

    // name
    p.expect(tkName)
    member.name = p.t.S.S
    p.next()
    if p.t.T == tkSemi {
        return member
    }

    // 解析默认值
    if p.t.T != tkEq {
        p.parseErr("expect ; or =")
    }
    if p.t.T == tkTVector || p.t.T == tkTMap || p.t.T == tkName {
        p.parseErr("vector, map, custom type can't set default value")
    }
    p.next()
    p.parseStructMemberDefault(member)
    p.expect(tkSemi)

    return member
}

func (p *parse) parseStructMemberDefault(member *structMember) {
    member.defTy = p.t.T
    switch p.t.T {
    case tkInteger:
        if !isNumberType(member.ty.ty) {
            p.parseErr("expect a non-number")
        }
        member.defVal = p.t.S.S
    case tkFloat:
        if !isNumberType(member.ty.ty) {
            p.parseErr("expect a non-number") 
        }
        member.defVal = p.t.S.S
    case tkString:
        if member.ty.ty != tkTString {
            p.parseErr("expect a non-string")
        }
        member.defVal = `"` + p.t.S.S + `"`
    case tkTrue:
        if member.ty.ty != tkTBool {
            p.parseErr("expect a non-bool")
        }
        member.defVal = "true"
    case tkFalse:
        if member.ty.ty != tkTBool {
            p.parseErr("expect a non-bool")
        }
        member.defVal = "false"
    default:
        p.parseErr("unsupported default value type " + TokenMap[p.t.T])
    }
}

func (p *parse) analyzeDepend() {
    // 解析include
    for _, v := range p.includes_ {
        p2 := newParse(v)
        p.includes = append(p.includes, p2)
    }

    // 解析结构体中的自定义类型是否合法
    for _, v := range p.structs {
        for _, m := range v.members {
            p.checkCustomType(m.ty)
        }
    }
}

func (p *parse) checkCustomType(ty *varType) {
    if ty.ty == tkName {
        name := ty.customTyName
        if strings.Count(name, "::") == 0 {
            name = p.module + "::" + name
        }

        ty2,dependModule := p.findCustomType(name)
        if ty2 == tkName {
            p.parseErr("can't find " + name + " define")
        }
        ty.customTy = ty2
        if dependModule != p.module {
            if p.dependModule == nil {
                p.dependModule = make(map[string]bool)
            }
            p.dependModule[dependModule] = true
        }
    } else if ty.ty == tkTVector {
        p.checkCustomType(ty.typeK)
    } else if ty.ty == tkTMap {
        p.checkCustomType(ty.typeK)
        p.checkCustomType(ty.typeV)
    }
}

func (p *parse) findCustomType(name string) (TK,string) {
    for _, v := range p.structs {
        if p.module+"::"+v.name  == name {
            return tkStruct,p.module
        }
    }

    for _, v := range p.enums {
        if p.module+"::"+v.name  == name {
            return tkEnum,p.module
        }
    }

    for _, p2 := range p.includes {
        if ret,module := p2.findCustomType(name); ret != tkName {
            return ret,module
        }
    }

    return tkName,""
}

func newParse(file string) *parse {
    content, err := ioutil.ReadFile(file)
    if err != nil {
        panic("read file " + file + " err:" + err.Error())
    }

    p := &parse{file: file, lex: NewLexState(file, content)}
    p.parse()

    return p
}
