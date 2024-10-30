package service

import (
	"context"
	"errors"
	"github.com/Andras5014/webook/internal/domain"
	articleEvent "github.com/Andras5014/webook/internal/events/article"
	"github.com/Andras5014/webook/internal/repository/article"
	"github.com/Andras5014/webook/pkg/logx"
	"time"
)
//go:generate mockgen -destination=mocks/article.mock.go -package=svcmocks -source=./article.go
type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) (int64, error)
	PublishV1(ctx context.Context, article domain.Article) (int64, error)
	Withdraw(ctx context.Context, article domain.Article) (int64, error)
	List(ctx context.Context, id int64, offset int, limit int) ([]domain.Article, error)
	ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPubById(ctx context.Context, id, uid int64) (domain.Article, error)
}

type articleService struct {
	// 方案一
	repo article.Repository

	// 方案二 依靠二个不同repository 来解决跨表跨库问题
	readerRepo article.ReaderRepository
	authorRepo article.AuthorRepository

	logger logx.Logger

	producer articleEvent.Producer
}

func (a *articleService) ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]domain.Article, error) {
	return a.repo.ListPub(ctx, start, offset, limit)
}

func (a *articleService) GetPubById(ctx context.Context, id, uid int64) (domain.Article, error) {
	art, err := a.repo.GetPubById(ctx, id)
	if errors.Is(err, nil) {
		go func() {
			er := a.producer.ProduceReadEvent(ctx, articleEvent.ReadEvent{
				ArticleId: id,
				UserId:    uid,
			})
			if er != nil {
				a.logger.Error("文章阅读事件写入失败", logx.Any("article_id", id), logx.Any("uid", uid), logx.Error(er))
			}
		}()

	}
	return art, err
}

func NewArticleService(repo article.Repository, producer articleEvent.Producer, l logx.Logger) ArticleService {
	return &articleService{
		repo:     repo,
		producer: producer,
		logger:   l,
	}
}

func NewArticleServiceV1(readerRepo article.ReaderRepository, authorRepo article.AuthorRepository, l logx.Logger) ArticleService {
	return &articleService{
		readerRepo: readerRepo,
		authorRepo: authorRepo,
		logger:     l,
	}
}

func (a *articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	article.Status = domain.ArticleStatusUnPublished
	if article.Id > 0 {
		return article.Id, a.repo.Update(ctx, article)
	}
	return a.repo.Create(ctx, article)

}

func (a *articleService) Publish(ctx context.Context, article domain.Article) (int64, error) {
	//制作库
	article.Status = domain.ArticleStatusPublished
	return a.repo.Sync(ctx, article)
}
func (a *articleService) PublishV1(ctx context.Context, article domain.Article) (int64, error) {
	var (
		id  = article.Id
		err error
	)

	if article.Id > 0 {
		err = a.authorRepo.Update(ctx, article)
	} else {
		id, err = a.authorRepo.Create(ctx, article)
	}
	if err != nil {
		return 0, err
	}
	article.Id = id

	// 保存到线上库并重试处理
	err = a.retrySaveToReaderRepo(ctx, article, 3)
	if err != nil {
		a.logger.Error("保存到线上库失败", logx.Any("article_id", article.Id), logx.Error(err))
		return id, err
	}
	return id, nil

}
func (a *articleService) List(ctx context.Context, id int64, offset int, limit int) ([]domain.Article, error) {
	return a.repo.List(ctx, id, offset, limit)
}
func (a *articleService) GetById(ctx context.Context, id int64) (domain.Article, error) {
	return a.repo.GetById(ctx, id)
}

func (a *articleService) Withdraw(ctx context.Context, article domain.Article) (int64, error) {
	article.Status = domain.ArticleStatusPrivate
	return a.repo.SyncStatus(ctx, article)
}

// retrySaveToReaderRepo 重试保存到 readerRepo，最多重试指定次数
func (a *articleService) retrySaveToReaderRepo(ctx context.Context, art domain.Article, maxRetries int) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = a.readerRepo.Save(ctx, art)
		if errors.Is(err, nil) {
			return nil
		}
		a.logger.Error("保存到 readerRepo 失败", logx.Any("article_id", art.Id), logx.Any("retry", i), logx.Error(err))
	}
	return err
}
