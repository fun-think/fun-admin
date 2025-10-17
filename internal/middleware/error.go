package middleware

import (
	"fun-admin/api/v1"
	"fun-admin/pkg/errors"
	"fun-admin/pkg/logger"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorMiddleware 全局错误处理中间件
func ErrorMiddleware(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录panic堆栈信息
				stack := debug.Stack()
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("stack", string(stack)),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
				)

				// 返回服务器错误响应
				if !c.Writer.Written() {
					v1.HandleServerError(c, errors.ErrInternalError)
				}
				c.Abort()
			}
		}()

		c.Next()

		// 处理业务层返回的错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			logger.Error("Business error",
				zap.Error(err),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
			)

			if !c.Writer.Written() {
				v1.HandleError(c, err)
			}
		}
	}
}

// NotFoundHandler 404处理器
func NotFoundHandler(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Warn("Route not found",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", c.ClientIP()),
		)

		v1.HandleNotFound(c)
	}
}

// MethodNotAllowedHandler 405处理器
func MethodNotAllowedHandler(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Warn("Method not allowed",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", c.ClientIP()),
		)

		err := errors.New(errors.CodeBadRequest, "方法不允许")
		v1.HandleError(c, err)
	}
}
