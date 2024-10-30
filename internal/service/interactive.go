package service

import (
	"context"
	"github.com/Andras5014/webook/internal/domain"
	"github.com/Andras5014/webook/internal/repository"
	"github.com/ecodeclub/ekit/slice"
	"golang.org/x/sync/errgroup"
)
//go:generate mockgen -destination=mocks/interactive.mock.go -package=svcmocks -source=interactive.go
type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, id int64) error
	Like(ctx context.Context, biz string, id int64, uid int64) error
	CancelLike(ctx context.Context, biz string, id int64, uid int64) error
	Collect(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	Get(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error)
	GetByIds(ctx context.Context, biz string, ids []int64) (map[int64]domain.Interactive, error)
}

type interactiveService struct {
	repo repository.InteractiveRepository
}

func (i *interactiveService) GetByIds(ctx context.Context, biz string, ids []int64) (map[int64]domain.Interactive, error) {
	intrs, err := i.repo.GetByIds(ctx, biz, ids)
	if err != nil {
		return nil, err
	}
	return slice.ToMap(intrs, func(intr domain.Interactive) int64 {
		return intr.BizId
	}), nil
}

func (i *interactiveService) Get(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error) {
	var (
		intr domain.Interactive
		err  error
		eg   errgroup.Group
	)

	intr, err = i.repo.Get(ctx, biz, id)
	if err != nil {
		return domain.Interactive{}, nil
	}

	eg.Go(func() error {
		intr.Liked, err = i.repo.Liked(ctx, biz, id, uid)
		return err
	})

	eg.Go(func() error {
		intr.Collected, err = i.repo.Collected(ctx, biz, id, uid)
		return err
	})

	return intr, eg.Wait()

}

func (i *interactiveService) Collect(ctx context.Context, biz string, id int64, cid int64, uid int64) error {
	return i.repo.AddCollectionItem(ctx, biz, id, cid, uid)
}

func NewInteractiveService(repo repository.InteractiveRepository) InteractiveService {
	return &interactiveService{repo: repo}
}

func (i *interactiveService) IncrReadCnt(ctx context.Context, biz string, id int64) error {
	return i.repo.IncrReadCnt(ctx, biz, id)
}
func (i *interactiveService) Like(ctx context.Context, biz string, id int64, uid int64) error {
	return i.repo.IncrLike(ctx, biz, id, uid)
}

func (i *interactiveService) CancelLike(ctx context.Context, biz string, id int64, uid int64) error {
	return i.repo.DecrLike(ctx, biz, id, uid)
}
