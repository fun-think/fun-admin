package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"
)

// Manager 存储管理器
type Manager struct {
	defaultStorage Storage
	storages       map[string]Storage
}

// NewManager 创建存储管理器
func NewManager() *Manager {
	return &Manager{
		storages: make(map[string]Storage),
	}
}

// Register 注册存储实例
func (m *Manager) Register(name string, storage Storage) {
	m.storages[name] = storage
}

// SetDefault 设置默认存储
func (m *Manager) SetDefault(storage Storage) {
	m.defaultStorage = storage
}

// Get 获取存储实例
func (m *Manager) Get(name string) (Storage, error) {
	if storage, exists := m.storages[name]; exists {
		return storage, nil
	}
	return nil, fmt.Errorf("存储实例 '%s' 不存在", name)
}

// Default 获取默认存储实例
func (m *Manager) Default() Storage {
	if m.defaultStorage == nil {
		// 如果没有设置默认存储，返回第一个注册的存储
		for _, storage := range m.storages {
			m.defaultStorage = storage
			break
		}
	}
	return m.defaultStorage
}

// NewFromConfig 根据配置创建存储实例
func NewFromConfig(config Config) (Storage, error) {
	switch config.Type {
	case "local":
		return NewLocalStorage(config.Extra["base_path"], config.Domain)
	case "oss":
		return NewOssStorage(config.Endpoint, config.AccessID, config.AccessKey, config.Bucket, config.Domain)
	// S3存储暂未实现
	// case "s3":
	// 	return NewS3Storage(config.Endpoint, config.Region, config.AccessID, config.AccessKey, config.Bucket, config.Domain)
	default:
		return nil, fmt.Errorf("不支持的存储类型: %s", config.Type)
	}
}

// Upload 上传文件（使用默认存储）
func (m *Manager) Upload(ctx context.Context, key string, reader interface{}, contentType string) (*FileInfo, error) {
	storage := m.Default()
	if storage == nil {
		return nil, fmt.Errorf("没有可用的存储实例")
	}

	var r io.Reader
	switch v := reader.(type) {
	case io.Reader:
		r = v
	case []byte:
		r = bytes.NewReader(v)
	default:
		return nil, fmt.Errorf("不支持的reader类型")
	}

	return storage.Upload(ctx, key, r, contentType)
}

// Download 下载文件（使用默认存储）
func (m *Manager) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	storage := m.Default()
	if storage == nil {
		return nil, fmt.Errorf("没有可用的存储实例")
	}

	return storage.Download(ctx, key)
}

// Delete 删除文件（使用默认存储）
func (m *Manager) Delete(ctx context.Context, key string) error {
	storage := m.Default()
	if storage == nil {
		return fmt.Errorf("没有可用的存储实例")
	}

	return storage.Delete(ctx, key)
}

// GetURL 获取文件访问URL（使用默认存储）
func (m *Manager) GetURL(ctx context.Context, key string, expire time.Duration) (string, error) {
	storage := m.Default()
	if storage == nil {
		return "", fmt.Errorf("没有可用的存储实例")
	}

	return storage.GetURL(ctx, key, expire)
}

// Exists 检查文件是否存在（使用默认存储）
func (m *Manager) Exists(ctx context.Context, key string) (bool, error) {
	storage := m.Default()
	if storage == nil {
		return false, fmt.Errorf("没有可用的存储实例")
	}

	return storage.Exists(ctx, key)
}

// GetSize 获取文件大小（使用默认存储）
func (m *Manager) GetSize(ctx context.Context, key string) (int64, error) {
	storage := m.Default()
	if storage == nil {
		return 0, fmt.Errorf("没有可用的存储实例")
	}

	return storage.GetSize(ctx, key)
}
