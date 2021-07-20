package delaytask

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"sync"
	"time"
)

var ctx = context.TODO()

type Options struct {
	Addr            string        							// redis链接地址
	Password        string        							// 访问密码
	DB              int           							// 指定数据库，默认：0
	Cache           MemoryCache   							// 缓存
	RefreshTime     time.Duration `default:"100"` 			// 刷新任务到执行队列时间（毫秒），默认：1000
	RefreshMaxNum   int           `default:"100000"`    	// 单次最大刷新任务数量，默认：1
	RefreshDuration time.Duration `default:"1800"`    		// 刷新任务的时间范围（秒），默认：3
}

type DelayTask struct {
	//订阅服务器实例
	cache MemoryCache
	//订阅列表
	pbFns sync.Map
	//读写锁
	lock sync.Mutex
	// 刷新任务到执行队列时间（毫秒），默认：1000
	refreshTime time.Duration
	// 单次最大刷新任务数量，默认：1
	refreshMaxNum int64
	//刷新任务的时间范围（秒），默认：3
	refreshDuration time.Duration
	//redis实例
	redis *redis.Client
}

type MemoryCache interface {
	GetOriginPoint() *redis.Client
	Subscribe(k string, pb func(message string))
	SubscribeAllEvents(pb func(message string))
}

func New(options Options) *DelayTask {
	ins := DelayTask{}

	//实例化redis连接池
	if options.Cache != nil {
		ins.cache = options.Cache
	} else {
		ins.cache = NewRedis(RedisOptions{
			Addr:     options.Addr,
			Password: options.Password,
			DB:       options.DB,
		})
	}

	ins.pbFns = sync.Map{}
	ins.lock = sync.Mutex{}
	ins.redis = ins.cache.GetOriginPoint()

	if ins.refreshTime == 0 {
		refreshTime, _ := strconv.Atoi(getStructTagDefaultValue(options, "RefreshTime"))
		ins.refreshTime = time.Duration(refreshTime) * time.Millisecond
	}

	if ins.refreshMaxNum == 0 {
		refreshMaxNum, _ := strconv.Atoi(getStructTagDefaultValue(options, "RefreshMaxNum"))
		ins.refreshMaxNum = int64(refreshMaxNum)
	}

	if ins.refreshDuration == 0 {
		refreshDuration, _ := strconv.Atoi(getStructTagDefaultValue(options, "RefreshDuration"))
		ins.refreshDuration = time.Duration(refreshDuration) * time.Second
	}

	//更新任务到待执行任务队列
	go func() {
		for {
			var startCount int64
			startCount = 0
			taskStatus := true

			for taskStatus {
				taskList, err := ins.redis.ZRangeWithScores(ctx, WaitQueue, startCount, startCount+ins.refreshMaxNum-1).Result()
				if err != nil {
					fmt.Printf("延时任务获取异常，错误：%s", err.Error())
					time.Sleep(5 * time.Second)
					continue
				}
				startCount += ins.refreshMaxNum
				if len(taskList) == 0 {
					taskStatus = false
					break
				}

				//分发任务到执行队列
				for _, task := range taskList {
					taskTime := int64(task.Score)
					currentTime := time.Now().Unix()

					//如果当前任务时间超过时间范围，则停止继续分发后面的任务
					if taskTime > currentTime+ int64(ins.refreshDuration.Seconds()) {
						goto jumpTask
					}

					if taskStr, ok := task.Member.(string); ok {
						taskInfo := QueueMsgSchema{}
						_ = json.Unmarshal([]byte(taskStr), &taskInfo)

						expireTime := time.Duration(taskTime - currentTime) * time.Second
						if expireTime <= 0 {
							expireTime = 1 * time.Millisecond
						}

						taskKey := generateKey(taskInfo.ID, taskInfo.NameSpace)
						taskContent := generateContentKey(taskInfo.ID, taskInfo.NameSpace)

						if _, err = ins.redis.Get(ctx, taskContent).Result(); err != nil {
							if err == redis.Nil {
								_ = ins.redis.Set(ctx, taskKey, taskInfo.Value, expireTime)
								_ = ins.redis.Set(ctx, taskContent, taskInfo.Value, -1)
							}
						}

						//从等待任务队列移除当前任务
						if _, err = ins.redis.ZRem(ctx, WaitQueue, task.Member).Result(); err != nil {
							fmt.Println(err)
						}
					}
				}
				time.Sleep(1 * time.Second)
			}

			jumpTask:
				time.Sleep(ins.refreshTime)
		}
	}()

	return &ins
}

// NewQueue 初始化新的任务队列
func (r *DelayTask) NewQueue(namespace string) *Queue {
	return NewQueue(QueueOptions{
		NameSpace: namespace,
		Cache:     r.cache,
	})
}
