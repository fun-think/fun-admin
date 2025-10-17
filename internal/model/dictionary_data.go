package model

import (
	"time"

	"gorm.io/gorm"
)

// DictionaryData 字典数据
type DictionaryData struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	TypeID    uint           `gorm:"not null;index" json:"type_id"`             // 字典类型ID
	Type      DictionaryType `gorm:"foreignKey:TypeID" json:"type"`             // 字典类型
	Label     string         `gorm:"size:100;not null" json:"label"`            // 字典标签
	Value     string         `gorm:"size:100;not null" json:"value"`            // 字典值
	Status    int            `gorm:"default:1;comment:1-启用,0-禁用" json:"status"` // 状态
	Remark    string         `gorm:"size:500" json:"remark"`                    // 备注
	Sort      int            `gorm:"default:0" json:"sort"`                     // 排序
	IsDefault bool           `gorm:"default:false" json:"is_default"`           // 是否默认值
	Ext1      string         `gorm:"size:255" json:"ext1"`                      // 扩展字段1
	Ext2      string         `gorm:"size:255" json:"ext2"`                      // 扩展字段2
	Ext3      string         `gorm:"size:255" json:"ext3"`                      // 扩展字段3
}

// TableName 指定表名
func (DictionaryData) TableName() string {
	return "admin_dictionary_data"
}
