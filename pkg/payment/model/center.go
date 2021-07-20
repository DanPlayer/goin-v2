package model

import "gorm.io/gorm"

type PayCenter struct {
	gorm.Model
	// 账户ID（UserID/CorpID）
	AccountId string `gorm:"index;type:varchar(40);not null;unique;comment:账户ID（UserID/CorpID）"`
	// 账户余额
	Balance uint `gorm:"type:int(12);default:0;comment:账户余额"`
	// 账户类型
	Type string `gorm:"type:varchar(40);default:group;comment:支付类型：group(团队)，person(个人)"`
}
