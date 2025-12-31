package service

import (
	"context"
	"fun-admin/internal/model"
	"fun-admin/internal/repository"
)

// OperationLogService 操作日志服务
type OperationLogService struct {
	*Service
	operationLogRepo repository.OperationLogRepository
}

// OperationLogServiceInterface 操作日志服务接口
type OperationLogServiceInterface interface {
	GetOperationLogs(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*model.OperationLog, int64, error)
	GetOperationLog(ctx context.Context, id uint) (*model.OperationLog, error)
	DeleteOperationLog(ctx context.Context, id uint) error
	DeleteOperationLogs(ctx context.Context, ids []uint) error
	ClearOperationLogs(ctx context.Context) error
	GetOperationLogStats(ctx context.Context) (*OperationLogStats, error)
}

// NewOperationLogService 创建操作日志服务
func NewOperationLogService(
	service *Service,
	operationLogRepo repository.OperationLogRepository,
) OperationLogServiceInterface {
	return &OperationLogService{
		Service:          service,
		operationLogRepo: operationLogRepo,
	}
}

// GetOperationLogs 获取操作日志列表
func (s *OperationLogService) GetOperationLogs(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*model.OperationLog, int64, error) {
	return s.operationLogRepo.GetOperationLogs(ctx, page, pageSize, filters)
}

// GetOperationLog 获取操作日志详情
func (s *OperationLogService) GetOperationLog(ctx context.Context, id uint) (*model.OperationLog, error) {
	return s.operationLogRepo.GetOperationLog(ctx, id)
}

// DeleteOperationLog 删除操作日志
func (s *OperationLogService) DeleteOperationLog(ctx context.Context, id uint) error {
	return s.operationLogRepo.DeleteOperationLog(ctx, id)
}

// DeleteOperationLogs 批量删除操作日志
func (s *OperationLogService) DeleteOperationLogs(ctx context.Context, ids []uint) error {
	return s.operationLogRepo.DeleteOperationLogs(ctx, ids)
}

// ClearOperationLogs 清空操作日志
func (s *OperationLogService) ClearOperationLogs(ctx context.Context) error {
	return s.operationLogRepo.ClearOperationLogs(ctx)
}

// CountStat 统计项
type CountStat struct {
	Key   string `json:"key"`
	Count int64  `json:"count"`
}

// OperationLogStats 日志统计
type OperationLogStats struct {
	Total        int64            `json:"total"`
	Today        int64            `json:"today"`
	LastSeven    int64            `json:"last_seven"`
	MethodCounts map[string]int64 `json:"method_counts"`
	StatusCounts map[string]int64 `json:"status_counts"`
	TopResources []CountStat      `json:"top_resources"`
	TopUsers     []CountStat      `json:"top_users"`
}

// GetOperationLogStats 返回统计信息
func (s *OperationLogService) GetOperationLogStats(ctx context.Context) (*OperationLogStats, error) {
	result, err := s.operationLogRepo.GetOperationLogStats(ctx)
	if err != nil {
		return nil, err
	}

	stats := &OperationLogStats{
		Total:        result.Total,
		Today:        result.Today,
		LastSeven:    result.LastSeven,
		MethodCounts: result.MethodCounts,
		StatusCounts: result.StatusCounts,
		TopResources: make([]CountStat, len(result.TopResources)),
		TopUsers:     make([]CountStat, len(result.TopUsers)),
	}
	for i, item := range result.TopResources {
		stats.TopResources[i] = CountStat{Key: item.Key, Count: item.Count}
	}
	for i, item := range result.TopUsers {
		stats.TopUsers[i] = CountStat{Key: item.Key, Count: item.Count}
	}
	return stats, nil
}
