package pdf

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/marstion/contract/cli"
)

type Pdf struct {
	Name string // pdf 文件名
}

func (p Pdf) Pages() (pages int, err error) {
	// pdfinfo hebeibingguan.pdf
	infoStr, err := cli.Cli("pdfinfo", []string{p.Name})
	if err != nil {
		return
	}
	reg := regexp.MustCompile(`Pages:\s+\d+`)
	reList := reg.FindStringSubmatch(infoStr)
	if len(reList) > 0 {
		lineList := strings.Split(reList[0], ` `)
		pages, err = strconv.Atoi(lineList[len(lineList)-1])
	} else {
		err = errors.New("not fount pages")
	}
	return
}

/*
PDF转换成图片, 生成图片规则 ${prefix}-%04d.jpg
可选参数:
parameters:

	first: 开始位置
	last: 结束位置
*/
func (p Pdf) ToJPG(prefix string, parameters map[string]interface{}) (err error) {
	// pdfimages -j -raw hebeibingguan.pdf prefix
	params := []string{}
	// first page to print
	if value, ok := parameters["first"]; ok {
		params = append(params, "-f")
		params = append(params, fmt.Sprint(value))
	}
	// last page to print
	if value, ok := parameters["last"]; ok {
		params = append(params, "-l")
		params = append(params, fmt.Sprint(value))
	}
	
	params = append(params, "-j")
	// <PDF-file> <PNG-root>
	params = append(params, p.Name)
	params = append(params, prefix)
	log.Printf("pdftopng, %#v", params)
	_, err = cli.Cli("pdfimages", params)
	return
}

/*
PDF转换成图片, 生成图片规则 ${prefix}-%06d.jpg
可选参数:
parameters:

	first: 开始位置
	last: 结束位置
*/
func (p Pdf) ToPNG(prefix string, parameters map[string]interface{}) (err error) {
	// pdftopng -raw hebeibingguan.pdf prefix
	params := []string{}
	// resolution, in DPI (default is 150)
	if value, ok := parameters["resolution"]; ok {
		params = append(params, "-r")
		params = append(params, fmt.Sprint(value))
	}
	// first page to print
	if value, ok := parameters["first"]; ok {
		params = append(params, "-f")
		params = append(params, fmt.Sprint(value))
	}
	// last page to print
	if value, ok := parameters["last"]; ok {
		params = append(params, "-l")
		params = append(params, fmt.Sprint(value))
	}
	// <PDF-file> <PNG-root>
	params = append(params, p.Name)
	params = append(params, prefix)
	log.Printf("pdftopng, %#v", params)
	_, err = cli.Cli("pdftopng", params)
	return
}
