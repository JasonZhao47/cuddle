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
	"strconv"
	"testing"
	"time"
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
			Id: 15,
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
		// 集成测试不需要mock

		before func(*testing.T)
		after  func(*testing.T)

		wantCode   int
		wantResult Result[int64]
	}{
		{
			name: "编辑文章",
			req: Article{
				Id:      20026,
				Topic:   "Title for testing",
				Content: "Content",
			},
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				var art dao.Article
				s.db.Where("id = ?", 20026).First(&art)
				assert.Equal(t, "Title for testing", art.Topic)
				assert.Equal(t, "Content", art.Content)
				assert.Equal(t, int64(15), art.AuthorId)
				assert.True(t, art.CTime > 0)
				assert.True(t, art.UTime > 0)
			},
			wantCode: 200,
			wantResult: Result[int64]{
				Data: 20026,
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

func (s *ArticleHandlerSuite) TestList() {
	// suite 拿取T的方法
	// test list all results
	t := s.T()
	testCases := []struct {
		name string
		req  Article

		before func(t *testing.T)
		after  func(t *testing.T)

		wantCode   int
		wantResult Result[[]Article]
	}{
		{
			name: "获取所有结果",
			req: Article{
				AuthorId: 1,
			},
			before: func(t *testing.T) {
				var arts []dao.Article
				for i := 0; i < 3; i++ {
					arts = append(arts, dao.Article{
						Id:       int64(i),
						AuthorId: 1,
						Topic:    "Test Topic " + strconv.Itoa(i),
						Content:  "Test Content " + strconv.Itoa(i),
						Status:   0,
						CTime:    time.Now().UnixMilli(),
						UTime:    time.Now().UnixMilli(),
					})
				}
				s.db.Create(arts)
			},
			after: func(t *testing.T) {

			},
			wantCode: 200,
			wantResult: Result[[]Article]{
				Code: 200,
				Msg:  "Success",
				Data: []Article{
					{
						Id:       0,
						AuthorId: 1,
						Topic:    "Test Topic 1",
						Content:  "Test Content 1",
					},
					{
						Id:       1,
						AuthorId: 1,
						Topic:    "Test Topic 2",
						Content:  "Test Content 2",
					},
					{
						Id:       2,
						AuthorId: 1,
						Topic:    "Test Topic 3",
						Content:  "Test Content 3",
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// before testing
			tc.before(t)
			// mock a request
			data, err := json.Marshal(tc.req)
			assert.NoError(t, err)
			// request: marshal / unmarshal
			req, err := http.NewRequest(http.MethodGet, "/articles/list", bytes.NewReader(data))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			// send a request
			// set up mock server
			s.server.ServeHTTP(recorder, req)
			code := recorder.Code
			assert.Equal(t, tc.wantCode, code)
			if code != http.StatusOK {
				return
			}
			var result Result[int64]
			err = json.Unmarshal(recorder.Body.Bytes(), &result)
			assert.NoError(t, err)
			// get result
			// compare
			assert.Equal(t, tc.wantResult, tc.wantCode)
			// after
			tc.after(t)
		})
	}
}

func (s *ArticleHandlerSuite) TearDownTest() {
	err := s.db.Exec("truncate table `articles`").Error
	assert.NoError(s.T(), err)
}

func TestArticleHandler(t *testing.T) {
	suite.Run(t, &ArticleHandlerSuite{})
}

type Article struct {
	Id       int64  `json:"id"`
	Topic    string `json:"topic"`
	Content  string `json:"content"`
	AuthorId int64  `json:"author_id"`
}
