package service

import (
	"context"
	"errors"
	v1 "fun-admin/api/v1"
	"fun-admin/internal/model"
	"fun-admin/internal/repository"

	"gorm.io/gorm"
)

type RoleService interface {
	GetRoles(ctx context.Context, req *v1.GetRoleListRequest) (*v1.GetRolesResponseData, error)
	RoleUpdate(ctx context.Context, req *v1.RoleUpdateRequest) error
	RoleCreate(ctx context.Context, req *v1.RoleCreateRequest) error
	RoleDelete(ctx context.Context, id uint) error
}

func NewRoleService(
	service *Service,
	roleRepository repository.RoleRepository,
	permissionRepository repository.PermissionRepository, // 添加权限仓库依赖
) RoleService {
	return &roleService{
		Service:              service,
		roleRepository:       roleRepository,
		permissionRepository: permissionRepository, // 注入权限仓库
	}
}

type roleService struct {
	*Service
	roleRepository       repository.RoleRepository
	permissionRepository repository.PermissionRepository // 添加权限仓库
}

func (s *roleService) RoleUpdate(ctx context.Context, req *v1.RoleUpdateRequest) error {
	return s.roleRepository.RoleUpdate(ctx, &model.Role{
		Name: req.Name,
		Sid:  req.Sid,
		BaseModel: model.BaseModel{
			ID: req.ID,
		},
	})
}

func (s *roleService) RoleCreate(ctx context.Context, req *v1.RoleCreateRequest) error {
	_, err := s.roleRepository.GetRoleBySid(ctx, req.Sid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.roleRepository.RoleCreate(ctx, &model.Role{
				Name: req.Name,
				Sid:  req.Sid,
			})
		} else {
			return err
		}
	}
	return nil
}

func (s *roleService) RoleDelete(ctx context.Context, id uint) error {
	old, err := s.roleRepository.GetRole(ctx, id)
	if err != nil {
		return err
	}

	// 使用permissionRepository而不是roleRepository来处理casbin相关操作
	if _, err := s.permissionRepository.DeleteRole(ctx, old.Sid); err != nil {
		return err
	}

	if err := s.permissionRepository.SavePolicy(ctx); err != nil {
		return err
	}

	return s.roleRepository.RoleDelete(ctx, id)
}

func (s *roleService) GetRoles(ctx context.Context, req *v1.GetRoleListRequest) (*v1.GetRolesResponseData, error) {
	list, total, err := s.roleRepository.GetRoles(ctx, req)
	if err != nil {
		return nil, err
	}
	data := &v1.GetRolesResponseData{
		List:  make([]v1.RoleDataItem, 0),
		Total: total,
	}
	for _, role := range list {
		data.List = append(data.List, v1.RoleDataItem{
			ID:        role.ID,
			Name:      role.Name,
			Sid:       role.Sid,
			UpdatedAt: role.UpdatedAt.Format("2006-01-02 15:04:05"),
			CreatedAt: role.CreatedAt.Format("2006-01-02 15:04:05"),
		})

	}
	return data, nil
}
