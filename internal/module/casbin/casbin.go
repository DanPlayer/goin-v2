package casbin

import (
	"flying-star/internal/config"
	"flying-star/internal/middleware"
	"flying-star/internal/module/casbin/router/application/v1"
	"flying-star/internal/module/common"
	"github.com/gin-gonic/gin"
)

// Init 注册权限模块
func Init(conf config.Options) common.ModuleOption {
	//配置业务模块路由规则
	return common.ModuleOption{
		Name: "casbin",
		ChildList: []common.ModuleChild{
			{
				Route:   "/v1/application/enforce",
				Method:  "POST",
				Auth:    []string{"admin"},
				Handles: []gin.HandlerFunc{middleware.AdminAuth(conf.Server.Token, true), v1.Enforce},
			},
			{
				Route:   "/v1/application/policy/all",
				Method:  "GET",
				Auth:    []string{"admin"},
				Handles: []gin.HandlerFunc{middleware.AdminAuth(conf.Server.Token, true), v1.AllPolicy},
			},
			{
				Route:   "/v1/application/policy/role",
				Method:  "GET",
				Auth:    []string{"admin"},
				Handles: []gin.HandlerFunc{middleware.AdminAuth(conf.Server.Token, true), v1.RolePolicy},
			},
			{
				Route:   "/v1/application/role/policy",
				Method:  "DELETE",
				Auth:    []string{"admin"},
				Handles: []gin.HandlerFunc{middleware.AdminAuth(conf.Server.Token, true), v1.RemoveRolesPolicy},
			},
			{
				Route:   "/v1/application/user/role",
				Method:  "POST",
				Auth:    []string{"admin"},
				Handles: []gin.HandlerFunc{middleware.AdminAuth(conf.Server.Token, true), v1.AddRolesForUser},
			},
			{
				Route:   "/v1/application/user/role",
				Method:  "PUT",
				Auth:    []string{"admin"},
				Handles: []gin.HandlerFunc{middleware.AdminAuth(conf.Server.Token, true), v1.UpdateRolesForUser},
			},
			{
				Route:   "/v1/application/user/role",
				Method:  "DELETE",
				Auth:    []string{"admin"},
				Handles: []gin.HandlerFunc{middleware.AdminAuth(conf.Server.Token, true), v1.DeleteRolesForUser},
			},
			{
				Route:   "/v1/application/user/role",
				Method:  "GET",
				Auth:    []string{"admin"},
				Handles: []gin.HandlerFunc{middleware.AdminAuth(conf.Server.Token, true), v1.UserRoles},
			},
			{
				Route:   "/v1/application/company/role",
				Method:  "GET",
				Handles: []gin.HandlerFunc{middleware.AdminAuth(conf.Server.Token, false), v1.CompanyRoles},
			},
		},
	}
}
