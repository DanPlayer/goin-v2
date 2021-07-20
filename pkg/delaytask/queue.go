package delaytask

import (
	"github.com/go-redis/redis/v8"
	"strconv"
	"sync"
	"time"
)

// QueueOptions 延时任务队列初始参数
type QueueOptions struct {
	NameSpace string 							// 命名空间
	Cache MemoryCache							// 缓存
	ConsumeMaxCount int64	`default:"10000"`	// 单次最大消费任务数量,默认：1
}

type Queue struct {
	basename string
	namespace string						// 命名空间
	cache MemoryCache						// 缓存
	mutex sync.Mutex						// 锁
	redis *redis.Client						// redis实例
	consumeMaxCount int64					// 单次最大消费任务数量,默认：1
	expiredHandleList []func(key string) bool	// 令牌过期处理
}

// NewQueue 初始化新的延时任务队列
func NewQueue(options QueueOptions) *Queue {
	queue := Queue{
		basename: "DelayTask",
		namespace: options.NameSpace,
		cache:     options.Cache,
		mutex:     sync.Mutex{},
		redis: options.Cache.GetOriginPoint(),
	}

	if queue.consumeMaxCount == 0 {
		consumeMaxCount, _ := strconv.Atoi(getStructTagDefaultValue(options, "ConsumeMaxCount"))
		queue.consumeMaxCount = int64(consumeMaxCount)
	}

	//处理全局过期事件
	queue.cache.SubscribeAllEvents(func(message string) {
		//过滤无效事件
		if !isUsableKey(message, queue.namespace) {
			return
		}
		//写入任务到待消费队列
		originKey := getOriginKey(message)
		contentKey := transformKey2Content(message)
		queue.redis.LPush(ctx, getConsumeQueueName(queue.namespace), originKey)
		queue.redis.Del(ctx, contentKey)
		for _, handle := range queue.expiredHandleList {
			if status := handle(originKey); status {
				_ = queue.Ack(originKey)
			}
		}
	})

	// 自动恢复当前实例历史未完成任务
	taskListKey, _, _ := queue.redis.Scan(ctx,0, getQueuePrefix(queue.namespace) + "*", 100000).Result()
	//批量恢复任务
	for _, taskKey := range taskListKey {
		curKey := transformContent2Key(taskKey)
		if content, err := queue.redis.Get(ctx, curKey).Result(); err != nil && err == redis.Nil {
			queue.redis.Set(ctx, curKey, content, 1 * time.Millisecond)
		}
	}

	return &queue
}

// Set 设置新的延时任务
func (r *Queue) Set(key string,  expire time.Duration) (err error) {
	expireTime := time.Now().Add(expire).Unix()
	_, err = r.redis.ZAdd(ctx, WaitQueue, &redis.Z{
		Score:  float64(expireTime),
		Member: QueueMsgSchema{
			NameSpace: r.namespace,
			ID:    key,
		},
	}).Result()
	return err
}

// SetWithTime 设置新的延时任务(设置指定时间)
func (r *Queue) SetWithTime(key string, expire time.Time) (err error) {
	expireTime := time.Now().Unix()
	if expire.After(time.Now()) {
		expireTime = expire.Unix()
	}
	_, err = r.redis.ZAdd(ctx, WaitQueue, &redis.Z{
		Score:  float64(expireTime),
		Member: QueueMsgSchema{
			NameSpace: r.namespace,
			ID:    key,
		},
	}).Result()
	return err
}

// Remove 移除指定延时任务
func (r *Queue) Remove(key string)  {
	r.redis.ZRem(ctx, WaitQueue, redis.Z{Member: QueueMsgSchema{
		NameSpace: r.namespace,
		ID:        key,
	}}.Member)

	r.redis.Del(ctx, generateKey(key, r.namespace))
	r.redis.Del(ctx, generateContentKey(key, r.namespace))

	r.redis.LRem(ctx, getConsumeQueueName(r.namespace), -1, key)
}

// Consume 消费到期任务
func (r *Queue) Consume() ([]string, error) {
	return r.redis.LRange(ctx, getConsumeQueueName(r.namespace),0, r.consumeMaxCount - 1).Result()
}

// Subscribe 订阅到期任务
func (r *Queue) Subscribe(fn func(key string) bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.expiredHandleList = append(r.expiredHandleList, fn)
}

// Ack 到期任务确认消费完成
func (r *Queue) Ack(key string) error {
	return r.redis.LRem(ctx, getConsumeQueueName(r.namespace), -1, key).Err()
}
