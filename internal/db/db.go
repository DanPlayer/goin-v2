package db

import (
	"flying-star/internal/config"
	"flying-star/pkg/db"
	"gorm.io/gorm"
	"time"
)

// RedisClient 初始化Redis客户端
var RedisClient = db.NewRedis(db.RedisOptions{
	Addr:     config.Info.Redis.Addr,
	Password: config.Info.Redis.Password,
	DB:       config.Info.Redis.DB,
})

// MySQLClient 初始化MySQL客户端
var MySQLClient = db.NewMySQL(db.Config{
	DNS: config.Info.MySQL.DSN,
})

// AutoMigrateApiServerTable 自动迁移数据表
func AutoMigrateApiServerTable(modes ...interface{}) {
	_ = MySQLClient.AutoMigrate(modes...)
}

type Model struct {
	ID        uint   `gorm:"primarykey"`
	Uid       string `gorm:"primarykey;default:uuid();"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
