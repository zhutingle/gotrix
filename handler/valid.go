package handler

import (
	"github.com/zhutingle/gotrix/global"
)

type Valid interface {
	Valid(param *Param, value interface{}) *global.GotrixError
}

type StringValid struct{}
type IntValid struct{}
type BoolValid struct{}
type ArrayValid struct{}

func (this StringValid) Valid(param *Param, value interface{}) *global.GotrixError {
	switch value.(type) {
	case string:
		var val string = value.(string)
		var length = len(val)
		if param.Must && length == 0 {
			return global.NewGotrixError(global.PARAM_MUST_EXISTS, param.Name)
		}
		if len(param.Len) > 0 && (length < param.min || length > param.max) {
			return global.NewGotrixError(global.PARAM_LENGTH_ERROR, param.Name, param.min, param.max)
		}
		break
	default:
		if value != nil {
			return global.NewGotrixError(global.STRING_PARAM_ERROR, param.Name)
		}
		if param.Must {
			return global.NewGotrixError(global.PARAM_MUST_EXISTS, param.Name)
		}
	}
	return nil
}

func (this IntValid) Valid(param *Param, value interface{}) *global.GotrixError {
	f, e := global.ToFloat64(value)
	if e == nil {
		var integerValue int = int(f)
		if len(param.Len) > 0 && (integerValue < param.min || integerValue > param.max) {
			return global.NewGotrixError(global.PARAM_VALUE_ERROR, param.Name, param.min, param.max)
		}
	} else {
		if value != nil {
			return global.NewGotrixError(global.INT_PARAM_ERROR, param.Name)
		}
		if param.Must {
			return global.NewGotrixError(global.PARAM_MUST_EXISTS, param.Name)
		}
	}
	return nil
}

func (this BoolValid) Valid(param *Param, value interface{}) *global.GotrixError {
	switch value.(type) {
	case bool:
		break
	default:
		if value != nil {
			return global.NewGotrixError(global.BOOL_PARAM_ERROR, param.Name)
		}
		if param.Must {
			return global.NewGotrixError(global.PARAM_MUST_EXISTS, param.Name)
		}
	}
	return nil
}

func (this ArrayValid) Valid(param *Param, value interface{}) *global.GotrixError {
	switch value.(type) {
	case []interface{}:
		break
	default:
		if value != nil {
			return global.NewGotrixError(global.ARRAY_PARAM_ERROR, param.Name)
		}
		if param.Must {
			return global.NewGotrixError(global.PARAM_MUST_EXISTS, param.Name)
		}
	}
	return nil
}
