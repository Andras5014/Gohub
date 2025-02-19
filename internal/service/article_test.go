package service

import (
	"context"
	"github.com/Andras5014/gohub/internal/domain"
	"github.com/Andras5014/gohub/internal/repository/article"
	artrepomocks "github.com/Andras5014/gohub/internal/repository/article/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

//func Test_articleService_Publish(t *testing.T) {
//	testCases := []struct {
//		name string
//		mock func(ctrl *gomock.Controller) article.Repository
//		article domain.Article
//		wantErr error
//		wantId int64
//		}{
//		{
//			name: "发布文章成功",
//			mock: func(ctrl *gomock.Controller) article.Repository {
//
//0
//			},
//			article: domain.Article{
//				Title: "我的标题",
//				Content: "我的内容",
//				Author: domain.Author{
//					Id:123,
//				},
//			},
//
//		},
//	}
//	}
//	for _, tc := range testCases {
//		ctrl :=gomock.NewController(t)
//		defer ctrl.Finish()
//		svc := NewArticleService(tc.mock(ctrl))
//		id,err := svc.Publish(context.Background(), tc.article)
//		assert.Equal(t, tc.wantErr, err)
//		assert.Equal(t, tc.wantId, id)
//	}
//}

func Test_articleService_PublishV1(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (article.ReaderRepository, article.AuthorRepository)
		article domain.Article
		wantErr error
		wantId  int64
	}{
		{
			name: "新建发布文章成功",
			mock: func(ctrl *gomock.Controller) (article.ReaderRepository, article.AuthorRepository) {
				artAuthor := artrepomocks.NewMockAuthorRepository(ctrl)
				artAuthor.EXPECT().Create(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				artReader := artrepomocks.NewMockReaderRepository(ctrl)
				artReader.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      1,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)
				return artReader, artAuthor
			},
			article: domain.Article{
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantErr: nil,
			wantId:  1,
		},
		//{
		//	name: "更新发布文章成功",
		//	mock: func(ctrl *gomock.Controller) (article.ReaderRepository,article.AuthorRepository) {}
		//}

	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			readerRepo, authorRepo := tc.mock(ctrl)
			svc := NewArticleServiceV1(readerRepo, authorRepo, nil)
			id, err := svc.PublishV1(context.Background(), tc.article)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantId, id)
		})
	}
}
