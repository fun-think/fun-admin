package model

// CrudItem keeps simple key/value data for the demo CRUD table.
type CrudItem struct {
	BaseModel
	Name   string `gorm:"size:255;not null" json:"name"`
	Value  string `gorm:"size:1024;not null" json:"value"`
	Remark string `gorm:"size:2048" json:"remark,omitempty"`
}

// TableName explicitly binds the model to the crud_items table.
func (CrudItem) TableName() string {
	return "crud_items"
}
