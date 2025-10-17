package handler

import (
	"fun-admin/internal/service"
	"fun-admin/pkg/admin/i18n"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// OperationLogHandler 操作日志处理器
type OperationLogHandler struct {
	*Handler
	operationLogService service.OperationLogServiceInterface
}

// NewOperationLogHandler 创建操作日志处理器
func NewOperationLogHandler(
	handler *Handler,
	operationLogService service.OperationLogServiceInterface,
) *OperationLogHandler {
	return &OperationLogHandler{
		Handler:             handler,
		operationLogService: operationLogService,
	}
}

// GetOperationLogs 获取操作日志列表
func (h *OperationLogHandler) GetOperationLogs(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 获取语言参数
	language := c.Query("language")
	if language == "" {
		language = "zh-CN"
	}

	// 构建过滤条件
	filters := make(map[string]interface{})

	// 添加搜索条件
	if keyword := c.Query("search_keyword"); keyword != "" {
		filters["keyword"] = keyword
	}

	// 添加用户筛选
	if userID := c.Query("user_id"); userID != "" {
		filters["user_id"] = userID
	}

	// 添加资源筛选
	if resource := c.Query("resource"); resource != "" {
		filters["resource"] = resource
	}

	// 添加操作类型筛选
	if action := c.Query("action"); action != "" {
		filters["action"] = action
	}

	// 添加方法筛选
	if method := c.Query("method"); method != "" {
		filters["method"] = method
	}

	// 添加路径筛选
	if path := c.Query("path"); path != "" {
		filters["path"] = path
	}

	// 添加IP筛选
	if ip := c.Query("ip"); ip != "" {
		filters["ip"] = ip
	}

	// 添加状态码筛选
	if statusCode := c.Query("status_code"); statusCode != "" {
		filters["status_code"] = statusCode
	}

	// 获取操作日志列表
	logs, total, err := h.operationLogService.GetOperationLogs(c, page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": i18n.Translate(language, "error.failed_to_get_data") + ": " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"items":     logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
		"message": "success",
	})
}

// GetOperationLog 获取单个操作日志详情
func (h *OperationLogHandler) GetOperationLog(c *gin.Context) {
	// 获取语言参数
	language := c.Query("language")
	if language == "" {
		language = "zh-CN"
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": i18n.Translate(language, "error.invalid_id"),
		})
		return
	}

	log, err := h.operationLogService.GetOperationLog(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": i18n.Translate(language, "error.log_not_found"),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    log,
		"message": "success",
	})
}

// DeleteOperationLog 删除操作日志
func (h *OperationLogHandler) DeleteOperationLog(c *gin.Context) {
	// 获取语言参数
	language := c.Query("language")
	if language == "" {
		language = "zh-CN"
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": i18n.Translate(language, "error.invalid_id"),
		})
		return
	}

	if err := h.operationLogService.DeleteOperationLog(c, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": i18n.Translate(language, "error.failed_to_delete_log") + ": " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": i18n.Translate(language, "message.deleted_successfully"),
	})
}

// BatchDeleteOperationLogs 批量删除操作日志
func (h *OperationLogHandler) BatchDeleteOperationLogs(c *gin.Context) {
	// 获取语言参数
	language := c.Query("language")
	if language == "" {
		language = "zh-CN"
	}

	var req struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": i18n.Translate(language, "error.invalid_request_data"),
		})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": i18n.Translate(language, "error.no_log_ids_provided"),
		})
		return
	}

	if err := h.operationLogService.DeleteOperationLogs(c, req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": i18n.Translate(language, "error.failed_to_batch_delete_logs") + ": " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": i18n.Translate(language, "message.batch_deleted_successfully"),
	})
}

// GetOperationLogStats 获取操作日志统计信息
func (h *OperationLogHandler) GetOperationLogStats(c *gin.Context) {
	// 获取语言参数
	language := c.Query("language")
	if language == "" {
		language = "zh-CN"
	}

	c.JSON(http.StatusNotImplemented, gin.H{
		"code":    501,
		"message": i18n.Translate(language, "error.not_implemented"),
	})
}
