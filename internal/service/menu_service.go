package service

import (
	"context"
	v1 "fun-admin/api/v1"
	"fun-admin/internal/model"
	"fun-admin/internal/repository"
	"fun-admin/pkg"
	"strings"

	"github.com/duke-git/lancet/v2/convertor"
	"go.uber.org/zap"
)

type MenuService interface {
	GetAdminMenus(ctx context.Context) (*v1.GetMenuResponseData, error)
	GetMenus(ctx context.Context, uid uint) (*v1.GetMenuResponseData, error)
	MenuUpdate(ctx context.Context, req *v1.MenuUpdateRequest) error
	MenuCreate(ctx context.Context, req *v1.MenuCreateRequest) error
	MenuDelete(ctx context.Context, id uint) error
}

func NewMenuService(
	service *Service,
	menuRepository repository.MenuRepository,
	userRepository repository.UserRepository,
	permissionService PermissionService,
) MenuService {
	return &menuService{
		Service:           service,
		menuRepository:    menuRepository,
		userRepository:    userRepository,
		permissionService: permissionService,
	}
}

type menuService struct {
	*Service
	menuRepository    repository.MenuRepository
	userRepository    repository.UserRepository
	permissionService PermissionService
}

func (s *menuService) MenuUpdate(ctx context.Context, req *v1.MenuUpdateRequest) error {
	return s.menuRepository.MenuUpdate(ctx, &model.Menu{
		Component:  req.Component,
		Icon:       req.Icon,
		KeepAlive:  req.KeepAlive,
		HideInMenu: req.HideInMenu,
		Locale:     req.Locale,
		Weight:     req.Weight,
		Name:       req.Name,
		ParentID:   req.ParentID,
		Path:       req.Path,
		Redirect:   req.Redirect,
		Title:      req.Title,
		URL:        req.URL,
		BaseModel: model.BaseModel{
			ID: req.ID,
		},
	})
}

func (s *menuService) MenuCreate(ctx context.Context, req *v1.MenuCreateRequest) error {
	return s.menuRepository.MenuCreate(ctx, &model.Menu{
		Component:  req.Component,
		Icon:       req.Icon,
		KeepAlive:  req.KeepAlive,
		HideInMenu: req.HideInMenu,
		Locale:     req.Locale,
		Weight:     req.Weight,
		Name:       req.Name,
		ParentID:   req.ParentID,
		Path:       req.Path,
		Redirect:   req.Redirect,
		Title:      req.Title,
		URL:        req.URL,
	})
}

func (s *menuService) MenuDelete(ctx context.Context, id uint) error {
	return s.menuRepository.MenuDelete(ctx, id)
}

func (s *menuService) GetMenus(ctx context.Context, uid uint) (*v1.GetMenuResponseData, error) {
	menuList, err := s.menuRepository.GetMenuList(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Error("GetMenuList error", zap.Error(err))
		return nil, err
	}
	data := &v1.GetMenuResponseData{
		List: make([]v1.MenuDataItem, 0),
	}
	// 获取权限的菜单
	permissionData, err := s.permissionService.GetUserPermissions(ctx, uid)
	if err != nil {
		return nil, err
	}

	menuPermMap := map[string]struct{}{}
	for _, permission := range permissionData.List {
		permissionParts := strings.Split(permission, pkg.PermSep)
		// 防呆设置，超管可以看到所有菜单
		if convertor.ToString(uid) == pkg.AdminUserID {
			if len(permissionParts) > 1 && strings.HasPrefix(permissionParts[0], pkg.MenuResourcePrefix) {
				menuPermMap[strings.TrimPrefix(permissionParts[0], pkg.MenuResourcePrefix)] = struct{}{}
			}
		} else {
			if len(permissionParts) > 1 && strings.HasPrefix(permissionParts[0], pkg.MenuResourcePrefix) {
				menuPermMap[strings.TrimPrefix(permissionParts[0], pkg.MenuResourcePrefix)] = struct{}{}
			}
		}
	}

	for _, menu := range menuList {
		if _, ok := menuPermMap[menu.Path]; ok {
			data.List = append(data.List, v1.MenuDataItem{
				ID:         menu.ID,
				Name:       menu.Name,
				Title:      menu.Title,
				Path:       menu.Path,
				Component:  menu.Component,
				Redirect:   menu.Redirect,
				KeepAlive:  menu.KeepAlive,
				HideInMenu: menu.HideInMenu,
				Locale:     menu.Locale,
				Weight:     menu.Weight,
				Icon:       menu.Icon,
				ParentID:   menu.ParentID,
				UpdatedAt:  menu.UpdatedAt.Format("2006-01-02 15:04:05"),
				URL:        menu.URL,
			})
		}
	}
	return data, nil
}

func (s *menuService) GetAdminMenus(ctx context.Context) (*v1.GetMenuResponseData, error) {
	menuList, err := s.menuRepository.GetMenuList(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Error("GetMenuList error", zap.Error(err))
		return nil, err
	}
	data := &v1.GetMenuResponseData{
		List: make([]v1.MenuDataItem, 0),
	}
	for _, menu := range menuList {
		data.List = append(data.List, v1.MenuDataItem{
			ID:         menu.ID,
			Name:       menu.Name,
			Title:      menu.Title,
			Path:       menu.Path,
			Component:  menu.Component,
			Redirect:   menu.Redirect,
			KeepAlive:  menu.KeepAlive,
			HideInMenu: menu.HideInMenu,
			Locale:     menu.Locale,
			Weight:     menu.Weight,
			Icon:       menu.Icon,
			ParentID:   menu.ParentID,
			UpdatedAt:  menu.UpdatedAt.Format("2006-01-02 15:04:05"),
			URL:        menu.URL,
		})
	}
	return data, nil
}
