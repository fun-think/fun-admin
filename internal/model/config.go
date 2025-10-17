package model

import (
	"time"

	"gorm.io/gorm"
)

// Config 系统设置模型
type Config struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Key       string         `gorm:"size:100;uniqueIndex;not null" json:"key"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	Value     string         `gorm:"type:text" json:"value"`
	Type      string         `gorm:"size:50;not null" json:"type"`  // 设置类型
	Group     string         `gorm:"size:50;not null" json:"group"` // 设置分组
	Sort      int            `gorm:"default:0" json:"sort"`         // 排序
	Remark    string         `gorm:"size:255" json:"remark"`        // 备注
}

// TableName 指定表名
func (Config) TableName() string {
	return "admin_config"
}
