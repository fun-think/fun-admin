package model

import (
	"time"
)

// OperationLog 操作日志模型
type OperationLog struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UserID       uint      `json:"user_id"`       // 操作用户ID
	UserName     string    `json:"user_name"`     // 操作用户名称
	IP           string    `json:"ip"`            // 操作IP
	Method       string    `json:"method"`        // HTTP方法 (GET, POST, PUT, DELETE)
	Path         string    `json:"path"`          // 请求路径
	UserAgent    string    `json:"user_agent"`    // 用户代理
	RequestData  string    `json:"request_data"`  // 请求数据
	ResponseData string    `json:"response_data"` // 响应数据
	StatusCode   int       `json:"status_code"`   // 状态码
	Duration     int64     `json:"duration"`      // 执行时长(毫秒)
	Description  string    `json:"description"`   // 操作描述
	Resource     string    `json:"resource"`      // 操作资源
	Action       string    `json:"action"`        // 操作类型 (create, update, delete, etc.)
}

// TableName 设置表名
func (OperationLog) TableName() string {
	return "admin_operation_log"
}
