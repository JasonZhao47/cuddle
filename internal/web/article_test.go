package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jasonzhao47/cuddle/internal/domain"
	"github.com/jasonzhao47/cuddle/internal/service"
	svcmock "github.com/jasonzhao47/cuddle/internal/service/mocks"
	"github.com/jasonzhao47/cuddle/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestArticleHandler_Detail(t *testing.T) {
	testCases := []struct {
		name       string
		mock       func(*gomock.Controller) (service.ArticleService, service.UserActivityService)
		reqBuilder func(*testing.T) *http.Request

		wantCode int
		wantBody string
	}{
		{
			name: "返回帖子内容",
			mock: func(ctrl *gomock.Controller) (service.ArticleService, service.UserActivityService) {
				svc := svcmock.NewMockArticleService(ctrl)
				svc.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(&domain.Article{
					Author: domain.Author{
						Id: 1,
					},
				}, nil)
				actSvc := svcmock.NewMockUserActivityService(ctrl)
				return svc, actSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {

				articleId := 1
				url := fmt.Sprintf("/articles/detail/%d", articleId)
				req, err := http.NewRequest(http.MethodGet, url, nil)
				req.Header.Set("Content-Type", "application/json")
				require.NoError(t, err)
				return req
			},
			wantCode: 200,
			wantBody: `{"code":0,"msg":"","data":null}`,
		},
		{
			name: "返回id错误",
			mock: func(ctrl *gomock.Controller) (service.ArticleService, service.UserActivityService) {
				svc := svcmock.NewMockArticleService(ctrl)
				actSvc := svcmock.NewMockUserActivityService(ctrl)
				return svc, actSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				articleId := "bcdedit"
				url := fmt.Sprintf("/articles/detail/%v", articleId)
				req, err := http.NewRequest(http.MethodGet, url, nil)
				req.Header.Set("Content-Type", "application/json")
				require.NoError(t, err)
				return req
			},
			wantCode: 200,
			wantBody: `{"code":4,"msg":"id参数错误","data":null}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc, actSvc := tc.mock(ctrl)
			hdl := NewArticleHandler(svc, actSvc, logger.NewLogger(zap.L()))

			recorder := httptest.NewRecorder()
			req := tc.reqBuilder(t)

			server := gin.Default()
			server.Use(func(c *gin.Context) {
				c.Set("user", UserClaim{Id: 1})
				c.Next()
			})
			hdl.RegisterRoutes(server)
			server.ServeHTTP(recorder, req)

			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantBody, recorder.Body.String())
		})
	}
}
