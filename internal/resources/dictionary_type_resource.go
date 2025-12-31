package resources

import (
	"context"

	"fun-admin/internal/model"
	"fun-admin/pkg/admin"
)

var dictStatusOptions = []admin.Option{
	{Value: "1", Label: "启用"},
	{Value: "0", Label: "禁用"},
}

var dictStatusBadgeMap = map[string]string{
	"1": "#52C41A",
	"0": "#F5222D",
}

var dictStatusEnumMap = map[string]string{
	"1": "启用",
	"0": "禁用",
}

// DictionaryTypeResource defines the schema for dictionary type management.
type DictionaryTypeResource struct {
	admin.BaseResource
}

// NewDictionaryTypeResource creates a new DictionaryTypeResource.
func NewDictionaryTypeResource() *DictionaryTypeResource {
	return &DictionaryTypeResource{}
}

// GetTitle returns the display title.
func (r *DictionaryTypeResource) GetTitle() string {
	return "字典类型"
}

// GetSlug returns the resource slug.
func (r *DictionaryTypeResource) GetSlug() string {
	return "admin_dictionary_type"
}

// GetModel returns the model behind this resource.
func (r *DictionaryTypeResource) GetModel() interface{} {
	return &model.DictionaryType{}
}

// GetFields describes the form and table fields.
func (r *DictionaryTypeResource) GetFields() []admin.Field {
	return []admin.Field{
		admin.NewIDField().Label("ID"),
		admin.NewTextField("name").Label("名称").Required(),
		admin.NewTextField("code").Label("编码").Required(),
		admin.NewSelectField("status").Label("状态").Required().SetOptions(dictStatusOptions),
		admin.NewNumberField("sort").Label("排序"),
		admin.NewTextareaField("remark").Label("备注").SetRows(3),
		admin.NewDateTimeField("created_at").Label("创建时间"),
		admin.NewDateTimeField("updated_at").Label("更新时间"),
	}
}

// GetActions exposes resource actions.
func (r *DictionaryTypeResource) GetActions() []admin.Action {
	return []admin.Action{
		admin.NewViewAction().Label("查看"),
		admin.NewEditAction().Label("编辑"),
		admin.NewDeleteAction().Label("删除"),
	}
}

// GetReadOnlyFields lists fields that must stay read-only.
func (r *DictionaryTypeResource) GetReadOnlyFields() []string {
	return []string{"id", "created_at", "updated_at"}
}

// GetColumns defines table columns.
func (r *DictionaryTypeResource) GetColumns() []*admin.Column {
	return []*admin.Column{
		admin.NewColumn("id", "ID", "number").SetSortable(true),
		admin.NewColumn("name", "名称", "text").SetSortable(true),
		admin.NewColumn("code", "编码", "text").SetSortable(true),
		admin.NewColumn("status", "状态", "badge").
			SetSortable(true).
			SetBadgeMap(dictStatusBadgeMap).
			SetEnumMap(dictStatusEnumMap),
		admin.NewColumn("sort", "排序", "number"),
		admin.NewColumn("created_at", "创建时间", "datetime").SetSortable(true),
	}
}

// GetFilters defines filterable fields.
func (r *DictionaryTypeResource) GetFilters() []*admin.Filter {
	return []*admin.Filter{
		{Name: "name", Label: "名称", Type: "text"},
		{Name: "code", Label: "编码", Type: "text"},
		{Name: "status", Label: "状态", Type: "select", Options: dictStatusOptions},
	}
}

// GetSearchableFields exposes search white list.
func (r *DictionaryTypeResource) GetSearchableFields() []string {
	return []string{"name", "code"}
}

// GetFilterableFields exposes filter white list.
func (r *DictionaryTypeResource) GetFilterableFields() []string {
	return []string{"name", "code", "status"}
}

// GetSortableFields exposes sortable fields.
func (r *DictionaryTypeResource) GetSortableFields() []string {
	return []string{"name", "code", "sort", "created_at"}
}

// GetDefaultOrder supplies the default sorting.
func (r *DictionaryTypeResource) GetDefaultOrder() (string, string) {
	return "sort", "ASC"
}

// IsHiddenInNavigation controls navigation visibility.
func (r *DictionaryTypeResource) IsHiddenInNavigation(ctx context.Context) bool {
	return false
}
