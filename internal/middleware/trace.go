package middleware

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// TraceMiddleware 追踪中间件，为每个请求生成唯一的追踪ID
func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从请求头获取追踪ID
		traceID := c.GetHeader("X-Trace-ID")

		// 如果没有提供，则生成一个新的
		if traceID == "" {
			traceID = generateTraceID()
		}

		// 设置到上下文中
		c.Set("trace_id", traceID)

		// 设置响应头
		c.Header("X-Trace-ID", traceID)

		c.Next()
	}
}

// generateTraceID 生成追踪ID
func generateTraceID() string {
	timestamp := time.Now().UnixNano()
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)

	return fmt.Sprintf("%d-%x", timestamp, randomBytes)
}
