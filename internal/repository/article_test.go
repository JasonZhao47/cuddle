package repository

import (
	"context"
	daomock "github.com/jasonzhao47/cuddle/internal/dao/mocks"
	"github.com/jasonzhao47/cuddle/internal/domain"
	"github.com/jasonzhao47/cuddle/internal/repository/dao"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestArticleRepository_Sync(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(controller *gomock.Controller) dao.ArticleDAO
		art     *domain.Article
		wantId  int64
		wantErr error
	}{
		{
			name: "创建并同步",
			mock: func(controller *gomock.Controller) dao.ArticleDAO {
				artDao := daomock.NewMockArticleDAO(controller)
				artDao.EXPECT().Sync(gomock.Any(), gomock.Any()).Return(int64(1), nil)
				return artDao
			},
			art: &domain.Article{
				Author: domain.Author{
					Id: 1,
				},
				Topic:   "标题",
				Content: "内容",
			},
			wantId:  1,
			wantErr: nil,
		},
		{
			name: "修改并同步",
			mock: func(controller *gomock.Controller) dao.ArticleDAO {
				artDao := daomock.NewMockArticleDAO(controller)
				artDao.EXPECT().Sync(gomock.Any(), gomock.Any()).Return(int64(1), nil)
				return artDao
			},
			art: &domain.Article{
				Id: 1,
				Author: domain.Author{
					Id: 1,
				},
				Topic:   "标题",
				Content: "内容",
			},
			wantId:  1,
			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			artDao := tc.mock(ctrl)
			repo := NewArticleRepository(artDao)
			id, err := repo.Sync(context.Background(), tc.art)
			assert.Equal(t, tc.wantId, id)
			assert.Equal(t, tc.wantErr, err)
		})
	}

}
