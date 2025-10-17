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
