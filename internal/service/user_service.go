package service

import (
	"context"
	"errors"
	v1 "fun-admin/api/v1"
	"fun-admin/internal/model"
	"fun-admin/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	GetUser(ctx context.Context, id uint) (*v1.GetUserResponseData, error)
	GetUsers(ctx context.Context, req *v1.GetUsersRequest) (*v1.GetUsersResponseData, error)
	UserUpdate(ctx context.Context, req *v1.UserUpdateRequest) error
	UserCreate(ctx context.Context, req *v1.UserCreateRequest) error
	UserDelete(ctx context.Context, id uint) error
}

func NewUserService(
	service *Service,
	userRepository repository.UserRepository,
	roleRepository repository.RoleRepository,
	permissionRepository repository.PermissionRepository, // 添加权限仓库依赖
) UserService {
	return &userService{
		Service:              service,
		userRepository:       userRepository,
		roleRepository:       roleRepository,
		permissionRepository: permissionRepository, // 注入权限仓库
	}
}

type userService struct {
	*Service
	userRepository       repository.UserRepository
	roleRepository       repository.RoleRepository
	permissionRepository repository.PermissionRepository // 添加权限仓库
}

func (s *userService) UserUpdate(ctx context.Context, req *v1.UserUpdateRequest) error {
	password := ""
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		password = string(hash)
	}
	return s.userRepository.UserUpdate(ctx, &model.User{
		Email:    req.Email,
		Nickname: req.Nickname,
		Password: password,
		Phone:    req.Phone,
		Username: req.Username,
		ID:       req.ID,
	})
}

func (s *userService) UserCreate(ctx context.Context, req *v1.UserCreateRequest) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.userRepository.UserCreate(ctx, &model.User{
		Email:    req.Email,
		Nickname: req.Nickname,
		Password: string(hash),
		Phone:    req.Phone,
		Username: req.Username,
	})
}

func (s *userService) UserDelete(ctx context.Context, id uint) error {
	return s.userRepository.UserDelete(ctx, id)
}

func (s *userService) GetUsers(ctx context.Context, req *v1.GetUsersRequest) (*v1.GetUsersResponseData, error) {
	list, total, err := s.userRepository.GetUsers(ctx, req)
	if err != nil {
		return nil, err
	}
	data := &v1.GetUsersResponseData{
		List:  make([]v1.UserDataItem, 0),
		Total: total,
	}
	for _, user := range list {
		data.List = append(data.List, v1.UserDataItem{
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			Email:     user.Email,
			ID:        user.ID,
			Nickname:  user.Nickname,
			Phone:     user.Phone,
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
			Username:  user.Username,
		})
	}
	return data, nil
}

func (s *userService) GetUser(ctx context.Context, id uint) (*v1.GetUserResponseData, error) {
	user, err := s.userRepository.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// 使用permissionRepository获取用户角色
	roles, err := s.permissionRepository.GetRolesForUser(ctx, uint64ToString(uint64(id)))
	if err != nil {
		return nil, err
	}

	data := &v1.GetUserResponseData{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		Email:     user.Email,
		Nickname:  user.Nickname,
		Phone:     user.Phone,
		Username:  user.Username,
		Roles:     roles,
	}
	return data, nil
}

func uint64ToString(v uint64) string {
	buf := make([]byte, 20)
	i := len(buf)
	for v >= 10 {
		i--
		buf[i] = byte(v%10 + '0')
		v /= 10
	}
	i--
	buf[i] = byte(v + '0')
	return string(buf[i:])
}
