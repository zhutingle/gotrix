package test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/zhutingle/gotrix/global"
)

func keepToolsTest() {
	fmt.Println("")
}

func TestToFloat64(t *testing.T) {
	zero := float64(0)
	one := float64(1)

	array := []interface{}{int(1), int8(1), int16(1), int32(1), int64(1), uint(1), uint8(1), uint16(1), uint32(1), uint64(1), float32(1), float64(1), "1"}
	for _, v := range array {
		value, e := global.ToFloat64(v)
		if value == one && e == nil {
			t.Log("正确")
		} else {
			t.Fatalf("类型为[%v]的值[%v]在调用 ToFloat64 时出错。", reflect.TypeOf(v), v)
		}
	}

	array = []interface{}{make(map[string]interface{}), make([]interface{}, 0), "abc"}
	for _, v := range array {
		value, e := global.ToFloat64(v)
		if value == zero && e != nil {
			t.Log("正确")
		} else {
			t.Fatalf("类型为[%v]的值[%v]在调用 ToFloat64 时出错。", reflect.TypeOf(v), v)
		}
	}
}

func TestToFloat64Must(t *testing.T) {
	zero := float64(0)
	one := float64(1)

	array := []interface{}{int(1), int8(1), int16(1), int32(1), int64(1), uint(1), uint8(1), uint16(1), uint32(1), uint64(1), float32(1), float64(1), "1"}
	for _, v := range array {
		value := global.ToFloat64Must(v)
		if value == one {
			t.Log("正确")
		} else {
			t.Fatalf("类型为[%v]的值[%v]在调用 ToFloat64Must 时出错。", reflect.TypeOf(v), v)
		}
	}

	array = []interface{}{make(map[string]interface{}), make([]interface{}, 0), "abc"}
	for _, v := range array {
		value := global.ToFloat64Must(v)
		if value == zero {
			t.Log("正确")
		} else {
			t.Fatalf("类型为[%v]的值[%v]在调用 ToFloat64Must 时出错。", reflect.TypeOf(v), v)
		}
	}
}

func TestToInt64(t *testing.T) {
	zero := int64(0)
	one := int64(1)

	array := []interface{}{int(1), int8(1), int16(1), int32(1), int64(1), uint(1), uint8(1), uint16(1), uint32(1), uint64(1), float32(1), float64(1), "1"}
	for _, v := range array {
		value, e := global.ToInt64(v)
		if value == one && e == nil {
			t.Log("正确")
		} else {
			t.Fatalf("类型为[%v]的值[%v]在调用 ToFloat64 时出错。", reflect.TypeOf(v), v)
		}
	}

	array = []interface{}{make(map[string]interface{}), make([]interface{}, 0), "abc"}
	for _, v := range array {
		value, e := global.ToInt64(v)
		if value == zero && e != nil {
			t.Log("正确")
		} else {
			t.Fatalf("类型为[%v]的值[%v]在调用 ToFloat64 时出错。", reflect.TypeOf(v), v)
		}
	}
}

func TestToInt64Must(t *testing.T) {
	zero := int64(0)
	one := int64(1)

	array := []interface{}{int(1), int8(1), int16(1), int32(1), int64(1), uint(1), uint8(1), uint16(1), uint32(1), uint64(1), float32(1), float64(1), "1"}
	for _, v := range array {
		value := global.ToInt64Must(v)
		if value == one {
			t.Log("正确")
		} else {
			t.Fatalf("类型为[%v]的值[%v]在调用 ToFloat64Must 时出错。", reflect.TypeOf(v), v)
		}
	}

	array = []interface{}{make(map[string]interface{}), make([]interface{}, 0), "abc"}
	for _, v := range array {
		value := global.ToInt64Must(v)
		if value == zero {
			t.Log("正确")
		} else {
			t.Fatalf("类型为[%v]的值[%v]在调用 ToFloat64Must 时出错。", reflect.TypeOf(v), v)
		}
	}
}
