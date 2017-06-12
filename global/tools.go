package global

import (
	"errors"
	"reflect"
	"strconv"
)

func ToFloat64(v interface{}) (float64, error) {
	zero := float64(0)
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
		return zero, errors.New("tools:This type cannot convert to float64.")
	}
}

func ToFloat64Must(v interface{}) float64 {
	r, _ := ToFloat64(v)
	return r
}

func ToInt64(v interface{}) (int64, error) {
	zero := int64(0)
	switch v.(type) {
	case string:
		value, e := strconv.ParseInt(v.(string), 10, 64)
		if e == nil {
			return value, nil
		} else {
			return zero, e
		}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return reflect.ValueOf(v).Convert(reflect.TypeOf(zero)).Int(), nil
	default:
		return zero, errors.New("tools:This type cannot convert to int64.")
	}
}

func ToInt64Must(v interface{}) int64 {
	r, _ := ToInt64(v)
	return r
}

func ToString(v interface{}) (string, error) {
	zero := ""
	switch v.(type) {
	case string:
		return v.(string), nil

	case int, int8, int16, int32, int64:
		return strconv.FormatInt(ToInt64Must(v), 10), nil

	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatInt(ToInt64Must(v), 10), nil

	case float32, float64:
		return strconv.FormatFloat(ToFloat64Must(v), 'g', 'e', 64), nil

	default:
		return zero, errors.New("tools:This type cannot convert to string.")
	}
}

func ToStringMust(v interface{}) string {
	r, _ := ToString(v)
	return r
}