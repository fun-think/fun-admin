package service

import (
	"context"
	"fmt"
	v1 "fun-admin/api/v1"
	"fun-admin/internal/model"
	"fun-admin/internal/repository"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginService interface {
	Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponseData, error)
}

func NewLoginService(
	service *Service,
	userRepository repository.UserRepository,
) LoginService {
	return &loginService{
		Service:        service,
		userRepository: userRepository,
	}
}

type loginService struct {
	*Service
	userRepository repository.UserRepository
}

func (s *loginService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponseData, error) {
	var user *model.User
	var err error

	// 根据登录类型处理不同逻辑
	if req.Type == "mobile" {
		// 手机号登录验证
		if !s.verifySMSCode(req.Mobile, req.Code) {
			return nil, fmt.Errorf("验证码错误或已过期")
		}
		user, err = s.userRepository.GetUserByMobile(ctx, req.Mobile)
	} else {
		// 用户名/邮箱登录
		user, err = s.userRepository.GetUserByUsername(ctx, req.Username)
	}

	if err != nil {
		return nil, err
	}

	// 如果是用户名登录，需要验证密码
	if req.Type != "mobile" {
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err != nil {
			return nil, err
		}
	}

	token, err := s.jwt.GenToken(user.ID, time.Now().Add(time.Hour*24*90))
	if err != nil {
		return nil, err
	}

	// 返回用户基本信息
	response := &v1.LoginResponseData{
		AccessToken: token,
		UserInfo: v1.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Email:    user.Email,
			Phone:    user.Phone,
		},
	}

	return response, nil
}

// verifySMSCode 验证短信验证码（模拟实现）
func (s *loginService) verifySMSCode(mobile, code string) bool {
	// 实际项目中这里应该查询数据库或调用短信服务API验证验证码
	// 这里为了演示，假设验证码是123456
	return code == "123456"
}
