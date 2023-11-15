package web

import (
	"github.com/gin-gonic/gin"
	"github.com/jasonzhao47/cuddle/internal/logger"
	"github.com/jasonzhao47/cuddle/internal/service"
	"net/http"
	"strconv"
)

type ArticleHandler struct {
	svc service.ArticleService
	l   logger.Logger
}

func NewArticleHandler(svc service.ArticleService, l logger.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
		l:   l,
	}
}

func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/articles")
	{
		ug.POST("/edit", h.Edit)
		ug.POST("/publish", h.Publish)
		ug.GET("/detail/:id", h.Detail)
		ug.POST("/list", h.List)
		ug.POST("/withdraw", h.Withdraw)
	}
}

func (h *ArticleHandler) Edit(ctx *gin.Context) {

}

func (h *ArticleHandler) Publish(ctx *gin.Context) {

}

func (h *ArticleHandler) Detail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return
	}
	h.svc.Detail(id)
	ctx.String(http.StatusOK, "")
}

func (h *ArticleHandler) List(ctx *gin.Context) {

}

func (h *ArticleHandler) Withdraw(ctx *gin.Context) {

}
