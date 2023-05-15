package ocr

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	Terrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ocr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
)

func EncodePath(path string) (output []byte, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	output = make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(output, data)
	return
}

type Tencent struct {
	SecretId  string
	SecretKey string
}

func (t Tencent) AdvertiseOCRbackup(filePath string) (rsp *ocr.AdvertiseOCRResponseParams, err error) {
	// 实例化一个认证对象，入参需要传入腾讯云账户 SecretId 和 SecretKey，此处还需注意密钥对的保密
	// 代码泄露可能会导致 SecretId 和 SecretKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考，建议采用更安全的方式来使用密钥，请参见：https://cloud.tencent.com/document/product/1278/85305
	// 密钥可前往官网控制台 https://console.cloud.tencent.com/cam/capi 进行获取
	credential := common.NewCredential(
		t.SecretId,
		t.SecretKey,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "ocr.ap-beijing.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, err := ocr.NewClient(credential, "ap-beijing", cpf)
	if err != nil {
		return
	}
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := ocr.NewAdvertiseOCRRequest()

	imgBS64, err := EncodePath(filePath)
	if err != nil {
		return
	}
	request.ImageBase64 = common.StringPtr(string(imgBS64))

	// 返回的resp是一个AdvertiseOCRResponse的实例，与请求对象对应
	response, err := client.AdvertiseOCR(request)

	if _, ok := err.(*Terrors.TencentCloudSDKError); ok {
		err = fmt.Errorf("an api error has returned: %s", err)
		return
	}
	if err != nil {
		return
	}

	rsp = response.Response
	return
}

func (t Tencent) GeneralHandwritingOCR(filePath string, parameters map[string]string) (rsp *ocr.GeneralHandwritingOCRResponseParams, err error) {
	// 实例化一个认证对象，入参需要传入腾讯云账户 SecretId 和 SecretKey，此处还需注意密钥对的保密
	// 代码泄露可能会导致 SecretId 和 SecretKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考，建议采用更安全的方式来使用密钥，请参见：https://cloud.tencent.com/document/product/1278/85305
	// 密钥可前往官网控制台 https://console.cloud.tencent.com/cam/capi 进行获取
	credential := common.NewCredential(
		t.SecretId,
		t.SecretKey,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "ocr.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := ocr.NewClient(credential, "ap-beijing", cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := ocr.NewGeneralHandwritingOCRRequest()

	imgBS64, err := EncodePath(filePath)
	if err != nil {
		return
	}
	request.ImageBase64 = common.StringPtr(string(imgBS64))

	// 判断是否只返回手写体
	if _, ok := parameters["only_hw"]; ok {
		request.Scene = common.StringPtr("only_hw")
	}

	// 返回的resp是一个GeneralHandwritingOCRResponse的实例，与请求对象对应
	// 循环， 避免出现： [TencentCloudSDKError] Code=FailedOperation.UnKnowError, Message=内部错误, RequestId=8b42ce0c-5f51-412e-8dff-9c7f3f4dfcfe
	var response *ocr.GeneralHandwritingOCRResponse
	for i := 1; i <= 5; i++ { // 五次内成功即返回， 不成功则最后返回错误
		response, err = client.GeneralHandwritingOCR(request)
		if err == nil {
			break
		}
	}
	if _, ok := err.(*Terrors.TencentCloudSDKError); ok {
		err = fmt.Errorf("an api error has returned: %s", err)
		return
	}
	if err != nil {
		panic(err)
	}
	// 输出json格式的字符串回包
	rsp = response.Response
	return
}
