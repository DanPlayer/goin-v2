package impl

import (
	"flying-star/internal/casbin"
	"flying-star/internal/module/casbin/model"
	"flying-star/utils"
)

type CasbinImpl struct{}

// Enforce 创建权限
func (i CasbinImpl) Enforce(sub, obj, act string) error {
	_, err := casbin.Casbin.AddNamedPolicy("p", sub, obj, act)
	return err
}

// GetFilterPolicy 获取角色的所有权限
func (i CasbinImpl) GetFilterPolicy(role string) [][]string {
	return casbin.Casbin.GetFilteredPolicy(0, role)
}

// BatchEnforce 批量新增权限
func (i CasbinImpl) BatchEnforce(enforces [][]interface{}) error {
	_, err := casbin.Casbin.BatchEnforce(enforces)
	return err
}

// RemoveNamedPolicy 删除角色权限
func (i CasbinImpl) RemoveNamedPolicy(sub, obj, act string) error {
	_, err := casbin.Casbin.RemoveNamedPolicy("p", sub, obj, act)
	return err
}

// RemoveFilteredNamedPolicy 批量删除角色权限
func (i CasbinImpl) RemoveFilteredNamedPolicy(enforces [][]string) error {
	_, err := casbin.Casbin.RemoveNamedPolicies("p", enforces)
	return err
}

// AddRoleForUser 新增用户角色
func (i CasbinImpl) AddRoleForUser(uid, role string) error {
	roleModel := model.Role{Code: role}
	info, err := roleModel.Info()
	if err != nil {
		return err
	}
	var casbinRole = role
	_, err = casbin.Casbin.AddRoleForUser(uid, casbinRole)
	if err != nil {
		return err
	}

	uerRole := model.UerRole{Uid: utils.GetUid(), UserID: uid, RoleID: info.Uid}
	return uerRole.Create()
}

// AddRolesForUser 新增用户多种角色
func (i CasbinImpl) AddRolesForUser(uid string, roles []string) error {
	for index, role := range roles {
		roleModel := model.Role{Code: role}
		info, err := roleModel.Info()
		if err != nil {
			continue
		}
		uerRole := model.UerRole{Uid: utils.GetUid(), UserID: uid, RoleID: info.Uid}
		if err = uerRole.Create(); err != nil {
			continue
		}

		roles[index] = role
	}
	_, err := casbin.Casbin.AddRolesForUser(uid, roles)
	return err
}

// DeleteRoleForUser 删除用户角色
func (i CasbinImpl) DeleteRoleForUser(uid, role string) error {
	_, err := casbin.Casbin.DeleteRoleForUser(uid, role)
	return err
}

// DeleteRolesForUser 清除用户所有角色
func (i CasbinImpl) DeleteRolesForUser(uid string) error {
	_, err := casbin.Casbin.DeleteRolesForUser(uid)
	return err
}

// GetPermissionsForUser 获取用户所有的权限
func (i CasbinImpl) GetPermissionsForUser(uid string) [][]string {
	return casbin.Casbin.GetPermissionsForUser(uid)
}

// GetRolesForUser 获取用户所有的角色
func (i CasbinImpl) GetRolesForUser(uid string) ([]string, error) {
	return casbin.Casbin.GetRolesForUser(uid)
}

func (i CasbinImpl) GetRolesInfoForUser(uid string) ([]model.Role, error) {
	role := model.Role{}
	return role.RolesInfoByUserID(uid)
}
