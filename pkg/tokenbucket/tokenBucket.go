package tokenbucket

import (
	"time"
)

type Options struct {
	//生成令牌数
	TokenNum uint16
	//超时时间
	TokenTimeout time.Duration
}

type TokenBucket struct {
	bucket  chan struct{}
	maxNum  int
	expires time.Duration
	timeout time.Duration
}

func NewTokenBucket(options Options) *TokenBucket {
	tokenBucket := TokenBucket{
		maxNum:  int(options.TokenNum),
		timeout: options.TokenTimeout,
	}
	bucket := make(chan struct{}, tokenBucket.maxNum)
	for i := 0; i < tokenBucket.maxNum; i++ {
		bucket <- struct{}{}
	}
	tokenBucket.bucket = bucket
	return &tokenBucket
}

// Get 申领令牌
func (t *TokenBucket) Get() bool {
	<-t.bucket
	return true
}

// Set 归还令牌
func (t *TokenBucket) Set() {
	if len(t.bucket) < t.maxNum {
		t.bucket <- struct{}{}
	}
}

// Len 剩余令牌数
func (t *TokenBucket) Len() int {
	return len(t.bucket)
}

// Count 令牌总数
func (t *TokenBucket) Count() int {
	return t.maxNum
}
