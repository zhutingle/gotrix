package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/zhutingle/gotrix/global"
)

type handleHttp struct {
}

func (this *handleHttp) handle(job *Job, cp *global.CheckedParams) (result interface{}, gErr *global.GotrixError) {

	handleUrl := sqlArgsReg.ReplaceAllStringFunc(job.Job, func(str string) string {
		var name = str[2 : len(str)-1]
		if _, ok := cp.V[name].(string); ok {
			return cp.V[name].(string)
		} else {
			bts, e := json.Marshal(cp.V[name])
			if e != nil {
				return ""
			} else {
				return string(bts)
			}
		}
	})

	resp, e := http.Get(handleUrl)
	if e != nil {
		log.Println(e)
		gErr = global.HTTPHANDLE_HTTP_GET_ERROR
		return
	}

	defer resp.Body.Close()
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		log.Println(e)
		gErr = global.HTTPHANDLE_HTTP_READ_BODY
		return
	}

	var returnInterface interface{}
	e = json.Unmarshal(body, &returnInterface)
	if e != nil {
		log.Println(e)
		gErr = global.HTTPHANDLE_ANALYZE_ERROR
		return
	}

	return returnInterface, nil
}
