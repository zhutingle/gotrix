package server

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/opesun/goquery"
	"github.com/zhutingle/gotrix/global"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// 对静态文件夹进行打包处理，输出到输出文件夹
func packageTarget(staticDir string, targetDir string) {

	// 判断输出文件夹是否存在，不存在则创建
	createDir(targetDir)

	b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")

	var htmlFiles map[string][]byte = make(map[string][]byte)
	var modelFiles map[string][]byte = make(map[string][]byte)
	var cssFiles map[string][]byte = make(map[string][]byte)
	var jsFiles map[string][]byte = make(map[string][]byte)
	var imgFiles map[string][]byte = make(map[string][]byte)

	// 遍历静态文件夹，并进行文件分类操作，并在输出文件夹中创建各子文件夹
	filepath.Walk(staticDir, func(longPath string, f os.FileInfo, err error) error {
		shortPath := strings.Replace(longPath, staticDir, "", -1)
		for _, c := range shortPath {
			if c == '/' || c == '\\' {
				shortPath = shortPath[1:]
			} else {
				break
			}
		}
		if f.IsDir() {
			createDir(filepath.Join(targetDir, filepath.FromSlash(path.Clean("/"+shortPath))))
			return nil
		}
		if strings.HasSuffix(shortPath, ".html") {
			// 目录中含有 model 就认为它是模版文件
			if strings.Contains(shortPath, "model") {
				modelFiles[shortPath], err = ioutil.ReadFile(longPath)
			} else {
				htmlFiles[shortPath], err = ioutil.ReadFile(longPath)
			}
		} else if strings.HasSuffix(shortPath, ".css") {
			cssFiles[shortPath], err = ioutil.ReadFile(longPath)
		} else if strings.HasSuffix(shortPath, ".js") {
			jsFiles[shortPath], err = ioutil.ReadFile(longPath)
		} else if strings.HasSuffix(shortPath, ".png") || strings.HasSuffix(shortPath, ".jpg") || strings.HasSuffix(shortPath, ".jpeg") {
			imgFiles[shortPath], err = ioutil.ReadFile(longPath)
		}
		if err != nil {
			log.Println("文件[", longPath, "]读取异常！")
		}
		return nil
	})

	// 对 img 文件进行 base64 处理，并在文件名后面加上 .cache
	var imgCacheFiles map[string][]byte = make(map[string][]byte)
	for shortPath, bs := range imgFiles {
		imgCacheFiles[shortPath+".cache"] = []byte(b64.EncodeToString(bs))
	}
	imgFiles = imgCacheFiles

	// 计算 html、css、js、img 文件的 MD5 值，不包括 js/potrix.cache.js
	var cacheFile string = "js/potrix.cache.js"
	var md5Map map[string]string = make(map[string]string)
	for _, files := range [](map[string][]byte){htmlFiles, cssFiles, jsFiles, imgFiles} {
		for shortPath, bs := range files {
			md5Map[shortPath] = b64.EncodeToString(global.Md5(bs))
		}
	}
	delete(md5Map, cacheFile)

	// 将各文件（不含 potrix.cache.js ）的 MD5 值合并成字符串，并写入 potrix.cache.js 再计算 potrix.cache.js 的 MD5 值
	md5Map[cacheFile] = b64.EncodeToString(global.Md5(bytes.Replace(jsFiles[cacheFile], []byte("(function (w) {"), createMd5String(md5Map), -1)))

	// 将各文件（包含 potrix.cache.js）的 MD5 值合并成字符串，并写入 potrix.cache.js
	jsFiles[cacheFile] = bytes.Replace(jsFiles[cacheFile], []byte("(function (w) {"), createMd5String(md5Map), -1)

	// 将模版文件中的 potrix.cache.js_stamp 替换成文件 potrix.cache.js 的 MD5 值
	for shortPath, bs := range modelFiles {
		modelFiles[shortPath] = bytes.Replace(bs, []byte("potrix.cache.js_stamp"), []byte(md5Map[cacheFile]), -1)
	}

	// 将各 html 文件按模版文件进行输出
	for shortPath, bs := range htmlFiles {
		htmlFiles[shortPath] = parseHtml(bs, func(modelName string) ([]byte, error) {
			return modelFiles[modelName], nil
		})
	}

	// 对各 html、js 文件进行 onload 替换处理
	onloadReg := regexp.MustCompile("<img[^<>]*?src=\"([^<>]*?)\"[^<>]*?/>")
	for _, files := range []map[string][]byte{htmlFiles, jsFiles} {
		for shortPath, bs := range files {
			files[shortPath] = onloadReg.ReplaceAllFunc(bs, func(match []byte) []byte {
				subMatch := onloadReg.FindSubmatch(match)[1]
				//if bytes.HasSuffix(subMatch, []byte(".png")) || bytes.HasSuffix(subMatch, []byte(".jpg")) {
				return bytes.Replace(match, subMatch, []byte("data:image/gif;base64,R0lGODlhAQABAIAAAP///wAAACH5BAEAAAAALAAAAAABAAEAAAICRAEAOw==\" cache=\""+string(subMatch)+"\" onload=\"javascript:P.img(this);"), -1)
				//} else {
				//	return match
				//}
			})
		}
	}

	// 对 html、js、css 进行压缩操作

	// 写入各文件至输出文件夹
	for _, files := range []map[string][]byte{htmlFiles, jsFiles, cssFiles, imgFiles} {
		for shortPath, bs := range files {
			writeToFile(filepath.Join(targetDir, filepath.FromSlash(path.Clean("/"+shortPath))), bs)
		}
	}
}

func createMd5String(md5Map map[string]string) []byte {
	var md5KeySort []string = make([]string, 0)
	for key, _ := range md5Map {
		md5KeySort = append(md5KeySort, key)
	}
	sort.Strings(md5KeySort)

	buffer := bytes.NewBuffer([]byte("(function(w) {\n    var F = {"))
	for i, length := 0, len(md5KeySort); i < length; i++ {
		buffer.WriteString(fmt.Sprintf("\"%s\":\"%s\"", md5KeySort[i], md5Map[md5KeySort[i]]))
		if i < length {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("};")

	return buffer.Bytes()
}

func createDir(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Println("文件夹[", path, "]不存在，")
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Fatalln("文件夹[", path, "]创建时出现异常[", err, "]！")
		}
		log.Println("文件夹[", path, "]创建成功。")
	}
}

func writeToFile(path string, bs []byte) {
	f, err := os.Create(path)
	if err != nil {
		log.Println("复制文件[", path, "]时出现异常：", err)
		return
	}

	defer f.Close()
	_, err = f.Write(bs)
	if err != nil {
		log.Println("写入文件[", path, "]时出现异常：", err)
		return
	}
}

func parseHtml(content []byte, getModel func(modelName string) ([]byte, error)) []byte {
	p, err := goquery.ParseString(string(content))
	if err != nil {
		log.Println("解析Html文件时出现异常:", err)
		return content
	}

	var html map[string]string = make(map[string]string)
	html["model"] = p.Find("html").Attr("model")
	html["res"] = p.Find("html").Attr("res")
	html["title"] = p.Find("title").Html()
	html["style"] = p.Find("style").Html()
	html["body"] = p.Find("body").Html()

	html["style"] = regexp.MustCompile("&gt;").ReplaceAllString(html["style"], ">")

	modelContent, err := getModel(html["model"])
	if err != nil {
		log.Println("获取模板内容时出现异常：", err)
		return content
	}

	newContent := regexp.MustCompile("\\$\\{\\w*\\}").ReplaceAllFunc(modelContent, func(bs []byte) []byte {
		return []byte(html[string(bs[2:len(bs)-1])])
	})

	return newContent
}

func parseHtmlFromFile(content []byte, dir string) []byte {

	return parseHtml(content, func(modelName string) ([]byte, error) {
		model, err := os.Open(filepath.Join(dir, filepath.FromSlash(path.Clean("/"+modelName))))
		if err != nil {
			log.Println("打开模版文件[", modelName, "]时出现异常：", err)
			return []byte(""), nil
		}
		modelContent, err := ioutil.ReadAll(model)
		if err != nil {
			log.Println("读取模版文件[", modelName, "]时出现异常：", err)
		}
		return modelContent, err
	})

}
