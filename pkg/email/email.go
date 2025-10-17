package email

import (
	"fmt"
	"gopkg.in/gomail.v2"
)

// EmailConfig 邮件配置
type EmailConfig struct {
	Host     string `json:"host"`      // SMTP服务器地址
	Port     int    `json:"port"`      // SMTP服务器端口
	Username string `json:"username"`  // 用户名
	Password string `json:"password"`  // 密码
	From     string `json:"from"`      // 发件人地址
	FromName string `json:"from_name"` // 发件人名称
}

// EmailMessage 邮件消息
type EmailMessage struct {
	To          []string `json:"to"`          // 收件人
	Subject     string   `json:"subject"`     // 主题
	Body        string   `json:"body"`        // 内容
	IsHTML      bool     `json:"is_html"`     // 是否HTML格式
	CC          []string `json:"cc"`          // 抄送
	BCC         []string `json:"bcc"`         // 密送
	Attachments []string `json:"attachments"` // 附件
}

// EmailService 邮件服务
type EmailService struct {
	config *EmailConfig
	dialer *gomail.Dialer
}

// NewEmailService 创建邮件服务
func NewEmailService(config *EmailConfig) *EmailService {
	dialer := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)

	return &EmailService{
		config: config,
		dialer: dialer,
	}
}

// Send 发送邮件
func (s *EmailService) Send(message *EmailMessage) error {
	if len(message.To) == 0 {
		return fmt.Errorf("收件人不能为空")
	}

	m := gomail.NewMessage()

	// 设置发件人
	from := s.config.From
	if from == "" {
		from = s.config.Username
	}

	fromName := s.config.FromName
	if fromName != "" {
		m.SetAddressHeader("From", from, fromName)
	} else {
		m.SetHeader("From", from)
	}

	// 设置收件人
	m.SetHeader("To", message.To...)

	// 设置抄送
	if len(message.CC) > 0 {
		m.SetHeader("Cc", message.CC...)
	}

	// 设置密送
	if len(message.BCC) > 0 {
		m.SetHeader("Bcc", message.BCC...)
	}

	// 设置主题
	m.SetHeader("Subject", message.Subject)

	// 设置内容
	if message.IsHTML {
		m.SetBody("text/html", message.Body)
	} else {
		m.SetBody("text/plain", message.Body)
	}

	// 添加附件
	for _, attachment := range message.Attachments {
		m.Attach(attachment)
	}

	// 发送邮件
	if err := s.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	return nil
}

// SendTest 发送测试邮件
func (s *EmailService) SendTest(to string) error {
	message := &EmailMessage{
		To:      []string{to},
		Subject: "Fun-Admin 邮件配置测试",
		Body:    "这是一封来自 Fun-Admin 的测试邮件，如果您收到这封邮件，说明邮件配置正确。",
		IsHTML:  false,
	}

	return s.Send(message)
}

// ValidateConfig 验证邮件配置
func (s *EmailService) ValidateConfig() error {
	if s.config.Host == "" {
		return fmt.Errorf("SMTP服务器地址不能为空")
	}

	if s.config.Port <= 0 || s.config.Port > 65535 {
		return fmt.Errorf("SMTP服务器端口无效")
	}

	if s.config.Username == "" {
		return fmt.Errorf("用户名不能为空")
	}

	if s.config.Password == "" {
		return fmt.Errorf("密码不能为空")
	}

	return nil
}

// UpdateConfig 更新配置
func (s *EmailService) UpdateConfig(config *EmailConfig) {
	s.config = config
	s.dialer = gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)
}

// GetConfig 获取配置
func (s *EmailService) GetConfig() *EmailConfig {
	return s.config
}
