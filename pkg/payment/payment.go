package payment

import (
	"flying-star/pkg/payment/store"
	"net/http"
)

type Payment struct {
	engine Plugin
}

// Load 加载支付插件
func Load(plugin Plugin) *Payment {
	return &Payment{engine: plugin}
}

func (p *Payment) BindAccount(options BindAccountOptions) error {
	return p.engine.BindAccount(options)
}

func (p *Payment) RechargeBegin(options RechargeOptions) (RechargeSchema, error) {
	return p.engine.RechargeBegin(options)
}

func (p *Payment) RechargeEnd(c *http.Request) (userId string, err error) {
	return p.engine.RechargeEnd(c)
}

func (p *Payment) Transfer(options TransferOptions) error {
	return p.engine.Transfer(options)
}


type ProfileSchema struct {
	Remain uint	`json:"remain"`				// 剩余金额
	Consume uint `json:"consume"`			// 已消费金额
}

// GetProfile 获取指定用户数据概况
func (p *Payment) GetProfile(accountId string) (info ProfileSchema, err error) {
	store := p.engine.GetStore()

	profile, err := store.GetProfile(accountId)
	if err != nil {
		return info, err
	}
	return ProfileSchema(profile), nil
}

type RechargeRecordSchema struct {
	Price uint `json:"price"`			// 充值金额
	UserId string `json:"userId"`		// 充值用户
	PaidTime string `json:"paidTime"`   // 充值时间
}

// GetRechargeRecordList 获取指定账户下的充值记录
func (p *Payment) GetRechargeRecordList(accountId string, page, size int) (list []RechargeRecordSchema, err error) {
	store := p.engine.GetStore()
	dataList, err := store.RechargeRecordList(accountId, page, size)
	if err != nil {
		return nil, err
	}
	for _, data := range dataList {
		list = append(list, RechargeRecordSchema{
			Price:    data.Price,
			UserId:   data.UserID,
			PaidTime: data.PaidTime.Format("2006-01-02 15:04"),
		})
	}
	return list, nil
}

type TransferRecordOptionalParams struct {
	Module string `json:"module"`		// 业务模块
	ReferId string `json:"referId"`		// 外部关联ID
	Remark string `json:"remark"`		// 备注信息
}

type TransferRecordSchema struct {
	UserId string `json:"userId"`		// 操作用户
	Payee string `json:"payee"`			// 收款人
	PayeeName string `json:"payeeName"` // 收款人姓名
	PaidTime string `json:"paidTime"`   // 充值时间
	Price uint `json:"price"`			// 充值金额
	Module string `json:"module"`		// 业务模块
	Remark string `json:"remark"`		// 备注信息
	ReferId string `json:"referId"`		// 外部关联ID
}

// GetTransferRecordList 获取指定账户下的转账记录
func (p *Payment) GetTransferRecordList(accountId string, page, size int, options TransferRecordOptionalParams) (list []TransferRecordSchema, err error) {
	storeIns := p.engine.GetStore()
	dataList, err := storeIns.TransferRecordList(accountId, page, size, store.TransferRecordOptionalParams(options))
	if err != nil {
		return nil, err
	}
	for _, data := range dataList {
		list = append(list, TransferRecordSchema{
			UserId:   data.UserID,
			Payee:    data.Payee,
			PayeeName: data.PayeeName,
			Price:    data.Price,
			PaidTime: data.PaidTime.Format("2006-01-02 15:04"),
			Module:   data.Module,
			Remark:   data.Remark,
		})
	}
	return list, nil
}
