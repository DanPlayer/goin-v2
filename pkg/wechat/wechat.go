package wechat

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"
)

// Options 初始化参数
type Options struct {
	AppID      string
	AppSecret  string
	EventToken string      // 微信公众号回调请求验证Token(仅用作服务器验证)
	ExpireTime int         `default:"6000"` //AccessToken过期时间
	Cache      MemoryCache // 缓存
}

type MemoryCache interface {
	Set(k, v string, expires time.Duration) error
	Get(k string) (string, error)
}

type ErrorCode struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// Wechat 微信实例
type Wechat struct {
	options     Options // 初始化参数
	accessToken struct {
		token      string // 临时授权token
		expireTime int    // 有效截至时间
	}
	eventQueue sync.Map // 事件分发
	mutex      sync.Mutex
}

// New 初始化实例
func New(options Options) (w *Wechat, err error) {
	if options.ExpireTime == 0 {
		field, _ := reflect.TypeOf(options).FieldByName("ExpireTime")
		options.ExpireTime, _ = strconv.Atoi(field.Tag.Get("default"))
	}
	wechat := Wechat{
		options: options,
	}
	err = wechat.initAccessToken()
	if err != nil {
		return nil, err
	}

	//自动更新AccessToken
	go func() {
		for {
			currentTime := int(time.Now().Unix())
			if currentTime > w.accessToken.expireTime {
				err = wechat.initAccessToken()
				if err != nil {
					//TODO: 添加系统警告
					fmt.Printf("AccessToken更新失败:%v\n", err.Error())
				} else {
					fmt.Printf("AccessToken更新成功:%v\n", wechat.accessToken.token)
				}
			}
			time.Sleep(10 * time.Second)
		}
	}()
	return &wechat, nil
}
