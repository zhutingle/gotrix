package global

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

func init() {
	fmt.Println("Gotrix Global Configuration ...")

	filePath, _ := filepath.Abs(os.Args[0])
	lastIndexOfSeperator := strings.LastIndex(filePath, string(filepath.Separator))
	filePath = filePath[:lastIndexOfSeperator+1]

	if _, err := toml.DecodeFile(filePath+"gotrix.conf", &Config); err != nil {
		fmt.Println("Gotrix Global Configuration cause an error:")
		fmt.Println(err)
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
	LogFile  string
	Redis    Redis
	Database Database
	WxCert   WxCert
	V        []V
	M        map[string]interface{}
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

type V struct {
	Key string
	Val string
}
