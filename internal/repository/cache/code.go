package cache

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	ErrCodeSendTooMany        = fmt.Errorf("code send too many")
	ErrCodeVerifyTooManyTimes = fmt.Errorf("code verify too many times")
)

//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}
type RedisCodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		client: client,
	}
}

func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}

	switch res {
	case 0:
		return nil
	case -1:
		return ErrCodeSendTooMany
	default:
		return fmt.Errorf("set code failed")

	}
}

//	func (c *RedisCodeCache) Get(ctx context.Context, biz, phone string) (string, error) {
//		return c.client.Get(ctx, biz+phone).Result()
//	}
func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		return false, ErrCodeVerifyTooManyTimes
	default:
		return false, fmt.Errorf("verify code failed")
	}
}
func (c *RedisCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
