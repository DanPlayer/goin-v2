package wechat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TemplateSendParam struct {
	ToUser      string      `json:"touser"`      // openid
	TemplateId  string      `json:"template_id"` // 模板id
	Url         string      `json:"url"`         // 跳转链接
	MiniProgram MiniProgram `json:"miniprogram"` // 跳转小程序
	Data        interface{} `json:"data"`        // 模板消息数据（根据实际模板）
}

type MiniProgram struct {
	AppID    string `json:"appid"`
	PagePath string `json:"pagepath"`
}

// SendTemplateMsg 发送模板消息
// param params TemplateSendParam 模板消息
func (w *Wechat) SendTemplateMsg(params TemplateSendParam) (re ErrorCode, err error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%v", w.accessToken.token)
	paramsStr, _ := json.Marshal(params)
	response, err := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer([]byte(paramsStr)))
	if err != nil {
		return re, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return re, errors.New("发送模板消息失败，请检查网络是否异常")
	}
	data, _ := ioutil.ReadAll(response.Body)
	_ = json.Unmarshal(data, &re)
	if re.ErrCode != 0 {
		return re, errors.New("发送模板消息失败, errmsg: " + re.ErrMsg)
	}
	return
}
