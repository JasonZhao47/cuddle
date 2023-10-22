package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
}

func (m *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			return
		}
		session := sessions.Default(ctx)
		userId := session.Get("userId")
		if userId == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 不用refresh token实现简易刷新
		// maxTime = now - then < 30
		now := time.Now()

		updateTime := session.Get("updateTime")
		lastUpdateTime, ok := updateTime.(time.Time)
		if updateTime == nil || now.Sub(lastUpdateTime) > time.Second*300 || !ok {
			session.Set(updateTime, now)
			session.Set("userId", userId)
			err := session.Save()
			if err != nil {
				// logging
			}
		}
	}
}
