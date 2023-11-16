package intergration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jasonzhao47/cuddle/internal/intergration/startup"
	"github.com/jasonzhao47/cuddle/internal/repository/dao"
	"github.com/jasonzhao47/cuddle/internal/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
type ArticleHandlerSuite struct {
	suite.Suite
	db     *gorm.DB
	server *gin.Engine
}

func (s *ArticleHandlerSuite) SetupSuite() {
	s.db = startup.InitDB()
	hdl := startup.InitArticleHandler(dao.NewArticleGormDAO(s.db))
	server := gin.Default()
	server.Use(func(ctx *gin.Context) {
		ctx.Set(`user`, web.UserClaim{
			Id: 123,
		})
	})
	hdl.RegisterRoutes(server)
	s.server = server
}

func (s *ArticleHandlerSuite) TestEdit() {
	t := s.T()
	testCases := []struct {
		name string
		req  Article
		// 在难以mock的时候，自己定义一个模型

		before func(*testing.T)
		after  func(*testing.T)

		wantCode   int
		wantResult Result[int64]
	}{
		{
			name: "编辑文章",
			req: Article{
				Id:      1,
				Title:   "",
				Content: "",
			},
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				var art dao.Article
				s.db.Where("id = ?", art.Id).First(&art)
				assert.Equal(t, "Title for testing", art.Topic)
				assert.Equal(t, "Content", art.Content)
				assert.Equal(t, int64(15), art.AuthorId)
				assert.True(t, art.CTime > 0)
				assert.True(t, art.UTime > 0)
			},
			wantCode: 200,
			wantResult: Result[int64]{
				Data: 1,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)

			data, err := json.Marshal(tc.req)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/articles/edit", bytes.NewReader(data))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			s.server.ServeHTTP(recorder, req)
			code := recorder.Code
			assert.Equal(t, tc.wantCode, code)
			if code != http.StatusOK {
				return
			}
			var result Result[int64]
			err = json.Unmarshal(recorder.Body.Bytes(), &result)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResult, result)
			// deal with databases here
			tc.after(t)
		})
	}
}

func TestArticleHandler(t *testing.T) {
	suite.Run(t, &ArticleHandlerSuite{})
}

type Article struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
