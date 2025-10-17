package admin

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// APIGenerator 用于根据资源定义生成 RESTful API
type APIGenerator struct {
	engine          *gin.Engine
	resourceService interface {
		Create(ctx context.Context, resourceSlug string, data map[string]interface{}) (map[string]interface{}, error)
		Update(ctx context.Context, resourceSlug string, id interface{}, data map[string]interface{}) error
		Delete(ctx context.Context, resourceSlug string, id interface{}) error
		DeleteBatch(ctx context.Context, resourceSlug string, ids []interface{}) (int64, error)
		Restore(ctx context.Context, resourceSlug string, id interface{}) error
		ForceDelete(ctx context.Context, resourceSlug string, id interface{}) error
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
		Export(
			ctx context.Context,
			resourceSlug string,
			filters map[string]interface{},
			search map[string]interface{},
			orderBy string,
			orderDirection string,
			format string,
		) ([]byte, string, error)
	}
	resourceManager *ResourceManager
}

// NewAPIGenerator 创建一个新的 API 生成器
func NewAPIGenerator(
	engine *gin.Engine,
	resourceService interface {
		Create(ctx context.Context, resourceSlug string, data map[string]interface{}) (map[string]interface{}, error)
		Update(ctx context.Context, resourceSlug string, id interface{}, data map[string]interface{}) error
		Delete(ctx context.Context, resourceSlug string, id interface{}) error
		DeleteBatch(ctx context.Context, resourceSlug string, ids []interface{}) (int64, error)
		Restore(ctx context.Context, resourceSlug string, id interface{}) error
		ForceDelete(ctx context.Context, resourceSlug string, id interface{}) error
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
		Export(
			ctx context.Context,
			resourceSlug string,
			filters map[string]interface{},
			search map[string]interface{},
			orderBy string,
			orderDirection string,
			format string,
		) ([]byte, string, error)
	},
	resourceManager *ResourceManager,
) *APIGenerator {
	return &APIGenerator{
		engine:          engine,
		resourceService: resourceService,
		resourceManager: resourceManager,
	}
}

// RegisterResourceAPI 为指定资源注册 API 路由
func (ag *APIGenerator) RegisterResourceAPI(resource Resource) {
	slug := resource.GetSlug()
	if slug == "" {
		// 如果没有定义 slug，则使用标题的小写形式
		slug = lowercase(resource.GetTitle())
	}

	// 注册资源的 CRUD 路由
	resourceGroup := ag.engine.Group("/api/v1/" + slug)
	{
		resourceGroup.GET("/", ag.listHandler(resource, slug))
		resourceGroup.POST("/", ag.createHandler(resource, slug))
		resourceGroup.GET("/:id", ag.detailHandler(resource, slug))
		resourceGroup.PUT("/:id", ag.updateHandler(resource, slug))
		resourceGroup.DELETE("/:id", ag.deleteHandler(resource, slug))
		resourceGroup.POST("/:id/restore", ag.restoreHandler(resource, slug))
		resourceGroup.DELETE("/:id/force", ag.forceDeleteHandler(resource, slug))
		// 添加导出路由
		resourceGroup.GET("/export", ag.exportHandler(resource, slug))
		// 添加批量删除路由
		resourceGroup.DELETE("/", ag.deleteBatchHandler(resource, slug))
		// 通用动作执行路由（单条与批量）
		resourceGroup.POST("/actions/:action", ag.runActionHandler(resource, slug))
	}
}

// RegisterAdminAPI 注册管理API路由
func (ag *APIGenerator) RegisterAdminAPI() {
	// 注册获取所有资源和页面的路由
	ag.engine.GET("/api/admin/resources", ag.listResourcesHandler())

	// 为每个资源注册 API 路由
	resources := ag.resourceManager.GetResources()
	for _, resource := range resources {
		// 跳过特殊资源（如仪表板）
		if resource.GetSlug() == "dashboard" {
			continue
		}
		ag.RegisterResourceAPI(resource)
	}
}

// listResourcesHandler 处理获取所有资源列表请求
func (ag *APIGenerator) listResourcesHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取语言参数
		language := c.Query("language")
		if language == "" {
			language = "zh-CN"
		}

		resources := ag.resourceManager.GetResources()
		pages := ag.resourceManager.GetPages()

		// 构建资源列表
		resourceList := make([]map[string]interface{}, 0)
		for _, resource := range resources {
			// 跳过特殊资源（如仪表板）
			if resource.GetSlug() == "dashboard" {
				continue
			}

			resourceMap := map[string]interface{}{
				"title":      translateResourceTitle(resource.GetTitle(), language),
				"slug":       resource.GetSlug(),
				"model":      resource.GetModel(),
				"fields":     resource.GetFields(),
				"actions":    resource.GetActions(),
				"editable":   true,
				"creatable":  true,
				"viewable":   true,
				"deletable":  true,
				"exportable": true,
			}
			resourceList = append(resourceList, resourceMap)
		}

		// 构建页面列表
		pageList := make([]map[string]interface{}, 0)
		for _, page := range pages {
			pageMap := map[string]interface{}{
				"title": translateResourceTitle(page.GetTitle(), language),
				"slug":  page.GetSlug(),
				"type":  "page",
			}
			pageList = append(pageList, pageMap)
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": map[string]interface{}{
				"resources": resourceList,
				"pages":     pageList,
			},
		})
	}
}

// translateResourceTitle 翻译资源标题
func translateResourceTitle(title, language string) string {
	// 这里可以根据需要实现翻译逻辑
	// 目前直接返回标题，后续可以扩展为从资源映射中获取翻译
	key := strings.ToLower(strings.ReplaceAll(title, " ", "_"))
	// 简单的中英文映射示例
	if language == "en" {
		switch key {
		case "用户":
			return "Users"
		case "角色":
			return "Roles"
		case "部门":
			return "Departments"
		case "文章":
			return "Posts"
		case "分类":
			return "Categories"
		case "菜单":
			return "Menus"
		case "api":
			return "APIs"
		case "仪表板":
			return "Dashboard"
		default:
			return title
		}
	}
	return title
}

// listHandler 处理列表请求
func (ag *APIGenerator) listHandler(resource Resource, slug string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取分页参数
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 10
		}

		// 获取过滤参数
		filters := make(map[string]interface{})
		search := make(map[string]interface{})

		// 处理查询参数
		for key, values := range c.Request.URL.Query() {
			// 跳过分页参数
			if key == "page" || key == "page_size" || key == "order_by" || key == "order_direction" {
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

		// 调用服务获取数据
		results, total, err := ag.resourceService.List(c, slug, page, pageSize, filters, search, orderBy, orderDirection)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取数据失败: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"items":     results,
				"total":     total,
				"page":      page,
				"page_size": pageSize,
			},
		})
	}
}

// createHandler 处理创建请求
func (ag *APIGenerator) createHandler(resource Resource, slug string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData map[string]interface{}
		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求数据格式错误",
			})
			return
		}

		// 验证数据
		errors := ValidateResourceData(resource, requestData)
		if len(errors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "数据验证失败",
				"errors":  errors,
			})
			return
		}

		// 调用服务创建数据
		result, err := ag.resourceService.Create(c, slug, requestData)
		if err != nil {
			// 检查是否为验证错误
			if validationErr, ok := err.(*ValidationError); ok {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "数据验证失败",
					"errors":  validationErr.Errors,
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "创建失败: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "创建成功",
			"data":    result,
		})
	}
}

// detailHandler 处理详情请求
func (ag *APIGenerator) detailHandler(resource Resource, slug string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "缺少ID参数",
			})
			return
		}

		// 调用服务获取数据
		result, err := ag.resourceService.Get(c, slug, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取数据失败: " + err.Error(),
			})
			return
		}

		if result == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "数据不存在",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data":    result,
		})
	}
}

// updateHandler 处理更新请求
func (ag *APIGenerator) updateHandler(resource Resource, slug string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "缺少ID参数",
			})
			return
		}

		var requestData map[string]interface{}
		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求数据格式错误",
			})
			return
		}

		// 验证数据
		errors := ValidateResourceData(resource, requestData)
		if len(errors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "数据验证失败",
				"errors":  errors,
			})
			return
		}

		// 调用服务更新数据
		err := ag.resourceService.Update(c, slug, id, requestData)
		if err != nil {
			// 检查是否为验证错误
			if validationErr, ok := err.(*ValidationError); ok {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "数据验证失败",
					"errors":  validationErr.Errors,
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "更新失败: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "更新成功",
		})
	}
}

// deleteHandler 处理删除请求
func (ag *APIGenerator) deleteHandler(resource Resource, slug string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "缺少ID参数",
			})
			return
		}

		// 调用服务删除数据
		err := ag.resourceService.Delete(c, slug, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "删除失败: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "删除成功",
		})
	}
}

// deleteBatchHandler 处理批量删除请求
func (ag *APIGenerator) deleteBatchHandler(resource Resource, slug string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析请求体中的 IDs
		var requestData struct {
			IDs []interface{} `json:"ids"`
		}

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求数据格式错误",
			})
			return
		}

		if len(requestData.IDs) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请提供要删除的记录ID",
			})
			return
		}

		// 调用服务批量删除数据
		affected, err := ag.resourceService.DeleteBatch(c, slug, requestData.IDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "批量删除失败: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "批量删除成功",
			"data": map[string]interface{}{
				"affected": affected,
			},
		})
	}
}

// exportHandler 处理导出请求
func (ag *APIGenerator) exportHandler(resource Resource, slug string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取过滤参数
		filters := make(map[string]interface{})
		search := make(map[string]interface{})

		// 处理查询参数
		for key, values := range c.Request.URL.Query() {
			// 跳过导出格式参数
			if key == "format" {
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

		// 调用服务导出数据
		data, filename, err := ag.resourceService.Export(c, slug, filters, search, orderBy, orderDirection, format)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "导出失败: " + err.Error(),
			})
			return
		}

		// 设置响应头
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Header("Content-Length", strconv.Itoa(len(data)))

		// 输出文件内容
		c.Data(http.StatusOK, "application/octet-stream", data)
	}
}

// runActionHandler 执行动作（支持批量）
func (ag *APIGenerator) runActionHandler(resource Resource, slug string) gin.HandlerFunc {
	type reqBody struct {
		IDs    []interface{}          `json:"ids"`
		Params map[string]interface{} `json:"params"`
	}
	return func(c *gin.Context) {
		actionName := c.Param("action")
		if actionName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 action 参数"})
			return
		}
		exec, ok := any(resource).(interface {
			RunAction(ctx context.Context, actionName string, ids []interface{}, params map[string]interface{}) (interface{}, error)
		})
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "该资源未实现动作执行"})
			return
		}
		var body reqBody
		if err := c.ShouldBindJSON(&body); err != nil {
			// 允许无体，仅通过 URL 触发
			body.Params = map[string]interface{}{}
		}
		result, err := exec.RunAction(c, actionName, body.IDs, body.Params)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": result})
	}
}

// restoreHandler 恢复软删除
func (ag *APIGenerator) restoreHandler(resource Resource, slug string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少ID参数"})
			return
		}
		if err := ag.resourceService.Restore(c, slug, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "恢复失败: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "恢复成功"})
	}
}

// forceDeleteHandler 强制删除
func (ag *APIGenerator) forceDeleteHandler(resource Resource, slug string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少ID参数"})
			return
		}
		if err := ag.resourceService.ForceDelete(c, slug, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "强制删除失败: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "强制删除成功"})
	}
}

// lowercase 将字符串转换为小写（简单实现）
func lowercase(s string) string {
	result := make([]rune, len(s))
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			result[i] = r + ('a' - 'A')
		} else {
			result[i] = r
		}
	}
	return string(result)
}

// ValidationError 验证错误类型
type ValidationError struct {
	Errors map[string][]string
}

func (e *ValidationError) Error() string {
	return "validation failed"
}
