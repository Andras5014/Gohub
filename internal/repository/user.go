package repository

import (
	"context"
	"database/sql"
	"github.com/Andras5014/webook/internal/domain"
	"github.com/Andras5014/webook/internal/repository/cache"
	"github.com/Andras5014/webook/internal/repository/dao"
	"log"
	"time"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
}

type CacheUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDAO, cache cache.UserCache) UserRepository {
	return &CacheUserRepository{
		dao:   dao,
		cache: cache,
	}
}
func (r *CacheUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.domainToEntity(&u))
}

func (r *CacheUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(user), nil
}

func (r *CacheUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	user, err := r.cache.Get(ctx, id)
	if err == nil {
		return user, nil
	}
	//if errors.Is(err, cache.ErrKeyNotExist) {
	//	user, err := r.dao.FindById(ctx, id)
	//	if err != nil {
	//		return domain.User{}, err
	//	}
	//	err = r.cache.Set(ctx, domain.User{
	//		Id:       user.Id,
	//		Email:    user.Email,
	//		Password: user.Password,
	//	})
	//	if err != nil {
	//		return domain.User{}, err
	//	}
	//	return domain.User{
	//		Id:       user.Id,
	//		Email:    user.Email,
	//		Password: user.Password,
	//	}, nil
	//}
	// 其他错误
	// 加载 redis崩溃，保护数据库 限流
	// 不加载 redis 崩溃
	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u := r.entityToDomain(ue)
	err = r.cache.Set(ctx, u)
	if err != nil {
		log.Println("redis set err", err)
	}
	return u, nil
}

func (r *CacheUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	user, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(user), nil
}

func (r *CacheUserRepository) entityToDomain(user *dao.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Email:    user.Email.String,
		Phone:    user.Phone.String,
		Password: user.Password,
		CreateAt: time.UnixMilli(user.CreatedAt),
	}
}
func (r *CacheUserRepository) domainToEntity(user *domain.User) dao.User {
	return dao.User{
		Id:        user.Id,
		Email:     sql.NullString{String: user.Email, Valid: user.Email != ""},
		Phone:     sql.NullString{String: user.Phone, Valid: user.Phone != ""},
		Password:  user.Password,
		CreatedAt: user.CreateAt.UnixMilli(),
	}
}
