package service

import (
	"context"
	"fmt"
	"fun-admin/pkg/logger"
	"fun-admin/pkg/storage"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileService 文件服务
type FileService struct {
	logger     *logger.Logger
	storageMgr *storage.Manager
	uploadPath string
}

// NewFileService 创建文件服务
func NewFileService(logger *logger.Logger) *FileService {
	// 确保上传目录存在
	uploadPath := "storage/uploads"
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		logger.Error("创建上传目录失败: " + err.Error())
	}

	// 创建存储管理器
	storageMgr := storage.NewManager()

	// 注册本地存储
	localStorage := storage.NewLocalStorage(uploadPath, "")
	storageMgr.Register("local", localStorage)
	storageMgr.SetDefault(localStorage)

	return &FileService{
		logger:     logger,
		storageMgr: storageMgr,
		uploadPath: uploadPath,
	}
}

// UploadFile 上传文件
func (s *FileService) UploadFile(file *multipart.FileHeader, allowedTypes []string, maxSize int64) (*FileInfo, error) {
	return s.UploadFileWithContext(context.Background(), file, allowedTypes, maxSize)
}

// UploadFileWithContext 带上下文的文件上传
func (s *FileService) UploadFileWithContext(ctx context.Context, file *multipart.FileHeader, allowedTypes []string, maxSize int64) (*FileInfo, error) {
	// 检查文件大小
	if file.Size > maxSize {
		return nil, fmt.Errorf("文件大小超过限制，最大允许 %d 字节", maxSize)
	}

	// 检查文件类型
	if len(allowedTypes) > 0 {
		fileExt := strings.ToLower(filepath.Ext(file.Filename))
		allowed := false
		for _, ext := range allowedTypes {
			if strings.ToLower(ext) == fileExt {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, fmt.Errorf("不支持的文件类型，仅支持: %s", strings.Join(allowedTypes, ", "))
		}
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开上传文件失败: %w", err)
	}
	defer src.Close()

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().Unix(), ext)

	// 按日期组织文件路径
	datePath := time.Now().Format("2006/01/02")
	fileKey := filepath.Join(datePath, filename)

	// 使用存储管理器上传文件
	storageInfo, err := s.storageMgr.Upload(ctx, fileKey, src, "")
	if err != nil {
		return nil, fmt.Errorf("上传文件失败: %w", err)
	}

	// 转换为文件服务返回的格式
	fileInfo := &FileInfo{
		Name:        filename,
		Path:        fileKey,
		URL:         storageInfo.URL,
		Size:        storageInfo.Size,
		Ext:         ext,
		StorageType: "local", // 当前默认使用本地存储
	}

	return fileInfo, nil
}

// FileInfo 文件信息
type FileInfo struct {
	Name        string `json:"name"`         // 文件名
	Path        string `json:"path"`         // 文件路径
	URL         string `json:"url"`          // 访问URL
	Size        int64  `json:"size"`         // 文件大小
	Ext         string `json:"ext"`          // 文件扩展名
	StorageType string `json:"storage_type"` // 存储类型
	ContentType string `json:"content_type"` // 内容类型
	ETag        string `json:"etag"`         // 文件标识
}

// DeleteFile 删除文件
func (s *FileService) DeleteFile(fileKey string) error {
	return s.DeleteFileWithContext(context.Background(), fileKey)
}

// DeleteFileWithContext 带上下文的文件删除
func (s *FileService) DeleteFileWithContext(ctx context.Context, fileKey string) error {
	return s.storageMgr.Delete(ctx, fileKey)
}

// GetFileURL 获取文件访问URL
func (s *FileService) GetFileURL(fileKey string, expire time.Duration) (string, error) {
	return s.GetFileURLWithContext(context.Background(), fileKey, expire)
}

// GetFileURLWithContext 带上下文的获取文件URL
func (s *FileService) GetFileURLWithContext(ctx context.Context, fileKey string, expire time.Duration) (string, error) {
	return s.storageMgr.GetURL(ctx, fileKey, expire)
}

// FileExists 检查文件是否存在
func (s *FileService) FileExists(fileKey string) (bool, error) {
	return s.FileExistsWithContext(context.Background(), fileKey)
}

// FileExistsWithContext 带上下文的检查文件存在
func (s *FileService) FileExistsWithContext(ctx context.Context, fileKey string) (bool, error) {
	return s.storageMgr.Exists(ctx, fileKey)
}

// GetFileSize 获取文件大小
func (s *FileService) GetFileSize(fileKey string) (int64, error) {
	return s.GetFileSizeWithContext(context.Background(), fileKey)
}

// GetFileSizeWithContext 带上下文的获取文件大小
func (s *FileService) GetFileSizeWithContext(ctx context.Context, fileKey string) (int64, error) {
	return s.storageMgr.GetSize(ctx, fileKey)
}

// GetFileCategory 获取文件分类
func (s *FileService) GetFileCategory(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp":
		return "image"
	case ".mp4", ".avi", ".mov", ".wmv", ".flv", ".mkv":
		return "video"
	case ".mp3", ".wav", ".flac", ".aac", ".ogg":
		return "audio"
	case ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx":
		return "document"
	case ".zip", ".rar", ".7z", ".tar", ".gz":
		return "archive"
	default:
		return "other"
	}
}

// ValidateFileType 验证文件类型
func (s *FileService) ValidateFileType(filename string, allowedTypes []string) bool {
	if len(allowedTypes) == 0 {
		return true
	}

	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedType := range allowedTypes {
		if strings.ToLower(allowedType) == ext {
			return true
		}
	}
	return false
}

// FormatFileSize 格式化文件大小
func (s *FileService) FormatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(size)/float64(div), "KMGTPE"[exp])
}
