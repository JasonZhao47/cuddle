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
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "id参数错误",
		})
		h.l.Warn("查询帖子失败，id格式不对",
			logger.Int64("id", id),
			logger.Error(err))
		return
	}
	article, err := h.svc.GetById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "帖子未找到",
		})
		h.l.Error("查询帖子失败",
			logger.Int64("article_id", id),
			logger.Error(err))
		return
	}
	user := ctx.MustGet("user").(UserClaim)
	if user.Id != article.Author.Id {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Warn("用户id参数错误",
			logger.Int64("user_id", user.Id),
			logger.Int64("id", id))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: nil,
	})
}

func (h *ArticleHandler) List(ctx *gin.Context) {

}

func (h *ArticleHandler) Withdraw(ctx *gin.Context) {

}
