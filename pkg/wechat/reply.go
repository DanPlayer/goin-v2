package wechat

import (
	"encoding/xml"
	"time"
)

// TextMsgXML XML文本消息
type TextMsgXML struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string
	// 若不标记XMLName, 则解析后的xml名为该结构体的名称
	XMLName xml.Name `xml:"xml"`
}

// GenerateTextMsgByXML 生成XML格式的文本消息
func GenerateTextMsgByXML(from string, to string, msg string) []byte {
	data := TextMsgXML{
		ToUserName:   to,
		FromUserName: from,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      msg,
		XMLName:      xml.Name{},
	}
	b, _ := xml.Marshal(&data)
	return b
}
