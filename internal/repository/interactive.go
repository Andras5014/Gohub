package repository

import (
	"context"
	"errors"
	"github.com/Andras5014/webook/internal/domain"
	"github.com/Andras5014/webook/internal/repository/cache"
	"github.com/Andras5014/webook/internal/repository/dao"
	"github.com/Andras5014/webook/pkg/logx"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, id int64, uid int64) error
	IncrLike(ctx context.Context, biz string, id int64, uid int64) error
	DecrLike(ctx context.Context, biz string, id int64, uid int64) error
	AddCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	Get(ctx context.Context, biz string, id int64) (domain.Interactive, error)
	Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error)
}

type CacheInteractiveRepository struct {
	dao   dao.InteractiveDAO
	cache cache.InteractiveCache
	l     logx.Logger
}

func (c *CacheInteractiveRepository) Get(ctx context.Context, biz string, id int64) (domain.Interactive, error) {
	intrCache, err := c.cache.Get(ctx, biz, id)
	if errors.Is(err, nil) {
		return intrCache, nil
	}

	intrDao, err := c.dao.Get(ctx, biz, id)
	if err != nil {
		return domain.Interactive{}, err
	}

	// 回写缓存
	if errors.Is(err, nil) {
		res := c.toDomain(intrDao)
		err = c.cache.Set(ctx, biz, id, res)
		if err != nil {
			c.l.Error("redis缓存失败", logx.Any("err", err),
				logx.Any("biz", biz), logx.Any("bizId", id))
		}
		return res, nil
	}
	return intrCache, err
}

func (c *CacheInteractiveRepository) Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := c.dao.GetLikeInfo(ctx, biz, id, uid)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, dao.ErrRecordNotFound):
		return false, nil
	default:
		return false, err
	}
}

func (c *CacheInteractiveRepository) Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := c.dao.GetCollectInfo(ctx, biz, id, uid)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, dao.ErrRecordNotFound):
		return false, nil
	default:
		return false, err
	}
}

func (c *CacheInteractiveRepository) AddCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error {
	err := c.dao.InsertCollectionBiz(ctx, dao.UserCollectionBiz{
		Cid:   cid,
		Uid:   uid,
		Biz:   biz,
		BizId: id,
	})
	if err != nil {
		return err
	}
	return c.cache.IncrCollectCntIfPresent(ctx, biz, id)
}

func NewInteractiveRepository(dao dao.InteractiveDAO, cache cache.InteractiveCache, l logx.Logger) InteractiveRepository {
	return &CacheInteractiveRepository{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}

func (c *CacheInteractiveRepository) IncrReadCnt(ctx context.Context, biz string, id int64, uid int64) error {
	err := c.dao.IncrReadCnt(ctx, biz, id, uid)
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

func (c *CacheInteractiveRepository) toDomain(intrDao dao.Interactive) domain.Interactive {
	return domain.Interactive{
		ReadCnt:    intrDao.ReadCnt,
		LikeCnt:    intrDao.LikeCnt,
		CollectCnt: intrDao.CollectCnt,
	}
}
