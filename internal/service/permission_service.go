package service

import (
	"context"
	v1 "fun-admin/api/v1"
	"fun-admin/internal/repository"
	"fun-admin/pkg"
	"strings"
)

type PermissionService interface {
	GetUserPermissions(ctx context.Context, uid uint) (*v1.GetUserPermissionsData, error)
	GetRolePermissions(ctx context.Context, role string) (*v1.GetRolePermissionsData, error)
	UpdateRolePermission(ctx context.Context, req *v1.UpdateRolePermissionRequest) error

	// 新增的方法
	AddRoleForUser(ctx context.Context, user string, role string) error
	DeleteRoleForUser(ctx context.Context, user string, role string) error
	GetPermissionsForUser(ctx context.Context, user string) ([]string, error)
	GetAllRoles(ctx context.Context) ([]string, error)
}

func NewPermissionService(
	service *Service,
	roleRepository repository.RoleRepository,
	userRepository repository.UserRepository,
	permissionRepository repository.PermissionRepository, // 添加权限仓库依赖
) PermissionService {
	return &permissionService{
		Service:              service,
		roleRepository:       roleRepository,
		userRepository:       userRepository,
		permissionRepository: permissionRepository, // 注入权限仓库
	}
}

type permissionService struct {
	*Service
	roleRepository       repository.RoleRepository
	userRepository       repository.UserRepository
	permissionRepository repository.PermissionRepository // 添加权限仓库
}

func (s *permissionService) UpdateRolePermission(ctx context.Context, req *v1.UpdateRolePermissionRequest) error {
	permissions := map[string]struct{}{}
	for _, v := range req.List {
		perm := strings.Split(v, pkg.PermSep)
		if len(perm) == 2 {
			permissions[v] = struct{}{}
		}

	}

	// 使用permissionRepository而不是roleRepository来处理权限
	// 先删除所有权限
	_, err := s.permissionRepository.RemoveFilteredPolicy(ctx, 0, req.Role)
	if err != nil {
		return err
	}

	// 添加新权限
	for permission := range permissions {
		_, err = s.permissionRepository.AddPolicy(ctx, req.Role, permission)
		if err != nil {
			return err
		}
	}

	return s.permissionRepository.SavePolicy(ctx)
}

func (s *permissionService) GetUserPermissions(ctx context.Context, uid uint) (*v1.GetUserPermissionsData, error) {
	data := &v1.GetUserPermissionsData{
		List: []string{},
	}

	// 使用permissionRepository获取用户权限，而不是userRepository
	// 先获取用户角色
	roles, err := s.permissionRepository.GetRolesForUser(ctx, uint64ToString(uint64(uid)))
	if err != nil {
		return nil, err
	}

	// 获取角色权限
	var permissions [][]string
	for _, role := range roles {
		policies, err := s.permissionRepository.GetFilteredPolicy(ctx, 0, role)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, policies...)
	}

	for _, v := range permissions {
		if len(v) == 3 {
			data.List = append(data.List, strings.Join([]string{v[1], v[2]}, pkg.PermSep))
		}
	}
	return data, nil
}

func (s *permissionService) GetRolePermissions(ctx context.Context, role string) (*v1.GetRolePermissionsData, error) {
	data := &v1.GetRolePermissionsData{
		List: []string{},
	}

	// 使用permissionRepository获取角色权限
	list, err := s.permissionRepository.GetFilteredPolicy(ctx, 0, role)
	if err != nil {
		return nil, err
	}

	for _, v := range list {
		if len(v) == 3 {
			data.List = append(data.List, strings.Join([]string{v[1], v[2]}, pkg.PermSep))
		}
	}
	return data, nil
}

// 新增方法的实现

func (s *permissionService) AddRoleForUser(ctx context.Context, user string, role string) error {
	_, err := s.permissionRepository.AddRoleForUser(ctx, user, role)
	if err != nil {
		return err
	}
	return s.permissionRepository.SavePolicy(ctx)
}

func (s *permissionService) DeleteRoleForUser(ctx context.Context, user string, role string) error {
	_, err := s.permissionRepository.DeleteRoleForUser(ctx, user, role)
	if err != nil {
		return err
	}
	return s.permissionRepository.SavePolicy(ctx)
}

func (s *permissionService) GetPermissionsForUser(ctx context.Context, user string) ([]string, error) {
	policies, err := s.permissionRepository.GetPermissionsForUser(ctx, user)
	if err != nil {
		return nil, err
	}

	var permissions []string
	for _, policy := range policies {
		if len(policy) >= 3 {
			permissions = append(permissions, strings.Join([]string{policy[1], policy[2]}, pkg.PermSep))
		}
	}
	return permissions, nil
}

func (s *permissionService) GetAllRoles(ctx context.Context) ([]string, error) {
	return s.permissionRepository.GetAllRoles(ctx)
}
