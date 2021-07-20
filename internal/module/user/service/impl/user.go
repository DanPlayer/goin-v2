package impl

import (
	"flying-star/internal/config"
	"flying-star/internal/module/user/model"
	"flying-star/internal/module/user/rdb"
	"flying-star/internal/module/user/service/pojo"
)

type UserImpl struct{}

// Info 用户信息
func (i UserImpl) Info(uid string) (pojo.UserInfo, error) {
	// 查询用户信息
	user := model.User{Uid: uid}
	u, err := user.Info()
	if err != nil {
		return pojo.UserInfo{}, err
	}

	userInfo := pojo.UserInfo{
		Uid:      u.Uid,
		UserName: u.UserName,
		NickName: u.NickName,
		Avatar:   u.Avatar,
		Sex:      u.Sex,
		Mobile:   u.Mobile,
		MailBox:  u.MailBox,
		RoleCode: u.RoleCode,
	}

	return userInfo, nil
}

// InfoByUserNameAndPass 根据账号密码获取用户信息
func (i UserImpl) InfoByUserNameAndPass(userName, password string) (pojo.UserInfo, error) {
	user := model.User{UserName: userName, Password: password}
	u, err := user.Info()
	if err != nil {
		return pojo.UserInfo{}, err
	}

	return i.Info(u.Uid)
}

// UpdateUserRoleByUid 更新用户的角色
func (i UserImpl) UpdateUserRoleByUid(uid, roleCode string) error {
	user := model.User{Uid: uid}
	updateUser := model.User{RoleCode: roleCode}
	err := user.Update(updateUser)
	if err != nil {
		return err
	}
	tokenRdb := rdb.Token{Uid: uid}
	token, err := tokenRdb.Get()
	if err != nil {
		return err
	}
	if token == "" {
		return nil
	}

	userRdb := rdb.Info{}
	info, err := userRdb.Get(token, config.Info.Server.Token)
	if err != nil {
		return err
	}
	info.Roles = roleCode
	_ = info.Set(token, config.Info.Server.Token)

	return nil
}
