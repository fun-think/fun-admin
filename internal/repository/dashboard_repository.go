package repository

import (
	"context"
	"fun-admin/internal/model"
	"fun-admin/pkg/logger"
	"strings"

	"gorm.io/gorm"
)

// DashboardRepository 仪表板仓库接口
type DashboardRepository interface {
	GetUserCount(ctx context.Context) (int64, error)
	GetDepartmentCount(ctx context.Context) (int64, error)
	GetPostCount(ctx context.Context) (int64, error)
	GetRoleCount(ctx context.Context) (int64, error)
	GetRecentUserStats(ctx context.Context) ([]map[string]interface{}, error)
	GetPostStatusStats(ctx context.Context) (map[string]interface{}, error)
	GetDatabaseVersion(ctx context.Context) (string, error)
	GetDatabaseSize(ctx context.Context) (int64, error)
}

// dashboardRepository 仪表板仓库实现
type dashboardRepository struct {
	logger *logger.Logger
	db     *gorm.DB
}

// NewDashboardRepository 创建仪表板仓库
func NewDashboardRepository(
	logger *logger.Logger,
	db *gorm.DB,
) DashboardRepository {
	return &dashboardRepository{
		logger: logger,
		db:     db,
	}
}

// GetUserCount 获取用户总数
func (r *dashboardRepository) GetUserCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Count(&count).Error
	return count, err
}

// GetDepartmentCount 获取部门总数
func (r *dashboardRepository) GetDepartmentCount(ctx context.Context) (int64, error) {
	// departments 表在 internal/model 中不存在，删除相关代码
	// 返回 0 和 nil 错误表示没有这个功能
	return 0, nil
}

// GetPostCount 获取文章总数
func (r *dashboardRepository) GetPostCount(ctx context.Context) (int64, error) {
	// posts 表在 internal/model 中不存在，删除相关代码
	// 返回 0 和 nil 错误表示没有这个功能
	return 0, nil
}

// GetRoleCount 获取角色总数
func (r *dashboardRepository) GetRoleCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Role{}).Count(&count).Error
	return count, err
}

// GetRecentUserStats 获取最近用户注册统计
func (r *dashboardRepository) GetRecentUserStats(ctx context.Context) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// 根据数据库类型使用不同的SQL
	var query string
	dbType := r.db.Dialector.Name()
	
	switch dbType {
	case "sqlite":
		query = `
			SELECT 
				DATE(created_at) as date,
				COUNT(*) as count
			FROM ` + (&model.User{}).TableName() + `
			WHERE created_at >= DATE('now', '-7 days')
			GROUP BY DATE(created_at)
			ORDER BY date
		`
	case "mysql":
		query = `
			SELECT 
				DATE(created_at) as date,
				COUNT(*) as count
			FROM ` + (&model.User{}).TableName() + `
			WHERE created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)
			GROUP BY DATE(created_at)
			ORDER BY date
		`
	case "postgres":
		query = `
			SELECT 
				DATE(created_at) as date,
				COUNT(*) as count
			FROM ` + (&model.User{}).TableName() + `
			WHERE created_at >= CURRENT_DATE - INTERVAL '7 days'
			GROUP BY DATE(created_at)
			ORDER BY date
		`
	default:
		// 默认使用通用SQL（可能不适用于所有数据库）
		query = `
			SELECT 
				DATE(created_at) as date,
				COUNT(*) as count
			FROM ` + (&model.User{}).TableName() + `
			WHERE created_at >= CURRENT_DATE - INTERVAL 7 DAY
			GROUP BY DATE(created_at)
			ORDER BY date
		`
	}

	err := r.db.WithContext(ctx).Raw(query).Scan(&results).Error
	return results, err
}

// GetPostStatusStats 获取文章状态统计
func (r *dashboardRepository) GetPostStatusStats(ctx context.Context) (map[string]interface{}, error) {
	// posts 表在 internal/model 中不存在，删除相关代码
	// 返回空 map 和 nil 错误表示没有这个功能
	stats := make(map[string]interface{})
	return stats, nil
}

// GetDatabaseVersion 获取数据库版本
func (r *dashboardRepository) GetDatabaseVersion(ctx context.Context) (string, error) {
	var version string
	
	// 根据数据库类型使用不同的SQL
	dbType := r.db.Dialector.Name()
	
	switch dbType {
	case "sqlite":
		err := r.db.WithContext(ctx).Raw("SELECT sqlite_version()").Scan(&version).Error
		return version, err
	case "mysql":
		err := r.db.WithContext(ctx).Raw("SELECT VERSION()").Scan(&version).Error
		return version, err
	case "postgres":
		err := r.db.WithContext(ctx).Raw("SELECT version()").Scan(&version).Error
		// 只返回版本号部分
		if err == nil && strings.Contains(version, " ") {
			parts := strings.Split(version, " ")
			if len(parts) > 0 {
				version = parts[0]
			}
		}
		return version, err
	default:
		// 默认尝试通用方式
		err := r.db.WithContext(ctx).Raw("SELECT VERSION()").Scan(&version).Error
		return version, err
	}
}

// GetDatabaseSize 获取数据库大小
func (r *dashboardRepository) GetDatabaseSize(ctx context.Context) (int64, error) {
	var dbSize int64
	
	// 根据数据库类型使用不同的SQL
	dbType := r.db.Dialector.Name()
	
	switch dbType {
	case "sqlite":
		err := r.db.WithContext(ctx).Raw("SELECT page_count * page_size as size FROM pragma_page_count(), pragma_page_size()").Scan(&dbSize).Error
		return dbSize, err
	case "mysql":
		// MySQL需要查询information_schema来获取数据库大小
		err := r.db.WithContext(ctx).Raw(`
			SELECT SUM(data_length + index_length) as size 
			FROM information_schema.tables 
			WHERE table_schema = DATABASE()
		`).Scan(&dbSize).Error
		return dbSize, err
	case "postgres":
		// PostgreSQL需要查询pg_database_size函数
		err := r.db.WithContext(ctx).Raw(`
			SELECT pg_database_size(current_database()) as size
		`).Scan(&dbSize).Error
		return dbSize, err
	default:
		// 对于不支持的数据库，返回0
		return 0, nil
	}
}