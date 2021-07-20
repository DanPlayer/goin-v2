package wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// WxUserInfo 微信用户资料
type WxUserInfo struct {
	Subscribe      int    `json:"subscribe"`
	OpenId         string `json:"openid"`
	NickName       string `json:"nickname"`
	Sex            int    `json:"sex"`
	Language       string `json:"language"`
	City           string `json:"city"`
	Province       string `json:"province"`
	Country        string `json:"country"`
	HeadImgUrl     string `json:"headimgurl"`
	SubscribeTime  int    `json:"subscribe_time"`
	UnionId        string `json:"unionid"`
	Remark         string `json:"remark"`
	GroupId        string `json:"groupid"`
	TagIdList      []int  `json:"tagid_list"`
	SubscribeScene string `json:"subscribe_scene"`
	QrScene        string `json:"qr_scene"`
	QrSceneStr     string `json:"qr_scene_str"`
}

// GetUserInfo 根据openid获取微信用户资料
func (w *Wechat) GetUserInfo(openid string) (result WxUserInfo, err error) {
	if openid == "" {
		return result, errors.New("接口 GetUserInfo 调用失败：参数 openid 不能为空")
	}
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/info?access_token=%v&openid=%v&lang=zh_CN", w.accessToken.token, openid)
	response, err := http.Get(url)
	if err != nil {
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return result, errors.New("用户资料获取失败，请检查网络是否异常")
	}
	//解析并存储最新AccessToken
	data, _ := ioutil.ReadAll(response.Body)
	_ = json.Unmarshal(data, &result)
	return result, nil
}
