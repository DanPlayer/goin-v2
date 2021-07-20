package model

import (
	"gorm.io/gorm"
	"time"
)

type PayTransfer struct {
	gorm.Model
	// 转账记录ID
	Uid string `gorm:"type:varchar(40);not null;unique;comment:转账记录ID"`
	// 付款企业ID
	CorpID string `gorm:"type:varchar(40);not null;comment:付款企业ID"`
	// 收款用户ID
	Payee string `gorm:"type:varchar(40);not null;comment:收款用户ID"`
	//收款真实姓名
	PayeeName string `gorm:"type:varchar(64);default:'';comment:收款真实姓名"`
	// 付款渠道
	ChargeType string `gorm:"type:varchar(40);default:wx_pay;comment:支付类型：wx_pay(微信支付)"`
	// 业务模块
	Module string `gorm:"type:varchar(40);default:'';comment:业务模块：task(任务激励)"`
	// 付款金额
	Price uint `gorm:"type:int(12);not null;comment:付款金额"`
	// 付款时间
	PaidTime time.Time `gorm:"comment:付款时间"`
	//收款备注信息
	Remark string `gorm:"type:text;comment:收款备注信息"`
	// 操作用户ID
	UserID string `gorm:"type:varchar(40);not null;comment:操作用户ID"`
	//外部关联ID
	ReferID string `gorm:"type:varchar(40);default:'';comment:外部关联ID"`
}