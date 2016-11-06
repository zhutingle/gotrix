package global

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// 加密字符串的前缀
var ENCRYPT_PREFIX string = "GOTRIX_ENCRYPTED:"

/**
 * 初始化密码
 **/
func InitPassword() {

	filePath, _ := filepath.Abs(os.Args[0])
	lastIndexOfSeperator := strings.LastIndex(filePath, string(filepath.Separator))
	filePath = filePath[:lastIndexOfSeperator+1]

	filePath = filePath + "gotrix.conf"

	if isEncrypted, _ := isGotrixEncrypted(filePath); isEncrypted {
		fmt.Print("请输入您的密码:")
		fmt.Scanln(&(Config.Args.Password))
	} else {
		fmt.Print("请设置您的密码（第一次输入）:")
		var password1, password2 string
		fmt.Scanln(&password1)
		fmt.Print("请再次输入您的密码（第二次输入）:")
		fmt.Scanln(&password2)

		if password1 != password2 {
			panic("两次输入的密码不一致，程序已退出，请重试。")
		} else {
			Config.Args.Password = password1
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
		panic(fmt.Sprintf("配置文件[%s]不存在，或者没有读取权限。", filePath))
	}

	return strings.HasPrefix(string(bs), ENCRYPT_PREFIX), bs
}

/**
 * 若文件已加密，从加密文件中读取内容并返回
 * 若文件尚未加密，利用密码将文件加密，并返回文件内容
 */
func ReadConfigFile(filePath string) ([]byte, error) {

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
			if Config.Args.Decrypt {
				ioutil.WriteFile(filePath, decryptBs[16:], 0)
			}
			return decryptBs[16:], err
		} else {
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

		if !Config.Args.Decrypt {
			ioutil.WriteFile(filePath, buf.Bytes(), 0)
		}

		return bs, err
	}

}
