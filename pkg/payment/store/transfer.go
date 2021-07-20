package store

import (
	"errors"
	"flying-star/pkg/payment/model"
	"flying-star/pkg/payment/utils"
	"gorm.io/gorm"
	"time"
)

type TransferRecordSetOptions struct {
	//付款用户
	UserId string `json:"userId"`
	//付款企业
	CorpId string `json:"corpId"`
	//付款金额
	Price uint `json:"price"`
	//收款用户ID
	PayeeId string `json:"payeeId"`
	//收款真实姓名
	PayeeName string `json:"payeeName"`
	//业务模块
	Module string `json:"module"`
	//支付备注
	Remark string `json:"remark"`
	//外部关联ID
	ReferId string `json:"referId"`
}

// TransferRecordSet 添加转账记录
func (s *Store) TransferRecordSet(moduleName string, options TransferRecordSetOptions) (*gorm.DB, error) {
	ts := s.client.Begin()

	//扣减账户余额
	center := ts.Debug().Model(model.PayCenter{}).Where("account_id = ? and balance > ?", options.CorpId, options.Price).Update("balance", gorm.Expr("balance - ?", options.Price))

	if err := center.Error; err != nil {
		ts.Rollback()
		return nil, err
	}
	if center.RowsAffected < 1 {
		ts.Rollback()
		return nil, errors.New("余额不足，请充值")
	}

	record := model.PayTransfer{
		Uid:        utils.GetUid(),
		CorpID:     options.CorpId,
		Payee:      options.PayeeId,
		PayeeName:  options.PayeeName,
		ChargeType: moduleName,
		Price:      options.Price,
		PaidTime:   time.Now(),
		Module:     options.Module,
		Remark:     options.Remark,
		UserID:     options.UserId,
		ReferID:    options.ReferId,
	}
	if err := ts.Create(&record).Error; err != nil {
		ts.Rollback()
		return nil, err
	}
	return ts, nil
}

type TransferRecordOptionalParams struct {
	Module string `json:"module"`		// 业务模块
	ReferId string `json:"referId"`		// 外部关联ID
	Remark string `json:"remark"`		// 备注信息
}

// TransferRecordList 获取指定账户下的消费记录
func (s *Store) TransferRecordList(accountId string, page, size int, options TransferRecordOptionalParams) (list []model.PayTransfer, err error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 1000 {
		size = 20
	}
	sql := s.client.Model(model.PayTransfer{}).Where("corp_id = ?", accountId)
	if options.Module != "" {
		sql = sql.Where("module = ?", options.Module)
	}
	if options.ReferId != "" {
		sql = sql.Where("refer_id = ?", options.ReferId)
	}
	if options.Remark != "" {
		sql = sql.Where("remark = ?", options.Remark)
	}
	if err = sql.Limit(size).Offset((page - 1) * size).Order("created_at desc").Find(&list).Error; err != nil {
		return list, err
	}
	return list, nil
}

