package model

import (
	"time"

	"gorm.io/gorm"
)

// Api API模型
type Api struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Group     string         `gorm:"size:50;not null" json:"group"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Path      string         `gorm:"size:255;not null" json:"path"`
	Method    string         `gorm:"size:10;not null" json:"method"`
}

// TableName 指定表名
func (Api) TableName() string {
	return "admin_api"
}
