package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Andras5014/webook/internal/domain"
	"github.com/ecodeclub/ekit/queue"
	"github.com/ecodeclub/ekit/slice"
	"math"
	"time"
)

type RankingService interface {
	TopN(ctx context.Context, n int) error
}

type BatchRankingService struct {
	artSvc    ArticleService
	intrSvc   InteractiveService
	batchSize int
	scoreFunc func(likeCnt int64, updateTime time.Time) float64
}

func NewRankingService(artSvc ArticleService, intrSvc InteractiveService) RankingService {
	return &BatchRankingService{
		artSvc:    artSvc,
		intrSvc:   intrSvc,
		batchSize: 100,
		scoreFunc: func(likeCnt int64, updateTime time.Time) float64 {
			ms := time.Since(updateTime).Seconds()
			return float64(likeCnt+1) / math.Pow(ms+2, 1.5)
		},
	}
}
func (b *BatchRankingService) TopN(ctx context.Context, n int) error {
	arts, err := b.topN(ctx, n)
	if err != nil {
		return err
	}
	// todo 处理榜单
	fmt.Println(arts)
	return nil
}

func (b *BatchRankingService) topN(ctx context.Context, n int) ([]domain.Article, error) {
	// 计算一周内的文章
	startTime := time.Now().Add(time.Hour * 24 * 7)
	offset := 0
	type Score struct {
		art   domain.Article
		score float64
	}

	topN := queue.NewPriorityQueue[Score](n, func(src Score, dst Score) int {
		if src.score > dst.score {
			return 1
		} else if src.score < dst.score {
			return -1
		} else {
			return 0
		}
	})
	for {
		arts, err := b.artSvc.ListPub(ctx, startTime, offset, b.batchSize)
		if err != nil {
			return nil, err
		}
		ids := slice.Map[domain.Article, int64](arts, func(idx int, src domain.Article) int64 {
			return src.Id
		})
		intrs, err := b.intrSvc.GetByIds(ctx, "article", ids)
		if err != nil {
			return nil, err
		}

		// 计算score
		for _, art := range arts {
			intr, ok := intrs[art.Id]
			if !ok {
				continue
			}
			score := b.scoreFunc(intr.LikeCnt+2, art.UpdatedAt)
			rankScore := Score{art: art, score: score}

			err := topN.Enqueue(rankScore)
			if errors.Is(err, queue.ErrOutOfCapacity) {
				// 队满取出最小元素
				val, _ := topN.Dequeue()
				if val.score < score {
					_ = topN.Enqueue(rankScore)
				} else {
					_ = topN.Enqueue(val)

				}
			}
		}
		// 判断是否还有剩余
		if len(arts) < b.batchSize {
			break
		}
		offset += b.batchSize
	}
	res := make([]domain.Article, topN.Len())
	for i := topN.Len() - 1; i >= 0; i-- {
		val, err := topN.Dequeue()
		if err != nil {
			// 空队列
			break
		}
		res[i] = val.art
	}
	return res, nil

}
