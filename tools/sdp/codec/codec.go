package codec

import (
    "bytes"
    "encoding/binary"
    "fmt"
    "math"
    "io"
)

// sdp 支持的数据类型
const (
    SdpType_Integer_Positive = 0 
    SdpType_Integer_Negative = 1
    SdpType_Float = 2
    SdpType_Double = 3
    SdpType_String = 4
    SdpType_Vector = 5
    SdpType_Map = 6
    SdpType_StructBegin = 7
    SdpType_StructEnd = 8
)

type Packer struct {
    buf *bytes.Buffer
}

func (p *Packer) writeData(buf []byte) error {
    _, err := p.buf.Write(buf)
    return err
}

func (p *Packer) WriteHeader(tag uint32, ty uint32) error {
    if tag < 15 {
        data := (ty << 4) | tag
        return p.WriteNumber32(data)
    }
    data := (ty << 4) | 0x0F
    if err := p.WriteNumber32(data); err != nil {
        return err
    }
    return p.WriteNumber32(tag)
}

func (p *Packer) WriteNumber32(v uint32) error {
    var b [5]byte
    var bs []byte
    bs = b[:]
    n := binary.PutUvarint(bs, uint64(v))
    return p.writeData(bs[0:n])
}

func (p *Packer) WriteNumber64(v uint64) error {
    var b [10]byte
    var bs []byte
    bs = b[:]
    n := binary.PutUvarint(bs, v)
    return p.writeData(bs[0:n])
}

func (p *Packer) WriteBool(tag uint32, v bool) error {
    tmp := int8(0)
    if v {
        tmp = 1
    }
    return p.WriteInt8(tag, tmp)
}

func (p *Packer) WriteByte(tag uint32, v byte) error {
    return p.WriteUint32(tag, uint32(v))
}

func (p *Packer) WriteUint8(tag uint32, v uint8) error {
    return p.WriteUint32(tag, uint32(v))
}

func (p *Packer) WriteInt8(tag uint32, v int8) error {
    return p.WriteInt32(tag, int32(v))
}

func (p *Packer) WriteUint16(tag uint32, v uint16) error {
    return p.WriteUint32(tag, uint32(v))
}

func (p *Packer) WriteInt16(tag uint32, v int16) error {
    return p.WriteInt32(tag, int32(v))
}

func (p *Packer) WriteUint32(tag uint32, v uint32) error {
    if err := p.WriteHeader(tag, SdpType_Integer_Positive); err != nil {
        return err
    }
    return p.WriteNumber32(v)
}

func (p *Packer) WriteInt32(tag uint32, v int32) error {
    if v < 0 {
        if err := p.WriteHeader(tag, SdpType_Integer_Negative); err != nil {
            return err
        }
        return p.WriteNumber32(uint32(-v))
    }
    return p.WriteUint32(tag, uint32(v))
}

func (p *Packer) WriteUint64(tag uint32, v uint64) error {
    if err := p.WriteHeader(tag, SdpType_Integer_Positive); err != nil {
        return err
    }
    return p.WriteNumber64(v)
}

func (p *Packer) WriteInt64(tag uint32, v int64) error {
    if v < 0 {
        if err := p.WriteHeader(tag, SdpType_Integer_Negative); err != nil {
            return err
        }
        return p.WriteNumber64(uint64(-v))
    }
    return p.WriteUint64(tag, uint64(v))
}

func (p *Packer) WriteFloat32(tag uint32, v float32) error {
    if err := p.WriteHeader(tag, SdpType_Float); err != nil {
        return err
    }
    return p.WriteNumber32(math.Float32bits(v));
}

func (p *Packer) WriteFloat64(tag uint32, v float64) error {
    if err := p.WriteHeader(tag, SdpType_Double); err != nil {
        return err
    }
    return p.WriteNumber64(math.Float64bits(v));
}

func (p *Packer) WriteString(tag uint32, v string) error {
    if err := p.WriteHeader(tag, SdpType_String); err != nil {
        return err
    }
    if err := p.WriteNumber32(uint32(len(v))); err != nil {
        return err
    }
    _, err := p.buf.WriteString(v)
    return err
}

func (p *Packer) ToBytes() []byte {
    return p.buf.Bytes()
}

func (p *Packer) Grow(n int) {
    p.buf.Grow(n)
}

func (p *Packer) Reset() {
    p.buf.Reset()
}

type UnPacker struct {
    buf *bytes.Reader
}

func (up *UnPacker) readHeader() (n int, tag uint32, ty uint32, err error) {
    data, err := up.buf.ReadByte()
    if err != nil {
        return 0, 0, 0, err
    }

    n = 1
    ty = uint32(data) >> 4
    tag = uint32(data) & 0x0F
    if tag == 15 {
        n1 := 0
        n1, tag, err = up.ReadNumber32()
        if err != nil {
            up.buf.UnreadByte()
            return 0, 0, 0, err
        }
        n += n1
    }

    return n, tag, ty, err
}

func (up *UnPacker) unreadHeader(n int) error {
    _, err := up.buf.Seek(int64(-n), io.SeekCurrent)
    return err
}

func (up *UnPacker) skip(n uint32) error {
    _, err := up.buf.Seek(int64(n), io.SeekCurrent)
    return err
}

func (up *UnPacker) skipCurField() error {
    _, _, ty, err := up.readHeader()
    if err != nil {
        return err
    }

    return up.skipField(ty)
}

func (up *UnPacker) skipVector() error {
    // 读vector长度
    _, len, err := up.ReadNumber32()
    if err != nil {
        return err
    }

    for i := uint32(0); i < len; i++ {
        if err := up.skipCurField(); err != nil {
            return err
        }
    }

    return nil
}

func (up *UnPacker) skipMap() error {
    // 读map长度
    _, len, err := up.ReadNumber32()
    if err != nil {
        return err
    }

    for i := uint32(0); i < len; i++ {
        if err := up.skipCurField(); err != nil {
            return err
        }
        if err := up.skipCurField(); err != nil {
            return err
        }
    }

    return nil
}

func (up *UnPacker) SkipStruct() error {
    for {
        _, _, ty, err := up.readHeader()
        if err != nil {
            return err
        }
        if ty == SdpType_StructEnd {
            break
        }
        if err := up.skipField(ty); err != nil {
            return err
        }
    }

    return nil
}

func (up *UnPacker) skipField(ty uint32) error {
    switch ty {
    case SdpType_Integer_Positive,SdpType_Integer_Negative,SdpType_Float,SdpType_Double:
        if _, _, err := up.ReadNumber64(); err != nil {
            return err
        }
    case SdpType_String:
         _, len, err := up.ReadNumber32()
        if err != nil {
            return err
        }
        if err := up.skip(len); err != nil {
            return err
        }
    case SdpType_Vector:
        if err := up.skipVector(); err != nil {
            return err
        }
    case SdpType_Map:
        if err := up.skipMap(); err != nil {
            return err
        }
    case SdpType_StructBegin:
        if err := up.SkipStruct(); err != nil {
            return err
        }
    default:
        return fmt.Errorf("unknown type:%d", ty)
    }
    return nil
}

func (up *UnPacker) SkipToTag(tag uint32, require bool) (has bool, ty uint32, err error) {
    for {
        if up.buf.Len() == 0 {
            break
        }
        n, curTag, curTy, err := up.readHeader()
        if err != nil {
            return false, 0, fmt.Errorf("tag:%d read header err:%s", tag, err.Error())
        }
        
        if curTy == SdpType_StructEnd || curTag > tag {
            // 多读了一个header
            if err := up.unreadHeader(n); err != nil {
                return false, 0, fmt.Errorf("tag:%d unread header err:%s", tag, err.Error())
            }
            break
        }
        if curTag == tag {
            return true, curTy, nil
        }
        // 跳过不需要的field
        if err:= up.skipField(curTy); err != nil {
            return false, 0, fmt.Errorf("tag:%d skip not enough data, err:%s", tag, err.Error())
        }
    }

    if require {
        return false, 0, fmt.Errorf("tag:%d field not exist", tag)
    }
    return false, 0, nil
}

func (up *UnPacker) ReadNumber32() (n int, v uint32, err error) {
    size := up.buf.Len()
    u64, err := binary.ReadUvarint(up.buf)
    return (size-up.buf.Len()), uint32(u64), err
}

func (up *UnPacker) ReadNumber64() (n int, v uint64, err error) {
    size := up.buf.Len()
    u64, err := binary.ReadUvarint(up.buf)
    return (size-up.buf.Len()), u64, err
}

func (up *UnPacker) ReadBool(v *bool, tag uint32, require bool) error {
    i8 := int8(0)
    if *v {
        i8 = 1
    }
    err := up.ReadInt8(&i8, tag, require)
    if i8 == 1 {
        *v = true
    } else {
        *v = false
    }
    return err
}

func (up *UnPacker) ReadByte(v *byte, tag uint32, require bool) error {
    u8 := uint8(*v)
    err := up.ReadUint8(&u8, tag, require)
    *v = byte(u8)
    return err
}

func (up *UnPacker) ReadUint8(v *uint8, tag uint32, require bool) error {
    u32 := uint32(*v)
    err := up.ReadUint32(&u32, tag, require)
    *v = uint8(u32)
    return err
}

func (up *UnPacker) ReadInt8(v *int8, tag uint32, require bool) error {
    i32 := int32(*v)
    err := up.ReadInt32(&i32, tag, require)
    *v = int8(i32)
    return err
}

func (up *UnPacker) ReadUint16(v *uint16, tag uint32, require bool) error {
    u32 := uint32(*v)
    err := up.ReadUint32(&u32, tag, require)
    *v = uint16(u32)
    return err
}

func (up *UnPacker) ReadInt16(v *int16, tag uint32, require bool) error {
    i32 := int32(*v)
    err := up.ReadInt32(&i32, tag, require)
    *v = int16(i32)
    return err
}

func (up *UnPacker) ReadUint32(v *uint32, tag uint32, require bool) error {
    has, ty, err := up.SkipToTag(tag, require)
    if !has || err != nil {
        return err
    }
    if ty != SdpType_Integer_Positive {
        return fmt.Errorf("tag:%d got wrong type %d", tag, ty)
    }

    _, u32, err := up.ReadNumber32()
    if err != nil {
        return fmt.Errorf("tag:%d read u32 err:%s", tag, err.Error())
    }
    *v = u32
    return err
}

func (up *UnPacker) ReadInt32(v *int32, tag uint32, require bool) error {
    has, ty, err := up.SkipToTag(tag, require)
    if !has || err != nil {
        return err
    }
    if ty != SdpType_Integer_Positive && ty != SdpType_Integer_Negative {
        return fmt.Errorf("tag:%d got wrong type %d", tag, ty)
    }

    _, u32, err := up.ReadNumber32()
    if err != nil {
        return fmt.Errorf("tag:%d read i32 err:%s", tag, err.Error())
    }
    *v = int32(u32)
    if ty == SdpType_Integer_Negative {
        *v = -int32(u32)
    }
    return err
}

func (up *UnPacker) ReadUint64(v *uint64, tag uint32, require bool) error {
    has, ty, err := up.SkipToTag(tag, require)
    if !has || err != nil {
        return err
    }
    if ty != SdpType_Integer_Positive {
        return fmt.Errorf("tag:%d got wrong type %d", tag, ty)
    }

    _, u64, err := up.ReadNumber64()
    if err != nil {
        return fmt.Errorf("tag:%d read u64 err:%s", tag, err.Error())
    }
    *v = u64
    return err
}

func (up *UnPacker) ReadInt64(v *int64, tag uint32, require bool) error {
    has, ty, err := up.SkipToTag(tag, require)
    if !has || err != nil {
        return err
    }
    if ty != SdpType_Integer_Positive && ty != SdpType_Integer_Negative {
        return fmt.Errorf("tag:%d got wrong type %d", tag, ty)
    }

    _, u64, err := up.ReadNumber64()
    if err != nil {
        return fmt.Errorf("tag:%d read i64 err:%s", tag, err.Error())
    }
    *v = int64(u64)
    if ty == SdpType_Integer_Negative {
        *v = -int64(u64)
    }
    return err
}

func (up *UnPacker) ReadFloat32(v *float32, tag uint32, require bool) error {
    has, ty, err := up.SkipToTag(tag, require)
    if !has || err != nil {
        return err
    }
    if ty != SdpType_Float {
        return fmt.Errorf("tag:%d got wrong type %d", tag, ty)
    }

    _, u32, err := up.ReadNumber32()
    if err != nil {
        return fmt.Errorf("tag:%d read float32 err:%s", tag, err.Error())
    }
    *v = math.Float32frombits(u32)
    return err
}

func (up *UnPacker) ReadFloat64(v *float64, tag uint32, require bool) error {
    has, ty, err := up.SkipToTag(tag, require)
    if !has || err != nil {
        return err
    }
    if ty != SdpType_Double {
        return fmt.Errorf("tag:%d got wrong type %d", tag, ty)
    }

    _, u64, err := up.ReadNumber64()
    if err != nil {
        return fmt.Errorf("tag:%d read float64 err:%s", tag, err.Error())
    }
    *v = math.Float64frombits(u64)
    return err
}

func (up *UnPacker) ReadString(v *string, tag uint32, require bool) error {
    has, ty, err := up.SkipToTag(tag, require)
    if !has || err != nil {
        return err
    }
    if ty != SdpType_String {
        return fmt.Errorf("tag:%d got wrong type %d", tag, ty)
    }

    _, len, err := up.ReadNumber32()
    if err != nil {
        return err
    }
    if len == 0 {
        return nil
    }
    if up.buf.Len() < int(len) {
        return fmt.Errorf("tag:%d end of data", tag)
    }
    var bs = make([]byte, len, len)
    if _, err := up.buf.Read(bs); err != nil {
        return fmt.Errorf("tag:%d read string err:%s", tag, err.Error())
    }
    *v = string(bs)
    return err
}

func (up *UnPacker) Reset(buf []byte) {
    up.buf.Reset(buf)
}

func NewPacker() *Packer {
    return &Packer{buf : &bytes.Buffer{}}
}

func NewUnPacker(buf []byte) *UnPacker {
    return &UnPacker{buf: bytes.NewReader(buf)}
}
