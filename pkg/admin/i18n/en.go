package i18n

// EnglishResources 英文语言资源
var EnglishResources = map[string]string{
	// 通用翻译
	"dashboard":      "Dashboard",
	"resources":      "Resources",
	"create":         "Create",
	"edit":           "Edit",
	"delete":         "Delete",
	"save":           "Save",
	"cancel":         "Cancel",
	"search":         "Search",
	"reset":          "Reset",
	"export":         "Export",
	"batch_delete":   "Batch Delete",
	"actions":        "Actions",
	"confirm":        "Confirm",
	"are_you_sure":   "Are you sure?",
	"yes":            "Yes",
	"no":             "No",
	"select_all":     "Select All",
	" deselect_all":  "Deselect All",
	"selected_items": "%d items selected",

	// 字段类型
	"text":         "Text",
	"email":        "Email",
	"number":       "Number",
	"select":       "Select",
	"textarea":     "Textarea",
	"boolean":      "Boolean",
	"date":         "Date",
	"datetime":     "DateTime",
	"relationship": "Relationship",
	"file":         "File",

	// 资源标题
	"user":       "User",
	"department": "Department",
	"post":       "Post",
	"role":       "Role",

	// 用户资源字段
	"user.id":            "ID",
	"user.name":          "Name",
	"user.email":         "Email",
	"user.password":      "Password",
	"user.department_id": "Department",
	"user.created_at":    "Created At",
	"user.updated_at":    "Updated At",

	// 部门资源字段
	"department.id":         "ID",
	"department.name":       "Name",
	"department.created_at": "Created At",
	"department.updated_at": "Updated At",

	// 文章资源字段
	"post.id":         "ID",
	"post.title":      "Title",
	"post.slug":       "Slug",
	"post.status":     "Status",
	"post.content":    "Content",
	"post.image":      "Cover Image",
	"post.created_at": "Created At",
	"post.updated_at": "Updated At",

	// 角色资源字段
	"role.id":          "ID",
	"role.name":        "Name",
	"role.description": "Description",
	"role.created_at":  "Created At",
	"role.updated_at":  "Updated At",

	// 状态
	"draft":     "Draft",
	"published": "Published",
	"archived":  "Archived",

	// 错误消息
	"error.create_failed":     "Create failed",
	"error.update_failed":     "Update failed",
	"error.delete_failed":     "Delete failed",
	"error.load_failed":       "Load failed",
	"error.validation_failed": "Validation failed",

	// 成功消息
	"success.create": "Created successfully",
	"success.update": "Updated successfully",
	"success.delete": "Deleted successfully",

	// 验证消息
	"validation.required":   "%s is required",
	"validation.email":      "%s must be a valid email address",
	"validation.min_length": "%s must be at least %d characters",
	"validation.max_length": "%s must be at most %d characters",
	"validation.min_value":  "%s must be at least %d",
	"validation.max_value":  "%s must be at most %d",
}
