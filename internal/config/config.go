package config

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"path"
	"testing"
)

// Options 系统配置参数声明
type Options struct {
	// 服务通用配置
	Server struct {
		//服务名
		AppName string `json:"appName"`
		//服务端口
		Port int `json:"port"`
		//Token加密令牌
		Token string `json:"token"`
		//AES加密密钥
		AesKey string `json:"aesKey"`
		//是否开启DEBUG
		Debug bool `json:"debug"`
	} `json:"server"`

	//权限配置
	Casbin struct {
		SuperRole string `json:"superRole"`
		AdminRole string `json:"adminRole"`
		BasicRole string `json:"baseRole"`
	} `json:"casbin"`

	//文档权限配置
	DocAuth struct {
		Admin string `json:"admin"`
	} `json:"docAuth"`

	//MySQL数据库
	MySQL struct {
		DSN string `json:"dsn"`
	} `json:"mysql"`

	//Redis数据库
	Redis struct {
		Addr     string `json:"addr"`
		DB       int    `json:"db"`
		Password string `json:"password"`
	} `json:"redis"`
}

var Info Options
var currentEnv *string

func init() {
	currentEnv = flag.String("env", "dev", "-env=dev")
	testing.Init()
	flag.Parse()

	//解析配置信息到Info
	currentPath, _ := os.Getwd()
	fileName := "app." + *currentEnv + ".toml"

	viper.SetConfigType("toml")
	viper.AddConfigPath(currentPath)
	viper.SetConfigName(fileName)
	viper.AddConfigPath(path.Join(currentPath, "./config"))
	viper.OnConfigChange(func(in fsnotify.Event) {
		refreshConfig()
	})
	refreshConfig()
}

// 刷新系统配置文件内容
func refreshConfig() {
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		panic("基础配置文件加载失败，请检查/config文件是否设置正确。")
	}
	//Dev环境下默认开启DEBUG模式
	if err := viper.Unmarshal(&Info); err == nil && *currentEnv == "dev" {
		Info.Server.Debug = true
	}
}
