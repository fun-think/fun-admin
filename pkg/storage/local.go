package storage

import (
	"context"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LocalStorage 本地存储实现
type LocalStorage struct {
	basePath string
	domain   string
}

// NewLocalStorage 创建本地存储实例
func NewLocalStorage(basePath, domain string) (*LocalStorage, error) {
	// 确保基础路径存在
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("创建本地存储目录失败: %w", err)
	}

	return &LocalStorage{
		basePath: basePath,
		domain:   domain,
	}, nil
}

// Upload 上传文件
func (s *LocalStorage) Upload(ctx context.Context, key string, reader io.Reader, contentType string) (*FileInfo, error) {
	// 构建文件路径
	filePath := filepath.Join(s.basePath, key)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %w", err)
	}

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	// 复制文件内容
	size, err := io.Copy(file, reader)
	if err != nil {
		return nil, fmt.Errorf("写入文件失败: %w", err)
	}

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	// 自动检测内容类型
	if contentType == "" {
		contentType = mime.TypeByExtension(filepath.Ext(key))
	}

	return &FileInfo{
		Key:          key,
		Name:         filepath.Base(key),
		Size:         size,
		ContentType:  contentType,
		ETag:         fmt.Sprintf("%x", fileInfo.ModTime().UnixNano()),
		LastModified: fileInfo.ModTime(),
		URL:          s.getURL(key),
	}, nil
}

// Download 下载文件
func (s *LocalStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	filePath := filepath.Join(s.basePath, key)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}

	return file, nil
}

// Delete 删除文件
func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	filePath := filepath.Join(s.basePath, key)

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	return nil
}

// GetURL 获取文件访问URL
func (s *LocalStorage) GetURL(ctx context.Context, key string, expire time.Duration) (string, error) {
	return s.getURL(key), nil
}

// Exists 检查文件是否存在
func (s *LocalStorage) Exists(ctx context.Context, key string) (bool, error) {
	filePath := filepath.Join(s.basePath, key)

	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("检查文件存在失败: %w", err)
	}

	return true, nil
}

// GetSize 获取文件大小
func (s *LocalStorage) GetSize(ctx context.Context, key string) (int64, error) {
	filePath := filepath.Join(s.basePath, key)

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, fmt.Errorf("获取文件大小失败: %w", err)
	}

	return fileInfo.Size(), nil
}

// getURL 获取文件访问URL
func (s *LocalStorage) getURL(key string) string {
	// 处理Windows路径
	key = filepath.ToSlash(key)

	if s.domain != "" {
		// 移除域名末尾的斜杠
		domain := strings.TrimSuffix(s.domain, "/")
		return fmt.Sprintf("%s/%s", domain, key)
	}

	// 返回相对路径
	return fmt.Sprintf("/uploads/%s", key)
}
