package global

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

/*
#include<stdlib.h>

void hide() {
	system("stty -echo");
}

void show() {
	system("stty echo");
}

*/
import "C"

// 加密字符串的前缀
var ENCRYPT_PREFIX string = "GOTRIX_ENCRYPTED:"

/**
 * 初始化密码
 **/
func InitPassword() {

	if len(Config.Args.Password) > 0 {
		return
	}

	if isEncrypted, _ := isGotrixEncrypted(Config.Args.ConfigFile); isEncrypted {
		log.Print("请输入您的密码:")
		password := readPassword()
		if len(password) == 0 {
			panic("密码不能为空，程序已退出，请重试。")
		} else {
			Config.Args.Password = password
		}
	} else {
		log.Print("请设置您的密码（第一次输入）:")
		password1 := readPassword()
		log.Print("请再次输入您的密码（第二次输入）:")
		password2 := readPassword()

		if password1 != password2 {
			panic("两次输入的密码不一致，程序已退出，请重试。")
		} else if len(password1) == 0 {
			panic("密码不能为空，程序已退出，请重试。")
		} else {
			Config.Args.Password = password1
		}
	}

}

/**
 * 读取用户输入的密码
 */
func readPassword() string {
	//exec.Command("stty", "/echo").Run()
	//defer exec.Command("stty", "echo").Run()
	C.hide()
	defer C.show()

	var password string
	var b []byte = make([]byte, 1)
	for {
		os.Stdin.Read(b)
		if b[0] != 10 {
			password += string(b)
		} else {
			return password
		}
	}
}

/**
 * 判断当前是否处于加密状态
 */
func isGotrixEncrypted(filePath string) (isEncrypted bool, bs []byte) {

	var err error
	bs, err = ioutil.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("配置文件[%s]不存在，或者没有读取权限：%v", filePath, err))
	}

	return strings.HasPrefix(string(bs), ENCRYPT_PREFIX), bs
}

/**
 * 若文件已加密，从加密文件中读取内容并返回
 * 若文件尚未加密，利用密码将文件加密，并返回文件内容
 */
func ReadConfigFile(filePath string, callback func([]byte, error)) ([]byte, error) {

	var isEncrypted bool
	var bs []byte

	if isEncrypted, bs = isGotrixEncrypted(filePath); isEncrypted {
		// 若文件已经做了加密处理

		bs = bytes.Replace(bs, []byte(ENCRYPT_PREFIX), []byte(nil), -1)

		decryptBs, err := AesDecrypt(bs, []byte(Config.Args.Password), 256)
		if err != nil {
			return bs, err
		}

		md5Bytes := Md5(decryptBs[16:])

		if string(decryptBs[:16]) == string(md5Bytes) {
			bs = decryptBs[16:]
			if callback != nil {
				callback(bs, err)
			}
			if Config.Args.Decrypt {
				ioutil.WriteFile(filePath, bs, 0)
			}
			return bs, err
		} else {
			if callback != nil {
				callback(nil, errors.New("密码不正确"))
			}
			return nil, errors.New("密码不正确")
		}

	} else {
		// 若文件尚未做加密处理

		md5Bytes := Md5(bs)

		buf := bytes.NewBuffer(md5Bytes)
		buf.Write(bs)

		encryptBs, err := AesEncrypt(buf.Bytes(), []byte(Config.Args.Password), 256)
		if err != nil {
			return bs, err
		}

		buf = bytes.NewBuffer([]byte(ENCRYPT_PREFIX))
		buf.Write(encryptBs)

		if callback != nil {
			callback(bs, err)
		}

		if !Config.Args.Decrypt {
			ioutil.WriteFile(filePath, buf.Bytes(), 0)
		}

		return bs, err
	}

}
