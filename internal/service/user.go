package service

import (
	"context"
	"github.com/Andras5014/webook/internal/domain"
	"github.com/Andras5014/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}
func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(encryptedPassword)
	return svc.repo.Create(ctx, u)
}
