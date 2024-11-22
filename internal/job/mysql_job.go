package job

import (
	"context"
	"errors"
	"github.com/Andras5014/webook/internal/domain"
	"github.com/Andras5014/webook/internal/service"
	"github.com/Andras5014/webook/pkg/logx"
	"golang.org/x/sync/semaphore"
)

type Executor interface {
	Name() string
	Exec(ctx context.Context, j domain.Job) error
}

type LocalFuncExecutor struct {
	funcs map[string]func(ctx context.Context, j domain.Job) error
}

func NewLocalFuncExecutor() *LocalFuncExecutor {
	return &LocalFuncExecutor{
		funcs: make(map[string]func(ctx context.Context, j domain.Job) error),
	}
}

func (l *LocalFuncExecutor) Name() string {
	return "LocalFuncExecutor"
}

func (l *LocalFuncExecutor) Exec(ctx context.Context, j domain.Job) error {
	fn, ok := l.funcs[j.Executor]
	if !ok {
		return errors.New("没有对应的执行器" + j.Name)
	}
	return fn(ctx, j)
}

func (l *LocalFuncExecutor) RegisterFunc(name string, fn func(ctx context.Context, j domain.Job) error) {
	l.funcs[name] = fn
}

type Scheduler struct {
	executors map[string]Executor
	svc       service.JobService
	limiter   *semaphore.Weighted
	l         logx.Logger
}

func NewScheduler(svc service.JobService, l logx.Logger) *Scheduler {
	return &Scheduler{
		executors: make(map[string]Executor),
		svc:       svc,
		limiter:   semaphore.NewWeighted(100),
		l:         l,
	}
}

func (s *Scheduler) RegisterExecutor(e Executor) {
	s.executors[e.Name()] = e
}

func (s *Scheduler) Schedule(ctx context.Context) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := s.limiter.Acquire(ctx, 1); err != nil {
			s.l.Error("获取信号量失败", logx.Error(err))
			continue
		}
		j, err := s.svc.Preempt(ctx)
		if err != nil {
			s.l.Error("抢占任务失败", logx.Error(err))
			continue
		}

		exec, ok := s.executors[j.Executor]
		if !ok {
			s.l.Error("没有对应的执行器", logx.String("name", j.Name))
			continue
		}
		go func() {
			defer func() {
				s.limiter.Release(1)
				if er := j.CancelFunc(); er != nil {
					s.l.Error("取消任务失败", logx.Error(er))
				}
			}()

			er := exec.Exec(ctx, j)
			if er != nil {
				s.l.Error("执行任务失败", logx.Error(er))
				return
			}
			er = s.svc.ResetNextTime(ctx, j)
			if er != nil {
				s.l.Error("重置下次执行时间失败", logx.Error(er))
				return
			}
		}()
	}
}
