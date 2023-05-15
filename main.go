package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/marstion/contract/config"
	"github.com/marstion/contract/img"
	"github.com/marstion/contract/ocr"
	"github.com/marstion/contract/pdf"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	v20181119 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
)

var name string
var nameHelp string = "源PDF文件."

var url string = "https://wzh.renshiren.com/ocr"
var urlHelp string = "ocr API"

var templateFile string
var templateFileHelp string = "模板配置文件"

var staffFile string
var staffFileHelp string = "员工信息台账 txt"

var first int
var firstHelp string = "起始页码. (可留空)"

var last int
var lastHelp string = "结束页码. (可留空)"

var thread int = 2
var threadHelp string = "线程数, 默认2线程, 最大4线程"

func contractSplit() {
	log.Println("Load Config: ", templateFile)
	template, err := config.LoadTemplateConfig(templateFile)
	if err != nil {
		panic(err)
	}
	log.Println("Load staff: ", name)
	contract := Contract{
		name: name,
		pdf:  pdf.Pdf{Name: name},
	}
	log.Println("Load: ", staffFile)
	staffinfoList, err := config.LoadStaffConfig(staffFile)
	if err != nil {
		panic(err)
	}

	// 初始化
	log.Println("init")

	contractInit(&contract)

	// ocr
	log.Println("cnocr")

	contractCNOcr(&contract)
	// 按员工拆分
	log.Println("staff split")

	contractStaffSplit(&contract, &template)

	// 姓名匹配
	log.Println("withname")

	contractWithName(&contract, &staffinfoList, &template)

	// 按姓名拆分
	log.Println("staffsplitname")

	contractStaffSplitName(&contract)
}

func contractStaffSplitName(contract *Contract) {
	for _, conL := range contract.contractL {
		// 处理可能没有匹配到员工的情况
		if conL.name.Id == "" {
			conL.name.Id = "00"
			conL.name.Name = "unknow"
		}

		/*
			切割身份证.
		*/
		front := false
		back := false
		for _, page := range conL.idcard {
			// 获取页面文件名
			pageSrcName := contract.GetTempImage(page - first)
			ocrList := contract.ocrList[page-first]

			for _, result := range ocrList.cnocr.Results {
				var retext string
				var offset []float64
				var reList []string
				// 根据 身份证号码 定位

				// 身份证背面(人像面)
				if !back {
					retext = `\d{10,18}[xX\d]`
					offset = []float64{-230, -500, 200, 100}
					reList = ocr.ReText(retext, result.Text)
					if len(reList) > 0 {
						pageDestName := contract.GetIdCardName(conL.name, conL.idcard, "back")
						img.TrimmingForFile(pageSrcName, pageDestName, result.Position, offset)
						back = true
						continue
					}

					retext = `姓名`
					offset = []float64{-220, -100, 500, 500}
					reList = ocr.ReText(retext, result.Text)
					if len(reList) > 0 {
						pageDestName := contract.GetIdCardName(conL.name, conL.idcard, "back")
						img.TrimmingForFile(pageSrcName, pageDestName, result.Position, offset)
						back = true
						continue
					}
				}

				// 身份证正面(国徽面)
				if !front {
					retext = `有效期限`
					offset = []float64{-270, -400, 300, 100}
					reList = ocr.ReText(retext, result.Text)
					if len(reList) > 0 {
						pageDestName := contract.GetIdCardName(conL.name, conL.idcard, "front")
						img.TrimmingForFile(pageSrcName, pageDestName, result.Position, offset)
						front = true
						continue
					}

					retext = `公安局`
					offset = []float64{-400, -400, 400, 200}
					reList = ocr.ReText(retext, result.Text)
					if len(reList) > 0 {
						pageDestName := contract.GetIdCardName(conL.name, conL.idcard, "front")
						img.TrimmingForFile(pageSrcName, pageDestName, result.Position, offset)
						front = true
						continue
					}

					retext = `中华人民`
					offset = []float64{-570, -200, 350, 500}
					reList = ocr.ReText(retext, result.Text)
					if len(reList) > 0 {
						pageDestName := contract.GetIdCardName(conL.name, conL.idcard, "front")
						img.TrimmingForFile(pageSrcName, pageDestName, result.Position, offset)
						front = true
						continue
					}
				}
			}
		}

		// 切割合同
		pdfName := contract.GetPdfName(conL.name, conL.pages)

		var imageFiles []string
		for _, page := range conL.pages {
			imageFiles = append(imageFiles, contract.GetTempImage(page-first))
		}
		log.Println("SAVE to pdf:", pdfName)
		err := api.ImportImagesFile(imageFiles, pdfName, nil, nil)
		if err != nil {
			panic(err)
		}
	}
}

// 对合同命名
func contractWithName(contract *Contract, staffinfoList *config.Staff, template *config.Template) {
	for index, conL := range contract.contractL {
		// 取身份证首页 或 合同首页, 进行ocr
		var ocrIndex int
		if len(conL.idcard) > 0 {
			ocrIndex = conL.idcard[0]
		} else if len(conL.pages) > 0 {
			ocrIndex = conL.pages[0]
		}
		imgName := contract.GetTempImage(ocrIndex - first)

		var ten ocr.Tencent = ocr.Tencent{
			SecretId:  "AKID2twz01mR58pNYVFft686tU1PeptwnW5S",
			SecretKey: "DKDsufBUaO4atxmMyPefRHzUyIWaR72n",
		}
		var rsp *v20181119.GeneralHandwritingOCRResponseParams
		var err error

		rsp, err = ten.GeneralHandwritingOCR(imgName, nil)

		if err != nil {
			panic(err)
		}
		var staffNameKey string = "^[\u4e00-\u9fa5]{2,3}$"
		var staffName []string
		var staffIdCardKey string = `\d{11,18}[xX\d]`
		var staffIdCard []string
		var staffPhoneKey string = `\d{10,11}`
		var staffPhone []string

		for _, textDetections := range rsp.TextDetections {
			// 匹配姓名
			reRsp := ocr.ReText(staffNameKey, strings.Replace(*textDetections.DetectedText, " ", "", -1))
			if len(reRsp) > 0 {
				staffName = append(staffName, reRsp...)
				continue
			}
			// 匹配身份证
			reRsp = ocr.ReText(staffIdCardKey, strings.Replace(*textDetections.DetectedText, " ", "", -1))
			if len(reRsp) > 0 {
				staffIdCard = append(staffIdCard, reRsp...)
				continue
			}

			// 匹配电话号码
			reRsp = ocr.ReText(staffPhoneKey, strings.Replace(*textDetections.DetectedText, " ", "", -1))
			if len(reRsp) > 0 {
				staffPhone = append(staffPhone, reRsp...)
				continue
			}
		}
		fmt.Println(staffName, staffIdCard, staffPhone)

		// 按身份证匹配员工
		if name, ok := staffinfoList.ComparisonIdCard(staffIdCard); ok {
			contract.contractL[index].name = name
			continue
		}

		// 按姓名匹配
		if name, ok := staffinfoList.ComparisonName(staffName); ok {
			contract.contractL[index].name = name
			continue
		}

		// 按 手机号
		if name, ok := staffinfoList.ComparisonPhone(staffPhone); ok {
			contract.contractL[index].name = name
			continue
		}
	}
}

func contractStaffSplit(contract *Contract, template *config.Template) {
	IsIdCard := true
	conList := contractList{}
	for _, ocrObj := range contract.ocrList {
		switch {
		case IsIdCard && !ocrObj.isContractStart(template.CNOCR.ContractStart): // 前页为身份证页面且没有匹配到合同开始页
			log.Printf("page: %d isIdCard\n", ocrObj.index)
			conList.idcard = append(conList.idcard, ocrObj.index)

		case IsIdCard && ocrObj.isContractStart(template.CNOCR.ContractStart): // 前页为身份证页面且匹配到合同开始页
			IsIdCard = false
			log.Printf("page: %d isContractStart\n", ocrObj.index)
			conList.pages = append(conList.pages, ocrObj.index)

		case !IsIdCard && ocrObj.isContractEnd(template.CNOCR.ContractEnd): // 前页为合同页且匹配到合同结束页
			log.Printf("page: %d isContractEnd\n", ocrObj.index)
			conList.pages = append(conList.pages, ocrObj.index)
			IsIdCard = true

			contract.contractL = append(contract.contractL, conList)
			conList = contractList{}
		case !IsIdCard: // 前页为合同页， 且没匹配到合同页结尾
			log.Printf("page: %d isContract\n", ocrObj.index)
			conList.pages = append(conList.pages, ocrObj.index)
		default:
			log.Fatalf("page: %d unknow!\n", ocrObj.index)
		}
	}
	log.Printf("contractStaffSplit: %#v\n", contract)
}

func contractInit(contract *Contract) {
	err := contract.HashCrc32()
	if err != nil {
		panic(err)
	}
	log.Println("crc32: ", contract.crc32)

	pages, err := contract.pdf.Pages()
	if err != nil {
		panic(err)
	}
	log.Println("Pages: ", pages)

	if last <= 0 || last > pages {
		last = pages
	}
	if first <= 0 {
		first = 1
	} else if first > pages {
		first = pages
	}
	if thread > 13 {
		thread = 13
	}
	// 如果 last 小于 first, 退出
	if last < first {
		panic("last > first.")
	}
}

func contractCNOcr(contract *Contract) {
	// pdf to png
	namePrefix := contract.GetTempImagePrefix()
	log.Println("PDF to IMG.")
	contract.pdf.ToJPG(namePrefix, map[string]interface{}{"resolution": 200, "last": last, "first": first})

	// pdf cnocr
	var wg sync.WaitGroup
	ch := make(chan int, thread)

	log.Println("IMG to CNOCR.")
	for index := first; index <= last; index++ {
		ch <- index
		wg.Add(1)

		// 避免出错, 循环, 成功再退出
		go func(index int) {
			defer wg.Done()
			imgName := contract.GetTempImage(index - first)
			log.Println("OCR: ", imgName)
			for {
				o, err := ocr.OCRImg(imgName, url)
				if err != nil || o == nil {
					time.Sleep(5 * time.Second)
					continue
				}
				contract.ocrList = append(contract.ocrList, &OcrIndex{index: index, cnocr: o})
				<-ch
				return
			}
		}(index)
	}
	wg.Wait()
	// 对结果排序, 因多线程返回结果时间不一致, 存入结果后乱序
	sort.Sort(contract.ocrList)
}

// func contract
func main() {
	flag.StringVar(&name, "p", "", nameHelp)
	flag.StringVar(&url, "u", url, urlHelp)
	// flag.Float64Var(&threshold, "t", threshold, thresholdHelp)
	flag.StringVar(&staffFile, "s", "", staffFileHelp)
	flag.IntVar(&first, "f", first, firstHelp)
	flag.IntVar(&last, "l", last, lastHelp)
	flag.IntVar(&thread, "P", thread, threadHelp)
	flag.StringVar(&templateFile, "t", "", templateFileHelp)
	flag.Parse()

	contractSplit()
}
