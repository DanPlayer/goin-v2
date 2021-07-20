package model

import (
	"flying-star/internal/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	gorm.Model
	// 用户唯一标识符 UUID
	Uid string `gorm:"<-:create;type:varchar(40);not null;unique;comment:用户唯一标识符"`
	// 账户名
	UserName string `gorm:"index;type:varchar(32);unique;default:'';comment:账户名"`
	// 账户密码,MD5加密
	Password string `gorm:"->:false;<-:create;index;type:varchar(64);default:'';comment:账户密码,MD5加密"`
	// 用户昵称
	NickName string `gorm:"type:varchar(64);default:'';comment:用户昵称"`
	// 用户头像
	Avatar string `gorm:"type:varchar(500);default:'';comment:用户头像"`
	// 性别 0：未知 1：男 2：女
	Sex int `gorm:"type:tinyint(1);default:0;comment:性别 0：未知 1：男 2：女"`
	// 用户手机号
	Mobile string `gorm:"type:varchar(20);default:'';comment:用户手机号"`
	// 用户邮箱
	MailBox string `gorm:"type:varchar(32);default:'';comment:用户邮箱"`
	// 用户角色
	RoleCode string `gorm:"index;type:varchar(40);default:'basic';comment:用户所属角色"`
}

func (user *User) TableName() string {
	return "user"
}

// Info 用户信息
func (user *User) Info() (info User, err error) {
	err = db.MySQLClient.Where(&user).First(&info).Error
	return
}

// Find 用户列表
func (user *User) Find() (list []User, err error) {
	err = db.MySQLClient.Where(&user).Find(&list).Error
	return
}

// Upsert 创建用户
func (user *User) Upsert() error {
	return db.MySQLClient.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "uid"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"user_name", "password", "nick_name", "avatar", "sex",
			"mobile", "mail_box", "role_code",
		}),
	}).Create(&user).Error
}

// Update 更新用户信息
func (user *User) Update(data User) error {
	return db.MySQLClient.Model(&user).Where(&user).Updates(data).Error
}

// Delete 删除用户信息
func (user *User) Delete(real bool) error {
	if real {
		return db.MySQLClient.Unscoped().Where(&user).Delete(&user).Error
	} else {
		return db.MySQLClient.Where(&user).Delete(&user).Error
	}
}

// FindByIds 根据ids获取数据
func (user *User) FindByIds(ids []string) (users []User, err error) {
	err = db.MySQLClient.Model(&user).Where("uid in ?", ids).Find(&users).Error
	return
}