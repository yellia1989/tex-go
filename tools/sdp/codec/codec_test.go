package codec

import (
    "testing"
    "math"
)

func unpacker(p *Packer) *UnPacker {
    return NewUnPacker(p.ToBytes())
}

// 测试byte读写
func TestByte(t *testing.T) {
    for tag := 0; tag <= 250; tag++ {
        for i := 0; i <= math.MaxUint8; i++ {
            packer := NewPacker()
            err := packer.WriteByte(uint32(tag), byte(i))
            if err != nil {
                t.Error(err)
            }
            var b byte 
            err = unpacker(packer).ReadByte(&b, uint32(tag), true)
            if err != nil {
                t.Error(err)
            }
            if b != byte(i) {
                t.Errorf("real:%d expect:%d", b, i)
            }
        }
    }
}

// 测试bool读写
func TestBool(t *testing.T) {
    for tag := 0; tag <= 250; tag++ {
        var bb = [2]bool{true,false}
        for i := 0; i < len(bb); i++ {
            packer := NewPacker()
            err := packer.WriteBool(uint32(tag), bb[i])
            if err != nil {
                t.Error(err)
            }
            var b bool
            err = unpacker(packer).ReadBool(&b, uint32(tag), true)
            if err != nil {
                t.Error(err)
            }
            if b != bb[i] {
                t.Errorf("real:%t expect:%t", b, bb[i])
            }
        }
    }
}

// 测试int8读写
func TestInt8(t *testing.T) {
    for tag := 0; tag <= 250; tag++ {
        for i := math.MinInt8; i <= math.MaxInt8; i++ {
            packer := NewPacker()
            err := packer.WriteInt8(uint32(tag), int8(i))
            if err != nil {
                t.Error(err)
            }
            var v int8
            err = unpacker(packer).ReadInt8(&v, uint32(tag), true)
            if err != nil {
                t.Error(err)
            }
            if v != int8(i) {
                t.Errorf("real:%d expect:%d", v, i)
            }
        }
    }
}

// 测试uint8读写
func TestUint8(t *testing.T) {
    for tag := 0; tag <= 250; tag++ {
        for i := 0; i <= math.MaxUint8; i++ {
            packer := NewPacker()
            err := packer.WriteUint8(uint32(tag), uint8(i))
            if err != nil {
                t.Error(err)
            }
            var v uint8
            err = unpacker(packer).ReadUint8(&v, uint32(tag), true)
            if err != nil {
                t.Error(err)
            }
            if v != uint8(i) {
                t.Errorf("real:%d expect:%d", v, i)
            }
        }
    }
}

// 测试int16读写
func TestInt16(t *testing.T) {
    for tag := 0; tag <= 250; tag++ {
        for i := math.MinInt16; i <= math.MaxInt16; i++ {
            packer := NewPacker()
            err := packer.WriteInt16(uint32(tag), int16(i))
            if err != nil {
                t.Error(err)
            }
            var v int16
            err = unpacker(packer).ReadInt16(&v, uint32(tag), true)
            if err != nil {
                t.Error(err)
            }
            if v != int16(i) {
                t.Errorf("real:%d expect:%d", v, i)
            }
        }
    }
}

// 测试uint16读写
func TestUint16(t *testing.T) {
    for tag := 0; tag <= 250; tag++ {
        for i := 0; i <= math.MaxUint16; i++ {
            packer := NewPacker()
            err := packer.WriteUint16(uint32(tag), uint16(i))
            if err != nil {
                t.Error(err)
            }
            var v uint16
            err = unpacker(packer).ReadUint16(&v, uint32(tag), true)
            if err != nil {
                t.Error(err)
            }
            if v != uint16(i) {
                t.Errorf("real:%d expect:%d", v, i)
            }
        }
    }
}

// 测试int32读写
func TestInt32(t *testing.T) {
    for tag := 0; tag <= 250; tag++ {
        for i := math.MinInt32; i <= math.MinInt32+100; i++ {
            packer := NewPacker()
            err := packer.WriteInt32(uint32(tag), int32(i))
            if err != nil {
                t.Error(err)
            }
            var v int32
            err = unpacker(packer).ReadInt32(&v, uint32(tag), true)
            if err != nil {
                t.Error(err)
            }
            if v != int32(i) {
                t.Errorf("real:%d expect:%d", v, i)
            }
        }
        for i := math.MaxInt32; i >= math.MaxInt32-100; i-- {
            packer := NewPacker()
            err := packer.WriteInt32(uint32(tag), int32(i))
            if err != nil {
                t.Error(err)
            }
            var v int32
            err = unpacker(packer).ReadInt32(&v, uint32(tag), true)
            if err != nil {
                t.Error(err)
            }
            if v != int32(i) {
                t.Errorf("real:%d expect:%d", v, i)
            }
        }
    }
}

// 测试uint32读写
func TestUint32(t *testing.T) {
    for tag := 0; tag <= 250; tag++ {
        for i := math.MaxUint32-100; i <= math.MaxUint32; i++ {
            packer := NewPacker()
            err := packer.WriteUint32(uint32(tag), uint32(i))
            if err != nil {
                t.Error(err)
            }
            var v uint32
            err = unpacker(packer).ReadUint32(&v, uint32(tag), true)
            if err != nil {
                t.Error(err)
            }
            if v != uint32(i) {
                t.Errorf("real:%d expect:%d", v, i)
            }
        }
    }
}

// 测试int64读写
func TestInt64(t *testing.T) {
    for tag := 0; tag <= 250; tag++ {
        for i := math.MinInt64; i <= math.MinInt64+100; i++ {
            packer := NewPacker()
            err := packer.WriteInt64(uint32(tag), int64(i))
            if err != nil {
                t.Error(err)
            }
            var v int64
            err = unpacker(packer).ReadInt64(&v, uint32(tag), true)
            if err != nil {
                t.Error(err)
            }
            if v != int64(i) {
                t.Errorf("real:%d expect:%d", v, i)
            }
        }
        for i := math.MaxInt64; i >= math.MaxInt64-100; i-- {
            packer := NewPacker()
            err := packer.WriteInt64(uint32(tag), int64(i))
            if err != nil {
                t.Error(err)
            }
            var v int64
            err = unpacker(packer).ReadInt64(&v, uint32(tag), true)
            if err != nil {
                t.Error(err)
            }
            if v != int64(i) {
                t.Errorf("real:%d expect:%d", v, i)
            }
        }
    }
}

// 测试uint64读写
func TestUint64(t *testing.T) {
    for tag := 0; tag <= 250; tag++ {
        for i := uint64(math.MaxUint64-100); i < math.MaxUint64; i++ {
            packer := NewPacker()
            err := packer.WriteUint64(uint32(tag), i)
            if err != nil {
                t.Error(err)
            }
            var v uint64
            err = unpacker(packer).ReadUint64(&v, uint32(tag), true)
            if err != nil {
                t.Error(err)
            }
            if v != i {
                t.Errorf("real:%d expect:%d", v, i)
            }
        }
    }
}

// 测试float读写
func TestFloat(t *testing.T) {
    for tag := 0; tag <= 250; tag++ {
        for i := float32(math.SmallestNonzeroFloat32); i < math.SmallestNonzeroFloat32+100; i++ {
            packer := NewPacker()
            err := packer.WriteFloat32(uint32(tag), i)
            if err != nil {
                t.Error(err)
            }
            var v float32
            err = unpacker(packer).ReadFloat32(&v, uint32(tag), true)
            if err != nil {
                t.Error(err)
            }
            if v != i {
                t.Errorf("real:%f expect:%f", v, i)
            }
        }
        for i := float32(math.MaxFloat32-100); i < math.MaxFloat32; i++ {
            packer := NewPacker()
            err := packer.WriteFloat32(uint32(tag), i)
            if err != nil {
                t.Error(err)
            }
            var v float32
            err = unpacker(packer).ReadFloat32(&v, uint32(tag), true)
            if err != nil {
                t.Error(err)
            }
            if v != i {
                t.Errorf("real:%f expect:%f", v, i)
            }
        }
    }
}

// 测试double读写
func TestDouble(t *testing.T) {
    for tag := 0; tag <= 250; tag++ {
        for i := float64(math.SmallestNonzeroFloat64); i < math.SmallestNonzeroFloat64+100; i++ {
            packer := NewPacker()
            err := packer.WriteFloat64(uint32(tag), i)
            if err != nil {
                t.Error(err)
            }
            var v float64
            err = unpacker(packer).ReadFloat64(&v, uint32(tag), true)
            if err != nil {
                t.Error(err)
            }
            if v != i {
                t.Errorf("real:%f expect:%f", v, i)
            }
        }
        for i := float64(math.MaxFloat64-100); i < math.MaxFloat64; i++ {
            packer := NewPacker()
            err := packer.WriteFloat64(uint32(tag), i)
            if err != nil {
                t.Error(err)
            }
            var v float64
            err = unpacker(packer).ReadFloat64(&v, uint32(tag), true)
            if err != nil {
                t.Error(err)
            }
            if v != i {
                t.Errorf("real:%f expect:%f", v, i)
            }
        }
    }
}

// 测试string读写
func TestString(t *testing.T) {
    ss := [...]string {
        "hello world",
        "中文字符",
        "包括字符串结束符的字符串\x00\xab",
    }

    for tag := 1; tag < 250; tag++ {
        for i := 0; i < len(ss); i++ {
            packer := NewPacker()
            err := packer.WriteString(uint32(tag), ss[i])
            if err != nil {
                t.Error(err)
            }
            var v string
            err = unpacker(packer).ReadString(&v, uint32(tag), true)
            if err != nil {
                t.Error(err)
            }
            if v != ss[i] {
                t.Errorf("real:%s expect:%s", v, ss[i])
            }
        }
    }
}

// 测试uint32读写性能
func BenchmarkUint32(t *testing.B) {
    packer := NewPacker()

    for i := 0; i < 200; i++ {
        err := packer.WriteUint32(uint32(i), uint32(0xffffffff))
        if err != nil {
            t.Error(err)
        }
    }

    up := unpacker(packer)

    for i := 0; i < 200; i++ {
        var v uint32
        err := up.ReadUint32(&v, uint32(i), true)
        if err != nil {
            t.Error(err)
        }
        if v != uint32(0xffffffff) {
            t.Errorf("real:%d expect:%d", v, uint32(0xffffffff))
        }
    }
}

// 测试string读写性能
func BenchmarkString(t *testing.B) {
    packer := NewPacker()

    for i := 0; i < 200; i++ {
        err := packer.WriteString(uint32(i), "hahahahahahahahahahahahahahahahahahahaha")
        if err != nil {
            t.Error(err)
        }
    }

    up := unpacker(packer)

    for i := 0; i < 200; i++ {
        var v string
        err := up.ReadString(&v, uint32(i), true)
        if err != nil {
            t.Error(err)
        }
        if v != "hahahahahahahahahahahahahahahahahahahaha" {
            t.Error("no eq.")
        }
    }
}
