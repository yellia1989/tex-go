package main

import (
    "github.com/yellia1989/tex-go/tools/sdp/sdp2go/test2"
    "github.com/yellia1989/tex-go/tools/sdp/sdp2go/test"
    "testing"
    "github.com/yellia1989/tex-go/tools/sdp/codec"
    "reflect"
)

// 简单测试struct的读写
/*
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
	Ss  string          `json:"s"`
	Vi  []int32         `json:"vi"`
	Mi  map[int32]int32 `json:"mi"`
	Age Age             `json:"age"`
}
*/
func TestStructSimple(t *testing.T) {
    var ss test2.SimpleStruct

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
    ss.Age = test2.Age_10

    p := codec.NewPacker()
    if err := ss.WriteStruct(p); err != nil {
        t.Fatalf("write struct faild:%s", err)
    }

    up := codec.NewUnPacker(p.ToBytes())
    var ss2 test2.SimpleStruct
    ss2.ReadStruct(up)

    p.Reset()
    ss2.WriteStruct(p)

    if reflect.DeepEqual(ss, ss2) == false {
        t.Fatalf("%v != %v", ss, ss2)
    }
}

/*
type RequireStruct struct {
	Ss test2.SimpleStruct `json:"ss"`
}
*/
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

/*
type DefaultStruct struct {
	B  bool   `json:"b"`
	By int8   `json:"by"`
	S  int16  `json:"s"`
	I  int32  `json:"i"`
	L  int64  `json:"l"`
	Ss string `json:"ss"`
}
struct defaultStruct {
    0   optional bool b = true;
    1   optional byte by = 1;
    2   optional short s = 10;
    3   optional int i = 1;
    4   optional long l = 0x0FFFFFFFFFFFFFFF;
    5   optional string ss = "yellia";
};
*/
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
