package repository

import (
	"context"
	"fmt"
	"fun-admin/internal/model"
	"fun-admin/pkg/logger"
	"time"

	"gorm.io/gorm"
)

// OperationLogRepository 操作日志仓库接口
type OperationLogRepository interface {
	GetOperationLogs(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*model.OperationLog, int64, error)
	GetOperationLog(ctx context.Context, id uint) (*model.OperationLog, error)
	DeleteOperationLog(ctx context.Context, id uint) error
	DeleteOperationLogs(ctx context.Context, ids []uint) error
	ClearOperationLogs(ctx context.Context) error
	GetOperationLogStats(ctx context.Context) (*OperationLogStatsResult, error)
}

type operationLogRepository struct {
	logger *logger.Logger
	db     *gorm.DB
}

// NewOperationLogRepository 创建仓库实例
func NewOperationLogRepository(logger *logger.Logger, db *gorm.DB) OperationLogRepository {
	return &operationLogRepository{
		logger: logger,
		db:     db,
	}
}

// GetOperationLogs 获取日志列表
func (r *operationLogRepository) GetOperationLogs(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*model.OperationLog, int64, error) {
	var list []*model.OperationLog
	query := r.db.WithContext(ctx).Model(&model.OperationLog{}).Order("id DESC")

	if method, ok := filters["method"]; ok && method != "" {
		query = query.Where("method LIKE ?", "%"+fmt.Sprint(method)+"%")
	}
	if path, ok := filters["path"]; ok && path != "" {
		query = query.Where("path LIKE ?", "%"+fmt.Sprint(path)+"%")
	}
	if ip, ok := filters["ip"]; ok && ip != "" {
		query = query.Where("ip = ?", ip)
	}
	if userAgent, ok := filters["user_agent"]; ok && userAgent != "" {
		query = query.Where("user_agent LIKE ?", "%"+fmt.Sprint(userAgent)+"%")
	}
	if statusCode, ok := filters["status_code"]; ok && statusCode != "" {
		query = query.Where("status_code = ?", statusCode)
	}
	if userID, ok := filters["user_id"]; ok && userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if resource, ok := filters["resource"]; ok && resource != "" {
		query = query.Where("resource = ?", resource)
	}
	if action, ok := filters["action"]; ok && action != "" {
		query = query.Where("action = ?", action)
	}
	if keyword, ok := filters["keyword"]; ok && keyword != "" {
		kw := fmt.Sprintf("%%%s%%", keyword)
		query = query.Where("(path LIKE ? OR description LIKE ? OR user_name LIKE ?)", kw, kw, kw)
	}
	if from, ok := filters["created_at_from"]; ok && from != "" {
		query = query.Where("created_at >= ?", from)
	}
	if to, ok := filters["created_at_to"]; ok && to != "" {
		query = query.Where("created_at <= ?", to)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// GetOperationLog 获取单条日志
func (r *operationLogRepository) GetOperationLog(ctx context.Context, id uint) (*model.OperationLog, error) {
	var operationLog model.OperationLog
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&operationLog).Error
	return &operationLog, err
}

// DeleteOperationLog 删除日志
func (r *operationLogRepository) DeleteOperationLog(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.OperationLog{}).Error
}

// DeleteOperationLogs 批量删除
func (r *operationLogRepository) DeleteOperationLogs(ctx context.Context, ids []uint) error {
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&model.OperationLog{}).Error
}

// ClearOperationLogs 清空日志
func (r *operationLogRepository) ClearOperationLogs(ctx context.Context) error {
	return r.db.WithContext(ctx).Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.OperationLog{}).Error
}

// OperationLogStatsResult 聚合统计
type OperationLogStatsResult struct {
	Total        int64
	Today        int64
	LastSeven    int64
	MethodCounts map[string]int64
	StatusCounts map[string]int64
	TopResources []CountResult
	TopUsers     []CountResult
}

// CountResult 通用 Key/Count 结构
type CountResult struct {
	Key   string
	Count int64
}

// GetOperationLogStats 返回统计信息
func (r *operationLogRepository) GetOperationLogStats(ctx context.Context) (*OperationLogStatsResult, error) {
	stats := &OperationLogStatsResult{
		MethodCounts: make(map[string]int64),
		StatusCounts: make(map[string]int64),
	}

	base := r.db.WithContext(ctx).Model(&model.OperationLog{})
	if err := base.Count(&stats.Total).Error; err != nil {
		return nil, err
	}

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	seven := now.AddDate(0, 0, -7)

	if err := base.Where("created_at >= ?", todayStart).Count(&stats.Today).Error; err != nil {
		return nil, err
	}
	if err := base.Where("created_at >= ?", seven).Count(&stats.LastSeven).Error; err != nil {
		return nil, err
	}

	var methodRows []struct {
		Method string
		Count  int64
	}
	if err := r.db.WithContext(ctx).
		Model(&model.OperationLog{}).
		Select("method, COUNT(*) as count").
		Group("method").
		Scan(&methodRows).Error; err != nil {
		return nil, err
	}
	for _, row := range methodRows {
		stats.MethodCounts[row.Method] = row.Count
	}

	var statusRows []struct {
		StatusCode int
		Count      int64
	}
	if err := r.db.WithContext(ctx).
		Model(&model.OperationLog{}).
		Select("status_code, COUNT(*) as count").
		Group("status_code").
		Scan(&statusRows).Error; err != nil {
		return nil, err
	}
	for _, row := range statusRows {
		key := fmt.Sprintf("%d", row.StatusCode)
		stats.StatusCounts[key] = row.Count
	}

	var resources []CountResult
	if err := r.db.WithContext(ctx).
		Model(&model.OperationLog{}).
		Select("resource as key, COUNT(*) as count").
		Where("resource <> ''").
		Group("resource").
		Order("count DESC").
		Limit(5).
		Scan(&resources).Error; err == nil {
		stats.TopResources = resources
	}

	var users []CountResult
	if err := r.db.WithContext(ctx).
		Model(&model.OperationLog{}).
		Select("user_name as key, COUNT(*) as count").
		Where("user_name <> ''").
		Group("user_name").
		Order("count DESC").
		Limit(5).
		Scan(&users).Error; err == nil {
		stats.TopUsers = users
	}

	return stats, nil
}
