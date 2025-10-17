package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityMiddleware 安全中间件
func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置安全头
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'; media-src 'self'; object-src 'none'; child-src 'none'; worker-src 'none'; frame-ancestors 'none'; form-action 'self'; base-uri 'self'; manifest-src 'self'")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		// 移除服务器信息
		c.Header("Server", "")

		c.Next()
	}
}

// ContentSecurityPolicyMiddleware CSP中间件
func ContentSecurityPolicyMiddleware(policy string) gin.HandlerFunc {
	if policy == "" {
		policy = "default-src 'self'"
	}

	return func(c *gin.Context) {
		c.Header("Content-Security-Policy", policy)
		c.Next()
	}
}

// HTTPSRedirectMiddleware HTTPS重定向中间件
func HTTPSRedirectMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Forwarded-Proto") == "http" {
			c.Redirect(301, "https://"+c.Request.Host+c.Request.RequestURI)
			c.Abort()
			return
		}
		c.Next()
	}
}
