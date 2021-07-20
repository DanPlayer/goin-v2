package optimisticlock

import (
	"errors"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

var UpdateTimeout = errors.New("update timeout")

type Model struct {
	gorm.Model
	Version uint `gorm:"type:int(12);default:1;comment:版本号"`
}

func (r Model) BeforeUpdate(tx *gorm.DB) error {
	//校验当前操作版本号是否一致
	tx.Statement.Where("version = ?", r.Version)
	tx.Statement.SetColumn("Version", gorm.Expr("version + ?", 1))
	return nil
}

func (r Model) AfterUpdate(tx *gorm.DB) error {
	//确保更新操作执行成功
	if tx.Statement.DB.RowsAffected < 1 {
		retryTimeout := 5 * time.Second
		currentTime := time.Now()
		sql := tx.Dialector.Explain(tx.Statement.SQL.String(), tx.Statement.Vars...)
		oldVersion := "AND version = " + strconv.Itoa(int(r.Version))

		retry:
			version := uint(0)

			whereList := strings.Split(sql, "WHERE")
			if len(whereList) == 2 {
				//重试
				querySql := strings.Replace("SELECT version FROM "+ tx.Statement.Table +" WHERE" + whereList[1], oldVersion, "", -1)
				tx.Raw(querySql).Scan(&version)
				if version != 0 {
					newVersion := "AND version = " + strconv.Itoa(int(version))
					updateSql := strings.Replace(sql, oldVersion, newVersion, -1)
					if tx.Exec(updateSql).RowsAffected > 0 {
						return nil
					}
				}
			}

			if time.Now().Unix() < currentTime.Add(retryTimeout).Unix() {
				goto retry
			}

		return UpdateTimeout
	}
	return nil
}