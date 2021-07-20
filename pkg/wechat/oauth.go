package wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type OAuthOptions struct {
	AppID      string		// 公众号AppID
	AppSecret  string		// 公众号AppSecret
	RedirectUrl string		// 授权回调地址
}

type OAuth struct {
	appID      string		// 公众号AppID
	appSecret  string		// 公众号AppSecret
	redirectUrl string		// 授权回调地址
}


func NewOAuth(options OAuthOptions) *OAuth {
	return &OAuth{
		appID:       options.AppID,
		appSecret:   options.AppSecret,
		redirectUrl: options.RedirectUrl,
	}
}

// GetLink 获取用户默认授权地址
func (r *OAuth) GetLink(code string) string {
	return "https://open.weixin.qq.com/connect/oauth2/authorize?appid="+ r.appID +"&redirect_uri="+ url.QueryEscape(r.redirectUrl + "/" + code) +"&response_type=code&scope=snsapi_userinfo&state=STATE#wechat_redirect"
}

type AccessTokenSchema struct {
	ErrCode int `json:"errcode"`					// 错误码
	ErrMsg string `json:"errmsg"`					// 错误信息
	AccessToken string `json:"access_token"`		// 网页授权接口调用凭证,注意：此access_token与基础支持的access_token不同
	ExpiresIn int `json:"expires_in"`				// access_token接口调用凭证超时时间，单位（秒）
	RefreshToken string `json:"refresh_token"`		// 用户刷新access_token
	OpenId string `json:"openid"`					// 用户唯一标识，请注意，在未关注公众号时，用户访问公众号的网页，也会产生一个用户和公众号唯一的OpenID
	Scope string `json:"scope"`						// 用户授权的作用域，使用逗号（,）分隔
}

func (r *OAuth) GetAccessToken(code string) (res AccessTokenSchema, err error) {
	if code == "" {
		return res, errors.New("获取用户授权失败：参数 code 不能为空")
	}
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", r.appID, r.appSecret, code)
	response, err := http.Get(url)
	if err != nil {
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return res, errors.New("用户资料获取失败，请检查网络是否异常")
	}
	data, _ := ioutil.ReadAll(response.Body)
	_ = json.Unmarshal(data, &res)
	if res.ErrCode != 0 {
		return res, errors.New("无效的code")
	}
	fmt.Println(string(data))
	return res, nil
}

type UserInfoSchema struct {
	OpenId string `json:"openid"`					// 用户的唯一标识
	NickName string `json:"nickname"`				// 用户昵称
	Sex uint `json:"sex"`							// 用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	Province string `json:"province"`				// 用户个人资料填写的省份
	City string `json:"city"`						// 普通用户个人资料填写的城市
	Country string `json:"country"`					// 国家，如中国为CN
	HeadImgUrl string `json:"headimgurl"`			// 用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效
	Privilege []string `json:"privilege"`			// 用户特权信息，json 数组，如微信沃卡用户为（chinaunicom）
	UnionId string `json:"unionid"`					// 只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段

	ErrCode int `json:"errcode"`					// 错误码
	ErrMsg string `json:"errmsg"`					// 错误信息
}

func (r *OAuth) GetUserInfo(token, openId string) (res UserInfoSchema, err error) {
	if token == "" || openId == "" {
		return res, errors.New("获取用户信息失败：参数 token 或 openId 不能为空")
	}
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN", token, openId)
	response, err := http.Get(url)
	if err != nil {
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return res, errors.New("用户资料获取失败，请检查网络是否异常")
	}
	data, _ := ioutil.ReadAll(response.Body)
	_ = json.Unmarshal(data, &res)
	if res.ErrCode != 0 {
		return res, errors.New("无效的token")
	}
	fmt.Println(string(data))
	return res, nil
}