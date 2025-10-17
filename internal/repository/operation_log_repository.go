package repository

import (
	"context"
	"fun-admin/internal/model"
	"fun-admin/pkg/logger"

	"gorm.io/gorm"
)

// OperationLogRepository 操作日志仓库接口
type OperationLogRepository interface {
	GetOperationLogs(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*model.OperationLog, int64, error)
	GetOperationLog(ctx context.Context, id uint) (*model.OperationLog, error)
	DeleteOperationLog(ctx context.Context, id uint) error
	DeleteOperationLogs(ctx context.Context, ids []uint) error
	ClearOperationLogs(ctx context.Context) error
}

// operationLogRepository 操作日志仓库实现
type operationLogRepository struct {
	logger *logger.Logger
	db     *gorm.DB
}

// NewOperationLogRepository 创建操作日志仓库
func NewOperationLogRepository(
	logger *logger.Logger,
	db *gorm.DB,
) OperationLogRepository {
	return &operationLogRepository{
		logger: logger,
		db:     db,
	}
}

// GetOperationLogs 获取操作日志列表
func (r *operationLogRepository) GetOperationLogs(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*model.OperationLog, int64, error) {
	var list []*model.OperationLog
	db := r.db.WithContext(ctx).Model(&model.OperationLog{}).Order("id DESC")

	// 应用过滤条件
	if method, ok := filters["method"]; ok && method != "" {
		db = db.Where("method LIKE ?", "%"+method.(string)+"%")
	}
	if path, ok := filters["path"]; ok && path != "" {
		db = db.Where("path LIKE ?", "%"+path.(string)+"%")
	}
	if ip, ok := filters["ip"]; ok && ip != "" {
		db = db.Where("ip = ?", ip)
	}
	if userAgent, ok := filters["user_agent"]; ok && userAgent != "" {
		db = db.Where("user_agent LIKE ?", "%"+userAgent.(string)+"%")
	}
	if statusCode, ok := filters["status_code"]; ok && statusCode != "" {
		db = db.Where("status_code = ?", statusCode)
	}

	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

// GetOperationLog 获取操作日志详情
func (r *operationLogRepository) GetOperationLog(ctx context.Context, id uint) (*model.OperationLog, error) {
	var operationLog model.OperationLog
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&operationLog).Error
	return &operationLog, err
}

// DeleteOperationLog 删除操作日志
func (r *operationLogRepository) DeleteOperationLog(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.OperationLog{}).Error
}

// DeleteOperationLogs 批量删除操作日志
func (r *operationLogRepository) DeleteOperationLogs(ctx context.Context, ids []uint) error {
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&model.OperationLog{}).Error
}

// ClearOperationLogs 清空操作日志
func (r *operationLogRepository) ClearOperationLogs(ctx context.Context) error {
	return r.db.WithContext(ctx).Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.OperationLog{}).Error
}
