package main

import (
    "testing"
    "reflect"
    "fmt"
    "strings"
    "encoding/hex"
    "github.com/yellia1989/tex-go/tools/sdp/codec"
    "github.com/yellia1989/tex-go/tools/sdp/sdp2go/test"
    "github.com/yellia1989/tex-go/tools/sdp/sdp2go/Test2"
    "github.com/yellia1989/tex-go/tools/sdp/sdp2go/rpc"
)

// 简单测试struct的读写
func TestStructSimple(t *testing.T) {
    ss := test.NewSimpleStruct()

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
    ss2 := test.NewSimpleStruct()
    if err := ss2.ReadStruct(up); err != nil {
        t.Fatalf("read struct faild:%s", err)
    }

    if reflect.DeepEqual(ss, ss2) == false {
        t.Fatalf("%v != %v", ss, ss2)
    }
}

func TestRequireStruct(t *testing.T) {
    p := codec.NewPacker()

    up := codec.NewUnPacker(p.ToBytes())
    ss2 := test.NewRequireStruct()
    if err := ss2.ReadStruct(up); err != nil {
        t.Logf("read struct err:%s", err)
    }

    p.Reset()
    ss := test.NewRequireStruct()
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
    ss := test.NewDefaultStruct()
    ss.WriteStruct(p)

    up := codec.NewUnPacker(p.ToBytes())
    ss2 := test.NewDefaultStruct()
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

func TestPrintSdp(t *testing.T) {
    s := Test2.NewStudent()
    s.IUid = 1234567890
    s.SName = "学生1"
    s.IAge = 12
    s.MSecret = make(map[string]string)
    s.MSecret["yellia"] = "hello"
    s.MSecret["luo"] = "juan"

    cl := Test2.NewClass()
    cl.IId = 1001
    cl.SName = "c1"
    cl.VStudent = append(cl.VStudent, *s.Copy())
    cl.VStudent = append(cl.VStudent, *s.Copy())
    cl.VData = append(cl.VData, 'c')

    tc := Test2.NewTeacher()
    tc.IId = 1001
    tc.S1 = *s.Copy()
    cl.VTeacher = append(cl.VTeacher, *tc.Copy())

    sc := Test2.NewSchool()
    sc.MClass = make(map[uint32]Test2.Class)
    sc.MClass[1001] = *cl.Copy()
    sc.MClass[1002] = *cl.Copy()

    t.Log("\n"+codec.PrintSdp(sc))
}

func TestC(t *testing.T) {
    packer := codec.NewPacker()

    s := Test2.NewStudent()
    s.IUid = 1234567890
    s.SName = "学生1"
    s.IAge = 12
    s.WriteStructFromTag(packer, 15, true)

    cl := Test2.NewClass()
    cl.IId = 1001
    cl.SName = "c1"
    cl.VStudent = append(cl.VStudent, *s.Copy())
    cl.VData = append(cl.VData, 'c')

    tc := Test2.NewTeacher()
    tc.IId = 1001
    cl.VTeacher = append(cl.VTeacher, *tc.Copy())
    cl.WriteStructFromTag(packer, 16, true)

    right := "7F0F00D285D8CC044107E5ADA6E7949F31020C807F1000E9074102633152017000D285D8CC044107E5ADA6E7949F31020C8043016354017000E90773808080"
    real := fmt.Sprintf("%X", packer.ToBytes())

    if strings.Index(right, real) == -1 {
        fmt.Printf("right:%s\n,real:%s\n", right, real)
    }
}

func TestMailDataInfo(t *testing.T) {
    right := "70002e4101614313323032302d30342d33302031353a33313a32364401614501615a02000100025c0100010f140780"
    decoded, err := hex.DecodeString(right)
    if err != nil {
        t.Fatal(err)
    }

    mail := rpc.NewMailDataInfo()
    codec.StringToSdp(decoded, mail)

    t.Log(codec.PrintSdp(mail))

    real := hex.EncodeToString(codec.SdpToString(mail))
    if real != right {
        t.Fatalf("MailDataInfo test failed, real:%s vs input:%s", real, right)
    }
}
