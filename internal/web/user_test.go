package web

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jasonzhao47/cuddle/internal/domain"
	"github.com/jasonzhao47/cuddle/internal/service"
	svcmock "github.com/jasonzhao47/cuddle/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		reqBuilder func(t *testing.T) *http.Request

		wantBody string
		wantCode int
	}{
		{
			name: "Should signup",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), domain.User{
					Email:    "testmail123@gmail.com",
					Password: "test@123",
				}).Return(nil)
				codeSvc := svcmock.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewReader([]byte(`{
					"email": "testmail123@gmail.com",
					"password": "test@123",
					"confirmPassword": "test@123"}`))
				req, err := http.NewRequest(http.MethodPost, "/users/signup", body)
				req.Header.Set("Content-Type", "application/json")
				require.NoError(t, err)
				return req
			},
			wantBody: "注册成功",
			wantCode: 200,
		},
		{
			name: "Contains illegal mailbox info",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), gomock.Any()).Return(service.ErrDuplicateEmail)
				codeSvc := svcmock.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewReader([]byte(`{
					"email": "testmail123@gmail.com",
					"password": "test@123",
					"confirmPassword": "test@123"}`))
				req, err := http.NewRequest(http.MethodPost, "/users/signup", body)
				req.Header.Set("Content-Type", "application/json")
				require.NoError(t, err)
				return req
			},
			wantBody: "邮箱冲突，换一个邮箱",
			wantCode: 200,
		},
		{
			name: "Other errors",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), gomock.Any()).Return(errors.New("系统爆炸了"))
				codeSvc := svcmock.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewReader([]byte(`{
					"email": "test@gmail.com",
					"password": "test@123",
					"confirmPassword": "test@123"}`))
				req, err := http.NewRequest(http.MethodPost, "/users/signup", body)
				// 设置request上相关请求头
				req.Header.Set("Content-Type", "application/json")
				require.NoError(t, err)
				return req
			},
			wantBody: "系统错误",
			wantCode: 200,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userSvc, codeSvc := tc.mock(ctrl)
			handler := NewUserHandler(userSvc, codeSvc)

			server := gin.Default()
			handler.RegisterRoutes(server)

			recorder := httptest.NewRecorder()
			req := tc.reqBuilder(t)

			server.ServeHTTP(recorder, req)

			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantBody, recorder.Body.String())
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		reqBuilder func(t *testing.T) *http.Request

		wantBody string
		wantCode int
	}{
		{
			name: "Should login",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				userSvc.EXPECT().
					Login(gomock.Any(), "normaluser@gmail.com", "test@123").
					Return(domain.User{}, nil)
				codeSvc := svcmock.NewMockCodeService(ctrl)
				return userSvc, codeSvc

			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewReader([]byte(`
				{
					"email": "normaluser@gmail.com",
					"password": "test@123"
				}
				`))
				req, err := http.NewRequest(http.MethodPost, "/users/login", body)
				req.Header.Set("Content-Type", "application/json")
				assert.Equal(t, nil, err)
				return req
			},
			wantBody: "登录成功",
			wantCode: 200,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc, codeSvc := tc.mock(ctrl)
			userHdl := NewUserHandler(userSvc, codeSvc)

			server := gin.Default()
			userHdl.RegisterRoutes(server)

			recorder := httptest.NewRecorder()
			req := tc.reqBuilder(t)

			server.ServeHTTP(recorder, req)

			assert.Equal(t, tc.wantBody, recorder.Body.String())
			assert.Equal(t, tc.wantCode, recorder.Code)
		})
	}
}
