package cache

import (
	"context"
	"errors"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"sync"
	"time"
)

type LocalCodeCache struct {
	cache      *lru.Cache
	lock       sync.Mutex
	expiration time.Duration
}

func NewLocalCodeCache(cache *lru.Cache, expiration time.Duration) *LocalCodeCache {
	return &LocalCodeCache{
		cache:      cache,
		expiration: expiration,
	}
}

func (l *LocalCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	key := l.key(biz, phone)
	now := time.Now()
	val, ok := l.cache.Get(key)
	if !ok {
		l.cache.Add(key, codeItem{
			code:   code,
			cnt:    3,
			expire: now.Add(l.expiration),
		})
		return nil
	}
	itm, ok := val.(codeItem)
	if !ok {
		return errors.New("系统错误")
	}
	if itm.expire.Sub(now) > time.Minute*9 {
		//不到一分钟
		return ErrCodeSendTooMany
	}

	//重发
	l.cache.Add(key, codeItem{
		code:   code,
		cnt:    3,
		expire: now.Add(l.expiration),
	})
	return nil
}

func (l *LocalCodeCache) Verify(ctx context.Context, biz string, phone string, code string) (bool, error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	key := l.key(biz, phone)
	val, ok := l.cache.Get(key)
	if !ok {
		return false, ErrKeyNotExist
	}
	itm, ok := val.(codeItem)
	if !ok {
		return false, errors.New("系统错误")
	}
	if itm.cnt <= 0 {
		return false, ErrCodeVerifyTooManyTimes
	}
	itm.cnt--

	return itm.code == code, nil
}
func (l *LocalCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

type codeItem struct {
	code   string
	cnt    int
	expire time.Time
}
