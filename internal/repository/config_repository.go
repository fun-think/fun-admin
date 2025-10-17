package repository

import (
	"context"
	"fun-admin/internal/model"
)

type ConfigRepository interface {
	GetConfig(ctx context.Context, key string) (*model.Config, error)
	CreateConfig(ctx context.Context, config *model.Config) error
	UpdateConfig(ctx context.Context, config *model.Config) error
	GetConfigsByGroup(ctx context.Context, group string) ([]model.Config, error)
	SearchConfigs(ctx context.Context, keyword string, group string) ([]model.Config, error)
	GetConfigGroups(ctx context.Context) ([]string, error)
	DeleteConfig(ctx context.Context, key string) error
	GetConfigCountByKey(ctx context.Context, key string) (int64, error)
	CreateConfigIfNotExists(ctx context.Context, config *model.Config) error
}

type configRepository struct {
	*Repository
}

func NewConfigRepository(repo *Repository) ConfigRepository {
	return &configRepository{repo}
}

// GetConfig 获取系统设置
func (r *configRepository) GetConfig(ctx context.Context, key string) (*model.Config, error) {
	var config model.Config
	if err := r.DB(ctx).Where("key = ?", key).First(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

// CreateConfig 创建系统设置
func (r *configRepository) CreateConfig(ctx context.Context, config *model.Config) error {
	return r.DB(ctx).Create(config).Error
}

// UpdateConfig 更新系统设置
func (r *configRepository) UpdateConfig(ctx context.Context, config *model.Config) error {
	return r.DB(ctx).Save(config).Error
}

// GetConfigsByGroup 根据分组获取系统设置
func (r *configRepository) GetConfigsByGroup(ctx context.Context, group string) ([]model.Config, error) {
	var configs []model.Config
	if err := r.DB(ctx).Where("group = ?", group).Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// SearchConfigs 搜索系统设置
func (r *configRepository) SearchConfigs(ctx context.Context, keyword string, group string) ([]model.Config, error) {
	db := r.DB(ctx).Model(&model.Config{})
	if keyword != "" {
		db = db.Where("key LIKE ? OR name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if group != "" {
		db = db.Where("group = ?", group)
	}

	var configs []model.Config
	if err := db.Find(&configs).Error; err != nil {
		return nil, err
	}

	return configs, nil
}

// GetConfigGroups 获取设置分组列表
func (r *configRepository) GetConfigGroups(ctx context.Context) ([]string, error) {
	var groups []string
	if err := r.DB(ctx).Model(&model.Config{}).Distinct().Pluck("group", &groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// DeleteConfig 删除系统设置
func (r *configRepository) DeleteConfig(ctx context.Context, key string) error {
	return r.DB(ctx).Where("key = ?", key).Delete(&model.Config{}).Error
}

// GetConfigCountByKey 根据key获取配置数量
func (r *configRepository) GetConfigCountByKey(ctx context.Context, key string) (int64, error) {
	var count int64
	if err := r.DB(ctx).Model(&model.Config{}).Where("key = ?", key).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CreateConfigIfNotExists 如果配置不存在则创建
func (r *configRepository) CreateConfigIfNotExists(ctx context.Context, config *model.Config) error {
	var count int64
	if err := r.DB(ctx).Model(&model.Config{}).Where("key = ?", config.Key).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return r.DB(ctx).Create(config).Error
	}
	return nil
}
