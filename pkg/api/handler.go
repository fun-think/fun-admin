package api

import (
	"context"
	"fun-admin/api/v1"
	"fun-admin/pkg/errors"
	"fun-admin/pkg/logger"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// BaseHandler API基础处理器
type BaseHandler struct {
	logger    *logger.Logger
	validator *validator.Validate
}

// NewBaseHandler 创建基础处理器
func NewBaseHandler(logger *logger.Logger) *BaseHandler {
	return &BaseHandler{
		logger:    logger,
		validator: validator.New(),
	}
}

// Success 返回成功响应
func (h *BaseHandler) Success(c *gin.Context, data interface{}) {
	v1.HandleSuccess(c, data)
}

// Error 返回错误响应
func (h *BaseHandler) Error(c *gin.Context, err error) {
	v1.HandleError(c, err)
}

// BadRequest 返回请求参数错误
func (h *BaseHandler) BadRequest(c *gin.Context, message string) {
	err := errors.New(errors.CodeBadRequest, message)
	h.Error(c, err)
}

// ValidationError 返回验证错误
func (h *BaseHandler) ValidationError(c *gin.Context, err error) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		message := h.formatValidationError(validationErrors)
		v1.HandleValidationError(c, message)
		return
	}

	wrappedErr := errors.Wrap(err, errors.CodeValidationFailed, "数据验证失败")
	h.Error(c, wrappedErr)
}

// Unauthorized 返回未授权错误
func (h *BaseHandler) Unauthorized(c *gin.Context) {
	v1.HandleUnauthorized(c)
}

// Forbidden 返回权限不足错误
func (h *BaseHandler) Forbidden(c *gin.Context) {
	v1.HandleForbidden(c)
}

// NotFound 返回资源不存在错误
func (h *BaseHandler) NotFound(c *gin.Context) {
	v1.HandleNotFound(c)
}

// InternalError 返回服务器内部错误
func (h *BaseHandler) InternalError(c *gin.Context, err error) {
	h.logger.Error("Internal server error",
		zap.Error(err),
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
	)
	v1.HandleServerError(c, err)
}

// BindJSON 绑定JSON数据并验证
func (h *BaseHandler) BindJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return err
	}

	if err := h.validator.Struct(obj); err != nil {
		return err
	}

	return nil
}

// BindQuery 绑定查询参数并验证
func (h *BaseHandler) BindQuery(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return err
	}

	if err := h.validator.Struct(obj); err != nil {
		return err
	}

	return nil
}

// BindURI 绑定URI参数并验证
func (h *BaseHandler) BindURI(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindUri(obj); err != nil {
		return err
	}

	if err := h.validator.Struct(obj); err != nil {
		return err
	}

	return nil
}

// GetUserID 获取当前用户ID
func (h *BaseHandler) GetUserID(c *gin.Context) (uint, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, errors.ErrUnauthorized
	}

	if id, ok := userID.(uint); ok {
		return id, nil
	}

	return 0, errors.ErrUnauthorized
}

// GetUserIDParam 从URL参数获取用户ID
func (h *BaseHandler) GetUserIDParam(c *gin.Context) (uint, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, errors.New(errors.CodeBadRequest, "无效的用户ID")
	}
	return uint(id), nil
}

// GetPageParam 获取分页参数
func (h *BaseHandler) GetPageParam(c *gin.Context) (page, pageSize int) {
	page = 1
	pageSize = 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	return page, pageSize
}

// GetSortParam 获取排序参数
func (h *BaseHandler) GetSortParam(c *gin.Context, defaultSort string) (string, string) {
	sortBy := c.Query("sort_by")
	if sortBy == "" {
		sortBy = defaultSort
	}

	order := c.Query("order")
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	return sortBy, order
}

// GetSearchParam 获取搜索参数
func (h *BaseHandler) GetSearchParam(c *gin.Context) string {
	return c.Query("search")
}

// WithContext 创建带有用户信息的上下文
func (h *BaseHandler) WithContext(c *gin.Context) context.Context {
	ctx := c.Request.Context()

	// 添加用户ID到上下文
	if userID, exists := c.Get("user_id"); exists {
		ctx = context.WithValue(ctx, "user_id", userID)
	}

	// 添加追踪ID到上下文
	if traceID, exists := c.Get("trace_id"); exists {
		ctx = context.WithValue(ctx, "trace_id", traceID)
	}

	return ctx
}

// formatValidationError 格式化验证错误消息
func (h *BaseHandler) formatValidationError(errs validator.ValidationErrors) string {
	messages := make([]string, 0, len(errs))

	for _, err := range errs {
		message := h.getFieldErrorMessage(err)
		messages = append(messages, message)
	}

	if len(messages) == 0 {
		return "数据验证失败"
	}

	if len(messages) == 1 {
		return messages[0]
	}

	result := messages[0]
	for i := 1; i < len(messages); i++ {
		result += "; " + messages[i]
	}

	return result
}

// getFieldErrorMessage 获取字段验证错误消息
func (h *BaseHandler) getFieldErrorMessage(err validator.FieldError) string {
	field := err.Field()

	switch err.Tag() {
	case "required":
		return field + " 是必填项"
	case "email":
		return field + " 必须是有效的邮箱地址"
	case "min":
		return field + " 长度不能少于 " + err.Param() + " 个字符"
	case "max":
		return field + " 长度不能超过 " + err.Param() + " 个字符"
	case "len":
		return field + " 长度必须为 " + err.Param() + " 个字符"
	case "gte":
		return field + " 必须大于等于 " + err.Param()
	case "lte":
		return field + " 必须小于等于 " + err.Param()
	case "gt":
		return field + " 必须大于 " + err.Param()
	case "lt":
		return field + " 必须小于 " + err.Param()
	case "oneof":
		return field + " 必须是以下值之一: " + err.Param()
	case "unique":
		return field + " 已经存在"
	default:
		return field + " 格式不正确"
	}
}

// PaginationResponse 分页响应结构
type PaginationResponse struct {
	Data       interface{}            `json:"data"`
	Pagination PaginationMeta         `json:"pagination"`
	Meta       map[string]interface{} `json:"meta,omitempty"`
}

// PaginationMeta 分页元信息
type PaginationMeta struct {
	Page        int   `json:"page"`
	PageSize    int   `json:"page_size"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
	HasPrevious bool  `json:"has_previous"`
	HasNext     bool  `json:"has_next"`
}

// NewPaginationResponse 创建分页响应
func (h *BaseHandler) NewPaginationResponse(data interface{}, page, pageSize int, total int64, meta ...map[string]interface{}) *PaginationResponse {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	if totalPages == 0 {
		totalPages = 1
	}

	response := &PaginationResponse{
		Data: data,
		Pagination: PaginationMeta{
			Page:        page,
			PageSize:    pageSize,
			Total:       total,
			TotalPages:  totalPages,
			HasPrevious: page > 1,
			HasNext:     page < totalPages,
		},
	}

	if len(meta) > 0 {
		response.Meta = meta[0]
	}

	return response
}

// SuccessWithPagination 返回带分页的成功响应
func (h *BaseHandler) SuccessWithPagination(c *gin.Context, data interface{}, page, pageSize int, total int64, meta ...map[string]interface{}) {
	response := h.NewPaginationResponse(data, page, pageSize, total, meta...)
	h.Success(c, response)
}
