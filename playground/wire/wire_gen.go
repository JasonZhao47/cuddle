// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	repository2 "github.com/jasonzhao47/cuddle/playground/wire/repository"
	"github.com/jasonzhao47/cuddle/playground/wire/repository/dao"
)

// Injectors from wire.go:

func InitUserRepository() *repository2.UserRepository {
	db := repository2.InitDB()
	userDAO := dao.NewUserDAO(db)
	userRepository := repository2.NewUserRepository(userDAO)
	return userRepository
}
