package handler

import (
	"net/http"
	"strings"

	"fun-admin/internal/repository"
	"fun-admin/pkg/admin"
	"fun-admin/pkg/admin/i18n"
	"github.com/gin-gonic/gin"
)

// ResourceHandler 资源处理器
type ResourceHandler struct {
	repo *repository.ResourceRepository
}

// NewResourceHandler 创建资源处理器
func NewResourceHandler(repo *repository.ResourceRepository) *ResourceHandler {
	return &ResourceHandler{repo: repo}
}

// GlobalSearch 全局搜索
func (h *ResourceHandler) GlobalSearch(c *gin.Context) {
	keyword := c.Query("keyword")
	language := c.Query("language")
	if language == "" {
		language = "zh-CN"
	}

	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": i18n.Translate(language, "error.keyword_required"),
		})
		return
	}

	// 获取所有资源
	resources := admin.GlobalResourceManager.GetResources()

	// 执行全局搜索（简化实现）
	results := make([]map[string]interface{}, 0)
	for _, resource := range resources {
		// 为每个资源创建一个简单的搜索结果条目
		result := map[string]interface{}{
			"resource": resource.GetSlug(),
			"title":    translateResourceTitle(resource.GetTitle(), language),
			"keyword":  keyword,
		}
		results = append(results, result)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": results,
	})
}

// ListResources 获取所有资源列表
func (h *ResourceHandler) ListResources(c *gin.Context) {
	// 获取语言参数
	language := c.Query("language")
	if language == "" {
		language = "zh-CN"
	}

	resources := admin.GlobalResourceManager.GetResources()
	pages := admin.GlobalResourceManager.GetPages()

	// 构建资源列表
	resourceList := make([]map[string]interface{}, 0)
	for _, resource := range resources {
		// 跳过特殊资源（如仪表板）
		if resource.GetSlug() == "dashboard" {
			continue
		}

		// 导航可见性与元信息
		if v, ok := any(resource).(admin.HiddenInNavigation); ok {
			if v.IsHiddenInNavigation(c.Request.Context()) {
				continue
			}
		}
		icon := ""
		if v, ok := any(resource).(admin.NavigationIcon); ok {
			icon = v.GetNavigationIcon()
		}
		group := ""
		if v, ok := any(resource).(admin.NavigationGroup); ok {
			group = v.GetNavigationGroup()
		}
		sort := 0
		if v, ok := any(resource).(admin.NavigationSort); ok {
			sort = v.GetNavigationSort()
		}
		var badgeCount *int64
		if v, ok := any(resource).(admin.NavigationBadge); ok {
			if count, err := v.GetNavigationBadge(c.Request.Context()); err == nil {
				badgeCount = new(int64)
				*badgeCount = count
			}
		}
		// 能力开关（默认 true）
		caps := admin.FrontendCapabilities{Editable: true, Creatable: true, Viewable: true, Deletable: true, Exportable: true}
		if v, ok := any(resource).(admin.CapabilityProvider); ok {
			caps = v.GetFrontendCapabilities(c.Request.Context())
		}

		resourceMap := map[string]interface{}{
			"title":           translateResourceTitle(resource.GetTitle(), language),
			"slug":            resource.GetSlug(),
			"model":           resource.GetModel(),
			"fields":          resource.GetFields(),
			"columns":         resource.GetColumns(),
			"filters":         resource.GetFilters(),
			"actions":         resource.GetActions(),
			"nav_icon":        icon,
			"nav_group":       group,
			"nav_sort":        sort,
			"nav_badge_count": badgeCount,
			"editable":        caps.Editable,
			"creatable":       caps.Creatable,
			"viewable":        caps.Viewable,
			"deletable":       caps.Deletable,
			"exportable":      caps.Exportable,
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

// GetResource 获取单个资源详情
func (h *ResourceHandler) GetResource(c *gin.Context) {
	slug := c.Param("slug")

	// 获取语言参数
	language := c.Query("language")
	if language == "" {
		language = "zh-CN"
	}

	// 先尝试查找资源
	resource := admin.GlobalResourceManager.GetResourceBySlug(slug)
	if resource != nil {
		// 构建字段信息（仅基础信息）
		fields := resource.GetFields()
		fieldList := make([]map[string]interface{}, len(fields))
		for i, field := range fields {
			fieldList[i] = map[string]interface{}{
				"name":     field.GetName(),
				"label":    translateFieldLabel(field.GetLabel(), language),
				"type":     field.GetType(),
				"required": field.IsRequired(),
			}
		}

		// 操作信息
		actions := resource.GetActions()
		actionList := make([]map[string]interface{}, 0, len(actions))
		for _, action := range actions {
			// 可见性过滤
			if v, ok := any(action).(admin.ActionVisibility); ok {
				if !v.IsVisible(c) {
					continue
				}
			}
			// 动作扩展元信息：图标、颜色、确认、权限键、是否批量
			icon := ""
			if v, ok := any(action).(interface{ GetIcon() string }); ok {
				icon = v.GetIcon()
			}
			color := "default"
			if v, ok := any(action).(interface{ GetColor() string }); ok {
				color = v.GetColor()
			}
			confirm := ""
			if v, ok := any(action).(interface{ GetConfirm() string }); ok {
				confirm = v.GetConfirm()
			}
			permission := ""
			if v, ok := any(action).(interface{ GetPermission() string }); ok {
				permission = v.GetPermission()
			}
			isBulk := false
			if v, ok := any(action).(interface{ IsBulk() bool }); ok {
				isBulk = v.IsBulk()
			}
			item := map[string]interface{}{
				"name":       action.GetName(),
				"label":      translateActionLabel(action.GetLabel(), language),
				"primary":    action.IsPrimary(),
				"icon":       icon,
				"color":      color,
				"confirm":    confirm,
				"permission": permission,
				"bulk":       isBulk,
			}
			// 若带表单，返回 schema
			if v, ok := any(action).(admin.ActionWithForm); ok {
				item["form_fields"] = v.GetFormFields()
			}
			actionList = append(actionList, item)
		}

		resourceMap := map[string]interface{}{
			"title":   translateResourceTitle(resource.GetTitle(), language),
			"slug":    resource.GetSlug(),
			"model":   resource.GetModel(),
			"fields":  fieldList,
			"columns": resource.GetColumns(),
			"filters": resource.GetFilters(),
			"actions": actionList,
			// 预留：由资源可选接口/策略动态控制
			"editable":        true,
			"creatable":       true,
			"viewable":        true,
			"deletable":       true,
			"exportable":      true,
			"readonly_fields": resource.GetReadOnlyFields(),
		}

		// 表格与页面元信息：排序/过滤/搜索白名单与默认排序（用于前端渲染与占位）
		if v, ok := any(resource).(interface{ GetSortableFields() []string }); ok {
			resourceMap["sortable_fields"] = v.GetSortableFields()
		}
		if v, ok := any(resource).(interface{ GetFilterableFields() []string }); ok {
			resourceMap["filterable_fields"] = v.GetFilterableFields()
		}
		if v, ok := any(resource).(interface{ GetSearchableFields() []string }); ok {
			resourceMap["searchable_fields"] = v.GetSearchableFields()
		}
		if v, ok := any(resource).(interface{ GetDefaultOrder() (string, string) }); ok {
			f, d := v.GetDefaultOrder()
			resourceMap["default_order"] = map[string]string{"field": f, "direction": d}
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": resourceMap,
		})
		return
	}

	// 再尝试查找页面
	page := admin.GlobalResourceManager.GetPageBySlug(slug)
	if page != nil {
		pageMap := map[string]interface{}{
			"title": translateResourceTitle(page.GetTitle(), language),
			"slug":  page.GetSlug(),
			"type":  "page",
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": pageMap,
		})
		return
	}

	// 资源不存在
	c.JSON(http.StatusNotFound, gin.H{
		"code":    404,
		"message": i18n.Translate(language, "error.resource_not_found"),
	})
}

// translateResourceTitle 翻译资源标题
func translateResourceTitle(title, language string) string {
	// 这里可以根据需要实现翻译逻辑
	// 目前直接返回标题，后续可以扩展为从资源映射中获取翻译
	key := strings.ToLower(strings.ReplaceAll(title, " ", "_"))
	return i18n.Translate(language, key, title)
}

// translateFieldLabel 翻译字段标签
func translateFieldLabel(label, language string) string {
	key := strings.ToLower(strings.ReplaceAll(label, " ", "_"))
	return i18n.Translate(language, key, label)
}

// translateActionLabel 翻译操作标签
func translateActionLabel(label, language string) string {
	key := strings.ToLower(strings.ReplaceAll(label, " ", "_"))
	return i18n.Translate(language, key, label)
}
