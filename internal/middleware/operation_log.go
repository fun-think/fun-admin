package middleware

import (
	"bytes"
	"io"
	"strings"
	"time"

	"fun-admin/internal/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// OperationLogMiddleware 操作日志中间件
type OperationLogMiddleware struct {
	logger *zap.Logger
	db     *gorm.DB
}

// NewOperationLogMiddleware 创建操作日志中间件
func NewOperationLogMiddleware(logger *zap.Logger, db *gorm.DB) *OperationLogMiddleware {
	return &OperationLogMiddleware{
		logger: logger,
		db:     db,
	}
}

// Handle 处理操作日志记录
func (l *OperationLogMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只记录 admin 相关的 API 操作
		if !strings.HasPrefix(c.Request.URL.Path, "/api/admin") &&
			!strings.HasPrefix(c.Request.URL.Path, "/api/v1") {
			c.Next()
			return
		}

		// 记录请求开始时间
		startTime := time.Now()

		// 读取请求体内容
		var reqBody []byte
		if c.Request.Body != nil {
			reqBody, _ = io.ReadAll(c.Request.Body)
			// 重新设置请求体，因为 ReadAll 会消耗掉
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}

		// 设置操作描述、资源和动作类型
		l.setOperationInfo(c)

		// 继续处理请求
		c.Next()

		// 请求处理完成后记录日志
		go l.saveOperationLog(c, startTime, reqBody)
	}
}

// setOperationInfo 设置操作相关的信息
func (l *OperationLogMiddleware) setOperationInfo(c *gin.Context) {
	c.Set("description", l.getOperationDescription(c))
	c.Set("resource", l.getResourceFromPath(c.Request.URL.Path))
	c.Set("action", l.getActionFromMethod(c.Request.Method))
}

// getOperationDescription 根据请求路径和方法生成操作描述
func (l *OperationLogMiddleware) getOperationDescription(c *gin.Context) string {
	method := c.Request.Method
	path := c.Request.URL.Path

	// 根据路径和方法生成描述
	switch {
	case strings.Contains(path, "/create") || method == "POST":
		return "创建记录"
	case strings.Contains(path, "/update") || method == "PUT":
		return "更新记录"
	case strings.Contains(path, "/delete") || method == "DELETE":
		return "删除记录"
	case strings.Contains(path, "/export"):
		return "导出数据"
	default:
		return "查看记录"
	}
}

// getResourceFromPath 从路径中提取资源名称
func (l *OperationLogMiddleware) getResourceFromPath(path string) string {
	// 例如: /api/v1/users/1 -> users
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "v1" && i+1 < len(parts) {
			return parts[i+1]
		}
		if part == "admin" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return "unknown"
}

// getActionFromMethod 根据HTTP方法获取操作类型
func (l *OperationLogMiddleware) getActionFromMethod(method string) string {
	switch method {
	case "POST":
		return "create"
	case "PUT":
		return "update"
	case "DELETE":
		return "delete"
	case "GET":
		return "read"
	default:
		return "other"
	}
}

// getUserInfo 从上下文中获取用户信息
func (l *OperationLogMiddleware) getUserInfo(c *gin.Context) (uint, string) {
	user, exists := c.Get("user")
	if !exists {
		return 0, "Anonymous"
	}

	if userModel, ok := user.(*model.User); ok {
		return userModel.ID, userModel.Username
	}

	return 0, "Unknown"
}

// saveOperationLog 保存操作日志
func (l *OperationLogMiddleware) saveOperationLog(c *gin.Context, start time.Time, reqBody []byte) {
	// 计算执行时长
	duration := time.Since(start)

	// 获取用户信息
	userID, username := l.getUserInfo(c)

	// 构建操作日志记录
	logEntry := &model.OperationLog{
		UserID:      userID,
		UserName:    username,
		IP:          c.ClientIP(),
		Method:      c.Request.Method,
		Path:        c.Request.URL.Path,
		UserAgent:   c.Request.UserAgent(),
		RequestData: string(reqBody),
		StatusCode:  c.Writer.Status(),
		Duration:    duration.Milliseconds(),
		Description: c.GetString("description"),
		Resource:    c.GetString("resource"),
		Action:      c.GetString("action"),
	}

	// 截取部分请求数据避免过长
	if len(logEntry.RequestData) > 1000 {
		logEntry.RequestData = logEntry.RequestData[:1000] + "... (truncated)"
	}

	// 保存到数据库
	if err := l.db.Create(logEntry).Error; err != nil {
		l.logger.Error("Failed to save operation logger", zap.Error(err))
	}
}
