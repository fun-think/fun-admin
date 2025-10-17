package migrate

import (
	"fun-admin/pkg/admin"

	"gorm.io/gorm"
)

// MigrateAdminResources 创建 admin 资源相关的数据表
func MigrateAdminResources(db *gorm.DB, resourceManager *admin.ResourceManager) error {
	// 获取所有注册的资源
	resources := resourceManager.GetResources()

	// 为每个资源创建数据表
	for _, resource := range resources {
		err := createResourceTable(db, resource)
		if err != nil {
			return err
		}
	}

	return nil
}

// createResourceTable 为资源创建数据表
func createResourceTable(db *gorm.DB, resource admin.Resource) error {
	tableName := resource.GetSlug()
	if tableName == "" {
		tableName = lowercase(resource.GetTitle())
	}

	// 构建创建表的 SQL 语句
	// 这里简化处理，实际项目中应该根据字段类型创建合适的列

	sql := "CREATE TABLE IF NOT EXISTS " + tableName + " ("
	sql += "id INTEGER PRIMARY KEY AUTOINCREMENT, "
	sql += "created_at DATETIME, "
	sql += "updated_at DATETIME, "
	sql += "deleted_at DATETIME"

	// 根据字段定义添加列
	fields := resource.GetFields()
	for _, field := range fields {
		columnName := field.GetName()
		columnType := getColumnType(field.GetType())

		sql += ", " + columnName + " " + columnType

		// 添加 NOT NULL 约束
		if field.IsRequired() {
			sql += " NOT NULL"
		}
	}

	sql += ")"

	// 执行 SQL
	return db.Exec(sql).Error
}

// getColumnType 根据字段类型获取数据库列类型
func getColumnType(fieldType string) string {
	switch fieldType {
	case "text", "email", "select", "textarea":
		return "TEXT"
	case "number":
		return "INTEGER"
	case "boolean":
		return "BOOLEAN"
	case "date":
		return "DATE"
	case "datetime":
		return "DATETIME"
	case "relationship":
		return "INTEGER"
	default:
		return "TEXT"
	}
}

// lowercase 将字符串转换为小写
func lowercase(s string) string {
	result := make([]rune, len(s))
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			result[i] = r + ('a' - 'A')
		} else {
			result[i] = r
		}
	}
	return string(result)
}

// MigrateAdminResourcesGORM 使用 GORM 的方式创建 admin 资源相关的数据表
func MigrateAdminResourcesGORM(db *gorm.DB, resourceManager *admin.ResourceManager) error {
	// 获取所有注册的资源
	resources := resourceManager.GetResources()

	// 为每个资源创建数据表
	for _, resource := range resources {
		err := createResourceTableGORM(db, resource)
		if err != nil {
			return err
		}
	}

	return nil
}

// dynamicModelWithTableName 动态模型类型，TableName 由闭包返回
type dynamicModelWithTableName struct {
	ID uint `gorm:"primaryKey"`
}

func (m dynamicModelWithTableName) TableName() string { return mTableName }

var mTableName string

// createResourceTableGORM 使用 GORM 为资源创建数据表
func createResourceTableGORM(db *gorm.DB, resource admin.Resource) error {
	tableName := resource.GetSlug()
	if tableName == "" {
		tableName = lowercase(resource.GetTitle())
	}
	mTableName = tableName

	// 自动迁移表结构（仅示意）
	return db.AutoMigrate(&dynamicModelWithTableName{})
}
