package tiktok

const baseUrl = "https://open.douyin.com"

type Options struct {
	ClientKey    string `json:"clientKey"`
	ClientSecret string `json:"clientSecret"`
}

type TikTok struct {
	clientKey    string
	clientSecret string
}

// New 初始化Tiktok实例
func New(options Options) *TikTok {
	return &TikTok{
		clientKey:    options.ClientKey,
		clientSecret: options.ClientSecret,
	}
}

// ErrorSchema 接口响应错误
type ErrorSchema struct {
	Data struct{
		Description string `json:"description"`				// 错误描述
		ErrorCode int `json:"error_code"`					// 错误码
	} `json:"data"`
	Extra struct{
		LogId string `json:"logid"`							// 日志ID
		Now int `json:"now"`								// 当前时间
	} `json:"extra"`
}