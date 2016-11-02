package handler

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"unicode"
)

var funcReg *regexp.Regexp = regexp.MustCompile("^(\\w+)\\((.*)\\)$")
var argsReg *regexp.Regexp = regexp.MustCompile("((true)|(false)|(null)|(-?\\d+)|(-?\\d+\\.\\d+)|(\\'.*?\\')|(\\$\\{\\w+\\}))(?:,|$)")
var sqlArgsReg *regexp.Regexp = regexp.MustCompile("\\$\\{\\w+\\}")

var pHandleHttp Handle
var pHandleFunc Handle
var pHandleRedis Handle
var pHandleSql Handle

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
	Parent int `xml:"parent,attr"`
	Insert int `xml:"insert,attr"`
	Delete int `xml:"delete,attr"`
	Update int `xml:"update,attr"`
	Select int `xml:"select,attr"`
}

type Job struct {
	Result  string `xml:"result,attr"`
	Test    string `xml:"test,attr"`
	One     bool   `xml:"one,attr"`
	Job     string `xml:",innerxml"`
	handle  Handle
	testJob *Job
}

type Param struct {
	XMLName xml.Name
	Type    string `xml:"type,attr"`
	Name    string `xml:"name,attr"`
	Must    bool   `xml:"must,attr"`
	Len     string `xml:"len,attr"`
	Valid   Valid
	min     int64
	max     int64
}

var funcMap map[int]*Func
var sqlMap map[int]*Sql
var pageMap map[int]*Page

func readXmlFolder(simpleHandler SimpleHandler, folder string) {

	pHandleHttp = &handleHttp{}
	pHandleFunc = (&handleFunc{simpleHandler: simpleHandler}).init()
	pHandleRedis = (&handleRedis{}).init()
	pHandleSql = (&handleSql{}).init()

	funcMap = make(map[int]*Func)
	sqlMap = make(map[int]*Sql)
	pageMap = make(map[int]*Page)

	var result Result

	_, filename, _, _ := runtime.Caller(1)
	baseDir := regexp.MustCompile("src.*$").ReplaceAllString(filename, "")

	// 对文件夹进行遍历，读取所有XML文件
	filepath.Walk(baseDir+folder, func(path string, info os.FileInfo, err error) error {
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

	stringVaid := StringValid{}
	intValid := IntValid{}
	boolValid := BoolValid{}
	arrayValid := ArrayValid{}

	for i := 0; i < len(result.Func); i++ {
		v := result.Func[i]
		funcMap[v.Id] = &v

		// 设置各Param的参数检查器
		for i := 0; i < len(v.Param); i++ {
			switch v.Param[i].XMLName.Local {
			case "string":
				v.Param[i].Valid = stringVaid
				if len(v.Param[i].Len) > 0 {
					v.Param[i].min, _ = strconv.ParseInt(regexp.MustCompile("^\\d+").FindString(v.Param[i].Len), 10, 64)
					v.Param[i].max, _ = strconv.ParseInt(regexp.MustCompile("\\d+$").FindString(v.Param[i].Len), 10, 64)
				}
				break
			case "int":
				v.Param[i].Valid = intValid
				if len(v.Param[i].Len) > 0 {
					v.Param[i].min, _ = strconv.ParseInt(regexp.MustCompile("^\\d+").FindString(v.Param[i].Len), 10, 64)
					v.Param[i].max, _ = strconv.ParseInt(regexp.MustCompile("\\d+$").FindString(v.Param[i].Len), 10, 64)
				}
				break
			case "bool":
				v.Param[i].Valid = boolValid
				break
			case "array":
				v.Param[i].Valid = arrayValid
				break
			default:
				break
			}
		}

		// 去掉Job的<![CDATA[]]>标签
		cdataExp := regexp.MustCompile("^<!\\[CDATA\\[(.*?)\\]\\]>$")
		for j := 0; j < len(v.Jobs); j++ {
			flag := cdataExp.MatchString(v.Jobs[j].Job)
			if flag {
				v.Jobs[j].Job = cdataExp.FindAllStringSubmatch(v.Jobs[j].Job, -1)[0][1]
			}
			strings.TrimSpace(v.Jobs[j].Job)
		}

		// 解析Job标签，并给不同的Job标签添加不同的处理器
		for j := 0; j < len(v.Jobs); j++ {
			funcStrs := funcReg.FindAllStringSubmatch(v.Jobs[j].Job, -1)
			if len(funcStrs) == 0 {
				if strings.HasPrefix(v.Jobs[j].Job, "http") {
					v.Jobs[j].handle = pHandleHttp
				} else {
					v.Jobs[j].handle = pHandleSql
				}
			} else {
				funcName := funcStrs[0][1]
				// 首字母大写，表示调用的是本地方法
				if unicode.IsUpper(rune(funcName[0])) {
					v.Jobs[j].handle = pHandleFunc
				} else { // 首字母不是大写则表示调用 Redis 的接口
					v.Jobs[j].handle = pHandleRedis
				}
			}

			// 解析Test标签，将Test标签转换为testJob
			if len(v.Jobs[j].Test) > 0 {
				v.Jobs[j].testJob = &Job{Job: v.Jobs[j].Test}
			}
		}

	}
	for i := 0; i < len(result.Sql); i++ {
		sqlMap[result.Sql[i].Id] = &result.Sql[i]
	}
	for i := 0; i < len(result.Page); i++ {
		pageMap[result.Page[i].Id] = &result.Page[i]
	}

}

func readXmlFile(xmlFileName string, result *Result) {

	// 读取文件内容
	content, err := ioutil.ReadFile(xmlFileName)
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
