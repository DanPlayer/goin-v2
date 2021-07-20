package eventnotification

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Options 初始化参数
type Options struct {
	ExpireTime uint        `default:"60"`  // AccessToken过期时间
	MaxBuffer  uint        `default:"600"` // 消息最大缓冲数量
	Cache      MemoryCache // 缓存
}

type MemoryCache interface {
	Set(k, v string, expires time.Duration) error
	Get(ctx context.Context,k string) (string, error)
	Del(ctx context.Context, k string) error
	Scan(cursor uint64, match string, count int64) (keys []string, newCursor uint64, err error)
}

type Message struct {
	Name    string // 消息名称
	Content []byte // 消息内容
}

type EventNotification struct {
	maxBuffer  uint
	expireTime uint        // 过期时间
	queue      sync.Map    // 消息队列
	mutex      sync.Mutex  // 互斥锁
	cache      MemoryCache // 消息存储
}

type EventRegistryObj struct {
	uid        string // 内部唯一标识符
	registryId string // 消息注册码
	event      *EventNotification
	once       sync.Once
}

var ctx = context.TODO()

// New 初始化事件通知服务
func New(options Options) *EventNotification {
	if options.ExpireTime == 0 {
		field, _ := reflect.TypeOf(options).FieldByName("ExpireTime")
		expireTime, _ := strconv.Atoi(field.Tag.Get("default"))
		options.ExpireTime = uint(expireTime)
	}
	if options.MaxBuffer == 0 {
		field, _ := reflect.TypeOf(options).FieldByName("MaxBuffer")
		maxBuffer, _ := strconv.Atoi(field.Tag.Get("default"))
		options.MaxBuffer = uint(maxBuffer)
	}
	eventNotification := EventNotification{cache: options.Cache, expireTime: options.ExpireTime, maxBuffer: options.MaxBuffer}
	//启动消息中心服务
	go eventNotification.runMessageCenter()
	return &eventNotification
}

// Registry 注册
func (r *EventNotification) Registry(registryId string) *EventRegistryObj {
	registryId = r.formatTxt(registryId)

	registry := EventRegistryObj{}
	registry.uid = uuid.NewV4().String()
	registry.registryId = registryId
	registry.event = r
	registry.once = sync.Once{}

	memberQueue, ok := r.queue.Load(registryId)
	if !ok {
		newMemberQueue := &sync.Map{}
		r.queue.Store(registryId, newMemberQueue)
		memberQueue = newMemberQueue
	}
	if memberDic, ok := memberQueue.(*sync.Map); ok {
		memberDic.Store(registry.uid, make(chan Message, r.maxBuffer))
	}

	return &registry
}

// UnRegistry 取消注册
func (r *EventRegistryObj) UnRegistry() bool {
	memberQueue, ok := r.event.queue.Load(r.registryId)
	if !ok {
		return false
	}
	memberDic, ok := memberQueue.(*sync.Map)
	if !ok {
		return false
	}
	msgQueue, ok := memberDic.Load(r.uid)
	if !ok {
		return false
	}
	msgDic, ok := msgQueue.(chan Message)
	if !ok {
		return false
	}
	r.once.Do(func() {
		close(msgDic)
	})
	memberDic.Delete(r.uid)
	return true
}

// GetMessage 获取全局消息
func (r *EventRegistryObj) GetMessage() (ch <-chan Message) {
	memberQueue, ok := r.event.queue.Load(r.registryId)
	if !ok {
		return
	}
	memberDic, ok := memberQueue.(*sync.Map)
	if !ok {
		return
	}
	msgQueue, ok := memberDic.Load(r.uid)
	if !ok {
		return
	}
	msgDic, ok := msgQueue.(chan Message)
	if !ok {
		return
	}
	return msgDic
}

// PutMessage 推送全局消息
func (r *EventNotification) PutMessage(registryId string, msg Message) (err error) {
	return r.setOriginMessage(registryId, msg.Name, string(msg.Content))
}

//runMessageCenter 消息中心
func (r *EventNotification) runMessageCenter() {
	for {
		time.Sleep(200 * time.Millisecond)
		stateList := r.scanOriginMessage()
		if len(stateList) == 0 {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		// 读取最新全局事件
		for _, state := range stateList {
			values := strings.Split(state, ".")
			if len(values) != 2 {
				continue
			}
			keys := strings.Split(values[0], ":")
			if len(keys) != 2 {
				continue
			}
			registryID := keys[1]
			msgName := values[1]
			msgContent := r.getOriginMessage(registryID, msgName)
			if msgName == "" {
				continue
			}

			//广播消息
			message := Message{
				Name:    msgName,
				Content: []byte(msgContent),
			}

			//获取同一注册码下所有成员消息队列
			if memberQueue, ok := r.queue.Load(registryID); ok {
				//向所有成员分发消息
				if memberList, ok := memberQueue.(*sync.Map); ok {
					memberList.Range(func(key, member interface{}) bool {
						if member, ok := member.(chan Message); ok {
							member <- message
						}
						return true
					})
				}
			}

			//消息置为已消费
			_ = r.delOriginMessage(registryID, msgName)
		}
	}
}

// getOriginMessage 获取原始消息
func (r *EventNotification) getOriginMessage(namespace, state string) string {
	namespace = r.formatTxt(namespace)
	state = r.formatTxt(state)
	key := "event:" + namespace + "." + state
	res, err := r.cache.Get(ctx, key)
	if err != nil {
		return ""
	}
	return res
}

// setOriginMessage 设置原始消息
func (r *EventNotification) setOriginMessage(namespace, state, msg string) error {
	namespace = r.formatTxt(namespace)
	state = r.formatTxt(state)
	key := "event:"+namespace+"."+state
	expire := time.Duration(int64(r.expireTime))
	return r.cache.Set(key, msg, expire)
}

// delOriginMessage 删除原始消息
func (r *EventNotification) delOriginMessage(namespace, state string) error {
	namespace = r.formatTxt(namespace)
	state = r.formatTxt(state)
	return r.cache.Del(ctx, "event:" + namespace + "." + state)
}

// scanOriginMessage 扫描原始消息
func (r *EventNotification) scanOriginMessage() (msgList []string) {
	list, _, err := r.cache.Scan(0, "event:*", 1000)
	if err != nil {
		fmt.Printf("全局共享状态扫描失败：%v", err.Error())
	}
	return list
}

var txtReg = regexp.MustCompile(`[_:.@-]+`)
func (r *EventNotification) formatTxt(txt string) string {
	return txtReg.ReplaceAllString(txt, "")
}
