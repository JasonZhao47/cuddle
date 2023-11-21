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
		req  PageInfo

		before func(t *testing.T)
		after  func(t *testing.T)

		wantCode   int
		wantResult Result[[]Article]
	}{
		{
			name: "获取所有结果",
			req: PageInfo{
				Page:     1,
				PageSize: 10,
			},
			before: func(t *testing.T) {
				var arts []*dao.Article
				for i := 0; i < 3; i++ {
					arts = append(arts, &dao.Article{
						AuthorId: 15,
						Topic:    "Test Topic " + strconv.Itoa(i+1),
						Content:  "Test Content " + strconv.Itoa(i+1),
						CTime:    time.Now().UnixMilli(),
						UTime:    time.Now().UnixMilli(),
					})
				}
				s.db.Create(&arts)
			},
			after: func(t *testing.T) {

			},
			wantCode: 200,
			wantResult: Result[[]Article]{
				Data: []Article{
					{
						Id: 1,
						Author: Author{
							Id: 15,
						},
						Topic:   "Test Topic 1",
						Content: "Test Content 1",
					},
					{
						Id: 2,
						Author: Author{
							Id: 15,
						},
						Topic:   "Test Topic 2",
						Content: "Test Content 2",
					},
					{
						Id: 3,
						Author: Author{
							Id: 15,
						},
						Topic:   "Test Topic 3",
						Content: "Test Content 3",
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
			req, err := http.NewRequest(http.MethodPost, "/articles/list", bytes.NewReader(data))
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
			var result Result[[]Article]
			err = json.Unmarshal(recorder.Body.Bytes(), &result)
			assert.NoError(t, err)
			// get result
			// compare
			assert.Equal(t, tc.wantResult, result)
			// after
			tc.after(t)
		})
	}
}

func (s *ArticleHandlerSuite) TestArticlePublish() {
	t := s.T()
	testCases := []struct {
		name   string
		req    Article
		before func(t *testing.T)
		after  func(t *testing.T)

		wantCode   int
		wantResult Result[int64]
	}{
		{
			name: "新建并发布帖子",
			req: Article{
				Id:      1,
				Topic:   "今天天气不错",
				Content: "不错不错",
				Author: Author{
					Id: 15,
				},
			},
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				var art dao.Article
				s.db.Where("id = ?", 1).First(&art)
				assert.Equal(t, "Title for testing", art.Topic)
				assert.Equal(t, "Content", art.Content)
				assert.Equal(t, int64(15), art.AuthorId)
				assert.True(t, art.CTime > 0)
				assert.True(t, art.UTime > 0)
			},
			wantCode: 0,
			wantResult: Result[int64]{
				Data: 1,
			},
		},
		{
			name: "发布编辑过的帖子",
			req: Article{
				Id:      1,
				Topic:   "今天天气不错",
				Content: "不错不错",
				Author: Author{
					Id: 15,
				},
			},
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				var art dao.Article
				s.db.Where("id = ?", 1).First(&art)
				assert.Equal(t, "Title for testing", art.Topic)
				assert.Equal(t, "Content", art.Content)
				assert.Equal(t, int64(15), art.AuthorId)
				assert.True(t, art.CTime > 0)
				assert.True(t, art.UTime > 0)
			},
			wantCode: 0,
			wantResult: Result[int64]{
				Data: 1,
			},
		},
		{
			name: "发布别人创建的帖子，失败",
			req: Article{
				Id:      1,
				Topic:   "我今天吃了安妮意大利餐厅",
				Content: "披萨很好吃",
				Author:  Author{Id: 1},
			},
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				var art dao.Article
				s.db.Where("id = ?", 1).First(&art)
				assert.Equal(t, "Title for testing", art.Topic)
				assert.Equal(t, "Content", art.Content)
				assert.Equal(t, int64(15), art.AuthorId)
				assert.True(t, art.CTime > 0)
				assert.True(t, art.UTime > 0)
			},
			wantCode: 0,
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
			req, err := http.NewRequest(http.MethodPost, "/articles/publish", bytes.NewReader(data))
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			s.server.ServeHTTP(recorder, req)
			code := recorder.Code
			if code != http.StatusOK {
				return
			}
			assert.Equal(t, tc.wantCode, code)
			var result Result[int64]
			err = json.Unmarshal(recorder.Body.Bytes(), &result)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResult, result)
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
	Id      int64  `json:"id"`
	Topic   string `json:"topic"`
	Content string `json:"content"`
	Author  Author `json:"author"`
}

type PageInfo struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type Author struct {
	Id int64 `json:"id"`
}
