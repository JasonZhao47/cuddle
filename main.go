package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/http"
)

func main() {
	// all-in-one initialization for web server
	initViperV2()
	//initThirdParty()
	app := InitWebApp()
	initPrometheus()

	// run a health check
	app.server.GET("/health", func(ctx *gin.Context) {
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
	err = app.server.Run(c.Addr)
	if err != nil {
		panic(err)
	}
	// shutdown the server gracefully
}

// initViper is deprecated: Use V2 instead.
//func initViper() {
//	// load in minimum configs
//	// including server itself
//	viper.SetConfigName("dev")
//	viper.SetConfigType("yaml")
//	// relative directory for viper(or for go, precisely) is looking from working directory
//	viper.AddConfigPath("configs")
//	err := viper.ReadInConfig()
//	if err != nil {
//		panic(fmt.Errorf("fatal error config file: %w", err))
//	}
//}

func initViperV2() {
	cfile := pflag.String("config", "configs/dev.yaml", "path to config")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initPrometheus() {
	go func() {
		// 专门给 prometheus 用的端口
		http.Handle("/metrics", promhttp.Handler())
		_ = http.ListenAndServe(":8081", nil)
	}()
}
