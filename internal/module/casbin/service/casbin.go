package service

import (
	"flying-star/internal/module/casbin/model"
	"flying-star/internal/module/casbin/service/impl"
)

func NewCasbin() CasbinService {
	return new(impl.CasbinImpl)
}

// CasbinService 服务
type CasbinService interface {
	// Enforce 创建权限
	Enforce(sub, obj, act string) error
	// GetFilterPolicy 获取角色的所有权限
	GetFilterPolicy(role string) [][]string
	// BatchEnforce 批量新增权限
	BatchEnforce(enforces [][]interface{}) error
	// RemoveNamedPolicy 删除角色权限
	RemoveNamedPolicy(sub, obj, act string) error
	// RemoveFilteredNamedPolicy 批量删除角色权限
	RemoveFilteredNamedPolicy(enforces [][]string) error
	// AddRoleForUser 新增用户角色
	AddRoleForUser(uid, role string) error
	// AddRolesForUser 新增用户多种角色
	AddRolesForUser(uid string, roles []string) error
	// DeleteRoleForUser 删除用户角色
	DeleteRoleForUser(uid, role string) error
	// DeleteRolesForUser 清除用户所有角色
	DeleteRolesForUser(uid string) error
	// GetPermissionsForUser 获取用户所有的权限
	GetPermissionsForUser(uid string) [][]string
	// GetRolesForUser 获取用户所有的角色
	GetRolesForUser(uid string) ([]string, error)
	// GetRolesInfoForUser 获取用户所有角色信息
	GetRolesInfoForUser(uid string) ([]model.Role, error)
}
