package middleware

import (
	"flying-star/internal/config"
	"flying-star/internal/module/user/rdb"
	"flying-star/utils"
	"strconv"
	"time"
)

// MakeAdminToken 生产管理后台的token密钥
func MakeAdminToken(info rdb.Info) (token string, err error) {
	token = utils.Md5hex(config.Info.Server.Token + info.Uid + strconv.FormatInt(time.Now().Unix(), 10))
	// 设置token缓存
	tokenRdb := rdb.Token{Uid: info.Uid, Token: token}
	_ = tokenRdb.Set()
	_ = info.Set(token, config.Info.Server.Token)
	return
}

// ParseAdminToken 验证管理后台的token密钥
func ParseAdminToken(token, secret string) (id, role string, err error) {
	userRdb := rdb.Info{}
	info, err := userRdb.Get(token, secret)
	if err != nil {
		return
	}
	id = info.Uid
	role = info.Roles
	return
}
