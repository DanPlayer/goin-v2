package wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

//微信公众号授权Token
type accessToken struct {
	Token     string `json:"access_token"` // 获取到的凭证
	ExpiresIn int    `json:"expires_in"`   // 凭证有效时间，单位：秒
	ErrCode   int    `json:"errcode"`      // 错误码
	ErrMsg    string `json:"errmsg"`       // 错误信息
}

// GetAccessToken 获取当前AccessToken
func (w *Wechat) GetAccessToken() string {
	return w.accessToken.token
}

//生成新的AccessToken
func (w *Wechat) initAccessToken() (err error) {
	//如果缓存存在有效token则直接同步到内存中
	ok := w.syncAccessTokenFromCache()
	if ok {
		return nil
	}
	target := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%v&secret=%v", w.options.AppID, w.options.AppSecret)
	response, err := http.Get(target)
	if err != nil {
		return errors.New(fmt.Sprintf("%v: AccessToken 更新失败，请检查AppId等参数是否正确，或者当前IP是否在白名单", time.Now().Format("2016-01-02 15:04:05")))
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return errors.New("AccessToken 更新失败,请检查AppID等配置参数书否正确")
	}
	result := new(accessToken)
	//解析并存储最新AccessToken
	data, _ := ioutil.ReadAll(response.Body)
	_ = json.Unmarshal(data, &result)
	if result.Token == "" {
		return errors.New(fmt.Sprintf("AccessToken 获取失败: %v", result.ErrMsg))
	}
	w.accessToken.token = result.Token
	w.accessToken.expireTime = int(time.Now().Unix()) + w.options.ExpireTime
	//同步最新token到缓存
	w.syncAccessTokenToCache()
	return nil
}

//同步token到缓存
func (w *Wechat) syncAccessTokenToCache() bool {
	err := w.options.Cache.Set("AccessTokenContent", w.accessToken.token, time.Duration(w.options.ExpireTime))
	err = w.options.Cache.Set("AccessTokenExpireTime", strconv.Itoa(w.accessToken.expireTime), time.Duration(w.options.ExpireTime))
	if err != nil {
		return false
	}
	return true
}

//同步token到缓存
func (w *Wechat) syncAccessTokenFromCache() bool {
	accessToken, err := w.options.Cache.Get("AccessTokenContent")
	if err != nil {
		return false
	}
	expireTimeStr, err := w.options.Cache.Get("AccessTokenExpireTime")
	if err != nil {
		return false
	}
	w.accessToken.token = accessToken
	expireTime, _ := strconv.Atoi(expireTimeStr)
	w.accessToken.expireTime = expireTime
	return true
}
