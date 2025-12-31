package service

import (
	"context"
	"crypto/rand"
	"fmt"
	v1 "fun-admin/api/v1"
	"fun-admin/internal/model"
	"fun-admin/internal/repository"
	"math/big"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type LoginService interface {
	Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponseData, error)
	SendSMSCode(ctx context.Context, mobile string) error
}

func NewLoginService(
	service *Service,
	userRepository repository.UserRepository,
	loginRepository repository.LoginRepository,
) LoginService {
	return &loginService{
		Service:        service,
		userRepository: userRepository,
		loginRepository: loginRepository,
	}
}

type loginService struct {
	*Service
	userRepository  repository.UserRepository
	loginRepository repository.LoginRepository
}

func (s *loginService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponseData, error) {
	var user *model.User
	var err error

	// 根据登录类型处理不同逻辑
	if req.Type == "mobile" {
		// 手机号登录验证
		if !s.verifySMSCode(ctx, req.Mobile, req.Code) {
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

// SendSMSCode 发送短信验证码
func (s *loginService) SendSMSCode(ctx context.Context, mobile string) error {
	// 验证手机号格式
	if !s.isValidMobile(mobile) {
		return fmt.Errorf("手机号格式不正确")
	}

	// 生成6位加密安全的随机数字验证码
	code, err := s.generateSecureSMSCode()
	if err != nil {
		return fmt.Errorf("生成验证码失败: %v", err)
	}
	
	// 设置5分钟过期时间
	expiredAt := time.Now().Add(5 * time.Minute)
	
	// 保存验证码到数据库
	if err := s.loginRepository.SaveSMSCode(ctx, mobile, code, expiredAt); err != nil {
		return fmt.Errorf("保存验证码失败: %v", err)
	}
	
	// 实际项目中这里应该调用短信服务发送验证码
	// 这里为了演示，我们记录日志
	s.logger.Info("短信验证码", 
		zap.String("mobile", mobile),
		zap.String("code", code),
		zap.Time("expired_at", expiredAt))
	
	return nil
}

// generateSecureSMSCode 生成加密安全的6位数字验证码
func (s *loginService) generateSecureSMSCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// isValidMobile 验证手机号格式
func (s *loginService) isValidMobile(mobile string) bool {
	// 简单的手机号验证，实际项目中应该更严格
	if len(mobile) != 11 {
		return false
	}
	// 检查是否都是数字
	for _, char := range mobile {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

// verifySMSCode 验证短信验证码
func (s *loginService) verifySMSCode(ctx context.Context, mobile, code string) bool {
	// 获取有效的验证码记录
	smsCode, err := s.loginRepository.GetValidSMSCode(ctx, mobile, code)
	if err != nil {
		s.logger.Error("验证码验证失败", zap.Error(err))
		return false
	}
	
	// 标记验证码为已使用
	if err := s.loginRepository.MarkSMSCodeAsUsed(ctx, smsCode.ID); err != nil {
		s.logger.Error("标记验证码为已使用失败", zap.Error(err))
		// 这里我们仍然认为验证成功，只是记录错误
	}
	
	return true
}
