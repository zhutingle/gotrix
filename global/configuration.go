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
	Dir      Dir
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
	Password   string // 控制台参数 --password {{password}}
	Port       int64  // 端口号
	Debug      bool
}

type Dir struct {
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
			Config.Args.Password = args[i + 1]
			break
		default:
			break
		}
	}

	if len(Config.Args.ConfigFile) == 0 {
		panic("未指定配置文件。")
	}
}

func pathJoin(str string) string {
	return filepath.Join(Config.Dir.Base, filepath.FromSlash(path.Clean("/" + str)))
}

func InitConfiguration() {

	log.Println("开始读取配置文件 : ", Config.Args.ConfigFile)

	ReadConfigFile(Config.Args.ConfigFile, func(bs []byte, err error) {
		if err != nil {
			panic(fmt.Sprintf("读取配置文件出现一个异常：[%v]", err))
		}

		if _, err := toml.Decode(string(bs), &Config); err != nil {
			panic(fmt.Sprintf("读取全局配置文件时出现一个异常：[%v]", err))
		} else {
			Config.M = make(map[string]interface{}, 0)
			for _, v := range Config.V {
				Config.M[v.Key] = v.Val
			}
			log.Println("读取配置文件成功!")
		}

		Config.Dir.LogFile = pathJoin(Config.Dir.LogFile)
		Config.Dir.Temp = pathJoin(Config.Dir.Temp)
		Config.Dir.Func = pathJoin(Config.Dir.Func)
		Config.Dir.Static = pathJoin(Config.Dir.Static)
		Config.Dir.Target = pathJoin(Config.Dir.Target)

		InitArgs()
	})

}
