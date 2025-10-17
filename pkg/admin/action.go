package admin

// Action 定义操作接口
type Action interface {
	GetName() string
	GetLabel() string
	IsPrimary() bool
}

// ActionWithForm 可选接口：若实现则返回动作需要的参数表单字段
// 前端可根据字段渲染对话框表单
// 字段类型复用资源字段体系
// 注意：仅定义 schema，不含校验与执行
// 执行时通过 RunAction 的 params 传入
// 例如：重置密码动作可要求 { sendMail: boolean }
type ActionWithForm interface {
	GetFormFields() []Field
}

// ActionVisibility 可选接口：控制动作在当前上下文是否可见
// 可用于基于权限/记录状态动态隐藏动作
type ActionVisibility interface {
	IsVisible(ctx interface{}) bool
}

// BaseAction 是所有操作的基类
type BaseAction struct {
	name       string
	label      string
	primary    bool
	icon       string
	color      string
	confirm    string
	permission string
	bulk       bool
}

func (a *BaseAction) GetName() string {
	return a.name
}

func (a *BaseAction) GetLabel() string {
	return a.label
}

func (a *BaseAction) IsPrimary() bool {
	return a.primary
}

// 扩展元信息读取
func (a *BaseAction) GetIcon() string       { return a.icon }
func (a *BaseAction) GetColor() string      { return a.color }
func (a *BaseAction) GetConfirm() string    { return a.confirm }
func (a *BaseAction) GetPermission() string { return a.permission }
func (a *BaseAction) IsBulk() bool          { return a.bulk }

// 链式设置方法
func (a *BaseAction) Label(label string) *BaseAction {
	a.label = label
	return a
}

func (a *BaseAction) SetPrimary(primary bool) *BaseAction {
	a.primary = primary
	return a
}

func (a *BaseAction) Icon(icon string) *BaseAction      { a.icon = icon; return a }
func (a *BaseAction) Color(color string) *BaseAction    { a.color = color; return a }
func (a *BaseAction) Confirm(msg string) *BaseAction    { a.confirm = msg; return a }
func (a *BaseAction) Permission(key string) *BaseAction { a.permission = key; return a }
func (a *BaseAction) AsBulk() *BaseAction               { a.bulk = true; return a }

// 通用操作构造
func NewAction(name string) *BaseAction {
	return (&BaseAction{name: name, label: name}).SetPrimary(false)
}

// CreateAction 创建操作
type CreateAction struct {
	BaseAction
}

func NewCreateAction() *CreateAction {
	return &CreateAction{
		BaseAction: BaseAction{
			name:    "create",
			label:   "Create",
			primary: true,
		},
	}
}

// EditAction 编辑操作
type EditAction struct {
	BaseAction
}

func NewEditAction() *EditAction {
	return &EditAction{
		BaseAction: BaseAction{
			name:    "edit",
			label:   "Edit",
			primary: false,
		},
	}
}

// DeleteAction 删除操作
type DeleteAction struct {
	BaseAction
}

func NewDeleteAction() *DeleteAction {
	return &DeleteAction{
		BaseAction: BaseAction{
			name:    "delete",
			label:   "Delete",
			primary: false,
		},
	}
}

// ViewAction 查看操作
type ViewAction struct {
	BaseAction
}

func NewViewAction() *ViewAction {
	return &ViewAction{
		BaseAction: BaseAction{
			name:    "view",
			label:   "View",
			primary: false,
		},
	}
}

// RestoreAction 恢复软删除
type RestoreAction struct{ BaseAction }

func NewRestoreAction() *RestoreAction {
	return &RestoreAction{BaseAction: BaseAction{name: "restore", label: "Restore"}}
}

// ForceDeleteAction 强制删除（硬删）
type ForceDeleteAction struct{ BaseAction }

func NewForceDeleteAction() *ForceDeleteAction {
	return &ForceDeleteAction{BaseAction: BaseAction{name: "force_delete", label: "Force Delete"}}
}
