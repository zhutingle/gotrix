package handler

import (
	"encoding/xml"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/zhutingle/gotrix/global"
)

var funcReg *regexp.Regexp = regexp.MustCompile("^(\\w+)\\((.*)\\)$")
var argsReg *regexp.Regexp = regexp.MustCompile("((true)|(false)|(null)|(-?\\d+)|(-?\\d+\\.\\d+)|(\\'.*?\\')|(\\$\\{\\w+\\}))(?:,|$)")
var sqlArgsReg *regexp.Regexp = regexp.MustCompile("\\$\\{\\w+\\}")
var autoTagReg *regexp.Regexp = regexp.MustCompile("<auto>(.*?)</auto>")
var autoTagItemReg *regexp.Regexp = regexp.MustCompile("\\w+\\s*=\\s*\\$\\{\\w+\\}\\s*(?:,|$|(and))\\s*")

type Result struct {
	Sql  []Sql  `xml:"sql"`
	Func []Func `xml:"func"`
	Page []Page `xml:"page"`
}

type Func struct {
	Id      int     `xml:"id,attr"`
	Name    string  `xml:"name,attr"`
	Des     string  `xml:"des,attr"`
	Private bool    `xml:"private,attr"`
	Self    bool    `xml:"self,attr"`
	Cron    string  `xml:"cron,attr"`
	Jobs    []Job   `xml:"job"`
	Param   []Param `xml:",any"`
}

type Sql struct {
	Func
	Test string `xml:"test,attr"`
}

type Page struct {
	Func
	Parent int    `xml:"parent,attr"`
	Insert []Func `xml:"insert"`
	Delete []Func `xml:"delete"`
	Update []Func `xml:"update"`
	Select []Func `xml:"select"`
}

type Job struct {
	Result  string `xml:"result,attr"`
	Test    string `xml:"test,attr"`
	Type    string `xml:"type,attr"`
	Job     string `xml:",innerxml"`
	handle  Handle
	testJob *Job
	auto    bool
}

type Param struct {
	XMLName xml.Name
	Type    string `xml:"type,attr"`
	Name    string `xml:"name,attr"`
	Des     string `xml:"des,attr"`
	Must    string `xml:"must,attr"`
	Len     string `xml:"len,attr"`
	Form    string `xml:"form,attr"`
	Dict    string `xml:"dict,attr"` // 数据字典
	Valid   Valid
	min     int64
	max     int64
	must    bool
}

func readXmlBytes(content []byte) {

	var result Result

	err := xml.Unmarshal(content, &result)
	if err != nil {
		log.Fatal(err)
		return
	}

	dealWithResult(&result)
}

func readXmlFolder(folder string) {

	var result Result

	// 对文件夹进行遍历，读取所有XML文件
	filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if !info.IsDir() {
			match, _ := regexp.MatchString("\\.xml$", path)
			if match {
				readXmlFile(path, &result)
			}
		}
		return nil
	})

	dealWithResult(&result)
}

func dealWithResult(result *Result) {

	dealWithFuncs(result.Func)

	for i := 0; i < len(result.Sql); i++ {
		v := result.Sql[i]
		sqlMap[v.Id] = &v

		v.Param = dealWithParam(v.Param)
	}
	for i := 0; i < len(result.Page); i++ {
		v := result.Page[i]
		pageMap[v.Id] = &v

		v.Param = dealWithParam(v.Param)

		dealWithFuncs(v.Insert)
		dealWithFuncs(v.Delete)
		dealWithFuncs(v.Update)
		dealWithFuncs(v.Select)

	}
}

func dealWithFuncs(funcs []Func) {
	for i := 0; i < len(funcs); i++ {
		v := funcs[i]
		funcMap[v.Id] = &v
		funcNameMap[v.Name] = &v

		v.Param = dealWithParam(v.Param)
		dealWithJob(v.Jobs)
	}
}

func dealWithParam(params []Param) []Param {
	// 设置各Param的参数检查器
	for i := 0; i < len(params); i++ {

		params[i].must = (params[i].Must == "true")
		params[i].Type = params[i].XMLName.Local

		switch params[i].XMLName.Local {
		case "string":
			params[i].Valid = stringVaid
			if len(params[i].Len) > 0 {
				params[i].min, _ = strconv.ParseInt(regexp.MustCompile("^\\d+").FindString(params[i].Len), 10, 64)
				params[i].max, _ = strconv.ParseInt(regexp.MustCompile("\\d+$").FindString(params[i].Len), 10, 64)
			}
			break
		case "int":
			params[i].Valid = intValid
			if len(params[i].Len) > 0 {
				params[i].min, _ = strconv.ParseInt(regexp.MustCompile("^\\d+").FindString(params[i].Len), 10, 64)
				params[i].max, _ = strconv.ParseInt(regexp.MustCompile("\\d+$").FindString(params[i].Len), 10, 64)
			}
			break
		case "bool":
			params[i].Valid = boolValid
			break
		case "array":
			params[i].Valid = arrayValid
			break
		case "file":
			params[i].Valid = fileValid
			break
		default:
			params = append(params[:i], params[i+1:]...)
			i--
			break
		}
	}
	return params
}

func dealWithJob(jobs []Job) {
	// 去掉Job的<![CDATA[]]>标签
	cdataExp := regexp.MustCompile("^<!\\[CDATA\\[([\\w\\W]*?)\\]\\]>$")
	autoExp := regexp.MustCompile("<auto>.*?</auto>")
	for j := 0; j < len(jobs); j++ {
		flag := cdataExp.MatchString(jobs[j].Job)
		if flag {
			jobs[j].Job = cdataExp.FindAllStringSubmatch(jobs[j].Job, -1)[0][1]
		}
		jobs[j].Job = strings.TrimSpace(jobs[j].Job)
	}

	// 解析Job标签，并给不同的Job标签添加不同的处理器
	for j := 0; j < len(jobs); j++ {

		jobs[j].auto = autoExp.MatchString(jobs[j].Job)

		funcStrs := funcReg.FindAllStringSubmatch(jobs[j].Job, -1)
		if len(funcStrs) == 0 {
			if strings.HasPrefix(jobs[j].Job, "http") {
				jobs[j].handle = pHandleHttp
			} else {
				jobs[j].handle = pHandleSql
			}
		} else {
			funcName := funcStrs[0][1]
			// 首字母大写，表示调用的是本地方法
			if unicode.IsUpper(rune(funcName[0])) {
				jobs[j].handle = pHandleFunc
			} else {
				// 首字母不是大写则表示调用 Redis 的接口
				jobs[j].handle = pHandleRedis
			}
		}

		// 解析Test标签，将Test标签转换为testJob
		if len(jobs[j].Test) > 0 {
			jobs[j].testJob = &Job{Job: jobs[j].Test}
		}
	}
}

func readXmlFile(xmlFileName string, result *Result) {

	// 读取文件内容
	content, err := global.ReadConfigFile(xmlFileName, nil)
	//	content, err := ioutil.ReadFile(xmlFileName)
	if err != nil {
		log.Fatal(err)
		return
	}

	// 读取XML内容
	err = xml.Unmarshal(content, result)
	if err != nil {
		log.Fatal(err)
		return
	}
}
