package service

import (
	"context"
	"errors"
	"github.com/Andras5014/gohub/internal/domain"
	"github.com/Andras5014/gohub/internal/repository"
	"github.com/Andras5014/gohub/pkg/logx"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicateEmail    = repository.ErrUserDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("无效的邮箱或者密码")
)

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email string, password string) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	UpdateNonSensitiveInfo(ctx context.Context, u domain.User) error
	FindOrCreateByWechat(ctx context.Context, info domain.WeChatInfo) (domain.User, error)
}
type userService struct {
	repo   repository.UserRepository
	logger logx.Logger
}

func NewUserService(repo repository.UserRepository, logger logx.Logger) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) SignUp(ctx context.Context, u domain.User) error {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(encryptedPassword)
	return svc.repo.Create(ctx, u)
}

func (svc *userService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	user, err := svc.repo.FindByEmail(ctx, email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return user, nil
}

func (svc *userService) UpdateNonSensitiveInfo(ctx context.Context, u domain.User) error {
	return svc.repo.Update(ctx, u)
}
func (svc *userService) Profile(ctx context.Context, id int64) (domain.User, error) {
	return svc.repo.FindById(ctx, id)
}
func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	user, err := svc.repo.FindByPhone(ctx, phone)
	if !errors.Is(err, repository.ErrUserNotFound) {
		return user, err
	}

	u := domain.User{
		Phone: phone,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil {
		return user, err
	}
	// phone 脱敏 1762****454
	svc.logger.Info("创建用户", logx.Any("phone", u.Phone))

	//主从延迟会出现问题
	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *userService) FindOrCreateByWechat(ctx context.Context, info domain.WeChatInfo) (domain.User, error) {
	user, err := svc.repo.FindByWechat(ctx, info.OpenId)
	if !errors.Is(err, repository.ErrUserNotFound) {
		return user, err
	}

	u := domain.User{
		WeChatInfo: info,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil {
		return user, err
	}

	//主从延迟会出现问题
	return svc.repo.FindByWechat(ctx, info.OpenId)
}

//func (svc *userService) Edit(ctx context.Context, u domain.User) error {
//	return svc.repo.Update(ctx, u)
//}
