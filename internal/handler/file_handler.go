package handler

import (
	v1 "fun-admin/api/v1"
	"fun-admin/internal/service"
	"fun-admin/pkg/logger"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// FileHandler 文件处理器
type FileHandler struct {
	fileService *service.FileService
	logger      *logger.Logger
}

// NewFileHandler 创建文件处理器
func NewFileHandler(logger *logger.Logger, fileService *service.FileService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
		logger:      logger,
	}
}

// Upload 上传文件
func (h *FileHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		v1.HandleValidationError(c, "获取上传文件失败")
		return
	}

	allowedTypesStr := c.PostForm("allowed_types")
	var allowedTypes []string
	if allowedTypesStr != "" {
		allowedTypes = strings.Split(allowedTypesStr, ",")
	}

	maxSize := int64(10 * 1024 * 1024)
	if maxSizeStr := c.PostForm("max_size"); maxSizeStr != "" {
		if parsed, parseErr := strconv.ParseInt(maxSizeStr, 10, 64); parseErr == nil {
			maxSize = parsed
		}
	}

	storageType := c.PostForm("storage_type")
	pathPrefix := c.PostForm("path")

	fileInfo, err := h.fileService.UploadFileWithOptions(c, file, allowedTypes, maxSize, storageType, pathPrefix)
	if err != nil {
		v1.HandleError(c, err)
		return
	}

	v1.HandleSuccess(c, fileInfo)
}

// List 返回文件列表
func (h *FileHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	storageType := c.Query("storage_type")

	files, total, err := h.fileService.ListFiles(c, storageType, page, pageSize)
	if err != nil {
		v1.HandleError(c, err)
		return
	}

	v1.HandleSuccess(c, gin.H{
		"list":  files,
		"total": total,
	})
}

// Info 返回文件元信息
func (h *FileHandler) Info(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		v1.HandleValidationError(c, "key 参数不能为空")
		return
	}
	storageType := c.Query("storage_type")

	info, err := h.fileService.GetFileInfo(c, storageType, key)
	if err != nil {
		v1.HandleError(c, err)
		return
	}

	v1.HandleSuccess(c, info)
}

// Delete 删除文件
func (h *FileHandler) Delete(c *gin.Context) {
	var req struct {
		Keys        []string `json:"keys"`
		StorageType string   `json:"storage_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// 兼容通过 query 传递的 key
		key := c.Query("key")
		if key != "" {
			req.Keys = []string{key}
			req.StorageType = c.Query("storage_type")
		}
	}

	if len(req.Keys) == 0 {
		v1.HandleValidationError(c, "缺少要删除的文件 key")
		return
	}

	for _, key := range req.Keys {
		if err := h.fileService.DeleteFileWithContext(c, req.StorageType, key); err != nil {
			v1.HandleError(c, err)
			return
		}
	}

	v1.HandleSuccess(c, gin.H{
		"deleted": len(req.Keys),
	})
}
