package resources

import (
	"context"

	"fun-admin/internal/model"
	"fun-admin/pkg/admin"
)

// CrudTableResource defines the schema for the demo CRUD table.
type CrudTableResource struct {
	admin.BaseResource
}

// NewCrudTableResource creates a new CrudTableResource instance.
func NewCrudTableResource() *CrudTableResource {
	return &CrudTableResource{}
}

// GetTitle returns the display title for the resource.
func (r *CrudTableResource) GetTitle() string {
	return "增删改查表格"
}

// GetSlug returns the identifier used for routing and tables.
func (r *CrudTableResource) GetSlug() string {
	return "crud_items"
}

// GetModel returns the underlying GORM model.
func (r *CrudTableResource) GetModel() interface{} {
	return &model.CrudItem{}
}

// GetFields describes the form/table fields.
func (r *CrudTableResource) GetFields() []admin.Field {
	return []admin.Field{
		admin.NewIDField().Label("ID"),
		admin.NewTextField("name").Label("名").Required(),
		admin.NewTextField("value").Label("值").Required(),
		admin.NewTextareaField("remark").Label("备注"),
		admin.NewDateTimeField("created_at").Label("创建时间"),
		admin.NewDateTimeField("updated_at").Label("更新时间"),
	}
}

// GetActions exposes available actions for the resource.
func (r *CrudTableResource) GetActions() []admin.Action {
	return []admin.Action{
		admin.NewViewAction().Label("查看"),
		admin.NewEditAction().Label("编辑"),
		admin.NewDeleteAction().Label("删除"),
		admin.NewAction("reset_values").
			Label("重置值").
			AsBulk().
			Confirm("确认重置选中记录的值？"),
		admin.NewAction("bulk_delete").
			Label("批量删除").
			Color("danger").
			AsBulk().
			Confirm("确认批量删除选中的记录？"),
	}
}

// GetReadOnlyFields returns fields that must remain read-only.
func (r *CrudTableResource) GetReadOnlyFields() []string {
	return []string{"id", "created_at", "updated_at"}
}

// GetColumns defines table columns.
func (r *CrudTableResource) GetColumns() []*admin.Column {
	return []*admin.Column{
		admin.NewColumn("id", "ID", "number").SetSortable(true),
		admin.NewColumn("name", "名", "text").SetSortable(true),
		admin.NewColumn("value", "值", "text"),
		admin.NewColumn("remark", "备注", "text"),
		admin.NewColumn("created_at", "创建时间", "datetime").SetSortable(true),
	}
}

// GetFilters defines filterable fields.
func (r *CrudTableResource) GetFilters() []*admin.Filter {
	return []*admin.Filter{
		{Name: "name", Label: "名", Type: "text"},
		{Name: "value", Label: "值", Type: "text"},
		{Name: "remark", Label: "备注", Type: "text"},
	}
}

// GetSearchableFields returns fields that support search.
func (r *CrudTableResource) GetSearchableFields() []string {
	return []string{"name", "value", "remark"}
}

// GetFilterableFields defines explicit filters.
func (r *CrudTableResource) GetFilterableFields() []string {
	return []string{"name", "value", "remark"}
}

// GetSortableFields returns sortable field names.
func (r *CrudTableResource) GetSortableFields() []string {
	return []string{"name", "value", "created_at", "updated_at"}
}

// GetDefaultOrder supplies default sorting.
func (r *CrudTableResource) GetDefaultOrder() (string, string) {
	return "created_at", "DESC"
}

// IsHiddenInNavigation specifies whether to hide this resource in navigation.
func (r *CrudTableResource) IsHiddenInNavigation(ctx context.Context) bool {
	return false
}
