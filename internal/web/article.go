package web

import (
	"github.com/gin-gonic/gin"
	"github.com/jasonzhao47/cuddle/internal/domain"
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
	type Req struct {
		Topic   string `json:"topic"`
		Content string `json:"content"`
		Id      int64  `json:"id"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 登陆态
	uc := ctx.MustGet("user").(UserClaim)
	// 为什么需要id?
	id, err := h.svc.Save(ctx, &domain.Article{
		Id: req.Id,
		// if a pointer is used to access a struct
		// are all the sub structs also copied?
		// yes
		Author: domain.Author{
			Id: uc.Id,
		},
		Topic:   req.Topic,
		Content: req.Content,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "编辑失败",
			Data: nil,
		})
		h.l.Error("编辑存储失败了", logger.Int64("id", req.Id), logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: id,
	})
}

func (h *ArticleHandler) Publish(ctx *gin.Context) {
	// 不能修改别人的
	// 登陆态
	// 只能修改自己的
	type Req struct {
		Id      int64  `json:"id"`
		Topic   string `json:"topic"`
		Content string `json:"content"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		h.l.Warn("发布失败，参数不对", logger.Int64("id", req.Id), logger.Error(err))
		return
	}
	user := ctx.MustGet("user").(UserClaim)
	id, err := h.svc.Publish(ctx, &domain.Article{
		Id: req.Id,
		Author: domain.Author{
			Id: user.Id,
		},
		Topic:   req.Topic,
		Content: req.Content,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "发布失败",
		})
		h.l.Error("发布失败", logger.Int64("id", req.Id), logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: id,
	})
	return
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
		// bad intention
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Warn("非法查询文章",
			logger.Int64("user_id", user.Id),
			logger.Int64("id", id))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: nil,
	})
}

func (h *ArticleHandler) List(ctx *gin.Context) {
	type Req struct {
		Page     int
		PageSize int
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 登陆态
	user := ctx.MustGet("user").(UserClaim)
	arts, err := h.svc.List(ctx, user.Id, req.Page, req.PageSize)
	// 不要在这里检测author了
	// 可以认为，在handler不必处理业务逻辑
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "帖子未找到",
		})
		h.l.Error("查询帖子错误",
			logger.Int64("user_id", user.Id),
			logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: arts,
	})
}

func (h *ArticleHandler) Withdraw(ctx *gin.Context) {

}
