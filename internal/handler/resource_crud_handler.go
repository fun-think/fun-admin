package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"

	"fun-admin/internal/service"
	"fun-admin/pkg/admin/i18n"

	"github.com/gin-gonic/gin"
)

type ResourceService interface {
	Create(ctx context.Context, resourceSlug string, data map[string]interface{}) (map[string]interface{}, error)
	Update(ctx context.Context, resourceSlug string, id interface{}, data map[string]interface{}) error
	Delete(ctx context.Context, resourceSlug string, id interface{}) error
	Get(ctx context.Context, resourceSlug string, id interface{}) (map[string]interface{}, error)
	List(
		ctx context.Context,
		resourceSlug string,
		page, pageSize int,
		filters map[string]interface{},
		search map[string]interface{},
		orderBy string,
		orderDirection string,
	) ([]map[string]interface{}, int64, error)
	RunAction(ctx context.Context, resourceSlug string, actionName string, ids []interface{}, params map[string]interface{}) (interface{}, error)
}

// ResourceCRUDHandler 资源 CRUD 处理器
type ResourceCRUDHandler struct {
	resourceService ResourceService
}

// NewResourceCRUDHandler 创建资源 CRUD 处理器
func NewResourceCRUDHandler(resourceService ResourceService) *ResourceCRUDHandler {
	return &ResourceCRUDHandler{
		resourceService: resourceService,
	}
}

// List 处理资源列表请求
func (h *ResourceCRUDHandler) List(c *gin.Context) {
	slug := c.Param("slug")

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // 限制最大页面大小
	}

	// 获取语言参数
	language := getLanguage(c)

	// 获取过滤参数
	filters := make(map[string]interface{})
	search := make(map[string]interface{})

	// 处理查询参数
	query := c.Request.URL.Query()
	for key, values := range query {
		// 跳过分页参数
		if key == "page" || key == "page_size" || key == "order_by" || key == "order_direction" || key == "language" {
			continue
		}

		// 如果参数名以 "search_" 开头，则作为搜索条件
		if len(values) > 0 && values[0] != "" {
			if len(key) > 7 && key[:7] == "search_" {
				fieldName := key[7:] // 去掉 "search_" 前缀
				search[fieldName] = values[0]
			} else {
				// 否则作为过滤条件
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

	// 调用服务获取数据
	results, total, err := h.resourceService.List(c, slug, page, pageSize, filters, search, orderBy, orderDirection)
	if err != nil {
		var notFoundErr *service.ResourceNotFoundError
		if errors.As(err, &notFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": i18n.Translate(language, "error.resource_not_found"),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": messageWithDebugError(i18n.Translate(language, "error.failed_to_get_data"), err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"items":     results,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
		"message": "success",
	})
}

// Create 处理资源创建请求
func (h *ResourceCRUDHandler) Create(c *gin.Context) {
	slug := c.Param("slug")

	// 获取语言参数
	language := getLanguage(c)

	var requestData map[string]interface{}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": i18n.Translate(language, "error.invalid_request_data"),
		})
		return
	}

	// 调用服务创建数据
	result, err := h.resourceService.Create(c, slug, requestData)
	if err != nil {
		var notFoundErr *service.ResourceNotFoundError
		if errors.As(err, &notFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": i18n.Translate(language, "error.resource_not_found"),
			})
			return
		}

		var validationErr *service.ValidationError
		if errors.As(err, &validationErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": i18n.Translate(language, "error.validation_failed"),
				"errors":  validationErr.Errors,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": messageWithDebugError(i18n.Translate(language, "error.failed_to_create_record"), err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    result,
		"message": i18n.Translate(language, "message.created_successfully"),
	})
}

// Get 处理资源详情请求
func (h *ResourceCRUDHandler) Get(c *gin.Context) {
	slug := c.Param("slug")
	id := c.Param("id")

	// 获取语言参数
	language := getLanguage(c)

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": i18n.Translate(language, "error.missing_id_parameter"),
		})
		return
	}

	// 调用服务获取数据
	result, err := h.resourceService.Get(c, slug, id)
	if err != nil {
		var notFoundErr *service.ResourceNotFoundError
		if errors.As(err, &notFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": i18n.Translate(language, "error.resource_not_found"),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": messageWithDebugError(i18n.Translate(language, "error.failed_to_get_data"), err),
		})
		return
	}

	if result == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": i18n.Translate(language, "error.record_not_found"),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    result,
		"message": "success",
	})
}

// Update 处理资源更新请求
func (h *ResourceCRUDHandler) Update(c *gin.Context) {
	slug := c.Param("slug")
	id := c.Param("id")

	// 获取语言参数
	language := getLanguage(c)

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": i18n.Translate(language, "error.missing_id_parameter"),
		})
		return
	}

	var requestData map[string]interface{}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": i18n.Translate(language, "error.invalid_request_data"),
		})
		return
	}

	// 调用服务更新数据
	err := h.resourceService.Update(c, slug, id, requestData)
	if err != nil {
		var notFoundErr *service.ResourceNotFoundError
		if errors.As(err, &notFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": i18n.Translate(language, "error.resource_not_found"),
			})
			return
		}

		var validationErr *service.ValidationError
		if errors.As(err, &validationErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": i18n.Translate(language, "error.validation_failed"),
				"errors":  validationErr.Errors,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": messageWithDebugError(i18n.Translate(language, "error.failed_to_update_record"), err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": i18n.Translate(language, "message.updated_successfully"),
	})
}

// Delete 处理资源删除请求
func (h *ResourceCRUDHandler) Delete(c *gin.Context) {
	slug := c.Param("slug")
	id := c.Param("id")

	// 获取语言参数
	language := getLanguage(c)

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": i18n.Translate(language, "error.missing_id_parameter"),
		})
		return
	}

	// 调用服务删除数据
	err := h.resourceService.Delete(c, slug, id)
	if err != nil {
		var notFoundErr *service.ResourceNotFoundError
		if errors.As(err, &notFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": i18n.Translate(language, "error.resource_not_found"),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": messageWithDebugError(i18n.Translate(language, "error.failed_to_delete_record"), err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": i18n.Translate(language, "message.deleted_successfully"),
	})
}

// RunAction 执行资源动作（包括批量）
func (h *ResourceCRUDHandler) RunAction(c *gin.Context) {
	slug := c.Param("slug")
	action := c.Param("action")
	language := getLanguage(c)

	var payload struct {
		IDs    []interface{}          `json:"ids"`
		Params map[string]interface{} `json:"params"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		if errors.Is(err, io.EOF) {
			payload.Params = map[string]interface{}{}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": i18n.Translate(language, "error.invalid_request_data"),
			})
			return
		}
	}
	if payload.Params == nil {
		payload.Params = map[string]interface{}{}
	}

	result, err := h.resourceService.RunAction(c, slug, action, payload.IDs, payload.Params)
	if err != nil {
		var notFoundErr *service.ResourceNotFoundError
		if errors.As(err, &notFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": i18n.Translate(language, "error.resource_not_found"),
			})
			return
		}

		if errors.Is(err, service.ErrActionNotSupported) {
			c.JSON(http.StatusNotImplemented, gin.H{
				"code":    http.StatusNotImplemented,
				"message": i18n.Translate(language, "error.action_not_supported"),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": messageWithDebugError(i18n.Translate(language, "error.failed_to_perform_action"), err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    result,
		"message": "success",
	})
}
