// parser is for struct parse, it is the main parse functions package
package parser

import (
	"encoding/binary"
	"errors"
	"os"
	"reflect"
)

// endians
const (
	LITTLE_ENDIAN = 0
	BIG_ENDIAN    = 1
)

// This is a wrapper function for simpler usage
func Parse(
	stream *os.File,
	typ reflect.Type,
	end uint,
) {
	ParseStructFromBinaryStream(stream, typ, end)
}

// main parse function dealing with file stream to parse arbitrary
// struct from that stream
// notice that this function can only parse linearly, high level
// should be implemented on your own since it knows nothing about
// high level knowledge
// params:
//      stream: file stream to read
//      typ: struct type, which is reflect.Type,
//           you can pass this like reflect.TypeOf(SampleStruct{})
//      end: endian, which is one of parser.LITTLE_ENDIAN or parser.BIG_ENDIAN
func ParseStructFromBinaryStream(
	stream *os.File,
	typ reflect.Type,
	end uint,
) (interface{}, error) {
	newStruct := reflect.Indirect(reflect.New(typ))
	for i := 0; i < newStruct.NumField(); i++ {
		field := newStruct.Field(i)
		parsed, ok := doParseValue(stream, field.Type(), field, end)
		if ok != nil {
			return nil, ok
		}
		field.Set(parsed)
	}
	return newStruct, nil
}

// parse an array value from a file stream, which is originally internal needed for
// ParseStructFromBinaryStream, but can be used alone as well
// params:
//      stream: file stream to read
//      typ: array type, which is reflect.Type
//      end: endian, which one of parser.LITTLE_ENDIAN or parser.BIG_ENDIAN
func ParseArrayFromBinaryStream(
	stream *os.File,
	typ reflect.Type,
	arr reflect.Value,
	end uint,
) (reflect.Value, error) {
	arrayLen := typ.Len()
	elemType := typ.Elem()
	for i := 0; i < arrayLen; i++ {
		parsed, ok := doParseValue(stream, elemType, arr, end)
		if ok != nil {
			return reflect.ValueOf(0), ok
		}
		arr.Index(i).Set(parsed)
	}
	return arr, nil
}

// parse a single value with type
func doParseValue(
	stream *os.File,
	typ reflect.Type,
	val reflect.Value,
	end uint,
) (reflect.Value, error) {
	kind := typ.Kind()
	switch {
	case kind == reflect.Invalid:
		return reflect.ValueOf(0), errors.New("struct field is not valid")
	case kind >= reflect.Chan && kind != reflect.Struct:
		return reflect.ValueOf(0), errors.New("struct field is variadic, unable to parse")
	case kind == reflect.Struct:
		newSubStruct, ok := ParseStructFromBinaryStream(stream, typ, end)
		if ok != nil {
			return reflect.ValueOf(0), ok
		}
		newParsed, assertRes := newSubStruct.(reflect.Value)
		if assertRes == false {
			return reflect.ValueOf(0), errors.New("type assertion of Value failed")
		}
		return newParsed, nil
	case kind == reflect.Array:
		newArray, ok := ParseArrayFromBinaryStream(stream, typ, val, end)
		if ok != nil {
			return reflect.ValueOf(0), ok
		}
		return newArray, nil
	default:
		// invariadic values only
		numberSize := typ.Size()
		buf := make([]byte, numberSize)
		stream.Read(buf)
		// WARNING! potential integer overflow! But I think it is ok here..
		// since buf should always be greater or equal than zero
		if numberSize > uintptr(len(buf)) {
			return reflect.ValueOf(0), errors.New("file content not long enough")
		}
		parsed, ok := convToInvariadicValue(buf, kind, end)
		if ok != nil {
			return reflect.ValueOf(0), ok
		}
		return parsed, nil
	}
}

func convToInvariadicValue(
	buf []byte,
	kind reflect.Kind,
	end uint,
) (reflect.Value, error) {
	// convert a byte slice to invariadic value
	// which means bool, int, int8, uint8 etc. that have no
	// variadic things inside.

	// this is assumed to be true out of this function

	// I do this because I have no way to know about exact type dynamically
	// to do coercion, but fortunately I only have to deal with types not that much
	switch kind {
	case reflect.Bool:
		// I don't know anything to convert a bool, just write one easily
		if buf[0] != 0 {
			return reflect.ValueOf(true), nil
		} else {
			return reflect.ValueOf(false), nil
		}
	case reflect.Int:
		if end == LITTLE_ENDIAN {
			return reflect.ValueOf(binary.LittleEndian.Uint64(buf)), nil
		} else {
			return reflect.ValueOf(binary.BigEndian.Uint64(buf)), nil
		}
	case reflect.Int8:
		return reflect.ValueOf(int8(buf[0])), nil
	case reflect.Int16:
		if end == LITTLE_ENDIAN {
			return reflect.ValueOf(int16(binary.LittleEndian.Uint16(buf))), nil
		} else {
			return reflect.ValueOf(int16(binary.BigEndian.Uint16(buf))), nil
		}
	case reflect.Int32:
		if end == LITTLE_ENDIAN {
			return reflect.ValueOf(int32(binary.LittleEndian.Uint32(buf))), nil
		} else {
			return reflect.ValueOf(int32(binary.BigEndian.Uint32(buf))), nil
		}
	case reflect.Int64:
		if end == LITTLE_ENDIAN {
			return reflect.ValueOf(int64(binary.LittleEndian.Uint64(buf))), nil
		} else {
			return reflect.ValueOf(int64(binary.BigEndian.Uint64(buf))), nil
		}
	case reflect.Uint:
		if end == LITTLE_ENDIAN {
			return reflect.ValueOf(uint(binary.LittleEndian.Uint64(buf))), nil
		} else {
			return reflect.ValueOf(uint(binary.BigEndian.Uint64(buf))), nil
		}
	case reflect.Uint8:
		return reflect.ValueOf(uint8(buf[0])), nil
	case reflect.Uint16:
		if end == LITTLE_ENDIAN {
			return reflect.ValueOf(binary.LittleEndian.Uint16(buf)), nil
		} else {
			return reflect.ValueOf(binary.BigEndian.Uint16(buf)), nil
		}
	case reflect.Uint32:
		if end == LITTLE_ENDIAN {
			return reflect.ValueOf(binary.LittleEndian.Uint32(buf)), nil
		} else {
			return reflect.ValueOf(binary.BigEndian.Uint32(buf)), nil
		}
	case reflect.Uint64:
		if end == LITTLE_ENDIAN {
			return reflect.ValueOf(binary.LittleEndian.Uint64(buf)), nil
		} else {
			return reflect.ValueOf(binary.BigEndian.Uint64(buf)), nil
		}
	case reflect.Uintptr:
		if end == LITTLE_ENDIAN {
			return reflect.ValueOf(uintptr(binary.LittleEndian.Uint64(buf))), nil
		} else {
			return reflect.ValueOf(uintptr(binary.BigEndian.Uint64(buf))), nil
		}
	case reflect.Float32:
		if end == LITTLE_ENDIAN {
			return reflect.ValueOf(float32(binary.LittleEndian.Uint32(buf))), nil
		} else {
			return reflect.ValueOf(float32(binary.BigEndian.Uint32(buf))), nil
		}
	case reflect.Float64:
		if end == LITTLE_ENDIAN {
			return reflect.ValueOf(float64(binary.LittleEndian.Uint64(buf))), nil
		} else {
			return reflect.ValueOf(float64(binary.BigEndian.Uint64(buf))), nil
		}
	default:
		return reflect.ValueOf(0), errors.New("type error")
	}
}
