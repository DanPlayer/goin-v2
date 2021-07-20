package wechat

import (
	"fmt"
)

type Media struct {
	VoiceUrl string `json:"voice_url"`
}

func (w *Wechat) GetMediaByID(mediaID string) (media Media, err error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/media/get?access_token=%v&media_id=%v", w.accessToken.token, mediaID)
	media.VoiceUrl = url
	return
}
