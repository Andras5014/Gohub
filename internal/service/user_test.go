package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Andras5014/webook/internal/domain"
	"github.com/Andras5014/webook/internal/repository"
	repomocks "github.com/Andras5014/webook/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func Test_userService_Login(t *testing.T) {
	testCases := []struct {
		name     string
		email    string
		password string
		mock     func(ctrl *gomock.Controller) repository.UserRepository

		ctx      context.Context
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "test1@qq.com").
					Return(domain.User{
						Email:    "test1@qq.com",
						Password: "$2a$10$U7hG2ztjaSxQcePCr1y0ROg6mXE6sbisj1DN/b6ylm2EWW.p06EB6",
					}, nil)
				return repo
			},
			email:    "test1@qq.com",
			password: "123456",

			wantUser: domain.User{
				Email:    "test1@qq.com",
				Password: "$2a$10$U7hG2ztjaSxQcePCr1y0ROg6mXE6sbisj1DN/b6ylm2EWW.p06EB6",
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "test1@qq.com").
					Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			email:    "test1@qq.com",
			password: "123456",

			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "test1@qq.com").
					Return(domain.User{}, errors.New("mock db error"))
				return repo
			},
			email:    "test1@qq.com",
			password: "123456",

			wantUser: domain.User{},
			wantErr:  errors.New("mock db error"),
		},
		{
			name: "密码错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "test1@qq.com").
					Return(domain.User{
						Email:    "test1@qq.com",
						Password: "xxx$2a$10$U7hG2ztjaSxQcePCr1y0ROg6mXE6sbisj1DN/b6ylm2EWW.p06EB6",
					}, nil)

				return repo
			},
			email:    "test1@qq.com",
			password: "123456",

			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewUserService(tc.mock(ctrl), nil)
			u, err := svc.Login(tc.ctx, tc.email, tc.password)
			fmt.Println(u, err)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)

		})
	}
}

func TestEncrypt(t *testing.T) {
	password, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	fmt.Println(string(password))
}
