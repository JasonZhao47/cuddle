package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jasonzhao47/cuddle/internal/domain/event/article"
	"github.com/robfig/cron/v3"
)

type App struct {
	server    *gin.Engine
	consumers []article.Consumer
	cron      *cron.Cron
}
