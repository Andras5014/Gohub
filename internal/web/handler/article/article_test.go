package article

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Andras5014/gohub/internal/domain"
	"github.com/Andras5014/gohub/internal/service"
	svcmocks "github.com/Andras5014/gohub/internal/service/mocks"
	"github.com/Andras5014/gohub/internal/web/result"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestArticleHandler_Publish(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) service.ArticleService

		reqBody string

		wantCode int
		wantRes  result.Result
	}{
		{
			name: "新建并发表",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				return svc
			},
			reqBody: `
{
	"title": "我的标题",
	"content": "我的内容"
}
`,
			wantCode: http.StatusOK,
			wantRes: result.Result{
				Code: 0,
				Data: float64(1),
				Msg:  "ok",
			},
		},
		{
			name: "publish失败",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(0), errors.New("publish失败"))
				return svc
			},
			reqBody: `
{
	"title": "我的标题",
	"content": "我的内容"
}
`,
			wantCode: http.StatusOK,
			wantRes: result.Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := gin.Default()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			server.Use(func(ctx *gin.Context) {
				ctx.Set("userId", int64(123))
			})
			h := NewArticleHandler(tc.mock(ctrl), nil)
			h.RegisterRoutes(server)
			req, err := http.NewRequest(http.MethodPost,
				"/articles/publish", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()

			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			var webRes result.Result
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantRes, webRes)
		})
	}
}
