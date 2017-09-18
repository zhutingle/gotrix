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

var b64 *base64.Encoding = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")

var StaticResource map[string]string = make(map[string]string)

// 将静态资源文件写入 Gotrix 代码中，方便其它项目自动引用
func packageToGotrix(staticDir string, fileName string) {

	buffer := bytes.NewBuffer(make([]byte, 0))
	buffer.Write([]byte("package server\n"))
	buffer.Write([]byte("func init() {\n"))

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
			return nil
		}
		fileContent, err := ioutil.ReadFile(longPath)
		if err != nil {
			log.Println("读取文件[", shortPath, "]时出错：", err)
			return err
		}
		buffer.Write([]byte("	StaticResource[\"" + shortPath + "\"] = \""))
		buffer.WriteString(b64.EncodeToString(fileContent))
		buffer.Write([]byte("\"\n"))
		return nil
	})

	buffer.Write([]byte("}\n"))

	writeToFile(path.Clean(staticDir+"/../../server/"+fileName), buffer.Bytes())

}

// 对静态文件夹进行打包处理，输出到输出文件夹
func packageTarget(staticDir string, targetDir string) {

	// 判断输出文件夹是否存在，不存在则创建
	createDir(targetDir)

	var htmlFiles map[string][]byte = make(map[string][]byte)
	var modelFiles map[string][]byte = make(map[string][]byte)
	var cssFiles map[string][]byte = make(map[string][]byte)
	var jsFiles map[string][]byte = make(map[string][]byte)
	var imgFiles map[string][]byte = make(map[string][]byte)
	var imgCacheFiles map[string][]byte = make(map[string][]byte)
	var ttfFiles map[string][]byte = make(map[string][]byte)
	var err error
	// 从代码中读取静态资源文件
	for shortPath, fileContentB64 := range StaticResource {
		if strings.HasSuffix(shortPath, ".html") {
			// 目录中含有 model 就认为它是模版文件
			if strings.Contains(shortPath, "model") {
				modelFiles[shortPath], err = b64.DecodeString(fileContentB64)
			} else {
				htmlFiles[shortPath], err = b64.DecodeString(fileContentB64)
			}
		} else if strings.HasSuffix(shortPath, ".css") {
			cssFiles[shortPath], err = b64.DecodeString(fileContentB64)
		} else if strings.HasSuffix(shortPath, ".js") {
			jsFiles[shortPath], err = b64.DecodeString(fileContentB64)
		} else if strings.HasSuffix(shortPath, ".png") || strings.HasSuffix(shortPath, ".jpg") || strings.HasSuffix(shortPath, ".jpeg") || strings.HasSuffix(shortPath, ".gif") {
			imgFiles[shortPath], err = b64.DecodeString(fileContentB64)
		} else if strings.HasSuffix(shortPath, ".ttf") {
			ttfFiles[shortPath], err = b64.DecodeString(fileContentB64)
		}
		if err != nil {
			log.Println("解码静态资源文件[", shortPath, "]时出错：", err)
		}
	}

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
		} else if strings.HasSuffix(shortPath, ".png") || strings.HasSuffix(shortPath, ".jpg") || strings.HasSuffix(shortPath, ".jpeg") || strings.HasSuffix(shortPath, ".gif") {
			imgFiles[shortPath], err = ioutil.ReadFile(longPath)
		} else if strings.HasSuffix(shortPath, ".ttf") {
			ttfFiles[shortPath], err = ioutil.ReadFile(longPath)
		}
		if err != nil {
			log.Println("文件[", longPath, "]读取异常！")
		}
		return nil
	})

	// 对 img 文件进行 base64 处理，并在文件名后面加上 .cache
	for shortPath, bs := range imgFiles {
		imgCacheFiles[shortPath+".cache"] = []byte(b64.EncodeToString(bs))
	}

	// 对各 html、js 文件进行 onload 替换处理
	onloadReg := regexp.MustCompile("<img[^<>]*?[^:]src=\"([^+<>!\\s]*?)\"[^<>]*?>")
	for _, files := range []map[string][]byte{htmlFiles, jsFiles} {
		for shortPath, bs := range files {
			files[shortPath] = onloadReg.ReplaceAllFunc(bs, func(match []byte) []byte {
				subMatch := onloadReg.FindSubmatch(match)[1]
				if bytes.HasSuffix(match, []byte("/>")) {
					return bytes.Replace(match, subMatch, []byte("data:image/gif;base64,R0lGODlhAQABAIAAAP///wAAACH5BAEAAAAALAAAAAABAAEAAAICRAEAOw==\" cache=\""+string(subMatch)+"\" onload=\"javascript:P.img(this);"), -1)
				} else {
					imgCacheFiles[string(subMatch)] = imgFiles[string(subMatch)]
					return match
				}
			})
		}
	}

	// 对 html、js、css 进行压缩操作

	// 计算 html、css、js、img 文件的 MD5 值，不包括 js/gotrix.cache.js
	var cacheFile string = "js/gotrix.cache.js"
	var md5Map map[string]string = make(map[string]string)
	for _, files := range [](map[string][]byte){htmlFiles, cssFiles, jsFiles, imgFiles} {
		for shortPath, bs := range files {
			md5Map[shortPath] = b64.EncodeToString(global.Md5(bs))
		}
	}
	delete(md5Map, cacheFile)

	// 将各文件（不含 gotrix.cache.js ）的 MD5 值合并成字符串，并写入 gotrix.cache.js 再计算 gotrix.cache.js 的 MD5 值
	md5Map[cacheFile] = b64.EncodeToString(global.Md5(bytes.Replace(jsFiles[cacheFile], []byte("(function (w) {"), createMd5String(md5Map), -1)))

	// 将各文件（包含 gotrix.cache.js）的 MD5 值合并成字符串，并写入 gotrix.cache.js
	jsFiles[cacheFile] = bytes.Replace(jsFiles[cacheFile], []byte("(function (w) {"), createMd5String(md5Map), -1)

	// 将模版文件中的 gotrix.cache.js_stamp 替换成文件 gotrix.cache.js 的 MD5 值
	for shortPath, bs := range modelFiles {
		modelFiles[shortPath] = bytes.Replace(bs, []byte("gotrix.cache.js_stamp"), []byte(md5Map[cacheFile]), -1)
	}

	// 将各 html 文件按模版文件进行输出
	for shortPath, bs := range htmlFiles {
		htmlFiles[shortPath] = parseHtml(shortPath, bs, func(modelName string) ([]byte, error) {
			return modelFiles[modelName], nil
		})
	}

	// 写入各文件至输出文件夹
	for _, files := range []map[string][]byte{htmlFiles, jsFiles, cssFiles, imgCacheFiles, ttfFiles} {
		for shortPath, bs := range files {
			writeToFile(filepath.Join(targetDir, filepath.FromSlash(path.Clean("/"+shortPath))), bs)
		}
	}
}

func createMd5String(md5Map map[string]string) []byte {
	var md5KeySort []string = make([]string, 0)
	for key := range md5Map {
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

	// 写入文件时先创建各文件夹
	createDir(filepath.Dir(path))

	// 写入文件内容
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

func parseHtml(fileName string, content []byte, getModel func(modelName string) ([]byte, error)) []byte {

	content = regexp.MustCompile("<func[\\w\\W]*?</func>").ReplaceAllFunc(content, func(match []byte) []byte {
		bs := bytes.NewBuffer([]byte(""))
		bs.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\" ?><funcs>")
		bs.Write(match)
		bs.WriteString("</funcs>")
		gotrixHandler.ReadXmlBytes(bs.Bytes())
		return []byte("")
	})

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
	if len(html["model"]) == 0 {
		log.Println("解析[", fileName, "]时，未指定模板文件。")
		return content
	}
	modelContent, err := getModel(html["model"])
	if err != nil {
		log.Println("解析[", fileName, "]时，获取模板[", html["model"], "]时出现异常：", err)
		return content
	}

	newContent := regexp.MustCompile("\\$\\{\\w*\\}").ReplaceAllFunc(modelContent, func(bs []byte) []byte {
		return []byte(html[string(bs[2:len(bs)-1])])
	})

	return newContent
}

func parseHtmlFromFile(fileName string, content []byte, dir string) []byte {

	return parseHtml(fileName, content, func(modelName string) ([]byte, error) {

		model, err := os.Open(filepath.Join(dir, filepath.FromSlash(path.Clean("/"+modelName))))
		if err != nil {
			fileResource := StaticResource[modelName]
			if fileResource != "" {
				return b64.DecodeString(fileResource)
			}
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

func parseFromB64(fileName string) []byte {

	content := StaticResource[fileName]
	if len(content) == 0 {
		return nil
	}

	bs, err := b64.DecodeString(content)
	if err != nil {
		log.Println("文件进行BASE64解码时出错")
	}

	if strings.HasSuffix(fileName, "html") {
		return parseHtml(fileName, bs, func(ModelName string) ([]byte, error) {
			return b64.DecodeString(StaticResource[ModelName])
		})
	}

	return bs

}
