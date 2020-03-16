package main

import (
    "testing"
    "reflect"
    "fmt"
    "strings"
    "github.com/yellia1989/tex-go/tools/sdp/codec"
    "github.com/yellia1989/tex-go/tools/sdp/sdp2go/test"
    "github.com/yellia1989/tex-go/tools/sdp/sdp2go/Test2"
)

// 简单测试struct的读写
func TestStructSimple(t *testing.T) {
    var ss test.SimpleStruct

    ss.B = true
    ss.By = 1
    ss.S = -1010
    ss.Us = 1010
    ss.I = -65535
    ss.Ui = 6553567
    ss.L = -1000000000
    ss.Ul = 1000000000
    ss.F = 0.12345678
    ss.D = -0.12345678
    ss.Ss = "yellia"
    ss.Vi = append(ss.Vi, 1)
    ss.Vi = append(ss.Vi, 2)
    ss.Mi = make(map[int32]int32)
    ss.Mi[1] = 1
    ss.Mi[2] = 2

    p := codec.NewPacker()
    if err := ss.WriteStruct(p); err != nil {
        t.Fatalf("write struct faild:%s", err)
    }

    up := codec.NewUnPacker(p.ToBytes())
    var ss2 test.SimpleStruct
    ss2.ReadStruct(up)

    p.Reset()
    ss2.WriteStruct(p)

    if reflect.DeepEqual(ss, ss2) == false {
        t.Fatalf("%v != %v", ss, ss2)
    }
}

func TestRequireStruct(t *testing.T) {
    p := codec.NewPacker()

    up := codec.NewUnPacker(p.ToBytes())
    var ss2 test.RequireStruct
    if err := ss2.ReadStruct(up); err != nil {
        t.Logf("read struct err:%s", err)
    }

    p.Reset()
    var ss test.RequireStruct
    if err := ss.WriteStruct(p); err != nil {
        t.Fatalf("write struct err:%s", err)
    }
    up.Reset(p.ToBytes())
    if err := ss2.ReadStruct(up); err != nil {
        t.Fatalf("read struct err:%s", err)
    }
}

func TestDefaultStruct(t *testing.T) {
    p := codec.NewPacker()
    var ss test.DefaultStruct
    ss.ResetDefault()
    ss.WriteStruct(p)

    up := codec.NewUnPacker(p.ToBytes())
    var ss2 test.DefaultStruct
    if err := ss2.ReadStruct(up); err != nil {
        t.Fatalf("read struct err:%s", err)
    }
    if ss2.B != true ||
        ss2.By != 1 ||
        ss2.S != 10 ||
        ss2.I != 1 ||
        ss2.L != 0x0FFFFFFFFFFFFFFF ||
        ss2.Ss != "yellia" {
        t.Fatalf("read struct %v", ss2)
    }
}

func TestC(t *testing.T) {
    packer := codec.NewPacker()

    s := Test2.Student{}
    s.ResetDefault()
    s.IUid = 1234567890;
    s.SName = "学生1";
    s.IAge = 12;
    s.WriteStructFromTag(packer, 15, true)

    cl := Test2.Class{}
    cl.ResetDefault()
    cl.IId = 1001
    cl.SName = "c1"
    cl.VStudent = append(cl.VStudent, s)
    cl.VData = append(cl.VData, 'c')

    tc := Test2.Teacher{}
    tc.ResetDefault()
    tc.IId = 1001
    cl.VTeacher = append(cl.VTeacher, tc)
    cl.WriteStructFromTag(packer, 16, true)

    right := "7F0F00D285D8CC044107E5ADA6E7949F31020C807F1000E9074102633152017000D285D8CC044107E5ADA6E7949F31020C8043016354017000E90773808080"
    real := fmt.Sprintf("%X", packer.ToBytes())

    if strings.Index(right, real) == -1 {
        fmt.Printf("right:%s\n,real:%s\n", right, real)
    }
}
