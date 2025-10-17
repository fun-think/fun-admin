package repository

import (
	"context"
	"errors"
	"fun-admin/internal/model"
	"fun-admin/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ProfileRepository 个人资料仓库接口
type ProfileRepository interface {
	GetProfile(ctx context.Context, userID uint) (*model.User, error)
	UpdateProfile(ctx context.Context, userID uint, updates map[string]interface{}) error
	UpdatePassword(ctx context.Context, userID uint, newPassword string) error
}

// profileRepository 个人资料仓库实现
type profileRepository struct {
	logger *logger.Logger
	db     *gorm.DB
}

// NewProfileRepository 创建个人资料仓库
func NewProfileRepository(
	logger *logger.Logger,
	db *gorm.DB,
) ProfileRepository {
	return &profileRepository{
		logger: logger,
		db:     db,
	}
}

// GetProfile 获取用户个人资料
func (r *profileRepository) GetProfile(ctx context.Context, userID uint) (*model.User, error) {
	var profile model.User
	if err := r.db.WithContext(ctx).First(&profile, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		r.logger.Error("获取用户个人资料失败", zap.Error(err))
		return nil, err
	}

	return &profile, nil
}

// UpdateProfile 更新用户个人资料
func (r *profileRepository) UpdateProfile(ctx context.Context, userID uint, updates map[string]interface{}) error {
	// 验证用户存在
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		r.logger.Error("验证用户失败", zap.Error(err))
		return err
	}

	// 更新用户资料
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		r.logger.Error("更新用户个人资料失败", zap.Error(err))
		return err
	}

	return nil
}

// UpdatePassword 更新用户密码
func (r *profileRepository) UpdatePassword(ctx context.Context, userID uint, newPassword string) error {
	// 验证用户存在
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		r.logger.Error("验证用户失败", zap.Error(err))
		return err
	}

	// 更新密码（此处未加密，实际项目中应该加密处理）
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Update("password", newPassword).Error; err != nil {
		r.logger.Error("更新密码失败", zap.Error(err))
		return err
	}

	return nil
}
