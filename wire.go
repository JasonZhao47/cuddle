//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/jasonzhao47/cuddle/internal/domain/event"
	"github.com/jasonzhao47/cuddle/internal/repository"
	"github.com/jasonzhao47/cuddle/internal/repository/cache"
	"github.com/jasonzhao47/cuddle/internal/repository/dao"
	"github.com/jasonzhao47/cuddle/internal/service"
	"github.com/jasonzhao47/cuddle/internal/web"
	"github.com/jasonzhao47/cuddle/ioc"
)

func InitWebApp() *App {
	// 按照从底到上来排列项目的依赖
	wire.Build(
		// ioc部分，公用组件集成 —— 数据库、缓存、日志、第三方
		ioc.InitRedis, ioc.InitDB, ioc.InitLogger, ioc.InitSaramaClient,

		// DAO
		dao.NewUserDAO,
		dao.NewArticleGormDAO,
		dao.NewUserActivityDAO,
		// cache
		cache.NewRedisUserCache,
		cache.NewRedisCodeCache,
		cache.NewArticleCache,
		cache.NewUserActivityCache,
		// repository
		repository.NewCacheUserRepository,
		repository.NewCodeRepository,
		repository.NewArticleRepository,
		repository.NewCacheUserActivityRepository,
		event.NewUserActivityEventConsumer,
		ioc.InitConsumers,
		// service
		ioc.InitSMSService,
		ioc.InitPrometheusService,
		service.NewUserService,
		service.NewSMSCodeService,
		service.NewArticleService,
		service.NewUserActivityService,
		// handler
		web.NewUserHandler,
		web.NewArticleHandler,
		// middleware
		ioc.GinMiddlewares,
		// Web服务器
		ioc.InitWebServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
