package pkg

// PermSep 权限分隔符
const PermSep = ","

// RBAC/Casbin 资源前缀与固定常量
const (
	MenuResourcePrefix = "menu:"
	ApiResourcePrefix  = "api:"
	AdminUserID        = "1"     // 超管用户ID(字符串形式用于比较)
	AdminRole          = "admin" // 超管角色SID
)
