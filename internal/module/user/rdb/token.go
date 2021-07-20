package rdb

import (
	"errors"
	"flying-star/internal/db"
)

const TokenKey = "token:"

type Token struct {
	Uid   string `json:"uid"`
	Token string `json:"token"`
}

func (rs Token) Get() (string, error) {
	if rs.Uid == "" {
		return "", errors.New("token缓存的必要属性没有设置")
	}
	return db.RedisClient.Get(ctx, TokenKey+rs.Uid)
}

func (rs Token) Set() error {
	if rs.Uid == "" || len(rs.Token) <= 0 {
		return errors.New("token缓存的必要属性没有设置")
	}
	return db.RedisClient.Set(TokenKey+rs.Uid, rs.Token, CacheTime)
}

func (rs Token) Del() error {
	if rs.Uid == "" || len(rs.Token) <= 0 {
		return errors.New("token缓存的必要属性没有设置")
	}
	return db.RedisClient.Del(ctx, TokenKey+rs.Uid)
}
