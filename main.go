package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	//	"os"
	"os/exec"
	//	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/zhutingle/gotrix/checker"
	"github.com/zhutingle/gotrix/global"
	"github.com/zhutingle/gotrix/handler"
	"github.com/zhutingle/gotrix/weichat"
)

func main() {

	//	gotrixHandler = handler.SimpleHandler{}
	//	gotrixHandler.Init()
	//
	//	checkedParams := &global.CheckedParams{Func: 1004, V: make(map[string]interface{})}
	//	response, err := gotrixHandler.Handle(checkedParams)
	//	log.Println(response)
	//	log.Println(err)

	if len(os.Args) > 1 {
		GotrixServer()
	} else {
		filePath, _ := filepath.Abs(os.Args[0])
		args := append([]string{filePath}, "")

		logFile, _ := os.Create(global.Config.LogFile)
		process, err := os.StartProcess(filePath, args, &os.ProcAttr{Files: []*os.File{logFile, logFile, logFile}})
		if err != nil {
			log.Println(err)
		}
		log.Println(process)
	}

}

var gotrixChecker global.Checker
var gotrixHandler global.Handler

func GotrixServer() {

	// -----杀掉原有实例---------------------------------------------------------
	c := exec.Command("netstat", "/nao")
	bs, err := c.Output()
	if err != nil {
		fmt.Println(err)
	}

	reg := regexp.MustCompile("TCP\\s*0\\.0\\.0\\.0:9080.*LISTENING\\s*(\\d*)")
	line := reg.Find(bs)

	if len(line) > 0 {
		reg2 := regexp.MustCompile("\\d*$")
		pid := reg2.Find(line)

		c1 := exec.Command("taskkill", "-f", "/pid", string(pid))
		err1 := c1.Run()
		if err1 != nil {
			fmt.Println(err1)
		}
	}
	// -----------------------------------------------------------------------

	gotrixChecker = checker.EncryptChecker{}
	gotrixHandler = handler.SimpleHandler{}
	gotrixHandler.Init()

	http.HandleFunc("/gotrix/", serverHandler)
	http.HandleFunc("/gotrix/wxpay.action", wxpayCallback)
	http.Handle("/", http.FileServer(http.Dir("src/github.com/zhutingle/gotrix/static")))

	err = http.ListenAndServe(":9080", nil)
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
	logBuffer.WriteString(strconv.FormatInt(int64(checkedParams.Func), 10))
	logBuffer.WriteString("\n----Param: ")
	logBuffer.WriteString(fmt.Sprint(checkedParams.V))

	// --------------------业务执行器--------------------
	var response interface{}
	var result []byte
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
	encryptResult, e := checker.AesEncrypt(buffer.Bytes(), checkedParams.Pass, 256)
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
	logBuffer.Write(result)
	logBuffer.WriteString("\n----Spend: ")
	logBuffer.WriteString(strconv.FormatInt((time.Now().UnixNano()-start)/1000000, 10))
	logBuffer.WriteString(" ms")
	logBuffer.WriteRune('\n')
	log.Println(logBuffer.String())

}

func wxpayCallback(w http.ResponseWriter, r *http.Request) {
	checkedParams := &global.CheckedParams{Func: 1001, V: make(map[string]interface{}, 0)}
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
