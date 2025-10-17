package admin

// Schema 用于定义数据库表结构
type Schema struct{}

// NewSchema 创建一个新的 Schema 实例
func NewSchema() *Schema {
	return &Schema{}
}

// CreateTable 创建新表
func (s *Schema) CreateTable(name string) *TableBuilder {
	return &TableBuilder{
		name:   name,
		fields: []Field{},
	}
}

// Table 获取现有表
func (s *Schema) Table(name string) *TableBuilder {
	return &TableBuilder{
		name:   name,
		fields: []Field{},
	}
}

// TableBuilder 表构建器
type TableBuilder struct {
	name   string
	fields []Field
}

// String 添加字符串字段
func (tb *TableBuilder) String(name string) *TableBuilder {
	field := NewTextField(name)
	tb.fields = append(tb.fields, field)
	return tb
}

// Text 添加文本字段
func (tb *TableBuilder) Text(name string) *TableBuilder {
	field := NewTextField(name)
	field.fieldType = "textarea"
	tb.fields = append(tb.fields, field)
	return tb
}

// Integer 添加整数字段
func (tb *TableBuilder) Integer(name string) *TableBuilder {
	field := NewNumberField(name)
	tb.fields = append(tb.fields, field)
	return tb
}

// Boolean 添加布尔字段
func (tb *TableBuilder) Boolean(name string) *TableBuilder {
	field := NewBooleanField(name)
	tb.fields = append(tb.fields, field)
	return tb
}

// GetFields 返回所有字段
func (tb *TableBuilder) GetFields() []Field {
	return tb.fields
}

// GetTableName 返回表名
func (tb *TableBuilder) GetTableName() string {
	return tb.name
}
