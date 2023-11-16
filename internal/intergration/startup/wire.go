//go:build wireinject

package startup

import (
	"github.com/google/wire"
	"github.com/jasonzhao47/cuddle/internal/repository"
	"github.com/jasonzhao47/cuddle/internal/repository/dao"
	"github.com/jasonzhao47/cuddle/internal/service"
	"github.com/jasonzhao47/cuddle/internal/web"
)

var (
	thirdPartyDep          = wire.NewSet(InitDB, InitLog)
	articleServiceProvider = wire.NewSet(dao.NewArticleGormDAO)
)

func InitArticleHandler(dao dao.ArticleDAO) *web.ArticleHandler {
	wire.Build(thirdPartyDep,
		repository.NewArticleRepository,
		service.NewArticleService,
		web.NewArticleHandler)
	return &web.ArticleHandler{}
}
