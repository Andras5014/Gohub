package job

import (
	"context"
	"github.com/Andras5014/gohub/internal/service"
	"github.com/Andras5014/gohub/pkg/logx"
	"github.com/Andras5014/gohub/pkg/redsyncx"
	"github.com/go-redsync/redsync/v4"
	"sync"
	"time"
)

type RankingJob struct {
	svc        service.RankingService
	timeout    time.Duration
	mutex      *redsyncx.AutoExtendMutex
	localMutex sync.Mutex
	l          logx.Logger
}

func NewRankingJob(svc service.RankingService, timeout time.Duration, mutex *redsync.Mutex, l logx.Logger) *RankingJob {
	return &RankingJob{
		svc:     svc,
		timeout: timeout,
		mutex:   redsyncx.NewAutoExtendMutex(mutex),
		l:       l,
	}
}

func (r *RankingJob) Name() string {
	return "ranking_job"
}
func (r *RankingJob) Run(ctx context.Context) error {
	r.localMutex.Lock()
	defer r.localMutex.Unlock()

	if r.mutex.Until().IsZero() {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := r.mutex.TryLockContext(ctx); err != nil {
			return err
		}

		// 拿到锁 自动续费
		go func() {
			err := r.mutex.AutoExtend(r.timeout/2, time.Second)
			if err != nil {
				r.l.Error("自动续期失败", logx.Error(err))
			}
		}()
	}

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	return r.svc.TopN(ctx, 100)
}

func (r *RankingJob) Close() error {
	_, err := r.mutex.Unlock()
	if err != nil {
		r.l.Error("释放锁失败", logx.Error(err))
	}
	return nil
}
