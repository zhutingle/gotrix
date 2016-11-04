package global

import (
	"fmt"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

func InitConfiguration() {
	fmt.Println("Gotrix 开始读取全局配置文件 ...")

	filePath, _ := filepath.Abs(os.Args[0])
	lastIndexOfSeperator := strings.LastIndex(filePath, string(filepath.Separator))
	filePath = filePath[:lastIndexOfSeperator+1]
	filePath = filePath + "gotrix.conf"

	bs, err := ReadConfigFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("读取配置文件出现一个异常：[%v]", err))
	}
	if _, err := toml.Decode(string(bs), &Config); err != nil {
		panic(fmt.Sprintf("Gotrix 读取全局配置文件时出现一个异常：[%v]", err))
	} else {
		Config.M = make(map[string]interface{}, 0)
		for _, v := range Config.V {
			Config.M[v.Key] = v.Val
		}
		fmt.Println("Gotrix Global Configuration success.")
	}

}

var Config Configuration

type Configuration struct {
	LogFile    string
	TempFolder string
	Redis      Redis
	Database   Database
	WxCert     WxCert
	Email      Email
	V          []V
	M          map[string]interface{}
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
