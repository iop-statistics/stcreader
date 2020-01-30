package stc

import (
	"errors"
	"reflect"
)

type fixedHeader struct {
	Code      uint16 // stc 文件的编号
	Length    uint16 // stc 文件长度-4
	DataCount uint16 // 数据数量
	TypeCount uint8  // 数据类型数量
}

type Header struct {
	fixedHeader
	Types []uint8 // 数据类型
	Index []IndexEntry
}

type IndexEntry struct {
	ID   uint32
	Addr uint32
}

type Row struct {
	Data []interface{}
}

const (
	DataTypeSByte = iota + 1
	DataTypeByte
	DataTypeShort
	DataTypeUShort
	DataTypeInt
	DataTypeUInt
	DataTypeLong
	DataTypeULong
	DataTypeFloat
	DataTypeDouble
	DataTypeString
)

func (r Row) Unmarshal(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("unmarshal: v should be a ptr and not nil")
	}

	pv := rv.Elem()
	if pv.Kind() != reflect.Struct {
		return errors.New("unmarshal: element of a should be a struct")
	}

	if pv.NumField() != len(r.Data) {
		return errors.New("unmarshal: struct's field count is not equal to data's count")
	}

	for i := 0; i < len(r.Data); i++ {
		fv := pv.Field(i)
		dv := reflect.ValueOf(r.Data[i])

		if fv.Kind() != dv.Kind() {
			return errors.New("unmarshal: field type not match")
		}

		fv.Set(dv)
	}

	return nil
}
