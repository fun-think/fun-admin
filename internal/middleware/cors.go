package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CORSConfig 用于配置跨域策略
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           time.Duration
}

// CORSMiddleware 返回带配置的跨域中间件
func CORSMiddleware(cfg CORSConfig) gin.HandlerFunc {
	defaultHeaders := []string{"Content-Type", "Authorization", "Accept", "Origin", "X-Requested-With"}
	defaultMethods := []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}

	allowedHeaders := joinValues(cfg.AllowedHeaders, defaultHeaders)
	allowedMethods := joinValues(cfg.AllowedMethods, defaultMethods)

	return func(c *gin.Context) {
		origin := strings.TrimSpace(c.GetHeader("Origin"))
		allowedOrigin, ok := resolveOrigin(origin, cfg.AllowedOrigins)
		if !ok {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Header("Access-Control-Allow-Origin", allowedOrigin)
		if cfg.AllowCredentials && origin != "" && allowedOrigin != "*" {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		c.Header("Access-Control-Allow-Headers", allowedHeaders)
		c.Header("Access-Control-Allow-Methods", allowedMethods)

		if cfg.MaxAge > 0 {
			c.Header("Access-Control-Max-Age", strconv.FormatInt(int64(cfg.MaxAge/time.Second), 10))
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func joinValues(current []string, fallback []string) string {
	values := current
	if len(values) == 0 {
		values = fallback
	}
	return strings.Join(values, ", ")
}

func resolveOrigin(origin string, allowed []string) (string, bool) {
	if origin == "" {
		if len(allowed) > 0 {
			return allowed[0], true
		}
		return "*", true
	}

	if len(allowed) == 0 {
		return origin, true
	}

	for _, item := range allowed {
		if item == "*" || strings.EqualFold(item, origin) {
			return origin, true
		}
	}

	return "", false
}
