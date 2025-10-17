package model

import (
	"time"

	"gorm.io/gorm"
)

// DictionaryType 字典类型
type DictionaryType struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Name      string         `gorm:"size:100;not null;uniqueIndex" json:"name"` // 字典名称
	Code      string         `gorm:"size:50;not null;uniqueIndex" json:"code"`  // 字典编码
	Status    int            `gorm:"default:1;comment:1-启用,0-禁用" json:"status"` // 状态
	Remark    string         `gorm:"size:500" json:"remark"`                    // 备注
	Sort      int            `gorm:"default:0" json:"sort"`                     // 排序
}

// TableName 指定表名
func (DictionaryType) TableName() string {
	return "admin_dictionary_type"
}
