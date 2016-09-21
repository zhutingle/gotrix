package handler

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/zhutingle/gotrix/ecdh"
	"github.com/zhutingle/gotrix/global"
	"github.com/zhutingle/gotrix/weichat"

	"github.com/tealeg/xlsx"
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
	orderId := int64(0)
	if _, b := args[0].(float64); b {
		orderId = int64(args[0].(float64))
	} else if _, b := args[0].(int64); b {
		orderId = args[0].(int64)
	}
	e := int(args[1].(int64))

	nowTime := time.Now()
	year, month, day := nowTime.Date()
	orderNo := int64(year*10000 + int(month)*100 + day)
	orderNo = orderNo*int64(math.Pow10(e)) + orderId
	response = strconv.FormatInt(orderNo, 10)

	return
}

func (this *handleFunc) WeichatSign(args []interface{}) (response interface{}, gErr *global.GotrixError) {
	mReq := args[0].(map[string]interface{})
	fmt.Println(mReq)
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
		if mReq[k] != nil {
			value := fmt.Sprintf("%v", mReq[k])
			if value != "" {
				signStrings = signStrings + k + "=" + value + "&"
			}
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

func (this *handleFunc) WeichatPay(args []interface{}) (response interface{}, gErr *global.GotrixError) {
	m := args[0].(map[string]interface{})
	orderReq := weichat.UnifyOrderReq{}
	orderReq.Appid = m["appid"].(string)
	orderReq.Body = m["body"].(string)
	orderReq.Mch_id = m["mch_id"].(string)
	orderReq.Nonce_str = m["nonce_str"].(string)
	orderReq.Notify_url = m["notify_url"].(string)
	orderReq.Trade_type = m["trade_type"].(string)
	orderReq.Spbill_create_ip = m["spbill_create_ip"].(string)
	orderReq.Total_fee = strconv.FormatInt(m["total_fee"].(int64), 10)
	orderReq.Out_trade_no = m["out_trade_no"].(string)
	orderReq.OpenId = m["openid"].(string)
	orderReq.Sign = m["sign"].(string)

	response = weichat.UnifiedOrder(orderReq)

	return
}

func (this *handleFunc) TimeStampSecond(args []interface{}) (response interface{}, gErr *global.GotrixError) {
	timeStampSecond := time.Now().Unix()
	response = strconv.FormatInt(timeStampSecond, 10)
	return
}

func (this *handleFunc) StringAdd(args []interface{}) (response interface{}, gErr *global.GotrixError) {
	var str string = ""
	for i, length := 0, len(args); i < length; i++ {
		str += fmt.Sprintf("%v", args[i])
	}
	response = str
	return
}

func (this *handleFunc) JsonToString(args []interface{}) (response interface{}, gErr *global.GotrixError) {
	bs, err := json.Marshal(args[0])
	if err != nil {
		gErr = global.JSON_TO_STRING_ERROR
		return
	}
	response = string(bs)
	return
}

func (this *handleFunc) Format(args []interface{}) (response interface{}, gErr *global.GotrixError) {
	format := args[0].(string)
	response = fmt.Sprintf(format, args[1:]...)
	return
}

func (this *handleFunc) Return(args []interface{}) (response interface{}, gErr *global.GotrixError) {
	gErr = &global.GotrixError{Status: 0, Msg: args[0].(string)}
	return
}

func (this *handleFunc) WxSendRedPack(args []interface{}) (response interface{}, gErr *global.GotrixError) {
	response = weichat.WxSendRedPack(args[0].(map[string]interface{}))
	return
}

func (this *handleFunc) TimeFormat(args []interface{}) (response interface{}, gErr *global.GotrixError) {
	layout := args[0].(string)
	response = time.Now().Format(layout)
	return
}

func (this *handleFunc) Max(args []interface{}) (response interface{}, gErr *global.GotrixError) {
	data := args[0].([]interface{})
	column := args[1].(string)

	var max int64 = math.MinInt64
	for i := 0; i < len(data); i++ {
		d := data[i].(map[string]interface{})
		v := d[column].(int64)
		if max < v {
			max = v
		}
	}
	response = max
	return
}

func (this *handleFunc) ToXls(args []interface{}) (response interface{}, gErr *global.GotrixError) {
	data := args[0].([]interface{})
	columnStr := args[1].(string)
	fileName := args[2].(string)

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet(fileName)
	if err != nil {
		fmt.Println(err.Error())
	}

	columns := strings.Split(columnStr, ",")
	row = sheet.AddRow()
	for i := 0; i < len(columns); i++ {
		cell = row.AddCell()
		cell.SetValue(columns[i])
	}

	for i := 0; i < len(data); i++ {
		row = sheet.AddRow()
		var d = data[i].(map[string]interface{})
		for j := 0; j < len(columns); j++ {
			cell = row.AddCell()
			cell.SetValue(d[columns[j]])
		}
	}

	err = file.Save(global.Config.TempFolder + fileName + ".xlsx")
	if err != nil {
		fmt.Println(err.Error())
	}

	return
}
