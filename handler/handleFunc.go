package handler

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/zhutingle/gotrix/ecdh"
	"github.com/zhutingle/gotrix/global"
)

type handleFunc struct {
}

func (this *handleFunc) handle(job *Job, cp *global.CheckedParams) (result interface{}, gErr *global.GotrixError) {

	funcStrs := funcReg.FindAllStringSubmatch(job.Job, -1)
	funcName := funcStrs[0][1]
	funcPara := funcStrs[0][2]

	args := make([]interface{}, 0)
	strs := argsReg.FindAllStringSubmatch(funcPara, -1)
	for j := 0; j < len(strs); j++ {
		str := strs[j][1]
		if str[0] == '$' {
			args = append(args, cp.V[str[2:len(str)-1]])
		} else if str[0] == '\'' {
			args = append(args, str[1:len(str)-1])
		} else {
			int64Value, err := strconv.ParseInt(str, 10, 64)
			if err == nil {
				args = append(args, int64Value)
				continue
			}
			float64Value, err := strconv.ParseFloat(str, 64)
			if err == nil {
				args = append(args, float64Value)
				continue
			}
			boolValue, err := strconv.ParseBool(str)
			if err == nil {
				args = append(args, boolValue)
				continue
			}
			if str == "null" {
				args = append(args, nil)
				continue
			}
		}
	}

	method := reflect.ValueOf(this).MethodByName(funcName)
	if !method.IsValid() {
		gErr = global.NewGotrixError(global.JOB_FUNC_NOT_FOUND, funcName)
		return
	}
	in := make([]reflect.Value, 1)
	in[0] = reflect.ValueOf(args)

	valus := method.Call(in)

	result = valus[0].Interface()
	if valus[1].Interface() != nil {
		gErr = valus[1].Interface().(*global.GotrixError)
	}

	return
}

func (this *handleFunc) Json(args []interface{}) (response interface{}, gErr *global.GotrixError) {
	if args[0] == nil {
		gErr = global.FUNC_JSON_NIL_ERROR
		return
	}
	var params = args[0].(map[string]interface{})
	if len(args) == 2 {
		response = params[args[1].(string)]
		return
	}
	var returnJson map[string]interface{} = make(map[string]interface{})
	for i := 1; i < len(args); i++ {
		returnJson[args[i].(string)] = params[args[i].(string)]
	}
	response = returnJson
	return
}

func (this *handleFunc) JsonSet(args []interface{}) (response interface{}, gErr *global.GotrixError) {
	if len(args) == 0 || args[0] == nil {
		gErr = global.FUNC_JSONSET_NIL_ERROR
		return
	}
	var params map[string]interface{}
	var i int
	if _, ok := args[0].(map[string]interface{}); ok {
		params = args[0].(map[string]interface{})
		i = 1
	} else {
		params = make(map[string]interface{})
		i = 0
	}
	for ; i <= len(args)-2; i = i + 2 {
		params[args[i].(string)] = args[i+1]
	}
	response = params
	return
}

func (this *handleFunc) Equal(args []interface{}) (response interface{}, gErr *global.GotrixError) {
	if args[0] == args[1] {
		if len(args) >= 3 {
			gErr = global.NewGotrixError(global.BLANK_ERROR, args[2])
		} else {
			gErr = global.INTERNAL_ERROR
		}
	}
	return
}

func (this *handleFunc) NotEqual(args []interface{}) (response interface{}, gErr *global.GotrixError) {
	if args[0] != args[1] {
		if len(args) >= 3 {
			gErr = global.NewGotrixError(global.BLANK_ERROR, args[2])
		} else {
			gErr = global.INTERNAL_ERROR
		}
	}
	return
}

func (this *handleFunc) LoginIn(args []interface{}) (response interface{}, gErr *global.GotrixError) {
	x := args[0].(string)
	y := args[1].(string)
	id := args[2].(int64)
	randKey := ecdh.Rand()
	S := ecdh.SecretKey(randKey)
	P := ecdh.PublicKey(randKey, x, y)

	var returnJson map[string]interface{} = make(map[string]interface{})
	returnJson["x"] = S.GetX().ToBigInteger().ToString(16)
	returnJson["y"] = S.GetY().ToBigInteger().ToString(16)

	session := ecdh.Rand()
	session.DMultiply(id)
	sessionHex := session.ToString(16)[:40]
	for len(sessionHex) < 40 {
		sessionHex += "0"
	}
	returnJson["session"] = sessionHex

	pass := P.GetX().ToBigInteger().Add(P.GetY().ToBigInteger())
	returnJson["pass"] = pass.ToString(16)

	response = returnJson
	return
}

func (this *handleFunc) Config(args []interface{}) (response interface{}, gErr *global.GotrixError) {
	key := args[0].(string)
	response = global.Config.M[key]
	return
}

func (this *handleFunc) RandString(args []interface{}) (response interface{}, gErr *global.GotrixError) {
	length := args[0].(int64)

	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := int64(0); i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	response = string(result)
	return
}

func (this *handleFunc) RandOrderNo(args []interface{}) (response interface{}, gErr *global.GotrixError) {
	orderId := int64(args[0].(float64))

	nowTime := time.Now()
	year, month, day := nowTime.Date()
	orderNo := int64(year*10000 + int(month)*100 + day)
	orderNo = orderNo*100000 + orderId
	response = orderNo

	return
}

func (this *handleFunc) WeiChatSign(args []interface{}) (response interface{}, gErr *global.GotrixError) {
	mReq := args[0].(map[string]interface{})
	key := args[1].(string)

	// 第一步：对 key 进行升序排序
	sorted_keys := make([]string, 0)
	for k, _ := range mReq {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)

	// 第二步：对 key = value 的键值对用 & 连接直接，略过空值
	var signStrings string
	for _, k := range sorted_keys {
		value := fmt.Sprintf("%v", mReq[k])
		if value != "" {
			signStrings = signStrings + k + "=" + value + "&"
		}
	}

	// 第三步：在键值对最后加上 key = API_KEY
	if key != "" {
		signStrings = signStrings + "key=" + key
	}

	//
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signStrings))
	cipherStr := md5Ctx.Sum(nil)
	upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))
	return upperSign, nil
}
