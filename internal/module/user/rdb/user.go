package rdb

import (
	"context"
	"encoding/json"
	"flying-star/internal/db"
)

const UserKey = "user:"

var ctx = context.Background()

type Info struct {
	// 用户唯一标识符 UUID
	Uid string `json:"uid"`
	// 账户名
	UserName string `json:"user_name"`
	// 用户昵称
	NickName string `json:"nick_name"`
	// 用户头像
	Avatar string `json:"avatar"`
	// 性别 0：未知 1：男 2：女
	Sex int `json:"sex"`
	// 用户手机号
	Mobile string `json:"mobile"`
	// 用户邮箱
	MailBox string `json:"mail_box"`
	// 角色
	Roles string `json:"roles"`
}

func (rs Info) Set(token, secret string) error {
	bytes, e := json.Marshal(rs)
	if e != nil {
		return e
	}
	return db.RedisClient.Set(UserKey+token+secret, string(bytes), CacheTime)
}

func (rs Info) Get(token, secret string) (re Info, err error) {
	s, err := db.RedisClient.Get(ctx, UserKey+token+secret)
	if err != nil {
		return re, err
	}
	if len(s) <= 0 {
		return Info{}, nil
	}
	err = json.Unmarshal([]byte(s), &re)
	return
}

func (rs Info) Del(token, secret string) error {
	return db.RedisClient.Del(ctx, UserKey+token+secret)
}

func (rs Info) Do(token, secret string) error {
	return db.RedisClient.Do(ctx, UserKey+token+secret, CacheTime)
}
