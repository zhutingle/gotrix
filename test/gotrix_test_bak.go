package test

import (
	"fmt"
	"io/ioutil"
	//	"log"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/zhutingle/gotrix/global"
	"github.com/zhutingle/gotrix/handler"
)

func insertDicti(gotrixHandler handler.SimpleHandler, t *testing.T, dictid float64, value float64, text string, des string) {
	checkedParams := &global.CheckedParams{Func: 100, V: make(map[string]interface{})}
	checkedParams.V["dictid"] = dictid
	checkedParams.V["value"] = value
	checkedParams.V["text"] = text
	checkedParams.V["des"] = des
	//response, err := gotrixHandler.GetPass([]byte("a65d8f37015e9aff1413855d1abaf1779e9005d1c19fbbbf2da447d0cf1dc1f7"))
	//	log.Println(response)
	fmt.Println(checkedParams.V)

	response, err := gotrixHandler.Handle(checkedParams)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(response)
	}
}

func TestSimpleHandler(t *testing.T) {

	gotrixHandler := handler.SimpleHandler{}
	gotrixHandler.Init()

	str := read3("C:\\Users\\Zhu\\Desktop\\p.txt")

	reg := regexp.MustCompile("[\u4e00-\u9fa5]*[区省市]")

	strs := reg.FindAllString(str, -1)
	var province string
	var indexMap map[int]string = make(map[int]string)
	var dataMap map[string][]string = make(map[string][]string)
	for i := 0; i < len(strs); i++ {
		if strs[i] == "副省级城市" || strs[i] == "省直辖县级市" || strs[i] == "地级市" || strs[i] == "县级市" || strs[i] == "自治区直辖县级市" {
			continue
		}
		if strings.HasSuffix(strs[i], "省") || strings.HasSuffix(strs[i], "区") {
			province = strs[i]
			dataMap[strs[i]] = make([]string, 0)
			indexMap[len(dataMap)+4] = strs[i]
		} else {
			dataMap[province] = append(dataMap[province], strs[i])
		}
	}
	var index float64 = 1
	for i := 0; i < len(indexMap); i++ {
		fmt.Printf("3 %d %s null\n", i+5, indexMap[i+5])
		insertDicti(gotrixHandler, t, 3, float64(i+5), indexMap[i+5], "")
		province = indexMap[i+5]
		array := dataMap[province]
		for j := 0; j < len(array); j++ {
			fmt.Printf("4 %d %s %d\n", j+1, array[j], i+5)
			insertDicti(gotrixHandler, t, 4, index, array[j], fmt.Sprintf("%d", i+5))
			index++
		}
	}

}

func read3(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}
