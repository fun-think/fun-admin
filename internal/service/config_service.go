package service

import (
	"context"
	"encoding/json"
	"fmt"
	"fun-admin/internal/model"
	"fun-admin/internal/repository"
	"fun-admin/pkg/email"
	"fun-admin/pkg/storage"
	"sync"
)

// EmailConfig 邮件配置
type EmailConfig struct {
	Host     string `json:"host"`      // SMTP服务器地址
	Port     int    `json:"port"`      // SMTP服务器端口
	Username string `json:"username"`  // 用户名
	Password string `json:"password"`  // 密码
	From     string `json:"from"`      // 发件人地址
	FromName string `json:"from_name"` // 发件人名称
}

// StorageConfig 存储配置
type StorageConfig struct {
	Type      string            `json:"type"`       // 存储类型
	Endpoint  string            `json:"endpoint"`   // 访问端点
	Region    string            `json:"region"`     // 区域
	AccessID  string            `json:"access_id"`  // 访问ID
	AccessKey string            `json:"access_key"` // 访问密钥
	Bucket    string            `json:"bucket"`     // 存储桶
	Domain    string            `json:"domain"`     // 自定义域名
	Extra     map[string]string `json:"extra"`      // 额外配置
}

// ConfigService 配置服务
type ConfigService struct {
	repo       *repository.Repository
	configRepo repository.ConfigRepository
	cache      map[string]string
	mu         sync.RWMutex
	emailSvc   *email.EmailService
	storageMgr *storage.Manager
}

// NewConfigService 创建配置服务
func NewConfigService(repo *repository.Repository, configRepo repository.ConfigRepository) *ConfigService {
	return &ConfigService{
		repo:       repo,
		configRepo: configRepo,
		cache:      make(map[string]string),
	}
}

// GetConfig 获取系统设置
func (s *ConfigService) GetConfig(ctx context.Context, key string) (string, error) {
	// 先从缓存获取
	s.mu.RLock()
	if value, exists := s.cache[key]; exists {
		s.mu.RUnlock()
		return value, nil
	}
	s.mu.RUnlock()

	// 从数据库获取
	config, err := s.configRepo.GetConfig(ctx, key)
	if err != nil {
		return "", fmt.Errorf("获取系统设置失败: %w", err)
	}

	// 更新缓存
	s.mu.Lock()
	s.cache[key] = config.Value
	s.mu.Unlock()

	return config.Value, nil
}

// SetConfig 设置系统设置
func (s *ConfigService) SetConfig(ctx context.Context, key, name, value, setType, group string) error {
	// 从数据库获取
	config, err := s.configRepo.GetConfig(ctx, key)
	if err != nil {
		// 如果记录不存在，创建新记录
		newConfig := model.Config{
			Key:   key,
			Name:  name,
			Value: value,
			Type:  setType,
			Group: group,
		}
		if err := s.configRepo.CreateConfig(ctx, &newConfig); err != nil {
			return fmt.Errorf("创建系统设置失败: %w", err)
		}
	} else {
		// 如果记录存在，更新值
		config.Name = name
		config.Value = value
		config.Type = setType
		config.Group = group
		if err := s.configRepo.UpdateConfig(ctx, config); err != nil {
			return fmt.Errorf("更新系统设置失败: %w", err)
		}
	}

	// 更新缓存
	s.mu.Lock()
	s.cache[key] = value
	s.mu.Unlock()

	return nil
}

// GetConfigByGroup 根据分组获取系统设置
func (s *ConfigService) GetConfigByGroup(ctx context.Context, group string) (map[string]string, error) {
	configs, err := s.configRepo.GetConfigsByGroup(ctx, group)
	if err != nil {
		return nil, fmt.Errorf("获取分组系统设置失败: %w", err)
	}

	result := make(map[string]string, len(configs))
	for _, config := range configs {
		result[config.Key] = config.Value
	}

	return result, nil
}

// GetEmailConfig 获取邮件配置
func (s *ConfigService) GetEmailConfig(ctx context.Context) (*EmailConfig, error) {
	configJSON, err := s.GetConfig(ctx, "email_config")
	if err != nil {
		return nil, err
	}

	if configJSON == "" {
		// 返回默认配置
		return &EmailConfig{
			Host: "smtp.example.com",
			Port: 587,
		}, nil
	}

	var config EmailConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("解析邮件配置失败: %w", err)
	}

	return &config, nil
}

// SetEmailConfig 设置邮件配置
func (s *ConfigService) SetEmailConfig(ctx context.Context, config *EmailConfig) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化邮件配置失败: %w", err)
	}

	return s.SetConfig(ctx, "email_config", "邮件配置", string(configJSON), "json", "email")
}

// GetStorageConfig 获取存储配置
func (s *ConfigService) GetStorageConfig(ctx context.Context) (*StorageConfig, error) {
	configJSON, err := s.GetConfig(ctx, "storage_config")
	if err != nil {
		return nil, err
	}

	if configJSON == "" {
		// 返回默认配置
		return &StorageConfig{
			Type: "local",
			Extra: map[string]string{
				"base_path": "storage/uploads",
			},
		}, nil
	}

	var config StorageConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("解析存储配置失败: %w", err)
	}

	return &config, nil
}

// SetStorageConfig 设置存储配置
func (s *ConfigService) SetStorageConfig(ctx context.Context, config *StorageConfig) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化存储配置失败: %w", err)
	}

	return s.SetConfig(ctx, "storage_config", "存储配置", string(configJSON), "json", "storage")
}

// GetEmailService 获取邮件服务
func (s *ConfigService) GetEmailService(ctx context.Context) (*email.EmailService, error) {
	if s.emailSvc != nil {
		return s.emailSvc, nil
	}

	config, err := s.GetEmailConfig(ctx)
	if err != nil {
		return nil, err
	}

	emailConfig := &email.EmailConfig{
		Host:     config.Host,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
		From:     config.From,
		FromName: config.FromName,
	}

	s.emailSvc = email.NewEmailService(emailConfig)
	return s.emailSvc, nil
}

// TestEmailConfig 测试邮件配置
func (s *ConfigService) TestEmailConfig(ctx context.Context, config *EmailConfig, testEmail string) error {
	emailConfig := &email.EmailConfig{
		Host:     config.Host,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
		From:     config.From,
		FromName: config.FromName,
	}

	emailSvc := email.NewEmailService(emailConfig)

	// 验证配置
	if err := emailSvc.ValidateConfig(); err != nil {
		return fmt.Errorf("邮件配置验证失败: %w", err)
	}

	// 发送测试邮件
	if err := emailSvc.SendTest(testEmail); err != nil {
		return fmt.Errorf("发送测试邮件失败: %w", err)
	}

	return nil
}

// GetStorageManager 获取存储管理器
func (s *ConfigService) GetStorageManager(ctx context.Context) (*storage.Manager, error) {
	if s.storageMgr != nil {
		return s.storageMgr, nil
	}

	config, err := s.GetStorageConfig(ctx)
	if err != nil {
		return nil, err
	}

	storageMgr := storage.NewManager()

	// 根据配置创建存储实例
	storageConfig := &storage.Config{
		Type:      config.Type,
		Endpoint:  config.Endpoint,
		Region:    config.Region,
		AccessID:  config.AccessID,
		AccessKey: config.AccessKey,
		Bucket:    config.Bucket,
		Domain:    config.Domain,
		Extra:     config.Extra,
	}

	storageInstance, err := storage.NewFromConfig(*storageConfig)
	if err != nil {
		return nil, fmt.Errorf("创建存储实例失败: %w", err)
	}

	storageMgr.Register("default", storageInstance)
	storageMgr.SetDefault(storageInstance)

	s.storageMgr = storageMgr
	return s.storageMgr, nil
}

// ClearCache 清除配置缓存
func (s *ConfigService) ClearCache() {
	s.mu.Lock()
	s.cache = make(map[string]string)
	s.mu.Unlock()

	// 清除邮件服务缓存
	s.emailSvc = nil

	// 清除存储管理器缓存
	s.storageMgr = nil
}

// InitDefaultConfigs 初始化默认设置
func (s *ConfigService) InitDefaultConfigs(ctx context.Context) error {
	defaultConfigs := []struct {
		key     string
		name    string
		value   string
		setType string
		group   string
	}{
		{
			key:     "site_name",
			name:    "站点名称",
			value:   "Fun-Admin",
			setType: "string",
			group:   "site",
		},
		{
			key:     "site_description",
			name:    "站点描述",
			value:   "基于 Go + Vue3 的现代化管理后台",
			setType: "string",
			group:   "site",
		},
		{
			key:     "site_keywords",
			name:    "站点关键词",
			value:   "Fun-Admin,Go,Vue3,管理后台",
			setType: "string",
			group:   "site",
		},
		{
			key:     "site_logo",
			name:    "站点Logo",
			value:   "/logo.png",
			setType: "string",
			group:   "site",
		},
		{
			key:     "site_favicon",
			name:    "站点图标",
			value:   "/favicon.ico",
			setType: "string",
			group:   "site",
		},
		{
			key:     "register_enabled",
			name:    "允许注册",
			value:   "false",
			setType: "boolean",
			group:   "system",
		},
		{
			key:     "login_captcha_enabled",
			name:    "登录验证码",
			value:   "true",
			setType: "boolean",
			group:   "system",
		},
		{
			key:     "file_upload_max_size",
			name:    "文件上传最大大小",
			value:   "10485760",
			setType: "number",
			group:   "file",
		},
		{
			key:     "file_allowed_types",
			name:    "允许上传的文件类型",
			value:   ".jpg,.jpeg,.png,.gif,.bmp,.pdf,.doc,.docx,.xls,.xlsx",
			setType: "string",
			group:   "file",
		},
	}

	for _, config := range defaultConfigs {
		newConfig := &model.Config{
			Key:   config.key,
			Name:  config.name,
			Value: config.value,
			Type:  config.setType,
			Group: config.group,
		}

		if err := s.configRepo.CreateConfigIfNotExists(ctx, newConfig); err != nil {
			continue
		}
	}

	return nil
}

// SearchConfigs 搜索系统设置
func (s *ConfigService) SearchConfigs(ctx context.Context, keyword string, group string) ([]model.Config, error) {
	configs, err := s.configRepo.SearchConfigs(ctx, keyword, group)
	if err != nil {
		return nil, fmt.Errorf("搜索系统设置失败: %w", err)
	}

	return configs, nil
}

// GetConfigGroups 获取设置分组列表
func (s *ConfigService) GetConfigGroups(ctx context.Context) ([]string, error) {
	groups, err := s.configRepo.GetConfigGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取设置分组失败: %w", err)
	}

	return groups, nil
}

// DeleteConfig 删除系统设置
func (s *ConfigService) DeleteConfig(ctx context.Context, key string) error {
	// 从数据库删除
	if err := s.configRepo.DeleteConfig(ctx, key); err != nil {
		return fmt.Errorf("删除系统设置失败: %w", err)
	}

	// 从缓存删除
	s.mu.Lock()
	delete(s.cache, key)
	s.mu.Unlock()

	return nil
}
