package global

import (
	"errors"
	"reflect"
	"strconv"
)

func ToFloat64(v interface{}) (float64, error) {
	var zero = float64(0)
	switch v.(type) {
	case string:
		value, e := strconv.ParseFloat(v.(string), 64)
		if e == nil {
			return value, nil
		} else {
			return zero, e
		}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return reflect.ValueOf(v).Convert(reflect.TypeOf(zero)).Float(), nil
	default:
		return zero, errors.New("This type cannot convert to float64.")
	}
}

func ToFloat64Must(v interface{}) float64 {
	var zero = float64(0)
	switch v.(type) {
	case string:
		value, e := strconv.ParseFloat(v.(string), 64)
		if e == nil {
			return value
		} else {
			return zero
		}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return reflect.ValueOf(v).Convert(reflect.TypeOf(zero)).Float()
	default:
		return zero
	}
}
