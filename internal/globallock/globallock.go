package globallock

import (
	"flying-star/internal/db"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

var RM *redsync.Redsync

func init()  {
	origin := db.RedisClient.GetOriginPoint()
	pool := goredis.NewPool(origin)
	RM = redsync.New(pool)
	fmt.Println(RM)
}