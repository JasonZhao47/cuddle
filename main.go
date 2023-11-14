package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

func main() {
	// all-in-one initialization for web server
	initViper()
	server := InitWebServer()
	// run a health check
	server.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "I'm still alive!")
	})
	// run on a port
	err := server.Run(viper.GetString("server.http.addr"))
	if err != nil {
		panic(err)
	}
	// shutdown the server gracefully
}

func initViper() {
	// load in minimum configs
	// including server itself
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	// relative directory for viper(or for go, precisely) is looking from working directory
	viper.AddConfigPath("configs")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
