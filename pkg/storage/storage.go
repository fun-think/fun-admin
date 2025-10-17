package storage

import (
	"context"
	"io"
	"time"
)

// Storage 云存储接口
type Storage interface {
	// Upload 上传文件
	Upload(ctx context.Context, key string, reader io.Reader, contentType string) (*FileInfo, error)

	// Download 下载文件
	Download(ctx context.Context, key string) (io.ReadCloser, error)

	// Delete 删除文件
	Delete(ctx context.Context, key string) error

	// GetURL 获取文件访问URL
	GetURL(ctx context.Context, key string, expire time.Duration) (string, error)

	// Exists 检查文件是否存在
	Exists(ctx context.Context, key string) (bool, error)

	// GetSize 获取文件大小
	GetSize(ctx context.Context, key string) (int64, error)
}

// FileInfo 文件信息
type FileInfo struct {
	Key          string    `json:"key"`           // 文件键
	Name         string    `json:"name"`          // 文件名
	Size         int64     `json:"size"`          // 文件大小
	ContentType  string    `json:"content_type"`  // 内容类型
	ETag         string    `json:"etag"`          // 文件标识
	LastModified time.Time `json:"last_modified"` // 最后修改时间
	URL          string    `json:"url"`           // 访问URL
}

// Config 存储配置
type Config struct {
	Type      string            `json:"type"`       // 存储类型: local, oss, s3
	Endpoint  string            `json:"endpoint"`   // 访问端点
	Region    string            `json:"region"`     // 区域
	AccessID  string            `json:"access_id"`  // 访问ID
	AccessKey string            `json:"access_key"` // 访问密钥
	Bucket    string            `json:"bucket"`     // 存储桶
	Domain    string            `json:"domain"`     // 自定义域名
	Extra     map[string]string `json:"extra"`      // 额外配置
}
