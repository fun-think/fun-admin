package resources

import (
	"fun-admin/internal/model"
	"fun-admin/pkg/admin"
)

// OperationLogResource 操作日志资源
type OperationLogResource struct {
	admin.BaseResource
}

// NewOperationLogResource 创建操作日志资源
func NewOperationLogResource() *OperationLogResource {
	return &OperationLogResource{}
}

// GetTitle 返回资源标题
func (r *OperationLogResource) GetTitle() string {
	return "操作日志"
}

// GetSlug 返回资源标识符
func (r *OperationLogResource) GetSlug() string {
	return "operation-logs"
}

// GetModel 返回关联的模型
func (r *OperationLogResource) GetModel() interface{} {
	return &model.OperationLog{}
}

// GetFields 返回字段定义
func (r *OperationLogResource) GetFields() []admin.Field {
	return []admin.Field{
		admin.NewIDField().Label("ID"),
		admin.NewTextField("user_id").Label("用户ID"),
		admin.NewTextField("username").Label("用户名"),
		admin.NewTextField("ip").Label("IP地址"),
		admin.NewTextField("method").Label("请求方法"),
		admin.NewTextField("path").Label("请求路径"),
		admin.NewTextareaField("user_agent").Label("用户代理"),
		admin.NewDateTimeField("created_at").Label("创建时间"),
	}
}

// GetActions 返回支持的操作
func (r *OperationLogResource) GetActions() []admin.Action {
	// 操作日志通常不允许编辑和删除
	viewAction := admin.NewViewAction()
	viewAction.Label("查看")

	return []admin.Action{
		viewAction,
	}
}

// GetReadOnlyFields 返回只读字段
func (r *OperationLogResource) GetReadOnlyFields() []string {
	return []string{"id", "user_id", "username", "ip", "method", "path", "user_agent", "created_at"}
}

// GetColumns 返回表格列定义
func (r *OperationLogResource) GetColumns() []*admin.Column {
	return []*admin.Column{
		admin.NewColumn("id", "ID", "number").SetSortable(true),
		admin.NewColumn("username", "用户名", "text").SetSortable(true),
		admin.NewColumn("ip", "IP地址", "text"),
		admin.NewColumn("method", "请求方法", "text"),
		admin.NewColumn("path", "请求路径", "text"),
		admin.NewColumn("user_agent", "用户代理", "text"),
		admin.NewColumn("created_at", "创建时间", "datetime").SetSortable(true),
	}
}

// GetFilters 返回过滤器定义
func (r *OperationLogResource) GetFilters() []*admin.Filter {
	return []*admin.Filter{
		{Name: "username", Label: "用户名", Type: "text"},
		{Name: "ip", Label: "IP地址", Type: "text"},
		{Name: "method", Label: "请求方法", Type: "select", Options: []admin.Option{
			{Label: "GET", Value: "GET"},
			{Label: "POST", Value: "POST"},
			{Label: "PUT", Value: "PUT"},
			{Label: "DELETE", Value: "DELETE"},
			{Label: "PATCH", Value: "PATCH"},
		}},
		{Name: "path", Label: "请求路径", Type: "text"},
	}
}
