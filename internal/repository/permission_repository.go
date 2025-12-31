package repository

import (
	"context"
	"fun-admin/pkg/logger"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	GetRolesForUser(ctx context.Context, userId string) ([]string, error)
	GetFilteredPolicy(ctx context.Context, fieldIndex int, fieldValues ...string) ([][]string, error)
	RemoveFilteredPolicy(ctx context.Context, fieldIndex int, fieldValues ...string) (bool, error)
	AddPolicy(ctx context.Context, params ...string) (bool, error)
	DeleteRole(ctx context.Context, role string) (bool, error)
	SavePolicy(ctx context.Context) error

	// 新增的方法
	AddRoleForUser(ctx context.Context, user string, role string) (bool, error)
	DeleteRoleForUser(ctx context.Context, user string, role string) (bool, error)
	GetPermissionsForUser(ctx context.Context, user string) ([][]string, error)
	GetAllRoles(ctx context.Context) ([]string, error)
	AddPermissionForUser(ctx context.Context, user string, permission ...string) (bool, error)
	DeletePermissionForUser(ctx context.Context, user string, permission ...string) (bool, error)
	GetUsersForRole(ctx context.Context, role string) ([]string, error) // 添加此方法
}

func NewPermissionRepository(
	logger *logger.Logger,
	db *gorm.DB,
	enforcer *casbin.SyncedEnforcer,
) PermissionRepository {
	return &permissionRepository{
		logger:   logger,
		db:       db,
		enforcer: enforcer,
	}
}

type permissionRepository struct {
	logger   *logger.Logger
	db       *gorm.DB
	enforcer *casbin.SyncedEnforcer
}

func (r *permissionRepository) GetRolesForUser(ctx context.Context, userId string) ([]string, error) {
	return r.enforcer.GetRolesForUser(userId)
}

func (r *permissionRepository) GetFilteredPolicy(ctx context.Context, fieldIndex int, fieldValues ...string) ([][]string, error) {
	return r.enforcer.GetFilteredPolicy(fieldIndex, fieldValues...)
}

func (r *permissionRepository) RemoveFilteredPolicy(ctx context.Context, fieldIndex int, fieldValues ...string) (bool, error) {
	return r.enforcer.RemoveFilteredPolicy(fieldIndex, fieldValues...)
}

func (r *permissionRepository) AddPolicy(ctx context.Context, params ...string) (bool, error) {
	return r.enforcer.AddPolicy(params)
}

func (r *permissionRepository) DeleteRole(ctx context.Context, role string) (bool, error) {
	return r.enforcer.DeleteRole(role)
}

func (r *permissionRepository) SavePolicy(ctx context.Context) error {
	return r.enforcer.SavePolicy()
}

// 新增方法的实现

func (r *permissionRepository) AddRoleForUser(ctx context.Context, user string, role string) (bool, error) {
	return r.enforcer.AddRoleForUser(user, role)
}

func (r *permissionRepository) DeleteRoleForUser(ctx context.Context, user string, role string) (bool, error) {
	return r.enforcer.DeleteRoleForUser(user, role)
}

func (r *permissionRepository) GetPermissionsForUser(ctx context.Context, user string) ([][]string, error) {
	return r.enforcer.GetPermissionsForUser(user)
}

func (r *permissionRepository) GetAllRoles(ctx context.Context) ([]string, error) {
	return r.enforcer.GetAllRoles()
}

func (r *permissionRepository) AddPermissionForUser(ctx context.Context, user string, permission ...string) (bool, error) {
	return r.enforcer.AddPermissionForUser(user, permission...)
}

func (r *permissionRepository) DeletePermissionForUser(ctx context.Context, user string, permission ...string) (bool, error) {
	return r.enforcer.DeletePermissionForUser(user, permission...)
}

func (r *permissionRepository) GetUsersForRole(ctx context.Context, role string) ([]string, error) {
	return r.enforcer.GetUsersForRole(role)
}