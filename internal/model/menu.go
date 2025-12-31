package model

// Menu 菜单模型
type Menu struct {
	BaseModel
	ParentID   uint   `gorm:"default:0" json:"parentId"`
	Path       string `gorm:"size:255" json:"path"`
	Title      string `gorm:"size:100;not null" json:"title"`
	Name       string `gorm:"size:100" json:"name"`
	Component  string `gorm:"size:100" json:"component"`
	Locale     string `gorm:"size:100" json:"locale"`
	Weight     int    `gorm:"default:0" json:"weight"`
	Icon       string `gorm:"size:50" json:"icon"`
	Redirect   string `gorm:"size:255" json:"redirect"`
	URL        string `gorm:"size:255" json:"url"`
	KeepAlive  bool   `gorm:"default:false" json:"keepAlive"`
	HideInMenu bool   `gorm:"default:false" json:"hideInMenu"`
	Type       string `gorm:"size:20;default:menu" json:"type"` // menu, button, api
	Permission string `gorm:"size:100" json:"permission"`       // 权限标识
}

// TableName 指定表名
func (Menu) TableName() string {
	return "admin_menu"
}
