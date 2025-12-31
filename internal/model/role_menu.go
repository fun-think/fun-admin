package model

import (
	"time"

	"gorm.io/gorm"
)

// RoleMenu 角色菜单关联模型
type RoleMenu struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	RoleID    uint           `gorm:"not null;index" json:"role_id"`
	MenuID    uint           `gorm:"not null;index" json:"menu_id"`
}

// TableName 指定表名
func (RoleMenu) TableName() string {
	return "admin_role_menu"
}
