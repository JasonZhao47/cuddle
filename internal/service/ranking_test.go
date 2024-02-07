package service

import (
	"context"
	"github.com/jasonzhao47/cuddle/internal/domain"
	svcmock "github.com/jasonzhao47/cuddle/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestBatchRankingService_GetTopN(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (UserActivityService, ArticleService)
		wantErr error
		wantRes []domain.PublishedArticle
	}{
		{
			name: "成功获取前三排名元素",
			mock: func(ctrl *gomock.Controller) (UserActivityService, ArticleService) {
				userActSvc := svcmock.NewMockUserActivityService(ctrl)
				articleSvc := svcmock.NewMockArticleService(ctrl)
				articleSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return([]domain.PublishedArticle{
					{Id: 1, Utime: now},
					{Id: 2, Utime: now},
				}, nil)
				articleSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return([]domain.PublishedArticle{
					{Id: 3, Utime: now},
					{Id: 4, Utime: now},
				}, nil)
				articleSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return([]domain.PublishedArticle{}, nil)
				userActSvc.EXPECT().GetReadByIds(gomock.Any(), "article", gomock.Any()).
					Return([]domain.UserActivity{
						{Id: 1, ReadCnt: 1},
						{Id: 2, ReadCnt: 2},
					}, nil)
				userActSvc.EXPECT().GetReadByIds(gomock.Any(), "article", gomock.Any()).
					Return([]domain.UserActivity{
						{Id: 3, ReadCnt: 3},
						{Id: 4, ReadCnt: 4},
					}, nil)
				userActSvc.EXPECT().GetReadByIds(gomock.Any(), "article", gomock.Any()).
					Return([]domain.UserActivity{}, nil)
				return userActSvc, articleSvc
			},
			wantErr: nil,
			wantRes: []domain.PublishedArticle{
				{Id: 4, Utime: now},
				{Id: 3, Utime: now},
				{Id: 2, Utime: now},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userActSvc, artSvc := tc.mock(ctrl)
			svc := &BatchRankingService{
				n:         3,
				batchSize: 2,
				artSvc:    artSvc,
				usrActSvc: userActSvc,
				scoreFn: func(readCnt int64, utime time.Time) float64 {
					return float64(readCnt)
				},
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
			defer cancel()
			res, err := svc.TopN(ctx)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)
		})
	}

}
