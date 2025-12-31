package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// SMSCode 短信验证码模型
type SMSCode struct {
	BaseModel
	Mobile      string    `gorm:"size:20;not null;index" json:"mobile"`      // 手机号
	Code        string    `gorm:"size:10;not null" json:"code"`              // 验证码
	ExpiredAt   time.Time `gorm:"not null" json:"expired_at"`                // 过期时间
	IsUsed      bool      `gorm:"default:false;not null" json:"is_used"`     // 是否已使用
}