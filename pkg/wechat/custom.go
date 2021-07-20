package wechat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Custom 客服属性
type Custom struct {
	KfAccount    string `json:"kf_account"`
	KfNick       string `json:"kf_nick"`
	KfId         string `json:"kf_id"`
	KfHeadimgurl string `json:"kf_headimgurl"`
}

// CustomParam 客服新增属性
type CustomParam struct {
	KfAccount string `json:"kf_account"`
	NickName  string `json:"nick_name"`
	Password  string `json:"password"`
}

type CustomText struct {
	ToUser  string            `json:"touser"`
	MsgType string            `json:"msgtype"`
	Text    CustomTextContent `json:"text"`
}

type CustomTextContent struct {
	Content string `json:"content"`
}

// CustomAdd 客服新增
func (w *Wechat) CustomAdd(params CustomParam) (re ErrorCode, err error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/customservice/kfaccount/add?access_token=%v", w.accessToken.token)
	paramsStr, _ := json.Marshal(params)
	response, err := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer([]byte(paramsStr)))
	if err != nil {
		return re, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return re, errors.New("客服新增失败，请检查网络是否异常")
	}
	data, _ := ioutil.ReadAll(response.Body)
	_ = json.Unmarshal(data, &re)
	if re.ErrCode != 0 {
		return re, errors.New("客服新增失败, errmsg: " + re.ErrMsg)
	}
	return
}

// CustomList 客服列表
func (w *Wechat) CustomList() (customs []Custom, err error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/customservice/getkflist?access_token=%v", w.accessToken.token)
	response, err := http.Get(url)
	if err != nil {
		return customs, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return customs, errors.New("获取客服列表失败，请检查网络是否异常")
	}
	data, _ := ioutil.ReadAll(response.Body)
	_ = json.Unmarshal(data, &customs)
	return
}

// SendCustomTextMsg 发送文本消息
// param msg 文本消息
// param touser 接收消息的用户（微信openid）
func (w *Wechat) SendCustomTextMsg(msg string, touser string) (re ErrorCode, err error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=%v", w.accessToken.token)
	var params CustomText
	params.ToUser = touser
	params.MsgType = "text"
	params.Text.Content = msg
	paramsStr, _ := json.Marshal(params)
	response, err := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer([]byte(paramsStr)))
	if err != nil {
		return re, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return re, errors.New("发送文本消息失败，请检查网络是否异常")
	}
	data, _ := ioutil.ReadAll(response.Body)
	_ = json.Unmarshal(data, &re)
	if re.ErrCode != 0 {
		return re, errors.New("发送文本消息失败, errmsg: " + re.ErrMsg)
	}
	return
}
