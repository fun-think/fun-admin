package model

import (
	"time"

	"gorm.io/gorm"
)

// Menu 菜单模型
type Menu struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	ParentID   uint           `gorm:"default:0" json:"parentId"`
	Path       string         `gorm:"size:255" json:"path"`
	Title      string         `gorm:"size:100;not null" json:"title"`
	Name       string         `gorm:"size:100" json:"name"`
	Component  string         `gorm:"size:100" json:"component"`
	Locale     string         `gorm:"size:100" json:"locale"`
	Weight     int            `gorm:"default:0" json:"weight"`
	Icon       string         `gorm:"size:50" json:"icon"`
	Redirect   string         `gorm:"size:255" json:"redirect"`
	URL        string         `gorm:"size:255" json:"url"`
	KeepAlive  bool           `gorm:"default:false" json:"keepAlive"`
	HideInMenu bool           `gorm:"default:false" json:"hideInMenu"`
}

// TableName 指定表名
func (Menu) TableName() string {
	return "admin_menu"
}
