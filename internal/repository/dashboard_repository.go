package repository

import (
	"context"
	"fun-admin/pkg/logger"

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
	err := r.db.WithContext(ctx).Table("users").Count(&count).Error
	return count, err
}

// GetDepartmentCount 获取部门总数
func (r *dashboardRepository) GetDepartmentCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("departments").Count(&count).Error
	return count, err
}

// GetPostCount 获取文章总数
func (r *dashboardRepository) GetPostCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("posts").Count(&count).Error
	return count, err
}

// GetRoleCount 获取角色总数
func (r *dashboardRepository) GetRoleCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("roles").Count(&count).Error
	return count, err
}

// GetRecentUserStats 获取最近用户注册统计
func (r *dashboardRepository) GetRecentUserStats(ctx context.Context) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as count
		FROM users 
		WHERE created_at >= DATE('now', '-7 days')
		GROUP BY DATE(created_at)
		ORDER BY date
	`

	err := r.db.WithContext(ctx).Raw(query).Scan(&results).Error
	return results, err
}

// GetPostStatusStats 获取文章状态统计
func (r *dashboardRepository) GetPostStatusStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var results []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}

	query := `
		SELECT 
			status,
			COUNT(*) as count
		FROM posts 
		GROUP BY status
	`

	err := r.db.WithContext(ctx).Raw(query).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		stats[result.Status] = result.Count
	}

	// 确保所有状态都有数据
	if _, ok := stats["draft"]; !ok {
		stats["draft"] = 0
	}
	if _, ok := stats["published"]; !ok {
		stats["published"] = 0
	}
	if _, ok := stats["archived"]; !ok {
		stats["archived"] = 0
	}

	return stats, nil
}

// GetDatabaseVersion 获取数据库版本
func (r *dashboardRepository) GetDatabaseVersion(ctx context.Context) (string, error) {
	var version string
	err := r.db.WithContext(ctx).Raw("SELECT sqlite_version()").Scan(&version).Error
	return version, err
}

// GetDatabaseSize 获取数据库大小
func (r *dashboardRepository) GetDatabaseSize(ctx context.Context) (int64, error) {
	var dbSize int64
	err := r.db.WithContext(ctx).Raw("SELECT page_count * page_size as size FROM pragma_page_count(), pragma_page_size()").Scan(&dbSize).Error
	return dbSize, err
}
