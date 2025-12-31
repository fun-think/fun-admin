package service

import (
	"context"
	"fmt"
	"fun-admin/pkg/logger"
	"fun-admin/pkg/storage"
	"io/fs"
	"mime/multipart"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// FileService 文件服务
type FileService struct {
	logger      *logger.Logger
	storageMgr  *storage.Manager
	defaultDisk string
	diskConfigs map[string]storage.Config
}

// NewFileService 创建文件服务
func NewFileService(logger *logger.Logger, conf *viper.Viper) *FileService {
	manager := storage.NewManager()
	defaultDisk := "local"
	if conf != nil && conf.GetString("storage.default") != "" {
		defaultDisk = conf.GetString("storage.default")
	}

	rawDisks := map[string]interface{}{}
	if conf != nil {
		rawDisks = conf.GetStringMap("storage.disks")
	}
	if len(rawDisks) == 0 {
		rawDisks = map[string]interface{}{
			"local": map[string]interface{}{
				"driver":    "local",
				"base_path": "storage/uploads",
				"domain":    "",
			},
		}
	}

	diskConfigs := make(map[string]storage.Config)

	for name, raw := range rawDisks {
		cfg := parseDiskConfig(raw)
		if cfg.Type == "" {
			cfg.Type = "local"
		}
		diskConfigs[name] = cfg

		var (
			store storage.Storage
			err   error
		)

		switch cfg.Type {
		case "oss":
			store, err = storage.NewOssStorage(cfg.Endpoint, cfg.AccessID, cfg.AccessKey, cfg.Bucket, cfg.Domain)
		default:
			basePath := cfg.Extra["base_path"]
			if basePath == "" {
				basePath = "storage/uploads"
			}
			store, err = storage.NewLocalStorage(basePath, cfg.Domain)
		}

		if err != nil {
			logger.Error("注册存储失败", zap.Error(err))
			continue
		}

		manager.Register(name, store)
		if name == defaultDisk || manager.Default() == nil {
			manager.SetDefault(store)
		}
	}

	if manager.Default() == nil {
		for name := range diskConfigs {
			if store, err := manager.Get(name); err == nil {
				manager.SetDefault(store)
				defaultDisk = name
				break
			}
		}
	}

	return &FileService{
		logger:      logger,
		storageMgr:  manager,
		defaultDisk: defaultDisk,
		diskConfigs: diskConfigs,
	}
}

// UploadFile 上传文件（使用默认存储）
func (s *FileService) UploadFile(file *multipart.FileHeader, allowedTypes []string, maxSize int64) (*FileInfo, error) {
	return s.UploadFileWithOptions(context.Background(), file, allowedTypes, maxSize, "", "")
}

// UploadFileWithContext 带上下文上传（默认存储）
func (s *FileService) UploadFileWithContext(ctx context.Context, file *multipart.FileHeader, allowedTypes []string, maxSize int64) (*FileInfo, error) {
	return s.UploadFileWithOptions(ctx, file, allowedTypes, maxSize, "", "")
}

// UploadFileWithOptions 上传并指定存储类型与路径前缀
func (s *FileService) UploadFileWithOptions(ctx context.Context, file *multipart.FileHeader, allowedTypes []string, maxSize int64, storageType, pathPrefix string) (*FileInfo, error) {
	if file == nil {
		return nil, fmt.Errorf("文件不能为空")
	}
	if maxSize > 0 && file.Size > maxSize {
		return nil, fmt.Errorf("文件大小超过限制，最大允许 %d 字节", maxSize)
	}
	if len(allowedTypes) > 0 {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		valid := false
		for _, allow := range allowedTypes {
			if strings.ToLower(allow) == ext {
				valid = true
				break
			}
		}
		if !valid {
			return nil, fmt.Errorf("不支持的文件类型，仅允许: %s", strings.Join(allowedTypes, ", "))
		}
	}

	store, diskName, err := s.resolveStorage(storageType)
	if err != nil {
		return nil, err
	}

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("读取上传文件失败: %w", err)
	}
	defer src.Close()

	key := s.buildFileKey(pathPrefix, file.Filename)
	fileInfo, err := store.Upload(ctx, key, src, file.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("上传文件失败: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	return &FileInfo{
		Name:        file.Filename,
		Path:        key,
		URL:         fileInfo.URL,
		Size:        fileInfo.Size,
		Ext:         ext,
		StorageType: diskName,
		ContentType: fileInfo.ContentType,
		ETag:        fileInfo.ETag,
	}, nil
}

// DeleteFile 删除默认存储的文件
func (s *FileService) DeleteFile(fileKey string) error {
	return s.DeleteFileWithContext(context.Background(), "", fileKey)
}

// DeleteFileWithContext 删除指定存储的文件
func (s *FileService) DeleteFileWithContext(ctx context.Context, storageType, fileKey string) error {
	if fileKey == "" {
		return fmt.Errorf("文件标识不能为空")
	}
	store, _, err := s.resolveStorage(storageType)
	if err != nil {
		return err
	}
	return store.Delete(ctx, fileKey)
}

// GetFileURL 获取默认存储的访问地址
func (s *FileService) GetFileURL(fileKey string, expire time.Duration) (string, error) {
	return s.GetFileURLWithContext(context.Background(), "", fileKey, expire)
}

// GetFileURLWithContext 获取指定存储的访问地址
func (s *FileService) GetFileURLWithContext(ctx context.Context, storageType, fileKey string, expire time.Duration) (string, error) {
	store, _, err := s.resolveStorage(storageType)
	if err != nil {
		return "", err
	}
	return store.GetURL(ctx, fileKey, expire)
}

// FileExists 检查默认存储文件是否存在
func (s *FileService) FileExists(fileKey string) (bool, error) {
	return s.FileExistsWithContext(context.Background(), "", fileKey)
}

// FileExistsWithContext 检查指定存储文件是否存在
func (s *FileService) FileExistsWithContext(ctx context.Context, storageType, fileKey string) (bool, error) {
	store, _, err := s.resolveStorage(storageType)
	if err != nil {
		return false, err
	}
	return store.Exists(ctx, fileKey)
}

// GetFileSize 获取默认存储文件大小
func (s *FileService) GetFileSize(fileKey string) (int64, error) {
	return s.GetFileSizeWithContext(context.Background(), "", fileKey)
}

// GetFileSizeWithContext 获取指定存储文件大小
func (s *FileService) GetFileSizeWithContext(ctx context.Context, storageType, fileKey string) (int64, error) {
	store, _, err := s.resolveStorage(storageType)
	if err != nil {
		return 0, err
	}
	return store.GetSize(ctx, fileKey)
}

// ListFiles 列出指定存储的文件（仅支持本地存储）
func (s *FileService) ListFiles(ctx context.Context, storageType string, page, pageSize int) ([]*FileInfo, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	diskName := storageType
	if diskName == "" {
		diskName = s.defaultDisk
	}
	cfg, ok := s.diskConfigs[diskName]
	if !ok || cfg.Type != "local" {
		return nil, 0, fmt.Errorf("当前存储不支持列出文件")
	}
	basePath := cfg.Extra["base_path"]
	if basePath == "" {
		basePath = "storage/uploads"
	}

	entries := make([]fileEntry, 0)
	err := filepath.WalkDir(basePath, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(basePath, path)
		if err != nil {
			return err
		}
		entries = append(entries, fileEntry{
			path:    filepath.ToSlash(rel),
			size:    info.Size(),
			modTime: info.ModTime(),
		})
		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].modTime.After(entries[j].modTime)
	})

	total := int64(len(entries))
	start := (page - 1) * pageSize
	if start > len(entries) {
		return []*FileInfo{}, total, nil
	}
	end := start + pageSize
	if end > len(entries) {
		end = len(entries)
	}

	store, _, err := s.resolveStorage(diskName)
	if err != nil {
		return nil, 0, err
	}

	list := make([]*FileInfo, 0, end-start)
	for _, entry := range entries[start:end] {
		url, _ := store.GetURL(ctx, entry.path, 0)
		list = append(list, &FileInfo{
			Name:        filepath.Base(entry.path),
			Path:        entry.path,
			URL:         url,
			Size:        entry.size,
			Ext:         filepath.Ext(entry.path),
			StorageType: diskName,
		})
	}

	return list, total, nil
}

// GetFileInfo 获取文件元信息
func (s *FileService) GetFileInfo(ctx context.Context, storageType, fileKey string) (*FileInfo, error) {
	if fileKey == "" {
		return nil, fmt.Errorf("文件标识不能为空")
	}
	store, diskName, err := s.resolveStorage(storageType)
	if err != nil {
		return nil, err
	}
	size, err := store.GetSize(ctx, fileKey)
	if err != nil {
		return nil, err
	}
	url, _ := store.GetURL(ctx, fileKey, 0)
	return &FileInfo{
		Name:        filepath.Base(fileKey),
		Path:        fileKey,
		URL:         url,
		Size:        size,
		Ext:         filepath.Ext(fileKey),
		StorageType: diskName,
	}, nil
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
	for _, allow := range allowedTypes {
		if strings.ToLower(allow) == ext {
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

// FileInfo 文件信息
type FileInfo struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	URL         string `json:"url"`
	Size        int64  `json:"size"`
	Ext         string `json:"ext"`
	StorageType string `json:"storage_type"`
	ContentType string `json:"content_type"`
	ETag        string `json:"etag"`
}

type fileEntry struct {
	path    string
	size    int64
	modTime time.Time
}

func (s *FileService) resolveStorage(storageType string) (storage.Storage, string, error) {
	target := storageType
	if target == "" {
		target = s.defaultDisk
	}
	if target == "" {
		return nil, "", fmt.Errorf("未配置可用的文件存储")
	}
	store, err := s.storageMgr.Get(target)
	if err != nil {
		if s.storageMgr.Default() != nil && target != s.defaultDisk {
			return s.storageMgr.Default(), s.defaultDisk, nil
		}
		return nil, "", fmt.Errorf("存储 %s 不存在", target)
	}
	return store, target, nil
}

func (s *FileService) buildFileKey(prefix, filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".dat"
	}
	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	base := time.Now().Format("2006/01/02")
	if prefix != "" {
		base = strings.Trim(prefix, "/")
	}
	return filepath.ToSlash(filepath.Join(base, name))
}

func parseDiskConfig(raw interface{}) storage.Config {
	cfg := storage.Config{
		Extra: map[string]string{},
	}
	data, _ := raw.(map[string]interface{})
	cfg.Type = stringify(data["driver"])
	cfg.Endpoint = stringify(data["endpoint"])
	cfg.Region = stringify(data["region"])
	cfg.AccessID = stringify(data["access_id"])
	cfg.AccessKey = stringify(data["access_key"])
	cfg.Bucket = stringify(data["bucket"])
	cfg.Domain = stringify(data["domain"])
	if base := stringify(data["base_path"]); base != "" {
		cfg.Extra["base_path"] = base
	}
	return cfg
}

func stringify(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case fmt.Stringer:
		return val.String()
	default:
		return fmt.Sprintf("%v", val)
	}
}
