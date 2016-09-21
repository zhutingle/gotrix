package handler

import (
	"fmt"
	"log"
	"strconv"

	"github.com/robfig/cron"
	"github.com/zhutingle/gotrix/global"
)

type SimpleHandler struct {
	cronManager *cron.Cron
}

type handle interface {
	handle(job *Job, cp *global.CheckedParams) (result interface{}, gErr *global.GotrixError)
}

func (this SimpleHandler) Init() {

	readXmlFolder("src/github.com/zhutingle/gotrix/func")

	this.cronTask()

}

func (this SimpleHandler) Handle(checkedParams *global.CheckedParams) (response interface{}, gErr *global.GotrixError) {
	// 1、从配置文件中取某个功能号对应的配置
	var handleFunc *Func = funcMap[checkedParams.Func]
	if handleFunc == nil {
		gErr = global.NewGotrixError(global.FUNC_NOT_EXISTS, checkedParams.Func)
		return
	}
	// 2、判断是否是私有方法，私有方法不能通过外部访问
	if checkedParams.Checked && handleFunc.Private {
		gErr = global.FUNC_PRIVATE_ERROR
		return
	}
	// 3、判断解密类型是否一致
	if handleFunc.Self != checkedParams.Self {
		gErr = global.FUNC_SELF_ERROR
		return
	}
	// 4、根据配置对参数进行验证
	for _, param := range handleFunc.Param {
		//	 参数验证
		gErr = param.Valid.Valid(&param, checkedParams.V[param.Name])
		if gErr != nil {
			return
		}
	}
	// 5、根据配置文件中的配置的业务，进行业务操作
	if len(handleFunc.Jobs) > 0 {
		return this.jobHandle(handleFunc, checkedParams)
	}

	return
}

func (this SimpleHandler) jobHandle(handleFunc *Func, checkedParams *global.CheckedParams) (response interface{}, gErr *global.GotrixError) {

	i, length := 0, len(handleFunc.Jobs)
	for ; i < length; i++ {

		job := &handleFunc.Jobs[i]

		if job.testJob != nil {
			_, err := pHandleFunc.handle(job.testJob, checkedParams)
			if err == nil {
				continue
			}
		}

		response, gErr = handleFunc.Jobs[i].handle.handle(job, checkedParams)

		if gErr != nil {
			log.Printf("指令[%s]在执行时出现异常：%s\n", job.Job, fmt.Sprint(gErr))
			return
		}

		if len(job.Result) == 0 {
			checkedParams.V[strconv.Itoa(i+1)] = response
		} else {
			checkedParams.V[job.Result] = response
		}

	}

	return
}

func (this SimpleHandler) GetPass(token []byte) (response interface{}, gErr *global.GotrixError) {
	var checkedParams *global.CheckedParams = &global.CheckedParams{Func: 0, V: make(map[string]interface{})}
	checkedParams.V["token"] = string(token)
	response, gErr = this.Handle(checkedParams)
	return
}

func (this SimpleHandler) GetSession(token []byte) (response interface{}, gErr *global.GotrixError) {
	var checkedParams *global.CheckedParams = &global.CheckedParams{Func: 1, V: make(map[string]interface{})}
	checkedParams.V["token"] = string(token)
	response, gErr = this.Handle(checkedParams)
	return
}
