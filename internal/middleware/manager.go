package middleware

import (
	"fun-admin/internal/repository"
	"fun-admin/pkg/logger"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// Manager 中间件管理器
type Manager struct {
	logger     *logger.Logger
	enforcer   *casbin.SyncedEnforcer
	repository *repository.Repository
	config     *viper.Viper
}

// NewManager 创建中间件管理器
func NewManager(logger *logger.Logger, db *gorm.DB, enforcer *casbin.SyncedEnforcer, repository *repository.Repository, config *viper.Viper) *Manager {
	return &Manager{
		logger:     logger,
		enforcer:   enforcer,
		repository: repository,
		config:     config,
	}
}

// SetupCORSMiddleware 设置CORS中间件
func (m *Manager) SetupCORSMiddleware() gin.HandlerFunc {
	cfg := CORSConfig{
		AllowedOrigins:   m.config.GetStringSlice("http.cors.allowed_origins"),
		AllowedMethods:   m.config.GetStringSlice("http.cors.allowed_methods"),
		AllowedHeaders:   m.config.GetStringSlice("http.cors.allowed_headers"),
		AllowCredentials: m.config.GetBool("http.cors.allow_credentials"),
		MaxAge:           m.config.GetDuration("http.cors.max_age"),
	}
	return CORSMiddleware(cfg)
}

// SetupAuthMiddleware 设置认证中间件
func (m *Manager) SetupAuthMiddleware() gin.HandlerFunc {
	return AuthMiddleware(m.enforcer)
}

// SetupOperationLogMiddleware 设置操作日志中间件
func (m *Manager) SetupOperationLogMiddleware() gin.HandlerFunc {
	// 修复：正确调用 OperationLogMiddleware
	operationLogMiddleware := NewOperationLogMiddleware(m.logger.Logger, m.repository.DB(nil))
	return operationLogMiddleware.Handle()
}

// SetupWebhookMiddleware 设置Webhook中间件
func (m *Manager) SetupWebhookMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Webhook验证逻辑
		c.Next()
	}
}

// SetupPermissionMiddleware 设置权限中间件
func (m *Manager) SetupPermissionMiddleware() gin.HandlerFunc {
	return PermissionMiddleware(m.enforcer)
}

// SetupRateLimitMiddleware 设置限流中间件
func (m *Manager) SetupRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddleware(60, time.Minute)
}

// SetupRequestLogMiddleware 设置请求日志中间件
func (m *Manager) SetupRequestLogMiddleware() gin.HandlerFunc {
	return RequestLogMiddleware(m.logger)
}

// SetupSecurityMiddleware 设置安全中间件
func (m *Manager) SetupSecurityMiddleware() gin.HandlerFunc {
	return SecurityMiddleware()
}

// SetupSignMiddleware 设置签名中间件
func (m *Manager) SetupSignMiddleware() gin.HandlerFunc {
	return SignMiddleware(m.logger, m.config)
}

// SetupTraceMiddleware 设置追踪中间件
func (m *Manager) SetupTraceMiddleware() gin.HandlerFunc {
	return TraceMiddleware()
}

// SetupErrorMiddleware 设置错误处理中间件
func (m *Manager) SetupErrorMiddleware() gin.HandlerFunc {
	return ErrorMiddleware(m.logger)
}

// SetupNotFoundHandler 设置404处理器
func (m *Manager) SetupNotFoundHandler(r *gin.Engine) {
	r.NoRoute(NotFoundHandler(m.logger))
}

// SetupMethodNotAllowedHandler 设置405处理器
func (m *Manager) SetupMethodNotAllowedHandler(r *gin.Engine) {
	r.NoMethod(MethodNotAllowedHandler(m.logger))
}
