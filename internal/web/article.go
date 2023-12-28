package web

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jasonzhao47/cuddle/internal/domain"
	"github.com/jasonzhao47/cuddle/internal/service"
	"github.com/jasonzhao47/cuddle/pkg/ginx"
	logger2 "github.com/jasonzhao47/cuddle/pkg/logger"
	"net/http"
	"strconv"
	"time"
)

type ArticleHandler struct {
	svc        service.ArticleService
	userActSvc service.UserActivityService
	l          logger2.Logger
	biz        string
}

func NewArticleHandler(svc service.ArticleService, userActSvc service.UserActivityService, l logger2.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc:        svc,
		userActSvc: userActSvc,
		l:          l,
		biz:        "article",
	}
}

func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	art := server.Group("/articles")
	{
		art.POST("/edit", h.Edit)
		art.POST("/publish", h.Publish)
		art.GET("/detail/:id", h.Detail)
		art.POST("/list", h.List)
		art.POST("/withdraw", h.Withdraw)
	}

	pub := server.Group("/pub")
	{
		pub.GET("/:id", h.PubDetail)
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
	id, err := h.svc.Save(ctx, domain.Article{
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
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "编辑失败",
			Data: nil,
		})
		h.l.Error("编辑存储失败了", logger2.Int64("id", req.Id), logger2.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: id,
	})
}

func (h *ArticleHandler) Publish(ctx *gin.Context) {
	// 不能修改别人的
	// 登陆态
	// 只能修改自己的
	type Req struct {
		Id      int64
		Topic   string `json:"topic"`
		Content string `json:"content"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		h.l.Warn("发布失败，参数不对", logger2.Int64("id", req.Id), logger2.Error(err))
		return
	}
	user := ctx.MustGet("user").(UserClaim)
	id, err := h.svc.Publish(ctx, domain.Article{
		Id: req.Id,
		Author: domain.Author{
			Id: user.Id,
		},
		Topic:   req.Topic,
		Content: req.Content,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "发布失败",
		})
		h.l.Error("发布失败", logger2.Int64("id", req.Id), logger2.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: id,
	})
	return
}

func (h *ArticleHandler) Detail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 4,
			Msg:  "id参数错误",
		})
		h.l.Warn("查询帖子失败，id格式不对",
			logger2.String("id", idStr),
			logger2.Error(err))
		return
	}
	article, err := h.svc.GetById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "帖子未找到",
		})
		h.l.Error("查询帖子失败",
			logger2.Int64("article_id", id),
			logger2.Error(err))
		return
	}
	user := ctx.MustGet("user").(UserClaim)
	if user.Id != article.Author.Id {
		// bad intention
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Warn("非法查询文章",
			logger2.Int64("user_id", user.Id),
			logger2.Int64("id", id))
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: nil,
	})
}

func (h *ArticleHandler) List(ctx *gin.Context) {
	type Req struct {
		Limit  int
		Offset int
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 登陆态
	user := ctx.MustGet("user").(UserClaim)
	arts, err := h.svc.List(ctx, user.Id, req.Limit, req.Offset)
	// 不要在这里检测author了
	// 可以认为，在handler不必处理业务逻辑
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "帖子未找到",
		})
		h.l.Error("查询帖子错误",
			logger2.Int64("user_id", user.Id),
			logger2.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: arts,
	})
}

func (h *ArticleHandler) Withdraw(ctx *gin.Context) {
	type Req struct {
		Id int64
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 登录态
	user := ctx.MustGet("user").(UserClaim)
	err := h.svc.WithDraw(ctx, user.Id, req.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "隐藏失败",
		})
		h.l.Error("隐藏帖子失败了", logger2.Int64("user_id", user.Id),
			logger2.Int64("id", req.Id),
			logger2.Error(err))
	}
	ctx.JSON(http.StatusOK, ginx.Result{})
}

func (h *ArticleHandler) PubDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 4,
			Msg:  "id参数错误",
		})
		h.l.Warn("查询帖子失败，id不对",
			logger2.String("id", idStr),
			logger2.Error(err))
		return
	}
	// 不需要登录就可以看
	art, err := h.svc.GetPubById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "查看失败",
		})
		h.l.Error("查看帖子失败了",
			logger2.String("id", idStr),
			logger2.Error(err))
	}
	// 在这增加阅读数
	go func() {
		newCtx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		err = h.userActSvc.IncrRead(newCtx, h.biz, art.Id)
		if err != nil {
			h.l.Error("阅读量增加失败了",
				logger2.String("id", idStr),
				logger2.String("biz", "pub"),
				logger2.Int64("biz_id", 1),
				logger2.Error(err))
		}
	}()
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: art,
	})
}
