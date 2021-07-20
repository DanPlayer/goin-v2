package model

import (
	"flying-star/internal/db"
	"gorm.io/gorm"
)

type UerRole struct {
	gorm.Model
	// 唯一标识
	Uid string `gorm:"type:varchar(40);not null;unique;comment:唯一标识"`
	// 用户唯一标识
	UserID string `gorm:"type:varchar(40);not null;comment:用户唯一标识"`
	// 角色唯一标识
	RoleID string `gorm:"type:varchar(40);not null;comment:角色唯一标识"`
}

func (ur UerRole) TableName() string {
	return "user_role"
}

// Create 创建用户角色
func (ur *UerRole) Create() error {
	return db.MySQLClient.Create(&ur).Error
}