package i18n

// ChineseResources 中文语言资源
var ChineseResources = map[string]string{
	// 通用翻译
	"dashboard":      "仪表板",
	"resources":      "资源",
	"create":         "创建",
	"edit":           "编辑",
	"delete":         "删除",
	"save":           "保存",
	"cancel":         "取消",
	"search":         "搜索",
	"reset":          "重置",
	"export":         "导出",
	"batch_delete":   "批量删除",
	"actions":        "操作",
	"confirm":        "确认",
	"are_you_sure":   "您确定吗？",
	"yes":            "是",
	"no":             "否",
	"select_all":     "全选",
	" deselect_all":  "取消全选",
	"selected_items": "已选择 %d 项",

	// 字段类型
	"text":         "文本",
	"email":        "邮箱",
	"number":       "数字",
	"select":       "选择",
	"textarea":     "文本域",
	"boolean":      "布尔值",
	"date":         "日期",
	"datetime":     "日期时间",
	"relationship": "关联",
	"file":         "文件",

	// 资源标题
	"user":       "用户",
	"department": "部门",
	"post":       "文章",
	"role":       "角色",

	// 用户资源字段
	"user.id":            "ID",
	"user.name":          "姓名",
	"user.email":         "邮箱",
	"user.password":      "密码",
	"user.department_id": "部门",
	"user.created_at":    "创建时间",
	"user.updated_at":    "更新时间",

	// 部门资源字段
	"department.id":         "ID",
	"department.name":       "名称",
	"department.created_at": "创建时间",
	"department.updated_at": "更新时间",

	// 文章资源字段
	"post.id":         "ID",
	"post.title":      "标题",
	"post.slug":       "别名",
	"post.status":     "状态",
	"post.content":    "内容",
	"post.image":      "封面图片",
	"post.created_at": "创建时间",
	"post.updated_at": "更新时间",

	// 角色资源字段
	"role.id":          "ID",
	"role.name":        "名称",
	"role.description": "描述",
	"role.created_at":  "创建时间",
	"role.updated_at":  "更新时间",

	// 状态
	"draft":     "草稿",
	"published": "已发布",
	"archived":  "已归档",

	// 错误消息
	"error.create_failed":     "创建失败",
	"error.update_failed":     "更新失败",
	"error.delete_failed":     "删除失败",
	"error.load_failed":       "加载失败",
	"error.validation_failed": "验证失败",

	// 成功消息
	"success.create": "创建成功",
	"success.update": "更新成功",
	"success.delete": "删除成功",

	// 验证消息
	"validation.required":   "%s 为必填项",
	"validation.email":      "%s 必须是有效的邮箱地址",
	"validation.min_length": "%s 长度不能少于 %d 个字符",
	"validation.max_length": "%s 长度不能超过 %d 个字符",
	"validation.min_value":  "%s 不能小于 %d",
	"validation.max_value":  "%s 不能大于 %d",
}
