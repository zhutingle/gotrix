package handler

import (
	"github.com/zhutingle/gotrix/global"
)

type Valid interface {
	Valid(param *Param, value interface{}) (v interface{}, gErr *global.GotrixError)
}

type StringValid struct {
	ZERO string
}
type IntValid struct {
	ZERO int64
}
type BoolValid struct {
	ZERO bool
}
type ArrayValid struct {
	ZERO []interface{}
}
type FileValid struct {
	
}

func (this StringValid) Valid(param *Param, value interface{}) (v interface{}, gErr *global.GotrixError) {
	if value == nil {
		if param.must {
			return nil, global.NewGotrixError(global.PARAM_MUST_EXISTS, param.Name)
		} else {
			return nil, nil
		}
	}
	if val, ok := value.(string); ok {
		var length = int64(len(val))
		if param.must && length == 0 {
			return this.ZERO, global.NewGotrixError(global.PARAM_MUST_EXISTS, param.Name)
		}
		if len(param.Len) > 0 && (length < param.min || length > param.max) {
			return this.ZERO, global.NewGotrixError(global.PARAM_LENGTH_ERROR, param.Name, param.min, param.max)
		}
	} else {
		return nil, global.NewGotrixError(global.STRING_PARAM_ERROR, param.Name)
	}
	return value, nil
}

func (this IntValid) Valid(param *Param, value interface{}) (v interface{}, gErr *global.GotrixError) {
	if value == nil {
		if param.must {
			return nil, global.NewGotrixError(global.PARAM_MUST_EXISTS, param.Name)
		} else {
			return nil, nil
		}
	}
	f, e := global.ToInt64(value)
	if e != nil {
		return this.ZERO, global.NewGotrixError(global.INT_PARAM_ERROR, param.Name)
	}
	if len(param.Len) > 0 && (f < param.min || f > param.max) {
		return this.ZERO, global.NewGotrixError(global.PARAM_VALUE_ERROR, param.Name, param.min, param.max)
	}
	return f, nil
}

func (this BoolValid) Valid(param *Param, value interface{}) (v interface{}, gErr *global.GotrixError) {
	switch value.(type) {
	case bool:
		break
	default:
		if value != nil {
			return this.ZERO, global.NewGotrixError(global.BOOL_PARAM_ERROR, param.Name)
		}
		if param.must {
			return this.ZERO, global.NewGotrixError(global.PARAM_MUST_EXISTS, param.Name)
		}
	}
	return value, nil
}

func (this ArrayValid) Valid(param *Param, value interface{}) (v interface{}, gErr *global.GotrixError) {
	switch value.(type) {
	case []interface{}:
		break
	default:
		if value != nil {
			return this.ZERO, global.NewGotrixError(global.ARRAY_PARAM_ERROR, param.Name)
		}
		if param.must {
			return this.ZERO, global.NewGotrixError(global.PARAM_MUST_EXISTS, param.Name)
		}
	}
	return value, nil
}

func (this FileValid) Valid(param *Param, value interface{}) (v interface{}, gErr *global.GotrixError) {

	return value, nil
}