package aliyun

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

type CommonResponse struct {
	RequestId string `json:"RequestId"`
	Data      string `json:"Data"`
}

type GetEmotionByTextResponse struct {
	Result struct {
		PositiveProb float64 `json:"positive_prob"`
		Sentiment    string  `json:"sentiment"`
		NeutralProb  float64 `json:"neutral_prob"`
		NegativeProb float64 `json:"negative_prob"`
	} `json:"result"`
	Success  bool   `json:"success"`
	TracerId string `json:"tracerId"`
}

// GetEmotionByText 从文本获取情感状态
func (y *AliYun) GetEmotionByText(text string) (ret GetEmotionByTextResponse, err error) {
	request := requests.NewCommonRequest()
	request.Domain = "alinlp.cn-hangzhou.aliyuncs.com"
	request.Version = "2020-06-29"
	request.ApiName = "GetSaChGeneral"
	request.QueryParams["ServiceCode"] = "alinlp"
	request.QueryParams["Text"] = text
	request.TransToAcsRequest()

	commonRequest, err := y.Client.ProcessCommonRequest(request)
	if err != nil {
		return
	}
	bytes := commonRequest.GetHttpContentBytes()
	var response CommonResponse
	_ = json.Unmarshal(bytes, &response)
	_ = json.Unmarshal([]byte(response.Data), &ret)
	return
}
