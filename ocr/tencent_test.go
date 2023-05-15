package ocr

import (
	"testing"
)

func TestAdvertiseOCR(t *testing.T) {
	var ten Tencent = Tencent{
		SecretId:  "xxxxxx",
		SecretKey: "xxxxx",
	}
	rsp, err := ten.GeneralHandwritingOCR("../data/temp/9d1903d7/image-000140.png", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v\n", rsp)
}
