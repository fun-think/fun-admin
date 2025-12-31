package handler

import (
	"net/http"
	"strconv"

	"fun-admin/internal/service"
	"fun-admin/pkg/admin/i18n"

	"github.com/gin-gonic/gin"
)

// OperationLogHandler 操作日志处理器
type OperationLogHandler struct {
	*Handler
	operationLogService service.OperationLogServiceInterface
}

// NewOperationLogHandler 创建处理器
func NewOperationLogHandler(
	handler *Handler,
	operationLogService service.OperationLogServiceInterface,
) *OperationLogHandler {
	return &OperationLogHandler{
		Handler:             handler,
		operationLogService: operationLogService,
	}
}

// GetOperationLogs 操作日志列表
func (h *OperationLogHandler) GetOperationLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	language := getLanguage(c)

	filters := make(map[string]interface{})
	if keyword := c.Query("search_keyword"); keyword != "" {
		filters["keyword"] = keyword
	}
	if userID := c.Query("user_id"); userID != "" {
		filters["user_id"] = userID
	}
	if resource := c.Query("resource"); resource != "" {
		filters["resource"] = resource
	}
	if action := c.Query("action"); action != "" {
		filters["action"] = action
	}
	if method := c.Query("method"); method != "" {
		filters["method"] = method
	}
	if path := c.Query("path"); path != "" {
		filters["path"] = path
	}
	if ip := c.Query("ip"); ip != "" {
		filters["ip"] = ip
	}
	if statusCode := c.Query("status_code"); statusCode != "" {
		filters["status_code"] = statusCode
	}
	if startTime := c.Query("start_time"); startTime != "" {
		filters["created_at_from"] = startTime
	}
	if endTime := c.Query("end_time"); endTime != "" {
		filters["created_at_to"] = endTime
	}

	logs, total, err := h.operationLogService.GetOperationLogs(c, page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": messageWithDebugError(i18n.Translate(language, "error.failed_to_get_data"), err),
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

// GetOperationLog 获取详情
func (h *OperationLogHandler) GetOperationLog(c *gin.Context) {
	language := getLanguage(c)

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

// DeleteOperationLog 删除日志
func (h *OperationLogHandler) DeleteOperationLog(c *gin.Context) {
	language := getLanguage(c)

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
			"message": messageWithDebugError(i18n.Translate(language, "error.failed_to_delete_log"), err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": i18n.Translate(language, "message.deleted_successfully"),
	})
}

// BatchDeleteOperationLogs 批量删除
func (h *OperationLogHandler) BatchDeleteOperationLogs(c *gin.Context) {
	language := getLanguage(c)

	var req struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": i18n.Translate(language, "error.invalid_request_data"),
		})
		return
	}

	if err := h.operationLogService.DeleteOperationLogs(c, req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": messageWithDebugError(i18n.Translate(language, "error.failed_to_batch_delete_logs"), err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": i18n.Translate(language, "message.batch_deleted_successfully"),
	})
}

// GetOperationLogStats 日志统计
func (h *OperationLogHandler) GetOperationLogStats(c *gin.Context) {
	language := getLanguage(c)

	stats, err := h.operationLogService.GetOperationLogStats(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": messageWithDebugError(i18n.Translate(language, "error.failed_to_get_data"), err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    stats,
		"message": "success",
	})
}

// ClearOperationLogs 清空日志
func (h *OperationLogHandler) ClearOperationLogs(c *gin.Context) {
	language := getLanguage(c)

	if err := h.operationLogService.ClearOperationLogs(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": messageWithDebugError(i18n.Translate(language, "error.failed_to_delete_log"), err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": i18n.Translate(language, "message.deleted_successfully"),
	})
}
