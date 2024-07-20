package repository

import (
	"context"
	"github.com/Andras5014/webook/internal/domain"
	"github.com/Andras5014/webook/internal/repository/dao"
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}
func (ur *UserRepository) Create(ctx context.Context, u domain.User) error {
	return ur.dao.Insert(ctx, &dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}
