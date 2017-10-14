package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/zhutingle/gotrix/checker"
	"github.com/zhutingle/gotrix/global"
	"github.com/zhutingle/gotrix/handler"
	"github.com/zhutingle/gotrix/weichat"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"time"
)

func GotrixServer() {

	global.InitArgs()
	global.InitPassword()
	global.InitConfiguration()
	global.InitArgs() // 命令行参数"强行覆盖"配置文件的[Args]参数

	for _, args := range os.Args {
		switch args {
		case "start":
			if global.Config.Args.Console {
				Server()
			} else {
				global.StartProcess()
			}
			break
		case "static":
			// 用于处理静态资源文件，将基础静态资源文件复制到 Gotrix 代码中，则其它项目默认就有这些静态资源文件
			packageToGotrix(global.Config.WEB.Static, "static.go")
			packageToGotrix(global.Config.WEB.Func, "func.go")
			break
		case "init":
			// TODO
			// 用于初始化基础数据表
			// 数据字典表 dict dicti
			// 用户表 user
			// 角色表 role
			// 功能表 func
			break
		}
	}

	//gotrixHandler = handler.SimpleHandler{}
	//gotrixHandler.Init()
	//checkedParams := &global.CheckedParams{Func: 100, V: make(map[string]interface{})}
	//response, err := gotrixHandler.Handle(checkedParams)
	//log.Println(response)
	//log.Println(err)

}

var gotrixChecker global.Checker
var gotrixHandler global.Handler

func Server() {

	// -----杀掉原有实例---------------------------------------------------------
	if runtime.GOOS == "windows" {
		c := exec.Command("netstat", "/ano")
		bs, err := c.Output()
		if err != nil {
			fmt.Println(err)
		}

		reg := regexp.MustCompile("TCP\\s*0\\.0\\.0\\.0:9080.*LISTENING\\s*(\\d*)")
		matches := reg.FindSubmatch(bs)

		if matches != nil && len(matches) > 0 {
			c1 := exec.Command("taskkill", "-f", "/pid", string(matches[1]))
			err1 := c1.Run()
			if err1 != nil {
				fmt.Println(err1)
			}
		}
	}
	// -----------------------------------------------------------------------

	gotrixChecker = checker.EncryptChecker{}
	gotrixHandler = handler.SimpleHandler{}
	gotrixHandler.Init()

	http.HandleFunc("/gotrix/", serverHandler)
	http.HandleFunc("/gotrix/push.action", pushHandler)
	http.HandleFunc("/gotrix/wxpay.action", wxpayCallback)

	// Debug 模式和非 Debug 模式的区别全写在这里，其它地方不允许写
	if global.Config.Args.Debug {
		http.Handle("/", http.FileServer(DevDir(global.Config.WEB.Static)))
	} else {
		packageTarget(global.Config.WEB.Static, global.Config.WEB.Target)
		http.Handle("/", http.FileServer(http.Dir(global.Config.WEB.Target)))
	}

	err := http.ListenAndServe(fmt.Sprint(":", global.Config.WEB.Port), nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

func writeError(w http.ResponseWriter, err *global.GotrixError) {
	if err.Status > 0 {
		w.Write([]byte(fmt.Sprintf("{\"status\":%d,\"msg\":\"%s\"}", err.Status, err.Msg)))
	} else {
		w.Write([]byte(err.Msg))
	}
}

func serverHandler(w http.ResponseWriter, r *http.Request) {
	var start = time.Now().UnixNano()
	var logBuffer bytes.Buffer
	// --------------------参数解析器--------------------
	checkedParams, gErr := gotrixChecker.Check(r, gotrixHandler)
	if gErr != nil {
		writeError(w, gErr)

		logBuffer.WriteString("\n----Error: ")
		logBuffer.WriteString(fmt.Sprint(gErr))
		logBuffer.WriteRune('\n')
		log.Println(logBuffer.String())
		return
	}

	logBuffer.WriteString("\n----Func: ")
	logBuffer.WriteString(checkedParams.Name)
	logBuffer.WriteString("\n----Param: ")
	logBuffer.WriteString(fmt.Sprint(checkedParams.V))

	// --------------------业务执行器--------------------
	var response interface{}
	response, gErr = gotrixHandler.Handle(checkedParams)
	if gErr != nil {
		writeError(w, gErr)

		logBuffer.WriteString("\n----Error: ")
		logBuffer.WriteString(fmt.Sprint(gErr))
		logBuffer.WriteRune('\n')
		log.Println(logBuffer.String())
		return
	}

	// --------------------结果输出器--------------------
	buffer := bytes.NewBufferString("{\"status\":0,\"msg\":\"成功\",\"data\":")
	str, _ := json.Marshal(response)
	buffer.Write(str)
	buffer.WriteString("}")
	encryptResult, e := global.AesEncrypt(buffer.Bytes(), checkedParams.Pass, 256)
	if e != nil {
		writeError(w, global.RETURN_DATE_ECNRYPT_ERROR)

		logBuffer.WriteString("\n----Error: ")
		logBuffer.WriteString(fmt.Sprint(e))
		logBuffer.WriteRune('\n')
		log.Println(logBuffer.String())
		return
	}
	w.Write(encryptResult)

	logBuffer.WriteString("\n----Result: ")
	logBuffer.Write(str)
	logBuffer.WriteString("\n----Spend: ")
	logBuffer.WriteString(strconv.FormatInt((time.Now().UnixNano()-start)/1000000, 10))
	logBuffer.WriteString(" ms")
	logBuffer.WriteRune('\n')
	log.Println(logBuffer.String())

}

func pushHandler(w http.ResponseWriter, r *http.Request) {

	cmd := exec.Command("git", "pull")
	cmd.Dir = global.Config.WEB.Base

	f, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(f))

	w.Write([]byte("success"))

}

func wxpayCallback(w http.ResponseWriter, r *http.Request) {
	checkedParams := &global.CheckedParams{Name: "wxpayCallback", V: make(map[string]interface{}, 0)}
	weichat, err := weichat.WxpayCallback(w, r)
	if err != nil {
		fmt.Println(err)
	}
	checkedParams.V["weichat"] = weichat
	fmt.Println(checkedParams.V)
	response, gErr := gotrixHandler.Handle(checkedParams)
	if gErr != nil {
		writeError(w, gErr)
	} else {
		writeError(w, &global.GotrixError{Status: 0, Msg: fmt.Sprintf("%v", response)})
	}
}
