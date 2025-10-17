package storage

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// OssStorage 阿里云OSS存储实现
type OssStorage struct {
	client *oss.Client
	bucket string
	domain string
}

// NewOssStorage 创建OSS存储实例
func NewOssStorage(endpoint, accessID, accessKey, bucket, domain string) (*OssStorage, error) {
	client, err := oss.New(endpoint, accessID, accessKey)
	if err != nil {
		return nil, fmt.Errorf("创建OSS客户端失败: %w", err)
	}

	return &OssStorage{
		client: client,
		bucket: bucket,
		domain: domain,
	}, nil
}

// Upload 上传文件
func (s *OssStorage) Upload(ctx context.Context, key string, reader io.Reader, contentType string) (*FileInfo, error) {
	bucket, err := s.client.Bucket(s.bucket)
	if err != nil {
		return nil, fmt.Errorf("获取存储桶失败: %w", err)
	}

	// 自动检测内容类型
	if contentType == "" {
		contentType = mime.TypeByExtension(key[strings.LastIndex(key, "."):])
	}

	// 设置上传选项
	options := []oss.Option{
		oss.ContentType(contentType),
	}

	// 上传文件
	err = bucket.PutObject(key, reader, options...)
	if err != nil {
		return nil, fmt.Errorf("上传文件失败: %w", err)
	}

	// 获取文件信息
	props, err := bucket.GetObjectDetailedMeta(key)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	// 解析文件大小
	contentLength, _ := strconv.ParseInt(props.Get("Content-Length"), 10, 64)

	// 解析最后修改时间
	lastModified, err := time.Parse(http.TimeFormat, props.Get("Last-Modified"))
	if err != nil {
		lastModified = time.Now()
	}

	return &FileInfo{
		Key:          key,
		Name:         key[strings.LastIndex(key, "/")+1:],
		Size:         contentLength,
		ContentType:  props.Get("Content-Type"),
		ETag:         strings.Trim(props.Get("ETag"), `"`),
		LastModified: lastModified,
		URL:          s.getURL(key),
	}, nil
}

// Download 下载文件
func (s *OssStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	bucket, err := s.client.Bucket(s.bucket)
	if err != nil {
		return nil, fmt.Errorf("获取存储桶失败: %w", err)
	}

	reader, err := bucket.GetObject(key)
	if err != nil {
		return nil, fmt.Errorf("下载文件失败: %w", err)
	}

	return reader, nil
}

// Delete 删除文件
func (s *OssStorage) Delete(ctx context.Context, key string) error {
	bucket, err := s.client.Bucket(s.bucket)
	if err != nil {
		return fmt.Errorf("获取存储桶失败: %w", err)
	}

	err = bucket.DeleteObject(key)
	if err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	return nil
}

// GetURL 获取文件访问URL
func (s *OssStorage) GetURL(ctx context.Context, key string, expire time.Duration) (string, error) {
	bucket, err := s.client.Bucket(s.bucket)
	if err != nil {
		return "", fmt.Errorf("获取存储桶失败: %w", err)
	}

	// 如果设置了自定义域名，直接返回
	if s.domain != "" {
		return fmt.Sprintf("%s/%s", strings.TrimSuffix(s.domain, "/"), key), nil
	}

	// 生成签名URL
	if expire > 0 {
		signedURL, err := bucket.SignURL(key, oss.HTTPGet, int64(expire/time.Second))
		if err != nil {
			return "", fmt.Errorf("生成签名URL失败: %w", err)
		}
		return signedURL, nil
	}

	// 返回公共URL
	return fmt.Sprintf("https://%s.%s/%s", s.bucket, s.client.Config.Endpoint, key), nil
}

// Exists 检查文件是否存在
func (s *OssStorage) Exists(ctx context.Context, key string) (bool, error) {
	bucket, err := s.client.Bucket(s.bucket)
	if err != nil {
		return false, fmt.Errorf("获取存储桶失败: %w", err)
	}

	_, err = bucket.GetObjectMeta(key)
	if err != nil {
		if IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("检查文件存在失败: %w", err)
	}

	return true, nil
}

// GetSize 获取文件大小
func (s *OssStorage) GetSize(ctx context.Context, key string) (int64, error) {
	bucket, err := s.client.Bucket(s.bucket)
	if err != nil {
		return 0, fmt.Errorf("获取存储桶失败: %w", err)
	}

	props, err := bucket.GetObjectDetailedMeta(key)
	if err != nil {
		return 0, fmt.Errorf("获取文件大小失败: %w", err)
	}

	contentLength, _ := strconv.ParseInt(props.Get("Content-Length"), 10, 64)
	return contentLength, nil
}

// getURL 获取文件访问URL
func (s *OssStorage) getURL(key string) string {
	if s.domain != "" {
		return fmt.Sprintf("%s/%s", strings.TrimSuffix(s.domain, "/"), key)
	}

	return fmt.Sprintf("https://%s.%s/%s", s.bucket, s.client.Config.Endpoint, key)
}

// IsNotExist 检查错误是否表示文件不存在
func IsNotExist(err error) bool {
	if err == nil {
		return false
	}

	// 检查是否为OSS错误
	if ossErr, ok := err.(oss.ServiceError); ok {
		return ossErr.StatusCode == 404
	}

	return false
}
