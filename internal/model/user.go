package model

// User 管理员用户模型
type User struct {
	BaseModel
	Username string `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password string `gorm:"size:255;not null" json:"-"`
	Nickname string `gorm:"size:50" json:"nickname"`
	Email    string `gorm:"size:100" json:"email"`
	Phone    string `gorm:"size:20" json:"phone"`
	Avatar   string `gorm:"size:255" json:"avatar"`
	Status   int    `gorm:"default:1" json:"status"` // 1:正常 2:禁用
}

// TableName 指定表名
func (User) TableName() string {
	return "admin_user"
}
