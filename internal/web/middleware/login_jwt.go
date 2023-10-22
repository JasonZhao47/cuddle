package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jasonzhao47/cuddle/internal/web"
	"net/http"
	"strings"
	"time"
)

type LoginJWTBuilder struct {
}

func (builder *LoginJWTBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			return
		}
		// 根据JWT验证登录信息

		// 解析Bearer的内容
		authCode := ctx.GetHeader("Authorization")
		if authCode == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// Bearer xxxxxxxxx
		seg := strings.Split(authCode, " ")
		if len(seg) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := seg[1]
		var userClaim web.UserClaim
		token, err := jwt.ParseWithClaims(tokenStr, &userClaim, func(token *jwt.Token) (interface{}, error) {
			return web.JWTKey, nil
		})
		if err != nil {
			// token解析有问题
			// 伪造的
			fmt.Println("Wrong token")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid {
			fmt.Println("Invalid token")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if userClaim.UserAgent != ctx.GetHeader("User-Agent") {
			// logger
			fmt.Println(`User-Agent changed from `, userClaim)
		}
		// then - now < 50
		// remaining
		expireTime := userClaim.ExpiresAt
		if expireTime.Sub(time.Now()) < time.Second*50 {
			// 刷新expires时间
			userClaim.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 5))
			// 用JWTKey生成一个新的key，发给客户端的header
			tokenStr, err := token.SignedString(web.JWTKey)
			if err != nil {
				// log: 会有啥问题？
			}
			ctx.Header("x-jwt-key", tokenStr)
		}
		// 方便后续的接口拿到UID作为业务使用
		ctx.Set("user", userClaim)
		// 解析完成之后刷新ExpireAt
	}
}
