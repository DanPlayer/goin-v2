package v1

import (
	"flying-star/internal/config"
	"flying-star/internal/middleware"
	"flying-star/internal/module/casbin/service"
	"flying-star/internal/validate"
	"flying-star/utils"
	"fmt"
	"github.com/gin-gonic/gin"
)

type EnforcePost struct {
	// 角色
	Sub string `json:"sub" validate:"required"`
	// 将要被访问的资源(功能路径)
	Obj string `json:"obj" validate:"required"`
	// 用户对资源的操作(GET,POST)
	Act string `json:"act" validate:"required"`
}

// Enforce
// @Summary 创建Casbin权限
// @Description 注册接口
// @Tags 权限
// @Accept json
// @Produce json
// @Param body body EnforcePost true "创建Casbin权限参数"
// @Success 200
// @Router /casbin/v1/application/enforce [post]
func Enforce(c *gin.Context) {
	var params EnforcePost
	if err := c.ShouldBindJSON(&params); err != nil {
		utils.OutParamErrorJson(c)
		return
	}

	if err := validate.ParseStruct(params); err != nil {
		utils.OutErrorJson(c, err)
		return
	}

	if params.Sub == config.Info.Casbin.SuperRole {
		utils.OutErrorJsonWithStr(c, "不允许创建超级管理员的权限，超级管理员已经是最高权限")
		return
	}

	if err := service.NewCasbin().Enforce(params.Sub, params.Obj, params.Act); err != nil {
		utils.OutErrorJson(c, err)
		return
	}
	utils.OutJsonOk(c, "创建Casbin权限成功")
}

// AllPolicy
// @Summary 获取所有权限
// @Description 获取所有权限
// @Tags 权限
// @Accept json
// @Produce json
// @Success 200 {object} [][]string
// @Router /casbin/v1/application/policy/all [get]
func AllPolicy(c *gin.Context) {
	utils.OutJson(c, service.NewCasbin().GetFilterPolicy(config.Info.Casbin.SuperRole))
}

// RolePolicy
// @Summary 获取角色所有权限
// @Description 获取角色所有权限
// @Tags 权限
// @Accept json
// @Produce json
// @Param role query string true "角色名称"
// @Success 200 {object} [][]string
// @Router /casbin/v1/application/policy/role [get]
func RolePolicy(c *gin.Context) {
	role := c.DefaultQuery("role", "")
	if role == "" {
		utils.OutErrorJsonWithStr(c, "角色名称不能为空")
		return
	}
	utils.OutJson(c, service.NewCasbin().GetFilterPolicy(role))
}

// RemoveRolesPolicy
// @Summary 批量删除角色权限
// @Description 批量删除角色权限
// @Tags 权限
// @Accept json
// @Produce json
// @Param body body [][]string true "需要删除的权限，格式：[{"superman","/user/info","GET"}]"
// @Success 200
// @Router /casbin/v1/application/role/policy [delete]
func RemoveRolesPolicy(c *gin.Context) {
	var policy [][]string
	_ = c.ShouldBindJSON(&policy)
	if err := validate.ParseStruct(policy); err != nil {
		utils.OutErrorJson(c, err)
		return
	}
	if err := service.NewCasbin().RemoveFilteredNamedPolicy(policy); err != nil {
		utils.OutErrorJson(c, err)
		return
	}
	utils.OutJsonOk(c, "删除角色权限成功")
}

type AddRolesForUserPost struct {
	Uid   string   `json:"uid" validate:"required"`
	Roles []string `json:"roles" validate:"required"`
}

// AddRolesForUser
// @Summary 新增用户角色
// @Description 新增用户角色
// @Tags 权限
// @Accept json
// @Produce json
// @Param body body AddRolesForUserPost true "新增用户角色参数"
// @Success 200
// @Router /casbin/v1/application/user/role [post]
func AddRolesForUser(c *gin.Context) {
	var params AddRolesForUserPost
	_ = c.ShouldBindJSON(&params)
	if err := validate.ParseStruct(params); err != nil {
		utils.OutErrorJson(c, err)
		return
	}

	if err := service.NewCasbin().AddRolesForUser(params.Uid, params.Roles); err != nil {
		utils.OutErrorJson(c, err)
		return
	}
	utils.OutJsonOk(c, "新增用户角色成功")
}

type UpdateRolesForUserPut struct {
	Uid   string   `json:"uid" validate:"required"`
	Roles []string `json:"roles" validate:"required"`
}

// UpdateRolesForUser
// @Summary 新增用户角色
// @Description 新增用户角色
// @Tags 权限
// @Accept json
// @Produce json
// @Param body body UpdateRolesForUserPut true "更新用户角色参数"
// @Success 200
// @Router /casbin/v1/application/user/role [put]
func UpdateRolesForUser(c *gin.Context) {
	var params UpdateRolesForUserPut
	_ = c.ShouldBindJSON(&params)
	if err := validate.ParseStruct(params); err != nil {
		utils.OutErrorJson(c, err)
		return
	}

	uid := middleware.GetLoginUid(c)

	casbinService := service.NewCasbin()

	// 先清除用户角色
	if err := casbinService.DeleteRolesForUser(uid); err != nil {
		utils.OutErrorJson(c, err)
		return
	}
	// 新增用户角色
	if err := casbinService.AddRolesForUser(params.Uid, params.Roles); err != nil {
		utils.OutErrorJson(c, err)
		return
	}
	utils.OutJsonOk(c, "更新用户角色成功")
}

type DeleteRolesForUserPost struct {
	Uid   string   `json:"uid" validate:"required"`
	Roles []string `json:"roles" validate:"required"`
}

// DeleteRolesForUser
// @Summary 批量删除用户角色
// @Description 批量删除用户角色
// @Tags 权限
// @Accept json
// @Produce json
// @Param body body DeleteRolesForUserPost true "批量删除用户角色参数"
// @Success 200
// @Router /casbin/v1/application/user/role [delete]
func DeleteRolesForUser(c *gin.Context) {
	var params DeleteRolesForUserPost
	_ = c.ShouldBindJSON(&params)
	if err := validate.ParseStruct(params); err != nil {
		utils.OutErrorJson(c, err)
		return
	}
	for index, _ := range params.Roles {
		err := service.NewCasbin().DeleteRoleForUser(params.Uid, params.Roles[index])
		if err != nil {
			fmt.Println(fmt.Sprintf("delete role for user validate_err: %s", err.Error()))
			return
		}
	}
	utils.OutJsonOk(c, "删除用户角色成功")
}

type RoleInfoResponse struct {
	Uid  string `json:"uid"`  // 角色唯一ID
	Name string `json:"name"` // 角色名称
	Code string `json:"code"` // 角色Code
}

// UserRoles
// @Summary 获取用户所有角色
// @Description 获取用户所有角色
// @Tags 权限
// @Accept json
// @Produce json
// @Param uid query string true "用户ID"
// @Success 200 {object} []string
// @Router /casbin/v1/application/user/role [get]
func UserRoles(c *gin.Context) {
	uid := c.DefaultQuery("uid", "")
	if uid == "" {
		utils.OutErrorJsonWithStr(c, "用户ID不能为空")
		return
	}
	roles, err := service.NewCasbin().GetRolesForUser(uid)
	if err != nil {
		utils.OutErrorJson(c, err)
		return
	}
	utils.OutJson(c, roles)
}
