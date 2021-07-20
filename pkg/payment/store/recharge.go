package store

import (
	"flying-star/pkg/payment/model"
	"flying-star/pkg/payment/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type RechargeRecordSetOptions struct {
	//付款用户
	UserId string `json:"userId"`
	//付款企业
	CorpId string `json:"corpId"`
	//付款金额
	Price uint `json:"price"`
	//支付备注
	Remark string `json:"remark"`
	//订单ID
	ChargeID string `json:"chargeId"`
}

// RechargeRecordSet 添加充值记录
func (s *Store) RechargeRecordSet(moduleName string, options RechargeRecordSetOptions) (*gorm.DB, error) {
	ts := s.client.Begin()
	record := model.PayRecharge{
		Uid:        utils.GetUid(),
		CorpID:     options.CorpId,
		ChargeType: moduleName,
		Price:      options.Price,
		PaidTime:   time.Now(),
		Remark:     options.Remark,
		UserID:     options.UserId,
		Paid:       0,
		ChargeID:  options.ChargeID,
	}
	if err := ts.Create(&record).Error; err != nil {
		ts.Rollback()
		return nil, err
	}

	//查询当前用户账户余额信息，如果不存在，则创建新的账户
	centerRecord := model.PayCenter{}

	if err := ts.Where("account_id = ?", options.CorpId).First(&centerRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			centerRecord.AccountId = options.CorpId
			if err = ts.Create(&centerRecord).Error; err != nil {
				ts.Rollback()
				return nil, err
			}
		} else {
			ts.Rollback()
			return nil, err
		}
	}

	return ts, nil
}

type RechargeRecordUpdateOptions struct {
	// 支付订单ID
	ChargeID string `json:"chargeId"`
	// 交易ID
	TransactionID string `json:"transactionId"`
	//支付金额
	Price uint `json:"price"`
}

// RechargeRecordUpdate 更新充值记录
func (s *Store) RechargeRecordUpdate(moduleName string, options RechargeRecordUpdateOptions) (*gorm.DB, string,error) {
	ts := s.client.Begin()

	//查询订单详情
	recordInfo := model.PayRecharge{}
	if err := ts.Where("charge_id = ? and charge_type = ?", options.ChargeID, moduleName).First(&recordInfo).Error; err != nil {
		//订单不存在
		if err == gorm.ErrRecordNotFound {
			return nil, "", nil
		}
		return nil, "", err
	}

	//订单支付金额异常
	if options.Price != recordInfo.Price {
		return nil, "", nil
	}

	record := model.PayRecharge{
		TransactionID: options.TransactionID,
		Paid:       1,
	}
	if err := ts.Model(model.PayRecharge{}).Where("charge_id = ? and charge_type = ?", options.ChargeID, moduleName).Updates(&record).Error; err != nil {
		ts.Rollback()
		return nil, "", err
	}

	//更新账户金额
	if err := ts.Clauses(clause.Locking{Strength: "UPDATE"}).Model(model.PayCenter{}).Where("account_id", recordInfo.CorpID).Update("balance", gorm.Expr("balance + ?", options.Price)).Error; err != nil {
		ts.Rollback()
		return nil, "", err
	}

	return ts, recordInfo.UserID, nil
}

// RechargeRecordList 获取指定账户下的充值列表
func (s *Store) RechargeRecordList(accountId string, page, size int) (list []model.PayRecharge, err error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 1000 {
		size = 20
	}
	if err = s.client.Model(model.PayRecharge{}).Where("corp_id = ? and paid = 1", accountId).Limit(size).Offset((page - 1) * size).Order("created_at desc").Find(&list).Error; err != nil {
		return list, err
	}
	return list, nil
}
