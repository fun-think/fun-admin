package admin

import "context"

// Resource 是所有管理资源的接口
type Resource interface {
	GetTitle() string
	GetSlug() string
	GetModel() interface{}
	GetFields() []Field
	GetActions() []Action
	GetReadOnlyFields() []string
	// 新增：表格列与筛选
	GetColumns() []*Column
	GetFilters() []*Filter
}

// Sortable 可选接口：声明允许排序的字段
type Sortable interface {
	GetSortableFields() []string
}

// Filterable 可选接口：声明允许精确过滤的字段
type Filterable interface {
	GetFilterableFields() []string
}

// Searchable 可选接口：声明允许模糊搜索的字段
type Searchable interface {
	GetSearchableFields() []string
}

// DefaultOrder 可选接口：声明默认排序
// 返回字段名与方向（"ASC"/"DESC"），方向大小写不敏感
type DefaultOrder interface {
	GetDefaultOrder() (field string, direction string)
}

// Exportable 可选接口：声明资源是否支持导出功能
type Exportable interface {
	IsExportable() bool
}

// NavigationMeta 定义资源在导航中的元信息（参考 Filament 导航能力）
// 通过可选接口提供，避免对现有资源实现造成破坏性变更
type NavigationMeta struct {
	Icon  string // 图标，例如 "heroicons-outline:user"
	Group string // 分组名，例如 "系统管理"
	Sort  int    // 排序（越小越靠前）
}

// NavigationIcon 可选接口：声明导航图标
type NavigationIcon interface {
	GetNavigationIcon() string
}

// NavigationGroup 可选接口：声明导航分组
type NavigationGroup interface {
	GetNavigationGroup() string
}

// NavigationSort 可选接口：声明导航排序
type NavigationSort interface {
	GetNavigationSort() int
}

// NavigationBadge 可选接口：返回导航徽章数
// 可按需基于上下文（如当前用户）动态计算
type NavigationBadge interface {
	GetNavigationBadge(ctx context.Context) (int64, error)
}

// HiddenInNavigation 可选接口：控制资源是否在导航中隐藏
type HiddenInNavigation interface {
	IsHiddenInNavigation(ctx context.Context) bool
}

// FrontendCapabilities 前端可见性能力开关
type FrontendCapabilities struct {
	Editable   bool
	Creatable  bool
	Viewable   bool
	Deletable  bool
	Exportable bool
}

// CapabilityProvider 可选接口：提供前端能力开关，支持按上下文动态控制
type CapabilityProvider interface {
	GetFrontendCapabilities(ctx context.Context) FrontendCapabilities
}

// ActionExecutor 可选接口：由资源实现具体动作的执行（含批量）
// actionName 对应前端声明的动作名；ids 可空（单动作）或多选（批量）；params 为动作表单参数
type ActionExecutor interface {
	RunAction(ctx context.Context, actionName string, ids []interface{}, params map[string]interface{}) (interface{}, error)
}
