package stc

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
)

type Reader struct {
	Header Header
	Index     []byte  // 索引
	cnt uint16
	r *bufio.Reader
}

func NewReader(r io.Reader) (*Reader, error) {
	rdr :=  &Reader{
		r: bufio.NewReader(r),
		cnt: 0,
	}

	err := rdr.readHeader()
	return rdr, err
}

func (r *Reader) readHeader() error {
	var header fixedHeader

	err := binary.Read(r.r, binary.LittleEndian, &header)
	if err != nil {
		return err
	}

	// TODO: 验证 header 正确性

	types := make([]uint8, header.TypeCount)
	err = binary.Read(r.r, binary.LittleEndian, &types)
	if err != nil {
		return err
	}

	length := int(math.Ceil(float64(header.DataCount) / 100.0) * 8) // count / 100 * 8

	idx := make([]byte, length)
	n, err := r.r.Read(idx)
	if n != length {
		return errors.New("index length not enough")
	}

	if err != nil {
		return err
	}

	r.Index = idx
	r.Header = Header{header, types}
	return nil
}

func (r *Reader) readRow() (Row, error) {
	row := Row{}

	for _, v := range r.Header.Types {
		var a interface{}
		switch v {
		case DataTypeSByte:
			a = new(uint8)
		case DataTypeByte:
			a = new(int8)
		case DataTypeShort:
			a = new(int16)
		case DataTypeUShort:
			a = new(uint16)
		case DataTypeInt:
			a = new(int32)
		case DataTypeUInt:
			a = new(uint32)
		case DataTypeLong:
			a = new(int64)
		case DataTypeULong:
			a = new(uint64)
		case DataTypeFloat:
			a = new(float32)
		case DataTypeDouble:
			a = new(float64)
		case DataTypeString:
			str, err := r.readString()
			if err != nil {
				return row, err
			}

			a = str
		default:
			return row, errors.New(fmt.Sprintf("unknown data type %d", v))
		}

		if v != DataTypeString {
			err := binary.Read(r.r, binary.LittleEndian, a)
			if err != nil {
				return row, err
			}

			row.Data = append(row.Data, reflect.ValueOf(a).Elem().Interface())
		} else {
			row.Data = append(row.Data, a)
		}


	}

	r.cnt++
	return row, nil
}

func (r *Reader) readString() (string, error) {
	var t uint8
	var length uint16

	err := binary.Read(r.r, binary.LittleEndian, &t)
	if err != nil {
		return "", err
	}

	err = binary.Read(r.r, binary.LittleEndian, &length)
	if err != nil {
		return "", err
	}

	bytes := make([]byte, length)
	err = binary.Read(r.r, binary.LittleEndian, &bytes)
	if err != nil {
		return "", err
	}

	// 应该不用处理t
	return string(bytes), nil
}

func (r *Reader) Read(v interface{}) error {
	if !r.HasNext() {
		return io.EOF
	}

	row, err := r.readRow()
	if err != nil {
		return err
	}

	return row.Unmarshal(v)
}

func (r *Reader) ReadAll(v interface{}) error {
	if !r.HasNext() {
		return io.EOF
	}

	rv := reflect.ValueOf(v)
	if  rv.Kind() != reflect.Ptr || rv.IsNil(){
		return errors.New("v should be a pointer of slice and not nil")
	}

	pv := rv.Elem()

	slice := reflect.MakeSlice(pv.Type(), int(r.Header.DataCount - r.cnt), int(r.Header.DataCount - r.cnt))

	for r.HasNext() {
		sv := slice.Index(int(r.cnt))
		r, err := r.readRow()
		if err != nil {
			return err
		}

		err = r.Unmarshal(sv.Addr().Interface())
		if err != nil {
			return err
		}
	}

	rv.Elem().Set(slice)
	return nil
}

func (r *Reader) ReadRaw() (Row, error) {
	return r.readRow()
}

func (r *Reader) ReadAllRaw() ([]Row, error) {
	var ret []Row
	for r.HasNext() {
		row, err := r.readRow()
		if err != nil {
			return ret, err
		}

		ret = append(ret, row)
	}

	return ret, nil
}

func (r *Reader) HasNext() bool {
	return r.cnt < r.Header.DataCount
}