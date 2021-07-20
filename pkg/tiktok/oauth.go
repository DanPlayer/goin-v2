package tiktok

import (
	"encoding/json"
	"flying-star/utils"
	"fmt"
	"net/url"
)

// AccessTokenSchema 用户授权凭证
type AccessTokenSchema struct {
	Data struct {
		AccessToken      string `json:"access_token"`       // 接口调用凭证
		Description      string `json:"description"`        // 错误码描述
		ErrorCode        int    `json:"error_code"`         // 错误码
		ExpiresIn        int    `json:"expires_in"`         // access_token接口调用凭证超时时间，单位（秒)
		OpenId           string `json:"open_id"`            // 授权用户唯一标识
		RefreshExpiresIn int    `json:"refresh_expires_in"` // refresh_token凭证超时时间，单位（秒)
		RefreshToken     string `json:"refresh_token"`      // 用户刷新access_token
		Scope            string `json:"scope"`              // 用户授权的作用域(Scope)，使用逗号（,）分隔，开放平台几乎几乎每个接口都需要特定的Scope
	} `json:"data"`
	Message string `json:"message"`
}

// RefreshTokenSchema 刷新用户授权凭证
type RefreshTokenSchema struct {
	Data struct {
		Description  string `json:"description"`   // 错误码描述
		ErrorCode    int    `json:"error_code"`    // 错误码
		ExpiresIn    int    `json:"expires_in"`    // 过期时间，单位（秒）
		RefreshToken string `json:"refresh_token"` // 用户刷新
	} `json:"data"`
	Message string `json:"message"`
}

// UserInfoSchema 用户公开信息
type UserInfoSchema struct {
	Data struct {
		Avatar       string `json:"avatar"`
		City         string `json:"city"`
		Country      string `json:"country"`
		Description  string `json:"description"`    // 错误码描述
		EAccountRole string `json:"e_account_role"` // 类型: * `EAccountM` - 普通企业号 * `EAccountS` - 认证企业号 * `EAccountK` - 品牌企业号
		ErrorCode    string `json:"error_code"`     // 错误码
		Gender       string `json:"gender"`         // 性别: * `0` - 未知 * `1` - 男性 * `2` - 女性
		Nickname     string `json:"nickname"`
		OpenId       string `json:"open_id"` // 用户在当前应用的唯一标识
		Province     string `json:"province"`
		UnionId      string `json:"union_id"` // 用户在当前开发者账号下的唯一标识（未绑定开发者账号没有该字段）
	} `json:"data"`
}

type ClientTokenSchema struct {
	Data struct {
		AccessToken string `json:"access_token"`
		Description string `json:"description"`
		ErrorCode   string `json:"error_code"`
		ExpiresIn   string `json:"expires_in"`
	} `json:"data"`
	Message string `json:"message"`
}

func (tk *TikTok) ClientToken() (res ClientTokenSchema, err error) {
	url := fmt.Sprintf("%s/oauth/client_token?client_key=%s&client_secret=%s&grant_type=client_credential", baseUrl, tk.clientKey, tk.clientSecret)
	body, err := utils.HttpGetBody(url)
	if err != nil {
		return res, err
	}
	_ = json.Unmarshal(body, &res)
	return res, nil
}

// AccessToken 获取抖音用户 AccessToken
func (tk *TikTok) AccessToken(code string) (res AccessTokenSchema, err error) {
	url := fmt.Sprintf("%s/oauth/access_token?client_key=%s&client_secret=%s&grant_type=authorization_code&code=%s", baseUrl, tk.clientKey, tk.clientSecret, code)
	body, err := utils.HttpGetBody(url)
	if err != nil {
		return res, err
	}
	_ = json.Unmarshal(body, &res)
	return res, nil
}

// RenewRefreshToken 刷新refresh_token的有效期；该接口适用于抖音授权
func (tk *TikTok) RenewRefreshToken(refreshToken string) (res RefreshTokenSchema, err error) {
	url := fmt.Sprintf("%s/oauth/renew_refresh_token?client_key=%s&refresh_token=%s", baseUrl, tk.clientKey, refreshToken)
	body, err := utils.HttpGetBody(url)
	if err != nil {
		return res, err
	}
	_ = json.Unmarshal(body, &res)
	return res, nil
}

// RefreshToken 使用 refreshToken 刷新 accessToken
func (tk *TikTok) RefreshToken(refreshToken string) (res AccessTokenSchema, err error) {
	url := fmt.Sprintf("%s/oauth/refresh_token?client_key=%s&grant_type=%s&refresh_token=%s", baseUrl, tk.clientKey, "refresh_token", refreshToken)
	body, err := utils.HttpGetBody(url)
	if err != nil {
		return res, err
	}
	_ = json.Unmarshal(body, &res)
	return res, nil
}

// UserInfo 获取用户的抖音公开信息，包含昵称、头像、性别和地区；适用于抖音
func (tk *TikTok) UserInfo(openId, accessToken string) (res UserInfoSchema, err error) {
	url := fmt.Sprintf("%s/oauth/userinfo/?open_id=%s&access_token=%s", baseUrl, openId, accessToken)
	body, err := utils.HttpGetBody(url)
	if err != nil {
		return res, err
	}
	_ = json.Unmarshal(body, &res)
	return res, nil
}

// AuthorizeUrl 用户权限授权地址
func (tk *TikTok) AuthorizeUrl(scope, redirectUri, state string) string {
	redirectUri = url.QueryEscape(redirectUri)
	return fmt.Sprintf("%s/platform/oauth/connect?client_key=%s&response_type=code&scope=%s&redirect_uri=%s&state=%s", baseUrl, tk.clientKey, scope, redirectUri, state)
}
