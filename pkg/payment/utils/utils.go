package utils

import uuid "github.com/satori/go.uuid"

// GetUid 获取唯一ID
//目前为了简便，直接采用了UUID v4生成，后期可以根据项目需要切换别的算法，例如雪花算法
func GetUid() string {
	return uuid.NewV4().String()
}
