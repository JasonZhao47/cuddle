// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jasonzhao47/cuddle/internal/repository"
	"github.com/jasonzhao47/cuddle/internal/repository/cache"
	"github.com/jasonzhao47/cuddle/internal/repository/dao"
	"github.com/jasonzhao47/cuddle/internal/service"
	"github.com/jasonzhao47/cuddle/internal/web"
	"github.com/jasonzhao47/cuddle/ioc"
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	cmdable := ioc.InitRedis()
	v := ioc.GinMiddlewares(cmdable)
	db := ioc.InitDB()
	userDAO := dao.NewUserDAO(db)
	userCache := cache.NewRedisUserCache(cmdable)
	userRepository := repository.NewCacheUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository)
	codeCache := cache.NewRedisCodeCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache)
	codeService := service.NewCodeService(codeRepository)
	userHandler := web.NewUserHandler(userService, codeService)
	engine := ioc.InitWebServer(v, userHandler)
	return engine
}