package model

import "gorm.io/gorm"

type PayBinding struct {
	gorm.Model
	// 用户ID
	UserID string `gorm:"type:varchar(40);not null;unique;comment:用户支付ID"`
	// 绑定ID
	OpenID string `gorm:"type:varchar(255);not null;unique;comment:用户开放ID"`
	//用户真实姓名
	RealName string `gorm:"type:varchar(64);default:'';comment:用户真实姓名"`
	//模块名
	ModuleName string `gorm:"type:varchar(64);default:'';comment:模块名"`
}