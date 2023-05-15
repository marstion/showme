package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Template struct {
	Note       string
	CNOCR      TemplateCNOCR
	TencentOCR TemplateTencentOCR
}

type TemplateCNOCR struct {
	ContractStart string // 合同开始页, 正则表达式
	ContractEnd   string // 合同结束页, 正则表达式

}

type TemplateTencentOCR struct {
	GeneralHandwritingOCR TemplateTencentOCRGeneralHandwritingOCR
}

type TemplateTencentOCRGeneralHandwritingOCR struct {
	Onlyhw bool `json:"only_hw"` // 开启进识别手写体, 忽略印刷体
}

func LoadTemplateConfig(name string) (c Template, err error) {
	jsonFile, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &c)

	return
}
