package middleware

import (
	v1 "fun-admin/api/v1"
	"fun-admin/pkg/jwt"
	"fun-admin/pkg/logger"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Jwt(j *jwt.JWT, logger *logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.Request.Header.Get("Authorization")
		if tokenString == "" {
			logger.WithContext(ctx).Warn("No token", zap.Any("data", map[string]interface{}{
				"url":    ctx.Request.URL,
				"params": ctx.Params,
			}))
			v1.HandleUnauthorized(ctx)
			ctx.Abort()
			return
		}

		// 移除Bearer前缀后再解析
		tokenString = strings.TrimSpace(tokenString)
		if strings.HasPrefix(strings.ToLower(tokenString), "bearer ") {
			tokenString = tokenString[7:] // 移除"Bearer "前缀
		}

		claims, err := j.ParseToken(tokenString)
		if err != nil {
			logger.WithContext(ctx).Error("token error", zap.Any("data", map[string]interface{}{
				"url":    ctx.Request.URL,
				"params": ctx.Params,
			}), zap.Error(err))
			v1.HandleUnauthorized(ctx)
			ctx.Abort()
			return
		}

		ctx.Set("claims", claims)
		recoveryLoggerFunc(ctx, logger)
		ctx.Next()
	}
}

func NoStrictAuth(j *jwt.JWT, logger *logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.Request.Header.Get("Authorization")
		if tokenString == "" {
			tokenString, _ = ctx.Cookie("accessToken")
		}
		if tokenString == "" {
			tokenString = ctx.Query("accessToken")
		}
		if tokenString == "" {
			ctx.Next()
			return
		}

		// 移除Bearer前缀后再解析
		tokenString = strings.TrimSpace(tokenString)
		if strings.HasPrefix(strings.ToLower(tokenString), "bearer ") {
			tokenString = tokenString[7:] // 移除"Bearer "前缀
		}

		claims, err := j.ParseToken(tokenString)
		if err != nil {
			ctx.Next()
			return
		}

		ctx.Set("claims", claims)
		recoveryLoggerFunc(ctx, logger)
		ctx.Next()
	}
}

func recoveryLoggerFunc(ctx *gin.Context, logger *logger.Logger) {
	if userInfo, ok := ctx.MustGet("claims").(*jwt.MyCustomClaims); ok {
		logger.WithValue(ctx, zap.Any("UserId", userInfo.UserId))
	}
}
