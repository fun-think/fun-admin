package model

// Role 角色模型
type Role struct {
	BaseModel
	Sid         string `gorm:"size:50;uniqueIndex;not null" json:"sid"` // 角色唯一标识
	Name        string `gorm:"size:50;not null" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	Status      int    `gorm:"default:1" json:"status"` // 1:正常 2:禁用
}

// TableName 指定表名
func (Role) TableName() string {
	return "admin_role"
}
