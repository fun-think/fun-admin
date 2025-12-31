package model

import (
	"time"

	"gorm.io/gorm"
)

// UserRole 用户角色关联模型
type UserRole struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	RoleID    uint           `gorm:"not null;index" json:"role_id"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "admin_user_role"
}
