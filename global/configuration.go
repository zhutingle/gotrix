package global

import (
	"net/mail"
	"time"

	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"path"
	"path/filepath"
)

var Config Configuration

type Configuration struct {
	Args     Args
	WEB      Dir
	CMS      Dir
	Redis    Redis
	Database Database
	WxCert   WxCert
	Email    Email
	V        []V
	M        map[string]interface{}
}

type Args struct {
	ConfigFile string // 配置文件
	Decrypt    bool   // 控制台参数 --decrypt
	Console    bool   // 控制台参数 --console
	Debug      bool

	password   string // 控制台参数 --password {{password}} ，该参数不允许在配置文件中配置，否则容易被窃取
}

type Dir struct {
	Port    int64  // 提供服务的端口号
	Base    string // 基目录，其它一切目录都基于此
	Temp    string // 临时文件夹
	LogFile string // 日志文件
	Func    string // 功能目录
	Static  string // 静态目录
	Target  string // 输出目录
}

type Redis struct {
	Ip          string
	Host        string
	Pass        string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

type Database struct {
	Url          string
	User         string
	Pass         string
	MaxOpenConns int
	MaxIdleConns int
}

type WxCert struct {
	Cert_pem    string
	Key_pem     string
	RootCA_Path string
}

type Email struct {
	Address  mail.Address
	SmtpUrl  string
	Identify string
	Username string
	Password string
	Host     string
}

type V struct {
	Key string
	Val string
}

func InitArgs() {

	args := append(os.Args, os.Environ()...)

	for i, arg := range args {
		switch arg {
		case "--config", "-f":
			Config.Args.ConfigFile = args[i + 1]
			break
		case "--decrypt", "-d":
			Config.Args.Decrypt = true
			break
		case "--console", "-c":
			Config.Args.Console = true
			break
		case "--password", "-p":
			Config.Args.password = args[i + 1]
			break
		default:
			break
		}
	}

	if len(Config.Args.ConfigFile) == 0 {
		// 未指定配置文件，从当前目录寻找 gotrix.conf
		log.Println("配置文件未配置，自动寻找默认配置文件 gotrix.conf")
		s, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			panic("获取当前目录失败")
		}
		Config.Args.ConfigFile = filepath.Join(s, filepath.FromSlash(path.Clean("/gotrix.conf")))
	}
}

func InitConfiguration() {

	log.Println("开始读取配置文件 : ", Config.Args.ConfigFile)

	// 根据配置文件的位置，设置默认参数，并按照如下约定自动设置其文件夹
	// project
	//     --WEB
	//     --WEB-FUNC
	//     --WEB-TEMP
	//     --WEB-TARGET
	//     --CMS
	//     --CMS-FUNC
	//     --CMS-TEMP
	//     --CMS-TARGET
	//     gotrix.conf
	//     web.log
	//     cms.log
	Config.WEB.Base = filepath.Dir(Config.Args.ConfigFile)
	Config.WEB.Port = 9080
	Config.WEB.Temp = "/WEB-TEMP"
	Config.WEB.LogFile = "/web.log"
	Config.WEB.Func = "/WEB-FUNC"
	Config.WEB.Static = "/WEB"
	Config.WEB.Target = "/WEB-TARGET"

	Config.CMS.Base = filepath.Dir(Config.Args.ConfigFile)
	Config.CMS.Port = 9088
	Config.CMS.Temp = "/CMS-TEMP"
	Config.CMS.LogFile = "/cms.log"
	Config.CMS.Func = "/CMS-FUNC"
	Config.CMS.Static = "/CMS"
	Config.CMS.Target = "/WEB-TARGET"

	ReadConfigFile(Config.Args.ConfigFile, func(bs []byte, err error) {
		if err != nil {
			panic(fmt.Sprintf("读取配置文件出现一个异常：[%v]", err))
		}

		if _, err := toml.Decode(string(bs), &Config); err != nil {
			panic(fmt.Sprintf("读取全局配置文件时出现一个异常：[%v]", err))
		}

		Config.M = make(map[string]interface{}, 0)
		for _, v := range Config.V {
			Config.M[v.Key] = v.Val
		}
		log.Println("读取配置文件成功!")

		Config.WEB.Temp = path.Clean(Config.WEB.Base + Config.WEB.Temp)
		Config.WEB.LogFile = path.Clean(Config.WEB.Base + Config.WEB.LogFile)
		Config.WEB.Func = path.Clean(Config.WEB.Base + Config.WEB.Func)
		Config.WEB.Static = path.Clean(Config.WEB.Base + Config.WEB.Static)
		Config.WEB.Target = path.Clean(Config.WEB.Base + Config.WEB.Target)

		Config.CMS.Temp = path.Clean(Config.CMS.Base + Config.CMS.Temp)
		Config.CMS.LogFile = path.Clean(Config.CMS.Base + Config.CMS.LogFile)
		Config.CMS.Func = path.Clean(Config.CMS.Base + Config.CMS.Func)
		Config.CMS.Static = path.Clean(Config.CMS.Base + Config.CMS.Static)
		Config.CMS.Target = path.Clean(Config.CMS.Base + Config.CMS.Target)

	})

}

func StartProcess() {
	filePath, _ := filepath.Abs(os.Args[0])
	logFile, err := os.Create(Config.WEB.LogFile)
	if err != nil {
		panic(fmt.Sprintf("创建日志文件时出现异常：%v", err))
	}
	log.Println("日志文件创建成功：", Config.WEB.LogFile)

	process, err := os.StartProcess(filePath, os.Args, &os.ProcAttr{Env: []string{"--console", "--password", Config.Args.password}, Files: []*os.File{logFile, logFile, logFile}})
	if err != nil {
		log.Println(err)
	}
	log.Println("新进程创建成功：", process)
}