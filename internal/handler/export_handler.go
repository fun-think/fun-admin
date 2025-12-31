package handler

import (
	"errors"
	"fun-admin/internal/service"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// ExportHandler 导出处理器
type ExportHandler struct {
	*Handler
	resourceService *service.ResourceService
	logger          *zap.Logger
}

// NewExportHandler 创建导出处理器
func NewExportHandler(
	handler *Handler,
	resourceService *service.ResourceService,
	logger *zap.Logger,
) *ExportHandler {
	return &ExportHandler{
		Handler:         handler,
		resourceService: resourceService,
		logger:          logger,
	}
}

// ExportData 导出数据
// @Summary 导出数据
// @Description 导出指定资源的数据为CSV或Excel格式
// @Tags export
// @Produce application/octet-stream
// @Param resource path string true "资源名称"
// @Param format query string false "导出格式 (csv/excel)" default(csv)
// @Param order_by query string false "排序字段"
// @Param order_direction query string false "排序方向 (ASC/DESC)"
// @Success 200 {file} file "导出的数据文件"
// @Router /api/v1/export/{resource} [get]
func (h *ExportHandler) ExportData(c *gin.Context) {
	slug := c.Param("resource")
	if slug == "" {
		h.logger.Warn("导出失败：缺少资源参数")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少资源参数",
		})
		return
	}

	// 获取过滤参数
	filters := make(map[string]interface{})
	search := make(map[string]interface{})

	// 处理查询参数
	for key, values := range c.Request.URL.Query() {
		// 跳过导出格式参数和其他系统参数
		if key == "format" || key == "order_by" || key == "order_direction" {
			continue
		}

		// 如果参数名以 "search_" 开头，则作为搜索条件
		if strings.HasPrefix(key, "search_") {
			fieldName := strings.TrimPrefix(key, "search_")
			if len(values) > 0 && values[0] != "" {
				search[fieldName] = values[0]
			}
		} else {
			// 否则作为过滤条件
			if len(values) > 0 && values[0] != "" {
				filters[key] = values[0]
			}
		}
	}

	// 获取排序参数
	orderBy := c.Query("order_by")
	orderDirection := c.Query("order_direction")
	if orderDirection != "" && orderDirection != "ASC" && orderDirection != "DESC" {
		orderDirection = "DESC" // 默认倒序
	}

	// 获取导出格式参数
	format := c.DefaultQuery("format", "csv")

	h.logger.Info("开始导出数据",
		zap.String("resource", slug),
		zap.String("format", format),
		zap.Int("filters_count", len(filters)),
		zap.Int("search_count", len(search)))

	// 调用服务导出数据
	data, filename, err := h.resourceService.Export(c, slug, filters, search, orderBy, orderDirection, format)
	if err != nil {
		h.logger.Error("导出数据失败",
			zap.String("resource", slug),
			zap.Error(err))

		// 检查是否为资源未找到错误
		var notFoundErr *service.ResourceNotFoundError
		if errors.As(err, &notFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "资源不存在",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": messageWithDebugError("导出失败", err),
		})
		return
	}

	h.logger.Info("导出数据成功",
		zap.String("resource", slug),
		zap.String("filename", filename),
		zap.Int("data_size", len(data)))

	// 设置响应头
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Length", strconv.Itoa(len(data)))

	// 输出文件内容
	c.Data(http.StatusOK, "application/octet-stream", data)
}
