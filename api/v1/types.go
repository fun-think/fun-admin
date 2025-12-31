package v1

type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"1234@gmail.com"`
	Password string `json:"password" binding:"required" example:"123456"`
	Type     string `json:"type" example:"account"` // 登录类型: account, mobile
	Mobile   string `json:"mobile" example:"13800138000"`
	Code     string `json:"code" example:"123456"`
}

type LoginResponse struct {
	Response
	Data LoginResponseData
}

type GetUsersRequest struct {
	Page     int    `form:"page" binding:"required" example:"1"`
	PageSize int    `form:"pageSize" binding:"required" example:"10"`
	ID       uint   `form:"id" binding:"" example:"1"`
	Username string `json:"username" binding:"" example:"张三"`
	Nickname string `json:"nickname" binding:"" example:"小Baby"`
	Phone    string `form:"phone" binding:"" example:"1858888888"`
	Email    string `form:"email" binding:"" example:"1234@gmail.com"`
}

type UserCreateRequest struct {
	Username string   `json:"username" binding:"required" example:"张三"`
	Nickname string   `json:"nickname" binding:"" example:"小Baby"`
	Password string   `json:"password" binding:"required" example:"123456"`
	Email    string   `json:"email" binding:"" example:"1234@gmail.com"`
	Phone    string   `form:"phone" binding:"" example:"1858888888"`
	Roles    []string `json:"roles" example:""`
}

type UserUpdateRequest struct {
	ID       uint     `json:"id"`
	Username string   `json:"username" binding:"required" example:"张三"`
	Nickname string   `json:"nickname" binding:"" example:"小Baby"`
	Password string   `json:"password" binding:"" example:"123456"`
	Email    string   `json:"email" binding:"" example:"1234@gmail.com"`
	Phone    string   `form:"phone" binding:"" example:"1858888888"`
	Roles    []string `json:"roles" example:""`
}

type UserDeleteRequest struct {
	ID uint `form:"id" binding:"required" example:"1"`
}

type MenuDataItem struct {
	ID         uint   `json:"id,omitempty"`         // 唯一id，使用整数表示
	ParentID   uint   `json:"parentId,omitempty"`   // 父级菜单的id，使用整数表示
	Weight     int    `json:"weight"`               // 排序权重
	Path       string `json:"path"`                 // 地址
	Title      string `json:"title"`                // 展示名称
	Name       string `json:"name,omitempty"`       // 同路由中的name，唯一标识
	Component  string `json:"component,omitempty"`  // 绑定的组件
	Locale     string `json:"locale,omitempty"`     // 本地化标识
	Icon       string `json:"icon,omitempty"`       // 图标，使用字符串表示
	Redirect   string `json:"redirect,omitempty"`   // 重定向地址
	KeepAlive  bool   `json:"keepAlive,omitempty"`  // 是否保活
	HideInMenu bool   `json:"hideInMenu,omitempty"` // 是否保活
	URL        string `json:"url,omitempty"`        // iframe模式下的跳转url，不能与path重复
	UpdatedAt  string `json:"updatedAt,omitempty"`  // 是否保活
}

type GetMenuResponseData struct {
	List []MenuDataItem `json:"list"`
}

type GetMenuResponse struct {
	Response
	Data GetMenuResponseData
}

type MenuCreateRequest struct {
	ParentID   uint   `json:"parentId,omitempty"`   // 父级菜单的id，使用整数表示
	Weight     int    `json:"weight"`               // 排序权重
	Path       string `json:"path"`                 // 地址
	Title      string `json:"title"`                // 展示名称
	Name       string `json:"name,omitempty"`       // 同路由中的name，唯一标识
	Component  string `json:"component,omitempty"`  // 绑定的组件
	Locale     string `json:"locale,omitempty"`     // 本地化标识
	Icon       string `json:"icon,omitempty"`       // 图标，使用字符串表示
	Redirect   string `json:"redirect,omitempty"`   // 重定向地址
	KeepAlive  bool   `json:"keepAlive,omitempty"`  // 是否保活
	HideInMenu bool   `json:"hideInMenu,omitempty"` // 是否保活
	URL        string `json:"url,omitempty"`        // iframe模式下的跳转url，不能与path重复

}

type MenuUpdateRequest struct {
	ID         uint   `json:"id,omitempty"`         // 唯一id，使用整数表示
	ParentID   uint   `json:"parentId,omitempty"`   // 父级菜单的id，使用整数表示
	Weight     int    `json:"weight"`               // 排序权重
	Path       string `json:"path"`                 // 地址
	Title      string `json:"title"`                // 展示名称
	Name       string `json:"name,omitempty"`       // 同路由中的name，唯一标识
	Component  string `json:"component,omitempty"`  // 绑定的组件
	Locale     string `json:"locale,omitempty"`     // 本地化标识
	Icon       string `json:"icon,omitempty"`       // 图标，使用字符串表示
	Redirect   string `json:"redirect,omitempty"`   // 重定向地址
	KeepAlive  bool   `json:"keepAlive,omitempty"`  // 是否保活
	HideInMenu bool   `json:"hideInMenu,omitempty"` // 是否保活
	URL        string `json:"url,omitempty"`        // iframe模式下的跳转url，不能与path重复
	UpdatedAt  string `json:"updatedAt"`
}

type MenuDeleteRequest struct {
	ID uint `form:"id"` // 唯一id，使用整数表示
}

type GetRoleListRequest struct {
	Page     int    `form:"page" binding:"required" example:"1"`
	PageSize int    `form:"pageSize" binding:"required" example:"10"`
	Sid      string `form:"sid" binding:"" example:"1"`
	Name     string `form:"name" binding:"" example:""`
}

type RoleDataItem struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Sid       string `json:"sid"`
	UpdatedAt string `json:"updatedAt"`
	CreatedAt string `json:"createdAt"`
}

type GetRolesResponseData struct {
	List  []RoleDataItem `json:"list"`
	Total int64          `json:"total"`
}

type GetRolesResponse struct {
	Response
	Data GetRolesResponseData
}

type RoleCreateRequest struct {
	Sid  string `form:"sid" binding:"required" example:"1"`
	Name string `form:"name" binding:"required" example:""`
}

type RoleUpdateRequest struct {
	ID   uint   `form:"id" binding:"required" example:"1"`
	Sid  string `form:"sid" binding:"required" example:"1"`
	Name string `form:"name" binding:"required" example:""`
}

type RoleDeleteRequest struct {
	ID uint `form:"id" binding:"required" example:"1"`
}

type PermissionCreateRequest struct {
	Sid  string `form:"sid" binding:"required" example:"1"`
	Name string `form:"name" binding:"required" example:""`
}

type GetApisRequest struct {
	Page     int    `form:"page" binding:"required" example:"1"`
	PageSize int    `form:"pageSize" binding:"required" example:"10"`
	Group    string `form:"group" binding:"" example:"权限管理"`
	Name     string `form:"name" binding:"" example:"菜单列表"`
	Path     string `form:"path" binding:"" example:"/v1/test"`
	Method   string `form:"method" binding:"" example:"GET"`
}

type ApiDataItem struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Method    string `json:"method"`
	Group     string `json:"group"`
	UpdatedAt string `json:"updatedAt"`
	CreatedAt string `json:"createdAt"`
}

type GetApisResponseData struct {
	List   []ApiDataItem `json:"list"`
	Total  int64         `json:"total"`
	Groups []string      `json:"groups"`
}

type GetApisResponse struct {
	Response
	Data GetApisResponseData
}

type ApiCreateRequest struct {
	Group  string `form:"group" binding:"" example:"权限管理"`
	Name   string `form:"name" binding:"" example:"菜单列表"`
	Path   string `form:"path" binding:"" example:"/v1/test"`
	Method string `form:"method" binding:"" example:"GET"`
}

type ApiUpdateRequest struct {
	ID     uint   `form:"id" binding:"required" example:"1"`
	Group  string `form:"group" binding:"" example:"权限管理"`
	Name   string `form:"name" binding:"" example:"菜单列表"`
	Path   string `form:"path" binding:"" example:"/v1/test"`
	Method string `form:"method" binding:"" example:"GET"`
}

type ApiDeleteRequest struct {
	ID uint `form:"id" binding:"required" example:"1"`
}

type GetUserPermissionsData struct {
	List []string `json:"list"`
}

type GetRolePermissionsRequest struct {
	Role string `form:"role" binding:"required" example:"admin"`
}

type GetRolePermissionsData struct {
	List []string `json:"list"`
}

type UpdateRolePermissionRequest struct {
	Role string   `form:"role" binding:"required" example:"admin"`
	List []string `form:"list" binding:"required" example:""`
}
