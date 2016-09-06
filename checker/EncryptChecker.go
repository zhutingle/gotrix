package checker

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/zhutingle/gotrix/global"
)

type EncryptChecker struct {
}

func (checker EncryptChecker) Check(r *http.Request, handler global.Handler) (checkedParams *global.CheckedParams, err *global.GotrixError) {

	checkedParams = &global.CheckedParams{V: make(map[string]interface{})}

	if r.Body == nil {
		err = global.NO_BODY_ERROR
		return
	}

	var reader io.Reader = r.Body
	b, e := ioutil.ReadAll(reader)
	if e != nil {
		err = global.READ_BODY_ERROR
		return
	}

	i, len := 0, len(b)
	for ; i < len; i++ {
		if b[i] == '=' {
			break
		}
	}
	if i == 0 || i >= len {
		err = global.BODY_SCHEME_ERROR
		return
	}

	var self = false
	var aesPass interface{}
	switch i {
	case 40: // session
		aesPass, err = handler.GetSession(b[:i])
		if err != nil {
			return
		}
		if aesPass == nil {
			err = global.USER_SESSION_NOT_EXISTES
			return
		}
		break
	case 64: // token
		aesPass, err = handler.GetPass(b[:i])
		if err != nil {
			return
		}
		if aesPass == nil {
			err = global.USER_NOT_EXISTES
			return
		}
		break
	case 32: // self 自解密
		aesPass = string(b[:31]) // 自解密
		self = true
		break
	default:
		err = global.NOT_SUPPORT_CONTENT_TYPE
		return
	}

	var pass []byte
	var userid int64
	switch aesPass.(type) {
	case string:
		pass = []byte(aesPass.(string))
		var decryptBytes []byte
		decryptBytes, e = AesDecrypt(b[i+1:], pass, 256)
		if e == nil {
			e = json.Unmarshal(decryptBytes, &(checkedParams.V))
		}
		break
	case map[string]interface{}:
		passJson := aesPass.(map[string]interface{})
		if passJson["id"] != nil {
			userid = int64(passJson["id"].(float64))
		}
		for _, val := range passJson {

			if _, ok := val.(string); ok {
				pass = []byte(val.(string))
				var decryptBytes []byte
				decryptBytes, e = AesDecrypt(b[i+1:], pass, 256)
				if e != nil {
					continue
				}

				e = json.Unmarshal(decryptBytes, &(checkedParams.V))
				if e != nil {
					continue
				}
			} else {
				continue
			}

			break
		}
		break
	}

	if e != nil {
		err = global.PASSWORD_ERROR
		return
	}

	var fun interface{} = checkedParams.V["func"]
	if fun == nil {
		err = global.FUNC_PARAM_MUST
		return
	}

	switch fun.(type) {
	case float64:
		checkedParams.Func = int(checkedParams.V["func"].(float64))
		checkedParams.V["token"] = string(b[:i])
		checkedParams.V["_ip"] = r.Header.Get("X-Forward-For")
		checkedParams.Pass = pass
		checkedParams.Checked = true
		checkedParams.Self = self
		break
	default:
		err = global.FUNC_PARAM_ERROR
		break
	}

	if userid > 0 {
		checkedParams.V["userid"] = userid
	}

	return
}
