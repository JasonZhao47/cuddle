//go:build wireinject

package startup

import (
	"github.com/google/wire"
	"github.com/jasonzhao47/cuddle/internal/repository"
	"github.com/jasonzhao47/cuddle/internal/repository/cache"
	"github.com/jasonzhao47/cuddle/internal/repository/dao"
	"github.com/jasonzhao47/cuddle/internal/service"
	"github.com/jasonzhao47/cuddle/internal/web"
)

var (
	thirdPartyDep          = wire.NewSet(InitDB, InitLog, InitRedis)
	articleServiceProvider = wire.NewSet(dao.NewArticleGormDAO, dao.NewUserActivityDAO)
)

func InitArticleHandler(dao dao.ArticleDAO, activityDao dao.UserActivityDAO) *web.ArticleHandler {
	wire.Build(
		thirdPartyDep,
		cache.NewArticleCache,
		cache.NewUserActivityCache,
		repository.NewArticleRepository,
		repository.NewCacheUserActivityRepository,
		service.NewArticleService,
		service.NewUserActivityService,
		web.NewArticleHandler)
	return &web.ArticleHandler{}
}
