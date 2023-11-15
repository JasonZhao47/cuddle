//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/jasonzhao47/cuddle/internal/repository"
	"github.com/jasonzhao47/cuddle/internal/repository/cache"
	"github.com/jasonzhao47/cuddle/internal/repository/dao"
	"github.com/jasonzhao47/cuddle/internal/service"
	"github.com/jasonzhao47/cuddle/internal/web"
	"github.com/jasonzhao47/cuddle/ioc"
)

func InitWebServer() *gin.Engine {
	// 按照从底到上来排列项目的依赖
	wire.Build(
		// ioc部分，公用组件集成 —— 数据库、缓存、日志
		ioc.InitRedis, ioc.InitDB, ioc.InitLogger,

		// DAO
		dao.NewUserDAO,
		dao.NewArticleGormDAO,
		// cache
		cache.NewRedisUserCache,
		cache.NewRedisCodeCache,
		// repository
		repository.NewCacheUserRepository,
		repository.NewCodeRepository,
		// service
		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,
		// handler
		web.NewUserHandler,
		web.NewArticleHandler,
		// middleware
		ioc.GinMiddlewares,
		// Web服务器
		ioc.InitWebServer,
	)
	return gin.Default()
}
