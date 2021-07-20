package user

import (
	"flying-star/internal/config"
	"flying-star/internal/middleware"
	"flying-star/internal/module/common"
	"flying-star/internal/module/user/router/application/v1"
	"github.com/gin-gonic/gin"
)

// Init 注册用户模块
func Init(conf config.Options) common.ModuleOption {
	//配置业务模块路由规则
	return common.ModuleOption{
		Name: "user",
		ChildList: []common.ModuleChild{
			{
				Route:   "/v1/application/info",
				Method:  "GET",
				Auth:    []string{"admin", "basic"},
				Handles: []gin.HandlerFunc{middleware.AdminAuth(conf.Server.Token, true), v1.Info},
			},
			{
				Route:   "/v1/application/login",
				Method:  "POST",
				Handles: []gin.HandlerFunc{v1.Login},
			},
			{
				Route:   "/v1/application/modify/role",
				Method:  "POST",
				Auth:    []string{"admin"},
				Handles: []gin.HandlerFunc{middleware.AdminAuth(conf.Server.Token, true), v1.ModifyUserRole},
			},
		},
	}
}
