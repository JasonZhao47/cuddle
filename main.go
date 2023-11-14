package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/http"
)

func main() {
	// all-in-one initialization for web server
	initViperV2()
	//initLogger()
	//initThirdParty()
	server := InitWebServer()
	// run a health check
	server.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "I'm still alive!")
	})
	// run on a port
	type config struct {
		Addr string `yaml:"addr"`
	}
	var c config
	c = config{
		Addr: ":8080",
	}
	err := viper.UnmarshalKey("server.http", &c)
	if err != nil {
		panic(err)
	}
	err = server.Run(c.Addr)
	if err != nil {
		panic(err)
	}
	// shutdown the server gracefully
}

// Deprecated: Use V2 instead
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

func initViperV2() {
	cfile := pflag.String("config", "config/dev.yaml", "path to config")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
