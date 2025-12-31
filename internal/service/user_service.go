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
		BaseModel: model.BaseModel{
			ID: req.ID,
		},
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
	
	// 收集用户ID列表
	userIDs := make([]uint, 0, len(list))
	for _, user := range list {
		userIDs = append(userIDs, user.ID)
	}
	
	// 批量获取用户角色信息
	usersRoles := make(map[uint][]string)
	if len(userIDs) > 0 {
		// 使用GetAllRoles方法获取所有角色关系，然后过滤出当前用户的角色
		// 这样可以避免多次调用GetRolesForUser造成的N+1问题
		allRoles, err := s.permissionRepository.GetAllRoles(ctx)
		if err != nil {
			return nil, err
		}
		
		// 创建用户ID到字符串ID的映射
		userIDMap := make(map[string]uint)
		for _, userID := range userIDs {
			userIDMap[uint64ToString(uint64(userID))] = userID
		}
		
		// 遍历所有角色关系，找出当前用户的
		for _, role := range allRoles {
			// 获取拥有该角色的所有用户
			usersForRole, err := s.permissionRepository.GetUsersForRole(ctx, role)
			if err != nil {
				continue
			}
			
			// 检查当前用户是否拥有该角色
			for _, userStrID := range usersForRole {
				if userID, exists := userIDMap[userStrID]; exists {
					usersRoles[userID] = append(usersRoles[userID], role)
				}
			}
		}
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
			Roles:     usersRoles[user.ID],
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
