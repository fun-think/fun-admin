package repository

import (
	"context"
	v1 "fun-admin/api/v1"
	"fun-admin/internal/model"
	"fun-admin/pkg/logger"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

type RoleRepository interface {
	GetRoles(ctx context.Context, req *v1.GetRoleListRequest) ([]*model.Role, int64, error)
	GetRoleBySid(ctx context.Context, sid string) (*model.Role, error)
	GetRole(ctx context.Context, id uint) (*model.Role, error)
	RoleUpdate(ctx context.Context, role *model.Role) error
	RoleCreate(ctx context.Context, role *model.Role) error
	RoleDelete(ctx context.Context, id uint) error
}

func NewRoleRepository(
	logger *logger.Logger,
	db *gorm.DB,
	enforcer *casbin.SyncedEnforcer,
) RoleRepository {
	return &roleRepository{
		logger:   logger,
		db:       db,
		enforcer: enforcer,
	}
}

type roleRepository struct {
	logger   *logger.Logger
	db       *gorm.DB
	enforcer *casbin.SyncedEnforcer
}

func (r *roleRepository) GetRole(ctx context.Context, id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.Where("id = ?", id).First(&role).Error
	return &role, err
}

func (r *roleRepository) GetRoleBySid(ctx context.Context, sid string) (*model.Role, error) {
	var role model.Role
	err := r.db.Where("sid = ?", sid).First(&role).Error
	return &role, err
}

func (r *roleRepository) GetRoles(ctx context.Context, req *v1.GetRoleListRequest) ([]*model.Role, int64, error) {
	var list []*model.Role
	db := r.db.Model(&model.Role{}).Order("id DESC")

	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	if req.Page > 0 && req.PageSize > 0 {
		db = db.Offset(int((req.Page - 1) * req.PageSize)).Limit(int(req.PageSize))
	}
	err = db.Find(&list).Error
	return list, total, err
}

func (r *roleRepository) RoleUpdate(ctx context.Context, role *model.Role) error {
	return r.db.Model(&model.Role{}).Where("id = ?", role.ID).Updates(role).Error
}

func (r *roleRepository) RoleCreate(ctx context.Context, role *model.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) RoleDelete(ctx context.Context, id uint) error {
	return r.db.Where("id = ?", id).Delete(&model.Role{}).Error
}
