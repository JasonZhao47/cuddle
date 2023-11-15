package service

import (
	"github.com/jasonzhao47/cuddle/internal/domain"
	"github.com/jasonzhao47/cuddle/internal/repository/dao"
)

type ArticleService interface {
	Detail(int64) *domain.Article
}

type articleService struct {
	dao dao.ArticleDAO
}

func NewArticleService(dao dao.ArticleDAO) ArticleService {
	return &articleService{dao: dao}
}

func (svc *articleService) Detail(id int64) *domain.Article {
	//TODO implement me
	panic("implement me")
}
