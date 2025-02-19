package cache

import (
	"context"
	"github.com/Andras5014/gohub/internal/repository/cache/redismocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRedisCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(*gomock.Controller) redis.Cmdable
		ctx     context.Context
		biz     string
		phone   string
		code    string
		wantErr error
	}{
		{
			name: "code存储成功",
			mock: func(controller *gomock.Controller) redis.Cmdable {
				mockCmdable := redismocks.NewMockCmdable(controller)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(0))
				mockCmdable.EXPECT().Eval(gomock.Any(),
					luaSetCode,
					[]string{"phone_code:register:12345678901"},
					[]any{"123456"},
				).Return(res)
				return mockCmdable
			},
			biz:     "register",
			phone:   "12345678901",
			code:    "123456",
			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			codeCache := NewCodeCache(tc.mock(ctrl))
			err := codeCache.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantErr, err)
		})
	}

}

/*
	a aaaa/bbbbb

b aaaacccbbb
c bbbbb
*/
