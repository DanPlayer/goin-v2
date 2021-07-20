package wechat

import "encoding/xml"

type Message struct {
	BaseMessage
	originData []byte
}

// BaseMessage 基础消息
type BaseMessage struct {
	ToUserName   string // 开发者微信号
	FromUserName string // 发送方帐号（一个OpenID）
	CreateTime   int    // 消息创建时间 （整型）
	MsgType      string // 消息类型
	Event        string // 事件类型
}

// Text 文本消息(text)
type Text struct {
	BaseMessage
	MsgId   int64  // 消息id，64位整型
	Content string // 文本消息内容
}

// Image 图片消息(image)
type Image struct {
	BaseMessage
	MsgId   int64  // 消息id，64位整型
	PicUrl  string // 图片链接（由系统生成）
	MediaId string // 图片消息媒体id，可以调用获取临时素材接口拉取数据
}

// Voice 语音消息(voice)
type Voice struct {
	BaseMessage
	MsgId       int64  // 消息id，64位整型
	MediaId     string // 图片消息媒体id，可以调用获取临时素材接口拉取数据
	Format      string // 语音格式：amr
	Recognition string // 语音识别结果，UTF8编码
}

// Video 视频消息(video)
type Video struct {
	BaseMessage
	MsgId        int64  // 消息id，64位整型
	MediaId      string // 图片消息媒体id，可以调用获取临时素材接口拉取 数据
	ThumbMediaId string // 视频消息缩略图的媒体id，可以调用多媒体文件下载接口拉取数据
}

// ShortVideo 小视频消息(shortvideo)
type ShortVideo struct {
	BaseMessage
	MsgId        int64  // 消息id，64位整型
	MediaId      string // 图片消息媒体id，可以调用获取临时素材接口拉取数据
	ThumbMediaId string // 视频消息缩略图的媒体id，可以调用多媒体文件下载接口拉取数据
}

// Location 地理位置消息(location)
type Location struct {
	BaseMessage
	MsgId     int64   // 消息id，64位整型
	LocationX float32 `xml:"Location_X"` // 地理位置纬度
	LocationY float32 `xml:"Location_Y"` // 地理位置经度
	Scale     float32 // 地图缩放大小
	Label     string  // 地理位置信息
}

// Link 链接消息(link)
type Link struct {
	BaseMessage
	MsgId       int64  // 消息id，64位整型
	Title       string // 消息标题
	Description string // 消息描述
	Url         string // 消息链接
}

// UnsubscribeEvent 取消关注事件(unsubscribe)
type UnsubscribeEvent struct {
	BaseMessage
}

// SubscribeEvent 关注或扫码关注事件(subscribe)
type SubscribeEvent struct {
	BaseMessage
	EventKey string // 事件KEY值，qrscene_为前缀，后面为二维码的参数值
	Ticket   string // 二维码的ticket，可用来换取二维码图片
}

// ScanEvent 扫描二维码事件(SCAN)
type ScanEvent struct {
	BaseMessage
	EventKey string // 事件KEY值，是一个32位无符号整数，即创建二维码时的二维码scene_id
	Ticket   string // 二维码的ticket，可用来换取二维码图片
}

// LocationEvent 上报地理位置事件(LOCATION)
type LocationEvent struct {
	BaseMessage
	Latitude  float32 // 地理位置纬度
	Longitude float32 // 地理位置经度
	Precision string  // 地理位置精度
}

// ClickEvent 自定义菜单点击拉取消息事件(CLICK)
type ClickEvent struct {
	BaseMessage
	EventKey string // 事件KEY值，与自定义菜单接口中KEY值对应
}

// ViewEvent 自定义菜单点击访问页面事件(VIEW)
type ViewEvent struct {
	BaseMessage
	EventKey string // 事件KEY值，设置的跳转URL
}

// NewMessage 初始化新的消息
func NewMessage(data []byte) (msg *Message) {
	_ = xml.Unmarshal(data, &msg)
	msg.originData = data
	return msg
}

// GetTextMessage 获取文本消息
func (m *Message) GetTextMessage() (msg Text) {
	_ = xml.Unmarshal(m.originData, &msg)
	return msg
}

// GetImageMessage 获取图片消息
func (m *Message) GetImageMessage() (msg Image) {
	_ = xml.Unmarshal(m.originData, &msg)
	return msg
}

// GetVoiceMessage 获取语音消息
func (m *Message) GetVoiceMessage() (msg Voice) {
	_ = xml.Unmarshal(m.originData, &msg)
	return msg
}

// GetLinkMessage 获取链接消息
func (m *Message) GetLinkMessage() (msg Link) {
	_ = xml.Unmarshal(m.originData, &msg)
	return msg
}

// GetSubscribeEventMessage 获取用户关注或扫码事件消息
func (m *Message) GetSubscribeEventMessage() (msg SubscribeEvent) {
	_ = xml.Unmarshal(m.originData, &msg)
	return msg
}

// GetScanEventMessage 获取用户扫描二维码事件消息
func (m *Message) GetScanEventMessage() (msg ScanEvent) {
	_ = xml.Unmarshal(m.originData, &msg)
	return msg
}
