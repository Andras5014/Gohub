package repository

import (
	"context"
	"github.com/Andras5014/webook/internal/repository/cache"
	"github.com/Andras5014/webook/internal/repository/dao"
	"github.com/Andras5014/webook/pkg/logx"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, id int64) error
	IncrLike(ctx context.Context, biz string, id int64, uid int64) error
	DecrLike(ctx context.Context, biz string, id int64, uid int64) error
}

type CacheInteractiveRepository struct {
	dao   dao.InteractiveDAO
	cache cache.InteractiveCache
	l     logx.Logger
}

func NewInteractiveRepository(dao dao.InteractiveDAO, cache cache.InteractiveCache, l logx.Logger) InteractiveRepository {
	return &CacheInteractiveRepository{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}

func (c *CacheInteractiveRepository) IncrReadCnt(ctx context.Context, biz string, id int64) error {
	err := c.dao.IncrReadCnt(ctx, biz, id)
	if err != nil {
		return err
	}
	return c.cache.IncrReadCntIfPresent(ctx, biz, id)
}

func (c *CacheInteractiveRepository) IncrLike(ctx context.Context, biz string, id int64, uid int64) error {
	err := c.dao.InsertLikeInfo(ctx, biz, id, uid)
	if err != nil {
		return err
	}
	return c.cache.IncrLikeCntIfPresent(ctx, biz, id)
}

func (c *CacheInteractiveRepository) DecrLike(ctx context.Context, biz string, id int64, uid int64) error {
	err := c.dao.DeleteLikeInfo(ctx, biz, id, uid)
	if err != nil {
		return err
	}
	return c.cache.DecrLikeCntIfPresent(ctx, biz, id)
}
