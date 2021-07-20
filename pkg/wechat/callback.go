package wechat

import (
	"crypto/sha1"
	"encoding/hex"
	"sort"
)

// ServiceVerificationOptions 微信公众号服务器验证
type ServiceVerificationOptions struct {
	Signature string `form:"signature"`
	TimeStamp string `form:"timestamp"`
	Nonce     string `form:"nonce"`
	EchoStr   string `form:"echostr"`
}

// ServiceVerification 校验服务器参数是否合法
func (w *Wechat) ServiceVerification(options ServiceVerificationOptions) (echo string, ok bool) {
	var paramsStr string
	//对参数进行字典排序
	paramsArr := []string{w.options.EventToken, options.TimeStamp, options.Nonce}
	sort.Strings(paramsArr)
	for _, list := range paramsArr {
		paramsStr += list
	}
	//进行SHA1加密
	sha := sha1.New()
	sha.Write([]byte(paramsStr))
	paramsStr = hex.EncodeToString(sha.Sum([]byte("")))

	if paramsStr == options.Signature {
		return options.EchoStr, true
	}
	return echo, false
}
