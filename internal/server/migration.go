package server

import (
	"context"
	"encoding/json"
	"fmt"
	"fun-admin/internal/migrate"
	"fun-admin/internal/model"
	"fun-admin/pkg"
	"fun-admin/pkg/admin"
	"fun-admin/pkg/logger"
	"fun-admin/pkg/sid"
	"net/http"
	"os"

	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RoleResource 表示角色资源关系模型
type RoleResource struct {
	ID       uint   `gorm:"primarykey"`
	RoleID   uint   `gorm:"column:role_id;index:idx_role_resource,unique" json:"roleId"`
	Resource string `gorm:"column:resource;size:255;index:idx_role_resource,unique" json:"resource"`
	Action   string `gorm:"column:action;size:50" json:"action"`
}

// TableName 指定表名
func (RoleResource) TableName() string {
	return "role_resources"
}

type MigrateServer struct {
	db  *gorm.DB
	log *logger.Logger
	sid *sid.Sid
	e   *casbin.SyncedEnforcer
}

func NewMigrateServer(
	db *gorm.DB,
	log *logger.Logger,
	sid *sid.Sid,
	e *casbin.SyncedEnforcer,
) *MigrateServer {
	return &MigrateServer{
		e:   e,
		db:  db,
		log: log,
		sid: sid,
	}
}

func (m *MigrateServer) Start(ctx context.Context) error {
	m.db.Migrator().DropTable(
		&model.User{},
		&model.Menu{},
		&model.Role{},
		&model.Api{},
		&RoleResource{},
	)
	if err := m.db.AutoMigrate(
		&model.User{},
		&model.Menu{},
		&model.Role{},
		&model.Api{},
		&RoleResource{},
	); err != nil {
		m.log.Error("user migrate error", zap.Error(err))
		return err
	}

	// 迁移 admin 资源表
	err := migrate.MigrateAdminResources(m.db, admin.GlobalResourceManager)
	if err != nil {
		m.log.Error("admin resources migrate error", zap.Error(err))
		return err
	}

	err = m.initialAdminUser(ctx)
	if err != nil {
		m.log.Error("initialAdminUser error", zap.Error(err))
	}

	err = m.initialMenuData(ctx)
	if err != nil {
		m.log.Error("initialMenuData error", zap.Error(err))
	}

	err = m.initialApisData(ctx)
	if err != nil {
		m.log.Error("initialApisData error", zap.Error(err))
	}

	err = m.initialRBAC(ctx)
	if err != nil {
		m.log.Error("initialRBAC error", zap.Error(err))
	}

	m.log.Info("AutoMigrate success")
	os.Exit(0)
	return nil
}

func (m *MigrateServer) Stop(ctx context.Context) error {
	m.log.Info("AutoMigrate stop")
	return nil
}

func (m *MigrateServer) initialAdminUser(ctx context.Context) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 创建管理员用户
	err = m.db.Create(&model.User{
		ID:       1,
		Username: "admin",
		Password: string(hashedPassword),
		Nickname: "Admin",
	}).Error

	if err != nil {
		return err
	}

	// 创建普通用户
	return m.db.Create(&model.User{
		ID:       2,
		Username: "user",
		Password: string(hashedPassword),
		Nickname: "运营人员",
	}).Error
}

func (m *MigrateServer) initialRBAC(ctx context.Context) error {
	roles := []model.Role{
		{Sid: pkg.AdminRole, Name: "超级管理员"},
		{Sid: "1000", Name: "运营人员"},
		{Sid: "1001", Name: "访客"},
	}

	if err := m.db.Create(&roles).Error; err != nil {
		return err
	}

	// 获取所有菜单和API
	menuList := make([]model.Menu, 0)
	if err := m.db.Find(&menuList).Error; err != nil {
		return err
	}

	apiList := make([]model.Api, 0)
	if err := m.db.Find(&apiList).Error; err != nil {
		return err
	}

	// 为管理员角色添加所有权限
	for _, item := range menuList {
		_, err := m.e.AddPermissionForUser(pkg.AdminRole, pkg.MenuResourcePrefix+item.Path, "read")
		if err != nil {
			m.log.Sugar().Info("为角色 %s 添加权限 %s:%s 失败: %v", pkg.AdminRole, pkg.MenuResourcePrefix+item.Path, "read", err)
		} else {
			fmt.Printf("为角色 %s 添加权限: %s %s\n", pkg.AdminRole, pkg.MenuResourcePrefix+item.Path, "read")
		}
	}

	for _, api := range apiList {
		_, err := m.e.AddPermissionForUser(pkg.AdminRole, pkg.ApiResourcePrefix+api.Path, api.Method)
		if err != nil {
			m.log.Sugar().Info("为角色 %s 添加权限 %s:%s 失败: %v", pkg.AdminRole, pkg.ApiResourcePrefix+api.Path, api.Method, err)
		} else {
			fmt.Printf("为角色 %s 添加权限: %s %s\n", pkg.AdminRole, pkg.ApiResourcePrefix+api.Path, api.Method)
		}
	}

	// 添加运营人员权限
	_, err := m.e.AddRoleForUser("2", "1000")
	if err != nil {
		m.log.Error("m.e.AddRoleForUser error", zap.Error(err))
		return err
	}

	// 为运营人员添加基础权限
	basicPermissions := []struct {
		resource string
		action   string
	}{
		{pkg.MenuResourcePrefix + "/profile/basic", "read"},
		{pkg.MenuResourcePrefix + "/profile/advanced", "read"},
		{pkg.MenuResourcePrefix + "/profile", "read"},
		{pkg.MenuResourcePrefix + "/dashboard", "read"},
		{pkg.MenuResourcePrefix + "/dashboard/workplace", "read"},
		{pkg.MenuResourcePrefix + "/dashboard/analysis", "read"},
		{pkg.MenuResourcePrefix + "/account/settings", "read"},
		{pkg.MenuResourcePrefix + "/account/center", "read"},
		{pkg.MenuResourcePrefix + "/account", "read"},
		{pkg.ApiResourcePrefix + "/v1/menus", http.MethodGet},
		{pkg.ApiResourcePrefix + "/v1/admin/user", http.MethodGet},
	}

	for _, perm := range basicPermissions {
		_, err := m.e.AddPermissionForUser("1000", perm.resource, perm.action)
		if err != nil {
			m.log.Sugar().Info("为角色 %s 添加权限 %s:%s 失败: %v", "1000", perm.resource, perm.action, err)
		} else {
			fmt.Printf("为角色 %s 添加权限: %s %s\n", "1000", perm.resource, perm.action)
		}
	}

	// 保存策略到数据库
	return m.e.SavePolicy()
}

func (m *MigrateServer) initialApisData(ctx context.Context) error {
	initialApis := []model.Api{
		{Group: "基础API", Name: "获取用户菜单列表", Path: "/v1/menus", Method: http.MethodGet},
		{Group: "基础API", Name: "获取管理员信息", Path: "/v1/admin/user", Method: http.MethodGet},

		{Group: "菜单管理", Name: "获取管理菜单", Path: "/v1/admin/menus", Method: http.MethodGet},
		{Group: "菜单管理", Name: "创建菜单", Path: "/v1/admin/menu", Method: http.MethodPost},
		{Group: "菜单管理", Name: "更新菜单", Path: "/v1/admin/menu", Method: http.MethodPut},
		{Group: "菜单管理", Name: "删除菜单", Path: "/v1/admin/menu", Method: http.MethodDelete},

		{Group: "权限模块", Name: "获取用户权限", Path: "/v1/admin/user/permissions", Method: http.MethodGet},
		{Group: "权限模块", Name: "获取角色权限", Path: "/v1/admin/role/permissions", Method: http.MethodGet},
		{Group: "权限模块", Name: "更新角色权限", Path: "/v1/admin/role/permission", Method: http.MethodPut},
		{Group: "权限模块", Name: "获取角色列表", Path: "/v1/admin/roles", Method: http.MethodGet},
		{Group: "权限模块", Name: "创建角色", Path: "/v1/admin/role", Method: http.MethodPost},
		{Group: "权限模块", Name: "更新角色", Path: "/v1/admin/role", Method: http.MethodPut},
		{Group: "权限模块", Name: "删除角色", Path: "/v1/admin/role", Method: http.MethodDelete},

		{Group: "权限模块", Name: "获取管理员列表", Path: "/v1/admin/users", Method: http.MethodGet},
		{Group: "权限模块", Name: "更新管理员信息", Path: "/v1/admin/user", Method: http.MethodPut},
		{Group: "权限模块", Name: "创建管理员账号", Path: "/v1/admin/user", Method: http.MethodPost},
		{Group: "权限模块", Name: "删除管理员", Path: "/v1/admin/user", Method: http.MethodDelete},

		{Group: "权限模块", Name: "获取API列表", Path: "/v1/admin/apis", Method: http.MethodGet},
		{Group: "权限模块", Name: "创建API", Path: "/v1/admin/api", Method: http.MethodPost},
		{Group: "权限模块", Name: "更新API", Path: "/v1/admin/api", Method: http.MethodPut},
		{Group: "权限模块", Name: "删除API", Path: "/v1/admin/api", Method: http.MethodDelete},
	}

	return m.db.Create(&initialApis).Error
}

func (m *MigrateServer) initialMenuData(ctx context.Context) error {
	menuList := make([]model.Menu, 0)
	err := json.Unmarshal([]byte(menuData), &menuList)
	if err != nil {
		m.log.Error("json.Unmarshal error", zap.Error(err))
		return err
	}

	return m.db.Create(&menuList).Error
}

var menuData = `[
  {
    "id": 1,
    "parentId": 0,
    "title": "仪表盘",
    "icon": "DashboardOutlined",
    "component": "RouteView",
    "redirect": "/dashboard/analysis",
    "path": "/dashboard",
    "name": "Dashboard",
    "locale": "menu.dashboard"
  },
  {
    "id": 2,
    "parentId": 0,
    "title": "分析页",
    "icon": "DashboardOutlined",
    "component": "/dashboard/analysis",
    "path": "/dashboard/analysis",
    "name": "DashboardAnalysis",
    "keepAlive": true,
    "locale": "menu.dashboard.analysis",
    "weight": 2
  },
  {
    "id": 15,
    "path": "/access",
    "component": "RouteView",
    "redirect": "/access/common",
    "title": "权限模块",
    "name": "Access",
    "parentId": 0,
    "icon": "ClusterOutlined",
    "locale": "menu.access",
    "weight": 1
  },
  {
    "id": 18,
    "parentId": 15,
    "path": "/access/admin",	
    "title": "管理员账号",
    "name": "accessAdmin",
    "component": "/access/admin",
    "locale": "menu.access.admin"
  },
  {
    "id": 51,
    "parentId": 15,
    "path": "/access/role",	
    "title": "角色管理",
    "name": "AccessRoles",
    "component": "/access/role",
    "locale": "menu.access.roles"
  },
  {
    "id": 52,
    "parentId": 15,
    "path": "/access/menu",	
    "title": "菜单管理",
    "name": "AccessMenu",
    "component": "/access/menu",
    "locale": "menu.access.menus"
  },
  {
    "id": 53,
    "parentId": 15,
    "path": "/access/api",	
    "title": "API管理",
    "name": "AccessAPI",
    "component": "/access/api",
    "locale": "menu.access.api"
  }
]`
