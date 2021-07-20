package api

import (
	"flying-star/internal/casbin"
	"flying-star/internal/config"
	"flying-star/internal/crontab"
	"flying-star/internal/db"
	"flying-star/internal/middleware"
	"flying-star/internal/module"
	casbinModel "flying-star/internal/module/casbin/model"
	userModel "flying-star/internal/module/user/model"
	"fmt"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"log"
)

// Run 服务启动
func Run() <-chan struct{} {
	ch := make(chan struct{})

	go func() {
		defer close(ch)
		// 自动迁移数据表
		db.AutoMigrateApiServerTable(
			&userModel.User{},
			&casbinModel.Role{},
			&casbinModel.UerRole{},
		)

		// 初始化casbin
		err := casbin.Dial(db.MySQLClient)
		if err != nil {
			log.Fatalln("casbin validate_err", err)
		}

		// http服务器
		gin.DisableConsoleColor()
		app := gin.New()
		app.Use(gin.Recovery())

		if config.Info.Server.Debug {
			gin.SetMode(gin.DebugMode)
			// 挂载中间件
			app.Use(middleware.CrossDomainForDebug())
		} else {
			gin.SetMode(gin.ReleaseMode)
			// 挂载中间件
			app.Use(middleware.CrossDomain())
		}
		//注册业务模块
		moduleList := module.Registry(config.Info)
		for _, option := range moduleList {
			group := app.Group(option.Name)
			{
				for _, child := range option.ChildList {
					//批量注册路由
					if child.Method == "GET" {
						group.GET(child.Route, child.Handles...)
					}
					if child.Method == "POST" {
						group.POST(child.Route, child.Handles...)
					}
					if child.Method == "PUT" {
						group.PUT(child.Route, child.Handles...)
					}
					if child.Method == "DELETE" {
						group.DELETE(child.Route, child.Handles...)
					}

					//初始化相关权限
					if len(child.Auth) == 0 {
						continue
					}
					for _, role := range child.Auth {
						if _, err := casbin.Casbin.AddNamedPolicy("p", role, "/"+option.Name+child.Route, child.Method); err != nil {
							log.Fatalln("casbin init policy validate_err:" + err.Error())
						}
					}
				}
			}
		}

		// 文档服务
		app.GET("/docs/*any", gin.BasicAuth(map[string]string{"admin": config.Info.DocAuth.Admin}), ginSwagger.WrapHandler(swaggerFiles.Handler))

		// 启动定时任务
		cron := crontab.InitCron()
		defer cron.Stop()

		// 启动服务
		if err := app.Run(fmt.Sprintf(":%v", config.Info.Server.Port)); err != nil {
			log.Println("app run validate_err", err)
		}
	}()

	return ch
}
