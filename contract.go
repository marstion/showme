package main

import (
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/marstion/contract/config"
	"github.com/marstion/contract/ocr"
	"github.com/marstion/contract/pdf"
)

const (
	mainDir    string = "data"  // 转换后存放主目录
	filePrefix string = "image" // 文件前缀
	fileSuffix string = "jpg"   // 文件后缀
)

type Contract struct {
	name string // 合同PDF扫描件

	pages int    // 合同总页数
	crc32 string // 合同PDF文件HASH

	pdf     pdf.Pdf // pdf 对象
	ocrList ocrList // ocr结果

	contractL []contractList // 合同按员工分组
}

type contractList struct {
	name   config.StaffInfo // 员工身份信息
	idcard []int            // 身份证页码
	pages  []int            // 合同页码
}

type ocrList []*OcrIndex
type OcrIndex struct {
	index int        // pdf 页码
	cnocr *ocr.CNOcr // ocr 结果
}

func (oi OcrIndex) isContractStart(key string) bool {
	for _, ocrLine := range oi.cnocr.Results {
		if len(ocr.ReText(key, ocrLine.Text)) > 0 {
			return true
		}
	}
	return false
}

func (oi OcrIndex) isContractEnd(key string) bool {
	for _, ocrLine := range oi.cnocr.Results {
		if len(ocr.ReText(key, ocrLine.Text)) > 0 {
			return true
		}
	}
	return false
}

// 实现排序
func (ol ocrList) Len() int           { return len(ol) }
func (ol ocrList) Less(i, j int) bool { return ol[i].index < ol[j].index }
func (ol ocrList) Swap(i, j int)      { ol[i], ol[j] = ol[j], ol[i] }

func NewContract(name string) (contract *Contract, err error) {
	contract = &Contract{
		name: name,
		pdf:  pdf.Pdf{Name: name},
	}

	err = contract.HashCrc32()
	if err != nil {
		return
	}
	contract.pages, err = contract.pdf.Pages()
	if err != nil {
		return
	}
	return
}

// 计算PDF文件HASH
func (c *Contract) HashCrc32() (err error) {
	os.Open(c.name)

	fileObj, err := ioutil.ReadFile(c.name)
	if err != nil {
		panic(err)
	}
	c.crc32 = strconv.FormatUint(uint64(crc32.ChecksumIEEE(fileObj)), 16)
	return
}

// 获取图片存放路劲和文件前缀， 自动创建目录
func (c Contract) GetTempImagePrefix() (name string) {
	dir := filepath.Join(mainDir, "temp", c.crc32)
	os.MkdirAll(dir, 0755)
	name = filepath.Join(dir, filePrefix)
	return
}

// 获取临时图片存放路劲和文件全称
func (c Contract) GetTempImage(index int) (name string) {
	name = fmt.Sprintf("%s-%04d.jpg", c.GetTempImagePrefix(), index)
	return
}

// 返回PDF正式文件存放目录
func (c Contract) GetPdfName(staff config.StaffInfo, pages []int) (name string) {
	dir := filepath.Join(mainDir, "complate", c.crc32)
	os.MkdirAll(dir, 0755)

	name = filepath.Join(dir, fmt.Sprintf("%s-%s-%s-%03d-%03d-%03d.pdf",
		staff.Id, staff.Name, staff.Idcard, len(pages), pages[0], pages[len(pages)-1]))
	return
}

// 返回身份证正反面文件存放路劲
func (c Contract) GetIdCardName(staff config.StaffInfo, pages []int, frontOback string) (name string) {
	dir := filepath.Join(mainDir, "complate", c.crc32)
	os.MkdirAll(dir, 0755)
	name = filepath.Join(dir,
		fmt.Sprintf("%s-%s-%s-%s-%03d-%03d-%03d.jpg",
			staff.Id, staff.Name, staff.Idcard, frontOback, len(pages), pages[0], pages[len(pages)-1]),
	)
	return
}
