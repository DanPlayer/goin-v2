package service

import (
	"flying-star/internal/module/user/service/impl"
	"flying-star/internal/module/user/service/pojo"
)

func NewUser() UserService {
	return new(impl.UserImpl)
}

type UserService interface {
	// Info 用户信息
	Info(uid string) (pojo.UserInfo, error)
	// InfoByUserNameAndPass 根据账号密码获取用户信息
	InfoByUserNameAndPass(userName, password string) (pojo.UserInfo, error)
	// UpdateUserRoleByUid 修改用户的角色
	UpdateUserRoleByUid(uid, roleCode string) error
}
