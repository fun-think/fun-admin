package service

import (
	"context"
	"fun-admin/internal/repository"
	"fun-admin/pkg/admin/i18n"
	"fun-admin/pkg/cache"
	"time"

	"go.uber.org/zap"
)

// DashboardService 仪表板服务
type DashboardService struct {
	*Service
	logger        *zap.Logger
	dashboardRepo repository.DashboardRepository
	cache         cache.CacheManager
}

// DashboardServiceInterface 仪表板服务接口
type DashboardServiceInterface interface {
	GetDashboardStats(ctx context.Context, language string) (map[string]interface{}, error)
	GetRecentUserStats(ctx context.Context, language string) ([]map[string]interface{}, error)
	GetPostStatusStats(ctx context.Context, language string) (map[string]interface{}, error)
	GetSystemInfo(ctx context.Context) (map[string]interface{}, error)
}

// NewDashboardService 创建仪表板服务
func NewDashboardService(
	service *Service,
	logger *zap.Logger,
	dashboardRepo repository.DashboardRepository,
	cache cache.CacheManager,
) DashboardServiceInterface {
	return &DashboardService{
		Service:       service,
		logger:        logger,
		dashboardRepo: dashboardRepo,
		cache:         cache,
	}
}

// GetDashboardStats 获取仪表板统计数据
func (s *DashboardService) GetDashboardStats(ctx context.Context, language string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 获取用户总数
	userCount, err := s.dashboardRepo.GetUserCount(ctx)
	if err != nil {
		return nil, err
	}
	stats["user_count"] = userCount

	// 获取部门总数 (已移除相关表)
	departmentCount, err := s.dashboardRepo.GetDepartmentCount(ctx)
	if err != nil {
		return nil, err
	}
	stats["department_count"] = departmentCount

	// 获取文章总数 (已移除相关表)
	postCount, err := s.dashboardRepo.GetPostCount(ctx)
	if err != nil {
		return nil, err
	}
	stats["post_count"] = postCount

	// 获取角色总数
	roleCount, err := s.dashboardRepo.GetRoleCount(ctx)
	if err != nil {
		return nil, err
	}
	stats["role_count"] = roleCount

	// 获取最近7天用户注册统计
	recentUsers, err := s.GetRecentUserStats(ctx, language)
	if err != nil {
		return nil, err
	}
	stats["recent_users"] = recentUsers

	// 获取文章状态统计 (已移除相关表)
	postStats, err := s.GetPostStatusStats(ctx, language)
	if err != nil {
		return nil, err
	}
	stats["post_stats"] = postStats

	// 获取系统信息
	systemInfo, err := s.GetSystemInfo(ctx)
	if err != nil {
		return nil, err
	}
	stats["system_info"] = systemInfo

	// 添加标签翻译
	stats["user_count_label"] = i18n.Translate(language, "dashboard.user_count")
	stats["department_count_label"] = i18n.Translate(language, "dashboard.department_count")
	stats["post_count_label"] = i18n.Translate(language, "dashboard.post_count")
	stats["role_count_label"] = i18n.Translate(language, "dashboard.role_count")

	return stats, nil
}

// GetRecentUserStats 获取最近用户注册统计
func (s *DashboardService) GetRecentUserStats(ctx context.Context, language string) ([]map[string]interface{}, error) {
	// 尝试从缓存获取数据
	cacheKey := "dashboard:recent_users:" + language
	if cached, err := s.cache.Get(ctx, cacheKey); err == nil && cached != nil {
		if data, ok := cached.([]map[string]interface{}); ok {
			return data, nil
		}
	}

	results, err := s.dashboardRepo.GetRecentUserStats(ctx)
	if err != nil {
		return nil, err
	}

	// 缓存数据
	s.cache.Set(ctx, cacheKey, results, time.Hour)

	return results, nil
}

// GetPostStatusStats 获取文章状态统计
func (s *DashboardService) GetPostStatusStats(ctx context.Context, language string) (map[string]interface{}, error) {
	// 尝试从缓存获取数据
	cacheKey := "dashboard:post_stats:" + language
	if cached, err := s.cache.Get(ctx, cacheKey); err == nil && cached != nil {
		if data, ok := cached.(map[string]interface{}); ok {
			return data, nil
		}
	}

	stats, err := s.dashboardRepo.GetPostStatusStats(ctx)
	if err != nil {
		return nil, err
	}

	// 添加标签翻译
	languageMap := map[string]string{
		"zh-CN": "zh-CN",
		"en":    "en-US",
	}[language]

	if languageMap == "" {
		languageMap = "zh-CN"
	}

	stats["published_label"] = i18n.Translate(languageMap, "dashboard.published")
	stats["draft_label"] = i18n.Translate(languageMap, "dashboard.draft")
	stats["archived_label"] = i18n.Translate(languageMap, "dashboard.archived")

	// 缓存数据
	s.cache.Set(ctx, cacheKey, stats, time.Hour)

	return stats, nil
}

// GetSystemInfo 获取系统信息
func (s *DashboardService) GetSystemInfo(ctx context.Context) (map[string]interface{}, error) {
	info := make(map[string]interface{})

	// 获取数据库版本
	version, err := s.dashboardRepo.GetDatabaseVersion(ctx)
	if err != nil {
		s.logger.Error("获取数据库版本失败", zap.Error(err))
	} else {
		info["database_version"] = version
	}

	// 获取数据库大小等信息
	dbSize, err := s.dashboardRepo.GetDatabaseSize(ctx)
	if err != nil {
		s.logger.Error("获取数据库大小失败", zap.Error(err))
	} else {
		info["database_size"] = dbSize
	}

	return info, nil
}