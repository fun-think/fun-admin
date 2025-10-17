package service

import (
	"context"
	"errors"
	"fun-admin/internal/model"
	"fun-admin/internal/repository"

	"go.uber.org/zap"
)

// ProfileService 个人资料服务
type ProfileService struct {
	*Service
	logger      *zap.Logger
	profileRepo repository.ProfileRepository
}

// ProfileServiceInterface 个人资料服务接口
type ProfileServiceInterface interface {
	GetProfile(ctx context.Context, userID uint) (*model.User, error)
	UpdateProfile(ctx context.Context, userID uint, updates map[string]interface{}) error
	UpdatePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error
}

// NewProfileService 创建个人资料服务
func NewProfileService(
	service *Service,
	logger *zap.Logger,
	profileRepo repository.ProfileRepository,
) ProfileServiceInterface {
	return &ProfileService{
		Service:     service,
		logger:      logger,
		profileRepo: profileRepo,
	}
}

// GetProfile 获取用户个人资料
func (s *ProfileService) GetProfile(ctx context.Context, userID uint) (*model.User, error) {
	return s.profileRepo.GetProfile(ctx, userID)
}

// UpdateProfile 更新用户个人资料
func (s *ProfileService) UpdateProfile(ctx context.Context, userID uint, updates map[string]interface{}) error {
	return s.profileRepo.UpdateProfile(ctx, userID, updates)
}

// UpdatePassword 更新用户密码
func (s *ProfileService) UpdatePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	// 验证用户存在
	user, err := s.profileRepo.GetProfile(ctx, userID)
	if err != nil {
		return err
	}

	// 验证旧密码（简化处理，实际应该使用密码哈希比较）
	if user.Password != oldPassword {
		return errors.New("旧密码不正确")
	}

	// 更新密码（此处未加密，实际项目中应该加密处理）
	return s.profileRepo.UpdatePassword(ctx, userID, newPassword)
}
