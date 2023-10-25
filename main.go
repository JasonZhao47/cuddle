package main

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jasonzhao47/cuddle/configs"
	"github.com/jasonzhao47/cuddle/internal/repository"
	"github.com/jasonzhao47/cuddle/internal/repository/cache"
	"github.com/jasonzhao47/cuddle/internal/repository/dao"
	"github.com/jasonzhao47/cuddle/internal/service"
	"github.com/jasonzhao47/cuddle/internal/web"
	"github.com/jasonzhao47/cuddle/internal/web/middleware"
	"github.com/jasonzhao47/cuddle/pkg/ginx/middleware/ratelimit"
	redis "github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

func main() {
	// init database
	// with configuration
	db := initDB()
	// init user Handlers
	// initialize a server
	redisClient := redis.NewClient(&redis.Options{
		Addr: configs.Config.Redis.Addr,
	})
	userCache := cache.NewUserCache(redisClient)
	codeCache := cache.NewRedisCodeCache(redisClient)
	server := initWebServer(redisClient)
	initUserHandlers(db, server, userCache, codeCache)
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

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(configs.Config.DB.DSN))

	if err != nil {
		panic(err)
	}
	// should init tables here at dao layer
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	// db.AutoMigrate
	return db
}

func initWebServer(redisClient redis.Cmdable) *gin.Engine {
	server := gin.Default()
	// CORS
	// 429 too much requests
	server.Use(cors.New(cors.Config{
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
	}))
	//server.Use(someMiddleware)
	//useSession(server)
	useJWT(server)
	// init redis client
	// 为啥这里要用builder
	// builder - 某个条件跟另外几个参数强烈耦合，否则就退化成了构造函数
	// 统一的话也都可以无脑builder，因为相当于上位替代版本
	// server.Use(middleware.LoginJWTMiddlewareBuilder{}.)
	server.Use(ratelimit.
		NewBuilder(redisClient, 1*time.Second, 3).
		Build())

	// other stuffs, JWT, session...
	// related to web layer
	// useJWT()
	// or use Session to store
	// or both
	return server
}

func initUserHandlers(db *gorm.DB, server *gin.Engine,
	userCache *cache.UserCache,
	codeCache cache.CodeCache) {
	// engines and database initialization
	// 切分的方向
	// dao
	userDAO := dao.NewUserDAO(db)
	// repo
	userRepo := repository.NewUserRepository(userDAO, userCache)
	codeRepo := repository.NewCodeRepository(codeCache)
	// service
	userService := service.NewUserService(userRepo)
	codeService := service.NewCodeService(codeRepo)
	// handler
	hdl := web.NewUserHandler(userService, codeService)
	// route
	hdl.RegisterRoutes(server)
}

func useSession(server *gin.Engine) {
	// stores the secret key for the algorithm
	loginMiddlewareBuilder := &middleware.LoginMiddlewareBuilder{}
	store := cookie.NewStore([]byte("secret_key"))
	// key in cookie
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

func useJWT(server *gin.Engine) {
	loginJWTBuilder := &middleware.LoginJWTBuilder{
		LoginPathRegExp: regexp.MustCompile(middleware.LoginPathPattern, regexp.None),
	}
	server.Use(loginJWTBuilder.Build())
}
