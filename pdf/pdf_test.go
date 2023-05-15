package pdf

import "testing"

var pdf Pdf = Pdf{Name: "../Scan.pdf"}

func TestPages(t *testing.T) {
	pages, err := pdf.Pages()
	if err != nil {
		t.Error(err)
	}
	t.Log(pages)
}

func TestToPng(t *testing.T) {
	// 测试切割所有
	err := pdf.ToPng("all", map[string]interface{}{"resolution": 300})
	if err != nil {
		t.Error("all: ", err)
	}
	err = pdf.ToPng("first12", map[string]interface{}{"first": 12})
	if err != nil {
		t.Error("first12: ", err)
	}
	err = pdf.ToPng("last3", map[string]interface{}{"last": 3})
	if err != nil {
		t.Error("first12: ", err)
	}
	err = pdf.ToPng("qujian4-11", map[string]interface{}{"first": 4, "last": 11})
	if err != nil {
		t.Error("qujian4-1: ", err)
	}
	err = pdf.ToPng("chaoguo32", map[string]interface{}{"first": 4, "last": 50})
	if err != nil {
		t.Error("chaoguo32: ", err)
	}
}
