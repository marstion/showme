package ocr

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type CNOcr struct {
	StatusCode int `json:"status_code"`
	Results    []CNOcrResult
}

type CNOcrResult struct {
	Text     string
	Score    float64
	Position [][]float64
}

func ReText(retext string, text string) (reList []string) {
	reg := regexp.MustCompile(retext)
	textL := strings.FieldsFunc(string(text), func(r rune) bool {
		return r == ':' || r == '：' || r == ')' || r == '）'
	})
	for _, t := range textL {
		reList = append(reList, reg.FindStringSubmatch(t)...)
	}
	return
}

// 图片文字提取
func OCRImg(name, url string) (o *CNOcr, err error) {
	form := new(bytes.Buffer)
	writer := multipart.NewWriter(form)
	fw, err := writer.CreateFormFile("image", filepath.Base(name))
	if err != nil {
		return
	}
	fd, err := os.Open(name)
	if err != nil {
		return
	}
	defer fd.Close()
	_, err = io.Copy(fw, fd)
	if err != nil {
		return
	}
	writer.Close()
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, form)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	json.Unmarshal(body, &o)
	return
}
