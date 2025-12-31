package resources

import (
	"context"

	"fun-admin/internal/model"
	"fun-admin/pkg/admin"
)

// DictionaryDataResource defines the schema for dictionary entry management.
type DictionaryDataResource struct {
	admin.BaseResource
}

// NewDictionaryDataResource creates a new DictionaryDataResource.
func NewDictionaryDataResource() *DictionaryDataResource {
	return &DictionaryDataResource{}
}

// GetTitle returns the display title.
func (r *DictionaryDataResource) GetTitle() string {
	return "字典数据"
}

// GetSlug returns the resource slug.
func (r *DictionaryDataResource) GetSlug() string {
	return "admin_dictionary_data"
}

// GetModel returns the model.
func (r *DictionaryDataResource) GetModel() interface{} {
	return &model.DictionaryData{}
}

// GetFields describes the form fields.
func (r *DictionaryDataResource) GetFields() []admin.Field {
	return []admin.Field{
		admin.NewIDField().Label("ID"),
		admin.NewRelationshipField("type_id", "admin_dictionary_type").
			Label("字典类型").
			Required().
			SetDisplayField("name"),
		admin.NewTextField("label").Label("标签").Required(),
		admin.NewTextField("value").Label("值").Required(),
		admin.NewSelectField("status").Label("状态").Required().SetOptions(dictStatusOptions),
		admin.NewBooleanField("is_default").Label("默认值"),
		admin.NewNumberField("sort").Label("排序"),
		admin.NewTextareaField("remark").Label("备注").SetRows(3),
		admin.NewTextField("ext1").Label("扩展字段1"),
		admin.NewTextField("ext2").Label("扩展字段2"),
		admin.NewTextField("ext3").Label("扩展字段3"),
		admin.NewDateTimeField("created_at").Label("创建时间"),
		admin.NewDateTimeField("updated_at").Label("更新时间"),
	}
}

// GetActions exposes actions.
func (r *DictionaryDataResource) GetActions() []admin.Action {
	return []admin.Action{
		admin.NewViewAction().Label("查看"),
		admin.NewEditAction().Label("编辑"),
		admin.NewDeleteAction().Label("删除"),
	}
}

// GetReadOnlyFields defines read-only fields.
func (r *DictionaryDataResource) GetReadOnlyFields() []string {
	return []string{"id", "created_at", "updated_at"}
}

// GetFilters defines filters.
func (r *DictionaryDataResource) GetFilters() []*admin.Filter {
	return []*admin.Filter{
		{Name: "type_id", Label: "字典类型", Type: "text"},
		{Name: "status", Label: "状态", Type: "select", Options: dictStatusOptions},
	}
}

// GetSearchableFields defines searchable fields.
func (r *DictionaryDataResource) GetSearchableFields() []string {
	return []string{"label", "value"}
}

// GetFilterableFields defines filterable fields.
func (r *DictionaryDataResource) GetFilterableFields() []string {
	return []string{"type_id", "status", "is_default"}
}

// GetSortableFields defines sortable fields.
func (r *DictionaryDataResource) GetSortableFields() []string {
	return []string{"label", "value", "sort", "created_at"}
}

// GetDefaultOrder supplies default ordering.
func (r *DictionaryDataResource) GetDefaultOrder() (string, string) {
	return "sort", "ASC"
}

// IsHiddenInNavigation controls navigation visibility.
func (r *DictionaryDataResource) IsHiddenInNavigation(ctx context.Context) bool {
	return false
}
