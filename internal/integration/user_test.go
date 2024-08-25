package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Andras5014/webook/internal/integration/startup"
	"github.com/Andras5014/webook/internal/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserHandler_e2e_SendLoginSMSCode(t *testing.T) {
	server := startup.InitWebServer()
	rdb := startup.InitRedis(startup.InitConfig())
	testCases := []struct {
		name     string
		before   func(t *testing.T)
		after    func(t *testing.T)
		reqBody  string
		wantCode int
		wantBody web.Result
	}{
		{
			name: "发送成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				val, err := rdb.Get(ctx, "phone_code:login:12345678902").Result()
				cancel()
				assert.NoError(t, err)
				assert.True(t, len(val) == 6)
			},
			reqBody:  `{"phone":"12345678902"}`,
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 200,
				Data: nil,
				Msg:  "发送成功",
			},
		},
		{
			name: "发送频繁",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

				_, err := rdb.Set(ctx, "phone_code:login:12345678902", "123456", 10*time.Minute).Result()
				cancel()
				assert.NoError(t, err)

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				val, err := rdb.Get(ctx, "phone_code:login:12345678902").Result()
				cancel()
				assert.NoError(t, err)
				assert.Equal(t, "123456", val)
			},
			reqBody:  `{"phone":"12345678902"}`,
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 200,
				Data: nil,
				Msg:  "发送太频繁",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			req, err := http.NewRequest(http.MethodPost, "/users/login_sms/code/send", bytes.NewBuffer([]byte(tc.reqBody)))

			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()

			server.ServeHTTP(resp, req)

			var webRes web.Result
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			require.Equal(t, tc.wantCode, resp.Code)
			require.Equal(t, tc.wantBody, webRes)
			tc.after(t)

		})
	}
}
