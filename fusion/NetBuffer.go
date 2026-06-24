package fusion

import (
	"admin/fusion/base"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"reflect"
)

var defByteOrder = binary.LittleEndian

type NetBuffer struct {
	Buffer bytes.Buffer
}

func (obj *NetBuffer) IsReadableEmpty() bool {
	return obj.Buffer.Len() <= 0
}

func (obj *NetBuffer) GetReadableSize() int {
	return obj.Buffer.Len()
}

func (obj *NetBuffer) GetReadableBytes() []byte {
	return obj.Buffer.Bytes()
}

func (obj *NetBuffer) Truncate(n int) {
	obj.Buffer.Truncate(n)
}

func (obj *NetBuffer) Grow(n int) {
	obj.Buffer.Grow(n)
}

func (obj *NetBuffer) Next(n int) []byte {
	return obj.Buffer.Next(n)
}

func (obj *NetBuffer) GetReader() io.Reader {
	return &obj.Buffer
}

func (obj *NetBuffer) GetWriter() io.Writer {
	return &obj.Buffer
}

func (obj *NetBuffer) Write(args ...interface{}) {
	for _, v := range args {
		switch v := v.(type) {
		case bool:
			obj.WriteBool(v)
		case float32:
			obj.WriteFloat32(v)
		case float64:
			obj.WriteFloat64(v)
		case int8:
			obj.WriteInt8(v)
		case uint8:
			obj.WriteUInt8(v)
		case int16:
			obj.WriteInt16(v)
		case uint16:
			obj.WriteUInt16(v)
		case int32:
			obj.WriteInt32(v)
		case uint32:
			obj.WriteUInt32(v)
		case int64:
			obj.WriteInt64(v)
		case uint64:
			obj.WriteUInt64(v)
		case string:
			obj.WriteString(v)
		case []byte:
			obj.WriteBytes(v)
		default:
			panic(errors.New("invalid type " + reflect.TypeOf(v).String()))
		}
	}
}

func (obj *NetBuffer) WriteBool(v bool) {
	Must(binary.Write(&obj.Buffer, defByteOrder, v))
}

func (obj *NetBuffer) WriteFloat32(v float32) {
	Must(binary.Write(&obj.Buffer, defByteOrder, v))
}

func (obj *NetBuffer) WriteFloat64(v float64) {
	Must(binary.Write(&obj.Buffer, defByteOrder, v))
}

func (obj *NetBuffer) WriteInt8(v int8) {
	Must(binary.Write(&obj.Buffer, defByteOrder, v))
}

func (obj *NetBuffer) WriteUInt8(v uint8) {
	Must(binary.Write(&obj.Buffer, defByteOrder, v))
}

func (obj *NetBuffer) WriteInt16(v int16) {
	Must(binary.Write(&obj.Buffer, defByteOrder, v))
}

func (obj *NetBuffer) WriteUInt16(v uint16) {
	Must(binary.Write(&obj.Buffer, defByteOrder, v))
}

func (obj *NetBuffer) WriteInt32(v int32) {
	Must(binary.Write(&obj.Buffer, defByteOrder, v))
}

func (obj *NetBuffer) WriteUInt32(v uint32) {
	Must(binary.Write(&obj.Buffer, defByteOrder, v))
}

func (obj *NetBuffer) WriteInt64(v int64) {
	Must(binary.Write(&obj.Buffer, defByteOrder, v))
}

func (obj *NetBuffer) WriteUInt64(v uint64) {
	Must(binary.Write(&obj.Buffer, defByteOrder, v))
}

func (obj *NetBuffer) WriteString(v string) {
	obj.WriteUInt16(uint16(len(v)))
	Must(obj.Buffer.WriteString(v))
}

func (obj *NetBuffer) WriteBytes(v []byte) {
	obj.WriteUInt16(uint16(len(v)))
	Must(obj.Buffer.Write(v))
}

func (obj *NetBuffer) AppendString(v string) {
	Must(obj.Buffer.WriteString(v))
}

func (obj *NetBuffer) AppendBytes(v []byte) {
	Must(obj.Buffer.Write(v))
}

func (obj *NetBuffer) Read(args ...interface{}) {
	for _, v := range args {
		switch v := v.(type) {
		case *bool:
			obj.ReadBool(v)
		case *float32:
			obj.ReadFloat32(v)
		case *float64:
			obj.ReadFloat64(v)
		case *int8:
			obj.ReadInt8(v)
		case *uint8:
			obj.ReadUInt8(v)
		case *int16:
			obj.ReadInt16(v)
		case *uint16:
			obj.ReadUInt16(v)
		case *int32:
			obj.ReadInt32(v)
		case *uint32:
			obj.ReadUInt32(v)
		case *int64:
			obj.ReadInt64(v)
		case *uint64:
			obj.ReadUInt64(v)
		case *string:
			obj.ReadString(v)
		case *[]byte:
			obj.ReadBytes(v)
		default:
			panic(errors.New("invalid type " + reflect.TypeOf(v).String()))
		}
	}
}

func (obj *NetBuffer) ReadBool(v *bool) {
	Must(binary.Read(&obj.Buffer, defByteOrder, v))
}

func (obj *NetBuffer) ReadFloat32(v *float32) {
	Must(binary.Read(&obj.Buffer, defByteOrder, v))
}

func (obj *NetBuffer) ReadFloat64(v *float64) {
	Must(binary.Read(&obj.Buffer, defByteOrder, v))
}

func (obj *NetBuffer) ReadInt8(v *int8) {
	Must(binary.Read(&obj.Buffer, defByteOrder, v))
}

func (obj *NetBuffer) ReadUInt8(v *uint8) {
	Must(binary.Read(&obj.Buffer, defByteOrder, v))
}

func (obj *NetBuffer) ReadInt16(v *int16) {
	Must(binary.Read(&obj.Buffer, defByteOrder, v))
}

func (obj *NetBuffer) ReadUInt16(v *uint16) {
	Must(binary.Read(&obj.Buffer, defByteOrder, v))
}

func (obj *NetBuffer) ReadInt32(v *int32) {
	Must(binary.Read(&obj.Buffer, defByteOrder, v))
}

func (obj *NetBuffer) ReadUInt32(v *uint32) {
	Must(binary.Read(&obj.Buffer, defByteOrder, v))
}

func (obj *NetBuffer) ReadInt64(v *int64) {
	Must(binary.Read(&obj.Buffer, defByteOrder, v))
}

func (obj *NetBuffer) ReadUInt64(v *uint64) {
	Must(binary.Read(&obj.Buffer, defByteOrder, v))
}

func (obj *NetBuffer) ReadString(v *string) {
	var n uint16
	obj.ReadUInt16(&n)
	if obj.Buffer.Len() >= int(n) {
		*v = string(obj.Buffer.Next(int(n)))
	} else {
		panic(io.ErrUnexpectedEOF)
	}
}

func (obj *NetBuffer) ReadBytes(v *[]byte) {
	var n uint16
	obj.ReadUInt16(&n)
	if obj.Buffer.Len() >= int(n) {
		*v = base.CopyBytes(obj.Buffer.Next(int(n)))
	} else {
		panic(io.ErrUnexpectedEOF)
	}
}
