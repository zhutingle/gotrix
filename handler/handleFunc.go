package handler

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/smtp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/zhutingle/gotrix/ecdh"
	"github.com/zhutingle/gotrix/global"
	"github.com/zhutingle/gotrix/weichat"

	"github.com/scorredoira/email"
	"github.com/tealeg/xlsx"
	"log"
	"regexp"
)

type handleFunc struct {
	simpleHandler SimpleHandler
	methodMap     map[string]func(args []interface{}) (response interface{}, gErr *global.GotrixError)
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

	if this.methodMap[funcName] == nil {
		gErr = global.NewGotrixError(global.JOB_FUNC_NOT_FOUND, funcName)
		return
	}
	result, gErr = this.methodMap[funcName](args)

	return
}

func (this *handleFunc) init() *handleFunc {

	this.methodMap = make(map[string]func(args []interface{}) (response interface{}, gErr *global.GotrixError))

	this.initJson()
	this.initTime()
	this.initJudge()
	this.initRand()
	this.initSpecial()
	this.initXlsx()
	this.initEmail()
	this.initHttp()
	this.initCall()
	this.initDebug()
	this.initGotrix()

	return this
}

/**
 * 定义所有与JSON操作相关的函数
 */
func (this *handleFunc) initJson() {
	// ToJson(string)
	// ToJson([]byte)
	// 将字符串转换为json格式
	this.methodMap["ToJson"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
		var param interface{} = args[0]
		var e error

		if _, ok := param.(string); ok {
			e = json.Unmarshal([]byte(param.(string)), &response)
		} else if _, ok := param.([]byte); ok {
			e = json.Unmarshal(param.([]byte), &response)
		}
		if e != nil {
			gErr = global.STRING_TO_JSON_ERROR
		}
		return
	}
	// ToString(map[string]interface{})
	// ToString([]interface{})
	// 将json格式转为字符串
	this.methodMap["ToString"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
		var jsonObject = args[0]

		bs, e := json.Marshal(jsonObject)
		if e != nil {
			gErr = global.JSON_TO_STRING_ERROR
		}

		response = string(bs)

		return
	}
	// Jget(map[string]interface{},string...)
	// 两个参数：           取JSON中的某个键的值，返回该值
	// 三个或以上参数：取JSON中的某些键的值，返回这些键和值组成的一个新的JSON
	this.methodMap["Jget"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
		if args[0] == nil {
			return nil, global.FUNC_JGET_NIL_ERROR
		}
		var params = args[0].(map[string]interface{})
		if len(args) == 2 {
			return params[args[1].(string)], nil
		}
		var returnJson map[string]interface{} = make(map[string]interface{})
		for i := 1; i < len(args); i++ {
			returnJson[args[i].(string)] = params[args[i].(string)]
		}
		return returnJson, nil
	}
	// Jset(map[string]interface{},(string,interface{})...)
	// 设置 JSON 中的某些键值对，并返回该 JSON
	// Jset((string,interface{})...)
	// 新建一个 JSON，并设置其中的健值对，并返回该 JSON
	this.methodMap["Jset"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
		if len(args) == 0 || args[0] == nil {
			return nil, global.FUNC_JSET_NIL_ERROR
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
		return params, nil
	}
	// Sprintf(string,interface{}...)
	// 调用fmt.Sprintf(string,interface{}...)格式化字符串
	this.methodMap["Sprintf"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
		format := args[0].(string)
		response = fmt.Sprintf(format, args[1:]...)
		return
	}
	// JAget([]interface{},int64)
	// 获取JSONArray中的某一个值，int64可以为负数，表示从后往前，JSONArray为nil时返回nil；超出JSONArray的范围时也返回  nil
	this.methodMap["JAget"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
		if args[0] == nil {
			return nil, nil
		}
		ja := args[0].([]interface{})
		index := args[1].(int64)
		if index < 0 {
			index = int64(len(ja)) + index
		}
		if index < 0 || index >= int64(len(ja)) {
			return nil, nil
		}

		return ja[index], nil
	}
}

/**
 * 定义所有与时间相关的函数
 */
func (this *handleFunc) initTime() {
	// Tsecond()
	// 返回UTC时间秒数的字符串
	this.methodMap["Tsecond"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
		timeStampSecond := time.Now().Unix()
		response = strconv.FormatInt(timeStampSecond, 10)
		return
	}
	// Tformat(string)
	// 按照Golang标准格式进行当前日期格式化
	this.methodMap["Tformat"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
		layout := args[0].(string)
		response = time.Now().Format(layout)
		return
	}
}

/**
 * 定义所有与判断相关的函数
 */
func (this *handleFunc) initJudge() {
	// Eq(interface{},interface{},string)
	// 第一个参数等于第二个参数时抛出异常
	// 第三个参数不为空时抛出第三个参数所示文字的异常，为空时抛出内部异常
	this.methodMap["Eq"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
		if args[0] == args[1] {
			if len(args) >= 3 {
				gErr = global.NewGotrixError(global.BLANK_ERROR, args[2])
			} else {
				gErr = global.INTERNAL_ERROR
			}
		}
		return
	}
	// Eq(interface{},interface{},string)
	// 第一个参数不等于第二个参数时抛出异常
	// 第三个参数不为空时抛出第三个参数所示文字的异常，为空时抛出内部异常
	this.methodMap["Neq"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
		if args[0] != args[1] {
			if len(args) >= 3 {
				gErr = global.NewGotrixError(global.BLANK_ERROR, args[2])
			} else {
				gErr = global.INTERNAL_ERROR
			}
		}
		return
	}
	this.methodMap["G"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
		if global.ToFloat64Must(args[0]) > global.ToFloat64Must(args[1]) {
			if len(args) >= 3 {
				gErr = global.NewGotrixError(global.BLANK_ERROR, args[2])
			} else {
				gErr = global.INTERNAL_ERROR
			}
		}
		return
	}
	this.methodMap["L"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
		if global.ToFloat64Must(args[0]) < global.ToFloat64Must(args[1]) {
			if len(args) >= 3 {
				gErr = global.NewGotrixError(global.BLANK_ERROR, args[2])
			} else {
				gErr = global.INTERNAL_ERROR
			}
		}
		return
	}
}

/**
 * 定义所有与随机相关的函数
 */
func (this *handleFunc) initRand() {
	// Rstring(int64)
	// 随机生成一个第一个参数所示长度的字符串
	this.methodMap["Rstring"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
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
}

/**
 * 定义所有特殊函数
 */
func (this *handleFunc) initSpecial() {
	this.methodMap["Return"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
		gErr = &global.GotrixError{Status: 0, Msg: args[0].(string)}
		return
	}
	this.methodMap["Config"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
		key := args[0].(string)
		response = global.Config.M[key]
		return
	}
	this.methodMap["LoginIn"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
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
	this.methodMap["WeichatSign"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
		mReq := args[0].(map[string]interface{})
		key := args[1].(string)

		// 第一步：对 key 进行升序排序
		sorted_keys := make([]string, 0)
		for k, _ := range mReq {
			if k != "sign" {
				sorted_keys = append(sorted_keys, k)
			}
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
	this.methodMap["WeichatPay"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
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
	this.methodMap["WxSendRedPack"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
		response = weichat.WxSendRedPack(args[0].(map[string]interface{}))
		return
	}
	this.methodMap["WxRequest"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
		response = weichat.WxRequest(args[0].(map[string]interface{}), args[1].(string))
		return
	}
	this.methodMap["WxCertRequest"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
		response = weichat.WxCertRequest(args[0].(map[string]interface{}), args[1].(string))
		return
	}
}

/**
 * 定义所有xlsx操作相关的函数
 */
func (this *handleFunc) initXlsx() {
	this.methodMap["ToXls"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
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

		filePath := global.Config.Dir.Temp + fileName + ".xlsx"
		err = file.Save(filePath)
		if err != nil {
			fmt.Println(err.Error())
		}

		return filePath, nil
	}
}

/**
 * 定义所有email操作相关的函数
 */
func (this *handleFunc) initEmail() {
	this.methodMap["SendEmail"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
		subject := args[0].(string)
		body := args[1].(string)
		receiver := args[2].(string)
		var err error

		emailConfig := global.Config.Email
		m := email.NewMessage(subject, body)
		m.From = emailConfig.Address
		m.To = strings.Split(receiver, ",")

		if len(args) == 4 {
			attach := args[3].(string)
			err = m.Attach(attach)
			if err != nil {
				fmt.Println(err)
			}
		}

		err = email.Send(emailConfig.SmtpUrl, smtp.PlainAuth(emailConfig.Identify, emailConfig.Username, emailConfig.Password, emailConfig.Host), m)
		if err != nil {
			fmt.Println(err)
		}
		return
	}
}

/**
 * 定义所有 Http 相关的扩展方法
 */
func (this *handleFunc) initHttp() {
	this.methodMap["HttpGet"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
		url := args[0].(string)

		resp, e := http.Get(url)
		if e != nil {
			fmt.Println(e)
			gErr = global.HTTPHANDLE_HTTP_GET_ERROR
			return
		}

		defer resp.Body.Close()
		body, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			fmt.Println(e)
			gErr = global.HTTPHANDLE_HTTP_READ_BODY
			return
		}

		return body, nil
	}

	this.methodMap["HttpPost"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
		url := args[0].(string)
		data := args[1]

		bs, e := json.Marshal(data)
		if e != nil {
			fmt.Println(e)
			gErr = global.HTTPHANDLE_HTTP_GET_ERROR
			return
		}

		reqBody := bytes.NewBuffer(bs)

		resp, e := http.Post(url, "application/json;charset=utf-8", reqBody)
		if e != nil {
			fmt.Println(e)
			gErr = global.HTTPHANDLE_HTTP_GET_ERROR
			return
		}

		defer resp.Body.Close()
		body, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			fmt.Println(e)
			gErr = global.HTTPHANDLE_HTTP_READ_BODY
			return
		}

		return body, nil
	}
}

/**
 * 定义所有调用其它相关的扩展方法
 */
func (this *handleFunc) initCall() {
	// 同步调用
	this.methodMap["SyncCall"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {

		var funcId = args[0].(int64)
		var param = args[1].(map[string]interface{})

		checkedParams := &global.CheckedParams{Func: int(funcId), V: param}
		return this.simpleHandler.Handle(checkedParams)
	}
	// 数组循环调用
	this.methodMap["JAeachCall"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
		if args[0] == nil {
			return nil, nil
		}
		ja := args[0].([]interface{})
		funcId, _ := global.ToInt64(args[1])

		result := make([]interface{}, 0)
		for _, v := range ja {
			checkedParams := &global.CheckedParams{Func: int(funcId), V: v.(map[string]interface{})}
			res, ge := this.simpleHandler.Handle(checkedParams)
			if ge == nil {
				result = append(result, res)
			} else {
				result = append(result, ge)
			}
		}
		return result, nil
	}
}

/**
 * 定义调试相关的扩展方法
 */
func (this *handleFunc) initDebug() {
	// 在控制台输出内容
	this.methodMap["Println"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
		log.Println(args...)
		return nil, nil
	}
}

/**
 * 定义操作框架本身的扩展方法
 */
func (this *handleFunc) initGotrix() {
	// 获取所有有权限的页面
	this.methodMap["GetPermissionedPages"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
		funcsStr := args[0].(string)
		funcsArray := strings.Split(funcsStr, ",")

		resultArray := make([]interface{}, 0)
		for i := 0; i < len(funcsArray); i++ {
			funcId, err := strconv.ParseInt(funcsArray[i], 10, 32)
			if err == nil && pageMap[int(funcId)] != nil {
				page := pageMap[int(funcId)]
				pageSimple := make(map[string]interface{})
				pageSimple["Id"] = page.Id
				pageSimple["Des"] = page.Des
				pageSimple["Parent"] = page.Parent
				resultArray = append(resultArray, pageSimple)
			}
		}

		return resultArray, nil
	}
	// 获取用户有权限的单个页面的详细信息
	this.methodMap["GetPermissionedPage"] = func(args []interface{}) (response interface{}, gErr *global.GotrixError) {
		funcsStr := args[0].(string)
		pageId := int(args[1].(int64))
		funcReg, err := regexp.Compile(fmt.Sprintf("(^%d,)|(,%d,)|(,%d$)", pageId, pageId, pageId))
		if err != nil {
			gErr = global.FUNC_PARAM_ERROR
			return
		}

		if !funcReg.MatchString(funcsStr) {
			gErr = global.NewGotrixError(global.NO_PERMISSION_PAGE_ERROR, pageId)
			return
		}

		var dictMap map[string]interface{} = make(map[string]interface{})

		var addDictMap = func(name string) *global.GotrixError {
			if dictMap[name] != nil {
				return nil
			}
			function := funcNameMap[name]
			var checkedParams *global.CheckedParams
			if function == nil {
				function = funcNameMap["dict"]
				if function == nil {
					return nil
				}
				checkedParams = &global.CheckedParams{Func: function.Id, V: make(map[string]interface{})}
				checkedParams.V["name"] = name
			} else {
				checkedParams = &global.CheckedParams{Func: function.Id, V: make(map[string]interface{})}
			}
			result, gotrixError := this.simpleHandler.Handle(checkedParams)
			if gotrixError != nil {
				return gotrixError
			}
			dictMap[name] = result
			return nil
		}

		var dealWithFuncs = func(funcs []Func) (returnFuncs []interface{}) {
			paramReg := regexp.MustCompile("\\$\\{\\w*\\}")
			returnFuncs = make([]interface{}, 0)
			for _, d := range funcs {
				tempMap := make(map[string]interface{})
				tempMap["Id"] = d.Id
				tempMap["Des"] = d.Des

				jobs := make([]string, 0)
				pagination := false
				for i, lenI := 0, len(d.Jobs); i < lenI; i++ {
					allParam := paramReg.FindAllString(d.Jobs[i].Job, -1)
					for j, lenJ := 0, len(allParam); j < lenJ; j++ {
						jobs = append(jobs, allParam[j][2:len(allParam[j])-1])
					}
					if d.Jobs[i].Type == "pagination" {
						pagination = true
					}
				}
				tempMap["Jobs"] = jobs
				tempMap["Pagination"] = pagination

				params := make([]interface{}, 0)
				for i, lenI := 0, len(d.Param); i < lenI; i++ {
					tempParam := make(map[string]interface{})
					tempParam["Name"] = d.Param[i].Name
					tempParam["Type"] = d.Param[i].Type
					if len(d.Param[i].Must) > 0 {
						tempParam["Must"] = d.Param[i].must
					}
					if len(d.Param[i].Des) > 0 {
						tempParam["Des"] = d.Param[i].Des
					}
					if len(d.Param[i].Form) > 0 {
						tempParam["Form"] = d.Param[i].Form
					}
					if len(d.Param[i].Len) > 0 {
						tempParam["Len"] = d.Param[i].Len
					}
					if len(d.Param[i].Dict) > 0 {
						addDictMap(d.Param[i].Dict)
						tempParam["Dict"] = d.Param[i].Dict
					}
					params = append(params, tempParam)
				}
				tempMap["Param"] = params

				returnFuncs = append(returnFuncs, tempMap)
			}
			return
		}

		page := pageMap[pageId]
		returnMap := make(map[string]interface{})
		returnMap["Id"] = page.Id
		returnMap["Name"] = page.Name
		returnMap["Des"] = page.Des
		returnMap["Param"] = page.Param
		for i, length := 0, len(page.Param); i < length; i++ {
			if len(page.Param[i].Dict) > 0 {
				addDictMap(page.Param[i].Dict)
			}
		}
		returnMap["Insert"] = dealWithFuncs(page.Insert)
		returnMap["Delete"] = dealWithFuncs(page.Delete)
		returnMap["Update"] = dealWithFuncs(page.Update)
		returnMap["Select"] = dealWithFuncs(page.Select)
		returnMap["Dict"] = dictMap
		fmt.Println(dictMap)

		return returnMap, nil
	}
}
