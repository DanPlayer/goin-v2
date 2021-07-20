package notification

import (
	"flying-star/internal/db"
	"flying-star/pkg/eventnotification"
)

var EventNotification *eventnotification.EventNotification

const (
	NoEvent      = "noEvent"
	// 可自定义
	//VideoReleaseSuccess = "videoReleaseSuccess"
	//RechargeSuccess="rechargeSuccess"
)

type WsEventSchema struct {
	Event   string      `json:"event"`   // 事件类型：noEvent(暂无消息), videoReleaseSuccess(视频发布成功), rechargeSuccess(充值成功)
	Message interface{} `json:"message"` // 事件消息内容
}

func init() {
	EventNotification = eventnotification.New(eventnotification.Options{Cache: db.RedisClient})
}
