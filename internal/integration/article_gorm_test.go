package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Andras5014/gohub/internal/domain"
	"github.com/Andras5014/gohub/internal/integration/startup"
	"github.com/Andras5014/gohub/internal/repository/dao/article"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ArticleTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

func (s *ArticleTestSuite) SetupSuite() {
	// 所有测试之前初始化
	s.server = gin.Default()
	s.server.Use(func(context *gin.Context) {
		var id int64
		id = 123
		context.Set("userId", id)
	})
	l := startup.InitLogger()
	config := startup.InitConfig()
	s.db = startup.InitDB(config, l)
	articleHdl := startup.InitArticleHandler()
	articleHdl.RegisterRoutes(s.server)
}

func (s *ArticleTestSuite) TearDownSuite() {
	s.db.Exec("truncate table articles")
}
func (s *ArticleTestSuite) TestEdit() {
	t := s.T()
	testCases := []struct {
		name     string
		art      Article
		before   func(t *testing.T)
		after    func(t *testing.T)
		wantCode int
		wantRes  Result[int64]
		wantErr  error
	}{
		{
			name: "新建帖子,保存成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				var art article.Article
				err := s.db.Where("id", 1).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.CreatedAt > 0)
				assert.True(t, art.UpdatedAt > 0)
				art.CreatedAt = 0
				art.UpdatedAt = 0

				assert.Equal(t, article.Article{
					Id:       1,
					Title:    "测试标题",
					Content:  "测试内容",
					AuthorId: 123,
					Status:   domain.ArticleStatusUnPublished.ToUint8(),
				}, art)
			},
			art: Article{
				Title:   "测试标题",
				Content: "测试内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{

				Msg:  "ok",
				Data: 1,
			},
			wantErr: nil,
		},
		{
			name: "修改已有帖子,保存成功",
			before: func(t *testing.T) {
				err := s.db.Create(&article.Article{
					Id:        2,
					Title:     "测试标题",
					Content:   "测试内容",
					AuthorId:  123,
					CreatedAt: int64(123),
					UpdatedAt: int64(123),
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var art article.Article
				err := s.db.Where("id", 2).First(&art).Error
				assert.NoError(t, err)
				//assert.True(t, article.CreatedAt.Int64 == int64(123))
				//assert.True(t, article.UpdatedAt.Int64 > 123)
				art.CreatedAt = 0
				art.UpdatedAt = 0

				assert.Equal(t, article.Article{
					Id:       2,
					Title:    "新的标题",
					Content:  "新的内容",
					AuthorId: 123,
					Status:   domain.ArticleStatusUnPublished.ToUint8(),
				}, art)

			},
			art: Article{
				Id:      2,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{

				Msg:  "ok",
				Data: 2,
			},
			wantErr: nil,
		},
		{
			name: "修改别人帖子",
			before: func(t *testing.T) {
				err := s.db.Create(&article.Article{
					Id:        3,
					Title:     "测试标题",
					Content:   "测试内容",
					AuthorId:  789,
					CreatedAt: int64(123),
					UpdatedAt: int64(123),
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var art article.Article
				err := s.db.Where("id", 3).First(&art).Error
				assert.NoError(t, err)
				//assert.True(t, article.CreatedAt.Int64 == int64(123))
				//assert.True(t, article.UpdatedAt.Int64 > 123)
				art.CreatedAt = 0
				art.UpdatedAt = 0

				assert.Equal(t, article.Article{
					Id:       3,
					Title:    "测试标题",
					Content:  "测试内容",
					AuthorId: 789,
				}, art)

			},
			art: Article{
				Id:      3,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Code: 5,
				Msg:  "系统错误",
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			reqBody, err := json.Marshal(tc.art)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/articles/edit", bytes.NewBuffer(reqBody))

			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()

			s.server.ServeHTTP(resp, req)

			var webRes Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			require.Equal(t, tc.wantCode, resp.Code)
			require.Equal(t, tc.wantRes, webRes)
			tc.after(t)

		})
	}
}
func TestArticle(t *testing.T) {
	suite.Run(t, new(ArticleTestSuite))
}

type Article struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func TestA(t *testing.T) {
	var ctx gin.Context
	ctx.Set("userId", 123)
	id := ctx.MustGet("userId")
	fmt.Println(id)
}

//type ResultTest[T any] struct {
//	result T
//}
//
//func testHandler(ctx *gin.Context) {
//	//拿唯一标识来
//	code := ctx.MustGet("bizcode")
//	switch {
//	case code == 1:
//		var result ResultTest[int]
//	}
//	ctx.ShouldBindJSON(&result)
//
//}
