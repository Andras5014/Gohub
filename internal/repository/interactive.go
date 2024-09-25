package repository

import (
	"context"
	"github.com/Andras5014/webook/internal/repository/cache"
	"github.com/Andras5014/webook/internal/repository/dao"
	"github.com/Andras5014/webook/pkg/logx"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, id int64) error
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
	return c.cache.IncrReadCnt(ctx, biz, id)
}
