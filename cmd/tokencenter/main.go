package main

import (
	"errors"
	"flying-star/internal/config"
	"flying-star/pkg/db"
	"flying-star/pkg/tokencenter"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func main()  {
	// RedisClient 初始化Redis客户端
	redisClient := db.NewRedis(db.RedisOptions{
		Addr:     config.Info.Redis.Addr,
		Password: config.Info.Redis.Password,
	})

	//初始化 RefreshToken 令牌中心
	tokenCenter := tokencenter.New(tokencenter.Options{
		UniqueCode: "refreshToken",
		Cache: redisClient,
		RefreshSW: true,
		RefreshStrategy: []time.Duration{15, 15, 15, 15, 15},
		RefreshHandle: func(token string) (string, error) {
			return strconv.Itoa(rand.Intn(10000000)), nil
		},
	})

	//监听令牌过期事件
	tokenCenter.SubscribeExpiredEvent(func(key string) {
		fmt.Printf("refresh token expired: %s \n", key)
	})

	//存储令牌
	_ = tokenCenter.Set("refresh_token", strconv.Itoa(rand.Intn(10000000)))

	//初始化 AccessToken 实例
	tokenCenter2 := tokencenter.New(tokencenter.Options{
		UniqueCode: "accessToken",
		Cache:           redisClient,
		RefreshSW: false,
		RefreshTime: 5,
		RefreshHandle: func(token string) (string, error) {
			refreshToken, err := tokenCenter.Get("refresh_token")
			if err != nil || refreshToken == "" {
				return "", errors.New("refresh token expired")
			}
			return strconv.Itoa(rand.Intn(10000000)), nil
		},
	})

	//监听令牌过期事件
	tokenCenter2.SubscribeExpiredEvent(func(key string) {
		fmt.Printf("access token expired: %s \n", key)
	})

	//存储令牌
	_ = tokenCenter2.Set("access_token", strconv.Itoa(rand.Intn(10000000)))
	select {}
}
