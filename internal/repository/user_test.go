package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Andras5014/webook/internal/domain"
	"github.com/Andras5014/webook/internal/repository/cache"
	cachemocks "github.com/Andras5014/webook/internal/repository/cache/mocks"
	"github.com/Andras5014/webook/internal/repository/dao"
	daomocks "github.com/Andras5014/webook/internal/repository/dao/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCacheUserRepository_FindById(t *testing.T) {
	now := time.Now()
	now = time.UnixMilli(now.UnixMilli())
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)

		ctx      context.Context
		id       int64
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "缓存未命中,查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				daoMock := daomocks.NewMockUserDAO(ctrl)
				cacheMock := cachemocks.NewMockUserCache(ctrl)
				cacheMock.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{}, cache.ErrKeyNotExist)
				daoMock.EXPECT().FindById(gomock.Any(), int64(1)).Return(&dao.User{
					Id: 1,
					Email: sql.NullString{
						String: "123@qq.com",
						Valid:  true,
					},
					Password: "123456",
					Phone: sql.NullString{
						String: "123456789",
						Valid:  true,
					},
					CreatedAt: sql.NullInt64{
						Int64: now.UnixMilli(),
						Valid: true,
					},
					UpdatedAt: sql.NullInt64{
						Int64: now.UnixMilli(),
						Valid: true,
					}}, nil)

				cacheMock.EXPECT().Set(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, u domain.User) error {
					// 确保忽略时间字段
					if u.Id != 1 || u.Email != "123@qq.com" || u.Password != "123456" || u.Phone != "123456789" {
						return fmt.Errorf("unexpected user data: %v", u)
					}
					return nil
				}).Return(nil)
				return daoMock, cacheMock
			},
			ctx: context.Background(),
			id:  1,
			wantUser: domain.User{
				Id:        1,
				Email:     "123@qq.com",
				Password:  "123456",
				Phone:     "123456789",
				CreatedAt: now,
			},
			wantErr: nil,
		},
		{
			name: "缓存命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				daoMock := daomocks.NewMockUserDAO(ctrl)
				cacheMock := cachemocks.NewMockUserCache(ctrl)
				cacheMock.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{Id: 1,
					Email:     "123@qq.com",
					Password:  "123456",
					Phone:     "123456789",
					CreatedAt: now}, nil)

				return daoMock, cacheMock
			},
			ctx: context.Background(),
			id:  1,
			wantUser: domain.User{
				Id:        1,
				Email:     "123@qq.com",
				Password:  "123456",
				Phone:     "123456789",
				CreatedAt: now,
			},
			wantErr: nil,
		},
		{
			name: "缓存未命中,db查询失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				daoMock := daomocks.NewMockUserDAO(ctrl)
				cacheMock := cachemocks.NewMockUserCache(ctrl)
				cacheMock.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{}, cache.ErrKeyNotExist)
				daoMock.EXPECT().FindById(gomock.Any(), int64(1)).Return(&dao.User{}, errors.New("db error"))

				return daoMock, cacheMock
			},
			ctx:      context.Background(),
			id:       1,
			wantUser: domain.User{},
			wantErr:  errors.New("db error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 测试用例
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			daoMock, cacheMock := tc.mock(ctrl)
			repo := NewUserRepository(daoMock, cacheMock)
			user, err := repo.FindById(tc.ctx, tc.id)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser.Id, user.Id)
			assert.Equal(t, tc.wantUser.Email, user.Email)
			assert.Equal(t, tc.wantUser.Password, user.Password)
			assert.Equal(t, tc.wantUser.Phone, user.Phone)
			// 因为时间戳比较会有问题，所以我们单独比较
			assert.True(t, tc.wantUser.CreatedAt.Equal(user.CreatedAt))
		})
	}
}
