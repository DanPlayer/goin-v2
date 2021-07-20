package wechat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//登录二维码参数
type qrcodeOptions struct {
	ExpireSeconds time.Duration `json:"expire_seconds"`
	ActionName    string        `json:"action_name"`
	ActionInfo    struct {
		Scene struct {
			SceneStr string `json:"scene_str"`
		} `json:"scene"`
	} `json:"action_info"`
}

type QrCode struct {
	//二维码图片地址
	Ticket string `json:"ticket"`
	//过期时间 当值是 0 时为永久二维码
	Expire int `json:"expire_seconds"`
	//二维码包含的数据
	Url string `json:"url"`
	//访问链接
	Link string `json:"link"`
}

// GetTempQrCode 获取临时登录二维码
func (w *Wechat) GetTempQrCode(sceneStr string, expireSeconds time.Duration) (qr QrCode, err error) {
	link := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=%v", w.accessToken.token)
	var params qrcodeOptions
	params.ExpireSeconds = expireSeconds
	params.ActionName = "QR_STR_SCENE"
	params.ActionInfo.Scene.SceneStr = sceneStr

	paramsStr, _ := json.Marshal(params)
	response, err := http.Post(link, "application/json;charset=utf-8", bytes.NewBuffer([]byte(paramsStr)))
	if err != nil {
		return qr, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return qr, errors.New("临时二维码获取失败，请检查网络是否异常")
	}
	data, _ := ioutil.ReadAll(response.Body)
	_ = json.Unmarshal(data, &qr)
	if qr.Ticket == "" {
		return qr, errors.New("临时二维码获取失败,请检查AccessToken是否过期")
	}
	qr.Link = fmt.Sprintf("https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=%v", qr.Ticket)
	return qr, nil
}
