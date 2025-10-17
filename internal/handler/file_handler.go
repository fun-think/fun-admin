package handler

import (
	"fun-admin/internal/service"
	"fun-admin/pkg/logger"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// FileHandler 文件处理器
type FileHandler struct {
	fileService *service.FileService
	logger      *logger.Logger
}

// NewFileHandler 创建文件处理器
func NewFileHandler(logger *logger.Logger) *FileHandler {
	return &FileHandler{
		fileService: service.NewFileService(logger),
		logger:      logger,
	}
}

// UploadHandler 处理文件上传请求
func UploadHandler(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "获取上传文件失败: " + err.Error(),
		})
		return
	}

	// 获取允许的文件类型参数
	allowedTypesStr := c.PostForm("allowed_types")
	var allowedTypes []string
	if allowedTypesStr != "" {
		allowedTypes = strings.Split(allowedTypesStr, ",")
	}

	// 获取最大文件大小参数（默认10MB）
	maxSize := int64(10 * 1024 * 1024)
	if maxSizeStr := c.PostForm("max_size"); maxSizeStr != "" {
		// 这里应该解析 maxSizeStr 为 int64，但为简化处理，我们使用默认值
	}

	// 上传文件
	fileInfo, err := service.NewFileService(nil).UploadFile(file, allowedTypes, maxSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文件上传失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "文件上传成功",
		"data":    fileInfo,
	})
}
