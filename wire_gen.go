// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/jasonzhao47/cuddle/internal/domain/event"
	"github.com/jasonzhao47/cuddle/internal/domain/event/article"
	"github.com/jasonzhao47/cuddle/internal/repository"
	"github.com/jasonzhao47/cuddle/internal/repository/cache"
	"github.com/jasonzhao47/cuddle/internal/repository/dao"
	"github.com/jasonzhao47/cuddle/internal/service"
	"github.com/jasonzhao47/cuddle/internal/web"
	"github.com/jasonzhao47/cuddle/ioc"
)

// Injectors from wire.go:

func InitWebApp() *App {
	cmdable := ioc.InitRedis()
	v := ioc.GinMiddlewares(cmdable)
	db := ioc.InitDB()
	userDAO := dao.NewUserDAO(db)
	userCache := cache.NewRedisUserCache(cmdable)
	userRepository := repository.NewCacheUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository)
	codeCache := cache.NewRedisCodeCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache)
	smsService := ioc.InitSMSService()
	decorator := ioc.InitPrometheusService(smsService)
	codeService := service.NewSMSCodeService(codeRepository, decorator)
	logger := ioc.InitLogger()
	userHandler := web.NewUserHandler(userService, codeService, logger)
	articleDAO := dao.NewArticleGormDAO(db)
	articleCache := cache.NewArticleCache(cmdable)
	articleRepository := repository.NewArticleRepository(articleDAO, articleCache)
	client := ioc.InitSaramaClient()
	syncProducer := ioc.InitSyncProducer(client)
	producer := article.NewSaramaSyncProducer(syncProducer)
	articleService := service.NewArticleService(articleRepository, producer)
	userActivityCache := cache.NewUserActivityCache(cmdable)
	userActivityDAO := dao.NewUserActivityDAO(db)
	userActivityRepository := repository.NewCacheUserActivityRepository(userActivityCache, userActivityDAO)
	userActivityService := service.NewUserActivityService(userActivityRepository)
	articleHandler := web.NewArticleHandler(articleService, userActivityService, logger)
	engine := ioc.InitWebServer(v, userHandler, articleHandler)
	userActivityEventConsumer := event.NewUserActivityEventConsumer(userActivityRepository, client, logger)
	v2 := ioc.InitConsumers(userActivityEventConsumer)
	app := &App{
		server:    engine,
		consumers: v2,
	}
	return app
}
