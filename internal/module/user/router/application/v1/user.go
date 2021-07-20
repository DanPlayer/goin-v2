package v1

import (
	"errors"
	"flying-star/internal/middleware"
	"flying-star/internal/module/user/rdb"
	userService "flying-star/internal/module/user/service"
	"flying-star/internal/module/user/service/pojo"
	"flying-star/internal/validate"
	"flying-star/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Info
// @Summary 用户信息
// @Description 用户信息
// @Tags 用户
// @Accept json
// @Produce json
// @Success 200 {object} pojo.UserInfo
// @Router /user/v1/application/info [get]
func Info(c *gin.Context) {
	uid := middleware.GetLoginUid(c)

	service := userService.NewUser()
	userInfo, err := service.Info(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.OutErrorJsonWithStr(c, "用户不存在")
			return
		}
		utils.OutErrorJsonWithStr(c, "服务异常，请稍后再试")
		return
	}
	utils.OutJson(c, userInfo)
}

type UserLogin struct {
	// 用户名
	UserName string `json:"user_name" validate:"required"`
	// 密码
	Password string `json:"password" validate:"required"`
}

type UserLoginResponse struct {
	// 用户信息
	Info pojo.UserInfo `json:"info"`
	// token
	Token string `json:"token"`
}

// Login
// @Summary 普通登录
// @Description 普通登录
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body UserLogin true "登陆参数"
// @Success 200 {object} UserLoginResponse
// @Router /user/v1/application/login [post]
func Login(c *gin.Context) {
	var params UserLogin
	_ = c.ShouldBindJSON(&params)
	if err := validate.ParseStruct(params); err != nil {
		utils.OutErrorJsonWithStr(c, "参数不正确或缺失")
		return
	}
	// 密码加密
	password := utils.EncryptUserPassword(params.Password)
	userInfo, err := userService.NewUser().InfoByUserNameAndPass(params.UserName, password)
	if err != nil {
		utils.OutErrorJsonWithStr(c, "账号密码错误")
		return
	}
	// 生产登录态
	userRdbInfo := rdb.Info{
		Uid:       userInfo.Uid,
		UserName:  userInfo.UserName,
		NickName:  userInfo.NickName,
		Avatar:    userInfo.Avatar,
		Sex:       userInfo.Sex,
		Mobile:    userInfo.Mobile,
		MailBox:   userInfo.MailBox,
		Roles:     userInfo.RoleCode,
		CompanyID: userInfo.CompanyID,
	}
	token, err := middleware.MakeAdminToken(userRdbInfo)
	if err != nil {
		utils.OutErrorJsonWithStr(c, "服务异常，请稍后再试")
		return
	}
	utils.OutJson(c, UserLoginResponse{Info: userInfo, Token: token})
}

type ModifyUserRoleParam struct {
	Uid      string `json:"uid" validate:"required"`      // 被修改的用户ID
	RoleCode string `json:"roleCode" validate:"required"` //修改为某个角色 basic:普通员工 admin:管理员
}

// ModifyUserRole
// @Summary 修改用户角色
// @Description 修改用户角色
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body ModifyUserRoleParam  true "修改用户角色权限参数"
// @Success 200
// @Router /user/v1/application/modify/role [post]
func ModifyUserRole(c *gin.Context) {
	var param ModifyUserRoleParam
	_ = c.ShouldBind(&param)
	if err := validate.ParseStruct(param); err != nil {
		utils.OutErrorJsonWithStr(c, "参数不正确或缺失")
		return
	}

	if param.Uid == middleware.GetLoginUid(c) {
		utils.OutErrorJsonWithStr(c, "不能修改自己的角色")
		return
	}

	if param.RoleCode != "basic" && param.RoleCode != "admin" {
		utils.OutErrorJsonWithStr(c, "角色参数不正确")
		return
	}

	service := userService.NewUser()
	if err := service.UpdateUserRoleByUid(param.Uid, param.RoleCode); err != nil {
		utils.OutErrorJsonWithStr(c, "修改角色失败，服务异常")
		return
	}

	utils.OutJsonOk(c, "修改成功")
}
