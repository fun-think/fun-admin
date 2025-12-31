package resources

import (
	"context"
	"fun-admin/internal/model"
	"fun-admin/pkg/admin"
)

// UserResource 用户资源定义
type UserResource struct {
	admin.BaseResource
}

// NewUserResource 创建用户资源
func NewUserResource() *UserResource {
	return &UserResource{}
}

// GetTitle 返回资源标题
func (r *UserResource) GetTitle() string {
	return "用户管理"
}

// GetSlug 返回资源标识符
func (r *UserResource) GetSlug() string {
	return "users"
}

// GetModel 返回关联的模型
func (r *UserResource) GetModel() interface{} {
	return &model.User{}
}

// GetFields 返回字段定义
func (r *UserResource) GetFields() []admin.Field {
	return []admin.Field{
		admin.NewIDField().Label("ID"),
		admin.NewTextField("username").Label("用户名").Required(),
		admin.NewTextField("nickname").Label("昵称"),
		admin.NewEmailField("email").Label("邮箱"),
		admin.NewTextField("phone").Label("手机号"),
		admin.NewSelectField("status").Label("状态").SetOptions([]admin.Option{
			{Label: "正常", Value: "1"},
			{Label: "禁用", Value: "2"},
		}).SetDefault("1"),
		admin.NewDateTimeField("created_at").Label("创建时间"),
		admin.NewDateTimeField("updated_at").Label("更新时间"),
	}
}

// GetColumns 返回列表列定义
func (r *UserResource) GetColumns() []*admin.Column {
	return []*admin.Column{
		admin.NewColumn("id", "ID", "number").SetWidth(80),
		admin.NewColumn("username", "用户名", "text"),
		admin.NewColumn("nickname", "昵称", "text"),
		admin.NewColumn("email", "邮箱", "text"),
		admin.NewColumn("phone", "手机号", "text"),
		admin.NewColumn("status", "状态", "badge").SetBadgeMap(map[string]string{
			"1": "正常",
			"2": "禁用",
		}),
		admin.NewColumn("created_at", "创建时间", "datetime"),
	}
}

// GetFilters 返回过滤器定义
func (r *UserResource) GetFilters() []*admin.Filter {
	return []*admin.Filter{
		{Name: "username", Label: "用户名", Type: "text"},
		{Name: "nickname", Label: "昵称", Type: "text"},
		{Name: "email", Label: "邮箱", Type: "text"},
		{Name: "phone", Label: "手机号", Type: "text"},
		{Name: "status", Label: "状态", Type: "select", Options: []admin.Option{
			{Label: "正常", Value: "1"},
			{Label: "禁用", Value: "2"},
		}},
	}
}

// GetActions 返回支持的操作
func (r *UserResource) GetActions() []admin.Action {
	return []admin.Action{
		admin.NewViewAction().Label("查看"),
		admin.NewEditAction().Label("编辑"),
		admin.NewDeleteAction().Label("删除"),
	}
}

// GetSearchableFields 返回可搜索字段
func (r *UserResource) GetSearchableFields() []string {
	return []string{"username", "nickname", "email", "phone"}
}

// GetFilterableFields 返回可过滤字段
func (r *UserResource) GetFilterableFields() []string {
	return []string{"username", "nickname", "email", "phone", "status"}
}

// GetReadOnlyFields 返回只读字段
func (r *UserResource) GetReadOnlyFields() []string {
	return []string{"id", "created_at", "updated_at"}
}

// GetFieldPermissions 返回字段级权限配置
func (r *UserResource) GetFieldPermissions(ctx context.Context) admin.FieldPermissions {
	return admin.FieldPermissions{
		Readable: []string{"id", "username", "nickname", "email", "phone", "status", "created_at", "updated_at"},
		Writable: []string{"username", "nickname", "email", "phone", "status"},
	}
}
