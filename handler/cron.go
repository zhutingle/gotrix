package handler

import (
	"log"
	"time"

	"github.com/robfig/cron"
	"github.com/zhutingle/gotrix/global"
)

type CronTask struct {
	fun     *Func
	handler global.Handler
}

func (this CronTask) Run() {
	startUnix := time.Now().Unix()
	log.Printf("定时任务[%d-%s]开始启动...\n", this.fun.Id, this.fun.Des)

	checkedParams := &global.CheckedParams{Func: this.fun.Id, V: make(map[string]interface{})}
	response, err := this.handler.Handle(checkedParams)
	if err != nil {
		log.Printf("定时任务[%d-%s]执行时出现异常：%v\n", this.fun.Id, this.fun.Des, err)
	} else {
		log.Printf("定时任务[%d-%s]执行成功返回：%v\n", this.fun.Id, this.fun.Des, response)
	}

	endUnix := time.Now().Unix()
	log.Printf("定时任务[%d-%s]执行完成...耗时[%d]秒\n", this.fun.Id, this.fun.Des, (endUnix - startUnix))
}

func (this SimpleHandler) cronTask() {
	// 初始化 cronManager
	if this.cronManager == nil {
		this.cronManager = cron.New()
	}
	// 遍历所有 Func 对象，取出所有配置有 cron 属性的 Func 对象，并将它加入定时任务
	for _, d := range funcMap {
		if len(d.Cron) > 0 {
			this.cronManager.AddJob(d.Cron, &CronTask{fun: d, handler: this})
		}
	}
	// 启动所有定时任务
	this.cronManager.Start()

}
