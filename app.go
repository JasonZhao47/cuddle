package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jasonzhao47/cuddle/internal/domain/event"
	"github.com/robfig/cron/v3"
)

type App struct {
	server    *gin.Engine
	consumers []event.Consumer
	cron      *cron.Cron
}
