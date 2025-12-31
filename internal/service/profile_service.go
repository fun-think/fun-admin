package service

import (
	"context"
	"errors"
	"fun-admin/internal/model"
	"fun-admin/internal/repository"

	"golang.org/x/crypto/bcrypt"
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
	// 验证用户存在并获取用户信息
	user, err := s.profileRepo.GetProfile(ctx, userID)
	if err != nil {
		return err
	}

	// 验证旧密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("旧密码不正确")
	}

	// 对新密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 更新密码
	return s.profileRepo.UpdatePassword(ctx, userID, string(hashedPassword))
}