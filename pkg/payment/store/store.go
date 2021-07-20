package store

import (
	"flying-star/pkg/payment/model"
	"gorm.io/gorm"
)

type Options struct {
	Client *gorm.DB
}

type Store struct {
	client *gorm.DB
}

func New(options Options) (*Store, error) {
	st := Store{
		client: options.Client,
	}
	//初始化相关表
	if err := st.client.AutoMigrate(model.PayBinding{}, model.PayCenter{}, model.PayRecharge{}, model.PayTransfer{}); err != nil {
		return nil, err
	}
	return &st, nil
}

type ProfileSchema struct {
	Remain uint				// 剩余金额
	Consume uint			// 已消费金额
}

// GetProfile 获取制定账户概况数据
func (s *Store) GetProfile(accountId string) (info ProfileSchema, err error) {
	balance := uint(0)
	consume := uint(0)

	if err = s.client.Model(model.PayCenter{}).Select("balance").Where("account_id = ?", accountId).Scan(&balance).Error; err != nil {
		return info, err
	}

	if err = s.client.Model(model.PayTransfer{}).Select("sum(price)").Where("corp_id = ?", accountId).Group("corp_id").Scan(&consume).Error; err != nil {
		return info, err
	}

	info.Remain = balance
	info.Consume = consume
	return info, nil
}