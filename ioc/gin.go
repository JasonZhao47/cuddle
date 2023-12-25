package ioc

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jasonzhao47/cuddle/internal/web"
	"github.com/jasonzhao47/cuddle/internal/web/middleware"
	prom "github.com/jasonzhao47/cuddle/pkg/ginx/middleware/prometheus"
	"github.com/jasonzhao47/cuddle/pkg/ginx/middleware/ratelimit"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitWebServer(
	middlewares []gin.HandlerFunc,
	userHandler *web.UserHandler,
	articleHandler *web.ArticleHandler) *gin.Engine {
	server := gin.Default()
	server.Use(middlewares...)
	userHandler.RegisterRoutes(server)
	articleHandler.RegisterRoutes(server)
	return server
}

func GinMiddlewares(cmd redis.Cmdable) []gin.HandlerFunc {
	// init redis client
	// 为啥这里要用builder
	// builder - 某个条件跟另外几个参数强烈耦合，否则就退化成了构造函数
	// 统一的话也都可以无脑builder，因为相当于上位替代版本
	// server.Use(middleware.LoginJWTMiddlewareBuilder{}.)
	// other stuffs, JWT, session...
	// related to web layer
	// useJWT()
	// or use Session to store
	// or both
	loginPathRegExp := regexp.MustCompile(middleware.LoginPathPattern, regexp.None)
	promBuilder := &prom.Builder{
		Namespace: "jason_zhao",
		Subsystem: "cuddle",
		Name:      "gin_http",
	}

	return []gin.HandlerFunc{
		// CORS
		// 429 too much requests
		corsHeader(),
		session(),
		promBuilder.BuildResponseTime(),
		promBuilder.BuildActiveRequests(),
		ratelimit.NewBuilder(cmd, time.Minute, 100).Build(),
		middleware.NewLoginJWTBuilder(loginPathRegExp).Build(),
	}
}

func corsHeader() gin.HandlerFunc {
	corsHdl := cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowCredentials: true,
		// AllowHeaders:
		// AllowMethods:
		ExposeHeaders: []string{"x-jwt-token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			// for production
			return strings.Contains(origin, "some_company.com")
		},
		// MaxAge: second precision
		MaxAge: 12 * time.Hour,
	})
	return corsHdl
}

func session() gin.HandlerFunc {
	// todo: update
	store := cookie.NewStore([]byte("secret_key"))
	return sessions.Sessions("ssid", store)
}

// deprecated
func useSession(server *gin.Engine) {
	// stores the secret key used for encryption algorithm
	loginMiddlewareBuilder := &middleware.LoginMiddlewareBuilder{}
	// key in cookie
	store := cookie.NewStore([]byte("secret_key"))
	// key in redis
	//store, err := redis.NewStore(16, "tcp",
	//	"localhost:6379",
	//	"",
	//	[]byte("9a0ff9e016a41dcd"),
	//	[]byte("898104dd3b97e4dc"))
	//if err != nil {
	//	panic(err)
	//}
	// only registers the session
	server.Use(sessions.Sessions("ssid", store))
	// check and protect each api
	server.Use(loginMiddlewareBuilder.Build())
}
