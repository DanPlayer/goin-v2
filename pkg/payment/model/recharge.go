package model

import (
	"gorm.io/gorm"
	"time"
)

type PayRecharge struct {
	gorm.Model
	// 充值ID
	Uid string `gorm:"type:varchar(40);not null;unique;comment:充值ID"`
	// 充值企业
	CorpID string `gorm:"type:varchar(40);not null;comment:充值企业"`
	// 支付金额
	Price uint `gorm:"type:int(12);not null;comment:支付金额"`
	// 支付状态
	Paid int `gorm:"type:tinyint(1);default:0;comment:支付状态 0：未支付 1：支付完成"`
	// 支付时间
	PaidTime time.Time `gorm:"comment:支付时间"`
	//充值备注
	Remark string `gorm:"type:text;comment:充值备注信息"`
	// 支付渠道
	ChargeType string `gorm:"type:varchar(40);default:wx_pay;comment:支付类型：wx_pay(微信支付)"`
	// 支付订单ID
	ChargeID string`gorm:"index;type:varchar(64);not null;comment:支付订单ID"`
	// 交易ID
	TransactionID string `gorm:"type:varchar(128);default:'';comment:订单交易ID"`
	// 操作用户ID
	UserID string `gorm:"type:varchar(40);not null;comment:操作用户ID"`
}

