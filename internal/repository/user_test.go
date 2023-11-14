package repository_test

//
//import (
//	"context"
//	"github.com/jasonzhao47/cuddle/internal/domain"
//	"github.com/jasonzhao47/cuddle/internal/repository/cache"
//	"github.com/jasonzhao47/cuddle/internal/repository/dao"
//	"github.com/jasonzhao47/cuddle/wire/repository"
//	"go.uber.org/mock/gomock"
//	"testing"
//)
//
//func TestCacheUserRepository_FindById(t *testing.T) {
//	testCases := []struct {
//		name string
//		mock func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)
//
//		// 输入
//		ctx context.Context
//		id  int64
//
//		// 输出
//		wantUser domain.User
//		wantErr  error
//	}{
//		{
//			name: "Should find user in cache",
//		},
//		{
//			name: "Should find user not in cache",
//		},
//		{
//			name: "User not found",
//		},
//	}
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//
//			userDao, userCache := tc.mock(ctrl)
//
//			repo := repository.NewUserRepository(userDao, userCache)
//
//		})
//
//	}
//}
