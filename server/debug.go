package server

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type DevDir string

func (d DevDir) Open(name string) (http.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) ||
		strings.Contains(name, "\x00") {
		return nil, errors.New("http: invalid character in file path")
	}
	dir := string(d)
	if dir == "" {
		dir = "."
	}
	f, err := os.Open(filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name))))
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(name, "html") {
		return newDevFile(f, dir), nil
	} else {
		return f, nil
	}
}

type DevFile struct {
	file     http.File
	fileInfo os.FileInfo
	bs       []byte
	cur      int64
}

func (f *DevFile) Close() error {
	return f.file.Close()
}

func (f *DevFile) Read(p []byte) (n int, err error) {
	cur := f.cur
	for i := 0; i < len(p) && f.cur < int64(len(f.bs)); {
		p[i] = f.bs[f.cur]
		i++
		f.cur++
	}
	n = int(f.cur - cur)
	if n == 0 {
		return n, io.EOF
	} else {
		return n, nil
	}
}

func (f *DevFile) Seek(offset int64, whence int) (int64, error) {
	if whence == io.SeekStart {
		f.cur = offset
	} else if whence == io.SeekCurrent {
		f.cur = f.cur + offset
	} else if whence == io.SeekEnd {
		f.cur = int64(len(f.bs)) + offset
	}

	if f.cur < 0 || f.cur > int64(len(f.bs)) {
		return 0, errors.New("超出边界")
	} else {
		return f.cur, nil
	}
}

func (f *DevFile) Readdir(count int) ([]os.FileInfo, error) {
	return f.file.Readdir(count)
}

func (f *DevFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *DevFile) Name() string {
	return f.fileInfo.Name()
}
func (f *DevFile) Size() int64 {
	return int64(len(f.bs))
}
func (f *DevFile) Mode() os.FileMode {
	return f.fileInfo.Mode()
}
func (f *DevFile) ModTime() time.Time {
	return f.fileInfo.ModTime()
}
func (f *DevFile) IsDir() bool {
	return f.fileInfo.IsDir()
}
func (f *DevFile) Sys() interface{} {
	return f.fileInfo.Sys()
}

func newDevFile(f http.File, dir string) *DevFile {
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		log.Print("读取文件时出现异常：", err)
	}

	fileInfo, err := f.Stat()
	if err != nil {
		log.Print("获取文件信息时出现异常：", err)
	}

	return &DevFile{file: f, fileInfo: fileInfo, bs: parseHtmlFromFile(bs, dir), cur: 0}
}
