package model

import (
	"time"

	"gorm.io/gorm"
)

// User 管理员用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Username  string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"size:255;not null" json:"password"`
	Nickname  string         `gorm:"size:50" json:"nickname"`
	Email     string         `gorm:"size:100" json:"email"`
	Phone     string         `gorm:"size:20" json:"phone"`
	Status    int            `gorm:"default:1" json:"status"` // 1:正常 2:禁用
}

// TableName 指定表名
func (User) TableName() string {
	return "admin_user"
}
