package user

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/Andras5014/gohub/internal/domain"
	"github.com/Andras5014/gohub/internal/service"
	svcmocks "github.com/Andras5014/gohub/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEncrypt(t *testing.T) {
	password := "andras"
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	err = bcrypt.CompareHashAndPassword(encrypted, []byte(password))
	assert.NoError(t, err)
}

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) service.UserService

		reqBody string

		wantCode int
		wantResp string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "480735106@qq.com",
					Password: "andras",
				}).Return(nil)
				return userSvc
			},
			reqBody:  `{"email":"480735106@qq.com","password":"andras","confirmPassword":"andras"}`,
			wantCode: http.StatusOK,
			wantResp: "注册成功",
		},
		{
			name: "参数不对 bind失败",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)

				return userSvc
			},
			reqBody:  `{"email":"480735106.com","password":"andras","confirmPassword":"andras"}`,
			wantCode: http.StatusOK,
			wantResp: `{"code":200,"msg":"参数错误"}`,
		},
		{
			name: "二次密码不同",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)

				return userSvc
			},
			reqBody:  `{"email":"480735106@qq.com","password":"andras","confirmPassword":"andras1"}`,
			wantCode: http.StatusOK,
			wantResp: `{"code":200,"msg":"两次密码不一致"}`,
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "480735106@qq.com",
					Password: "andras",
				}).Return(service.ErrUserDuplicateEmail)
				return userSvc
			},
			reqBody:  `{"email":"480735106@qq.com","password":"andras","confirmPassword":"andras"}`,
			wantCode: http.StatusOK,
			wantResp: `{"code":200,"msg":"邮箱冲突"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := gin.Default()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			h := NewUserHandler(tc.mock(ctrl), nil, nil, nil)
			h.RegisterRoutes(server)
			req, err := http.NewRequest(http.MethodPost,
				"/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()

			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)

			fmt.Println(resp.Body.String())
			fmt.Println(tc.wantResp)
			assert.Equal(t, tc.wantResp, resp.Body.String())

		})
	}
}

func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userSvc := svcmocks.NewMockUserService(ctrl)

	userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))

	err := userSvc.SignUp(context.Background(), domain.User{
		Email: "480735106@qq.com",
	})
	t.Log(err)
}
