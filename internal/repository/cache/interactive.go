package cache

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/Andras5014/webook/internal/domain"
	"github.com/redis/go-redis/v9"
)

const fieldReadCnt = "read_cnt"
const fieldLikeCnt = "like_cnt"
const fieldCollectCnt = "collect_cnt"

var (
	//go:embed lua/incr_cnt.lua
	luaIncrCnt string
)

type InteractiveCache interface {
	IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error
	IncrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
	DecrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
	IncrCollectCntIfPresent(ctx context.Context, biz string, bizId int64) error
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	Set(ctx context.Context, biz string, bizId int64, res domain.Interactive) error
}

type InteractiveRedisCache struct {
	client redis.Cmdable
}

func NewInteractiveCache(client redis.Cmdable) InteractiveCache {
	return &InteractiveRedisCache{client: client}
}
func (i *InteractiveRedisCache) IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	key := i.key(biz, bizId)
	return i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldReadCnt).Err()
}

func (i *InteractiveRedisCache) IncrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	key := i.key(biz, bizId)
	return i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldLikeCnt, 1).Err()
}

func (i *InteractiveRedisCache) DecrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	key := i.key(biz, bizId)
	return i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldLikeCnt, -1).Err()
}

func (i *InteractiveRedisCache) IncrCollectCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	//TODO implement me
	panic("implement me")
}

func (i *InteractiveRedisCache) Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
	//TODO implement me
	panic("implement me")
}

func (i *InteractiveRedisCache) Set(ctx context.Context, biz string, bizId int64, res domain.Interactive) error {
	//TODO implement me
	panic("implement me")
}

func (i *InteractiveRedisCache) key(biz string, id int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, id)
}
