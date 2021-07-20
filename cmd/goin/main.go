package main

import (
	_ "flying-star/docs"
	"flying-star/internal/api"
	"fmt"
)

// @title GoIn
// @version 2.0
// @description GoIn
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @contact.name GoIn
// @contact.url localhost:8088
// @host localhost:8088
func main() {
	// 启动api服务
	apiMsg := api.Run()

	for {
		select {
		case _, ok := <-apiMsg:
			if !ok {
				fmt.Println("API服务异常退出，请检查服务配置是否错误")
				return
			}
		}
	}
}
