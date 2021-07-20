package model

import (
	"flying-star/internal/db"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	// 角色唯一标识
	Uid string `gorm:"type:varchar(40);not null;unique;comment:角色唯一标识"`
	// 角色名称
	Name string `gorm:"type:varchar(40);not null;comment:角色名称"`
	// 角色标识码
	Code string `gorm:"type:varchar(40);not null;comment:角色标识码"`
}

func (role *Role) TableName() string {
	return "role"
}

func (role *Role) Find() (roles []Role, err error) {
	err = db.MySQLClient.Model(&role).Where(&role).Find(&roles).Error
	return
}

// Info 信息
func (role *Role) Info() (info Role, err error) {
	if err = db.MySQLClient.Where(&role).First(&info).Error; err != nil {
		return
	}
	return
}

// Create 创建用户角色
func (role Role) Create() error {
	return db.MySQLClient.Create(&role).Error
}

// RolesInfoByUserID 通过用户ID 获取该用户的角色信息
func (role *Role) RolesInfoByUserID(userID string) ([]Role, error) {
	var roles []Role
	err := db.MySQLClient.Model(&role).Where("uid IN ( SELECT role_id FROM user_role WHERE user_id = ? )", userID).Find(&roles).Error
	return roles, err
}
