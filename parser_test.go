package parser

import (
	"os"
	"reflect"
	"testing"
)

func TestParserInt(t *testing.T) {
	type MyStruct struct {
		Num uint64
	}
	f, ok := os.Create("tempStruct")
	if ok != nil {
		t.Error("unable to create struct, error is ", ok)
	}
	buf := []byte{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}
	f.Write(buf)
	f.Close()

	f, ok = os.Open("tempStruct")
	if ok != nil {
		t.Error("unable to open created struct, error is ", ok)
	}
	parsed, ok := ParseStructFromBinaryStream(f, reflect.TypeOf(MyStruct{}), BIG_ENDIAN)
	if ok != nil {
		t.Error("parse error ", ok)
	}
	val, assertComplete := parsed.(reflect.Value)
	if !assertComplete {
		t.Error("Value type assertion error")
	}

	realVal, assertComplete := val.Interface().(MyStruct)
	if !assertComplete {
		t.Error("real struct type assertion error")
	}

	if realVal.Num != 16045690984833335023 {
		t.Error("value is wrong")
	}

	realVal.Num = 1
	if realVal.Num != 1 {
		t.Error("error")
	}

	ok = os.Remove("tempStruct")
	if ok != nil {
		t.Error("temp file delete error")
	}
}
