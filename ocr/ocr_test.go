package ocr

import (
	"log"
	"testing"
)

func TestRequest(t *testing.T) {
	filename := "test.jpg"
	ocrURL := "https://di29.rushi.pw/ocr"
	ocr, err := OCRImg(filename, ocrURL)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%#v", ocr)
}
