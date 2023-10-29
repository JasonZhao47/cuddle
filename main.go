package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// all-in-one initialization for web server
	server := InitWebServer()
	// run a health check
	server.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "I'm still alive!")
	})
	// run on a port
	err := server.Run(":8089")
	if err != nil {
		panic(err)
	}
	// shutdown the server gracefully
}
