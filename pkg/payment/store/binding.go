package store

import (
	"flying-star/pkg/payment/model"
	"gorm.io/gorm/clause"
)

// BindingSet 存储外部账号关联关系
func (s *Store) BindingSet(userId, openId, realName, moduleName string) error {
	bindInfo := model.PayBinding{
		UserID:     userId,
		OpenID:     openId,
		RealName:   realName,
		ModuleName: moduleName,
	}
	return s.client.Clauses(clause.OnConflict{
		Columns:      []clause.Column{{Name: "user_id"}},
		DoUpdates:    clause.AssignmentColumns([]string{"open_id", "real_name"}),
	}).Create(&bindInfo).Error
}

// BindingGet 获取外部关联ID
func (s *Store) BindingGet(userId, moduleName string) (info model.PayBinding, err error) {
	bindInfo := model.PayBinding{
		UserID:     userId,
		ModuleName:     moduleName,
	}
	if err := s.client.Model(model.PayBinding{}).Where(bindInfo).First(&info).Error; err != nil {
		return info, err
	}
	return info, nil
}