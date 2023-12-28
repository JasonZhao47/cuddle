package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jasonzhao47/cuddle/internal/domain/event"
)

type App struct {
	server    *gin.Engine
	consumers []event.Consumer
}
