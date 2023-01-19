package ocr

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ocr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
)

// 腾讯云的 secretId 和 secretKey
const secretId = ""
const secretKey = ""
const region = "ap-beijing"
const start = "我单位开票信息如下"
const end = "我单位申请开具以下到期汇票的增值税普票"

var credential = common.NewCredential(secretId, secretKey)
var client, _ = ocr.NewClient(credential, region, profile.NewClientProfile())

type ResponseText struct {
	Response TextDetections `json:"Response"`
}

type TextDetections struct {
	TextDetections []ContentEntry
	Language       string `json:"Language"`
}

type ContentEntry struct {
	DetectedText string
}

func imageBase64(targetUrl string) string {
	imgByte, err := os.ReadFile(targetUrl)
	if err != nil {
		fmt.Println(err)
	}
	res := base64.StdEncoding.EncodeToString(imgByte)
	return res
}

func parse(targetUrl string) *ResponseText {

	request := ocr.NewGeneralBasicOCRRequest()
	request.ImageBase64 = common.StringPtr(imageBase64(targetUrl))

	isPdf, _ := regexp.MatchString(".[p|P][d|D][f|F]$", targetUrl)

	if isPdf {
		request.IsPdf = common.BoolPtr(true)
	}

	response, err := client.GeneralBasicOCR(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return nil
	}
	if err != nil {
		panic(err)
	}

	responseText := &ResponseText{}
	json.Unmarshal([]byte(response.ToJsonString()), responseText)
	return responseText
}

func GetTaxNumber(targetUrl string) []string {
	responseText := parse(targetUrl)
	contentEntrys := responseText.Response.TextDetections
	flag := false
	resultArray := make([]string, 0)
	for _, contentEntry := range contentEntrys {
		detectedText := contentEntry.DetectedText
		if !flag && strings.Contains(detectedText, start) {
			flag = true
			continue
		}
		if flag && strings.Contains(detectedText, end) {
			flag = false
			break
		}

		if flag {
			resultArray = append(resultArray, detectedText)
		}
	}
	return resultArray
}

// func main() {
// 	taxNum := GetTaxNumber("/Users/liws/Downloads/pdf和图片解析/117261673847775_.pic.jpg")
// 	fmt.Println(taxNum)
// }
