package email

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/gomail.v2"
)

// MailConfig 存储了发送邮件所需的配置信息
type MailConfig struct {
	Host               string
	Port               int
	UserName           string
	PasswordOrAuthCode string // 授权码或密码
	TlsVerify          bool
}

// MailService 是一个邮件发送服务的接口
type MailService interface {
	Send(to []string, subject, body string) error
}

// GoMailSender 是基于 GoMail 库实现的邮件发送服务
type GoMailSender struct {
	config *MailConfig
	dialer *gomail.Dialer
}

// NewGoMailSender 创建一个新的 NewGoMailSender 实例
func NewGoMailSender(config *MailConfig) (*GoMailSender, error) {
	dialer := gomail.NewDialer(config.Host, config.Port, config.UserName, config.PasswordOrAuthCode)
	dialer.TLSConfig = &tls.Config{
		ServerName:         config.Host,
		InsecureSkipVerify: !config.TlsVerify,
	}
	return &GoMailSender{
		config: config,
		dialer: dialer,
	}, nil
}

// Send 发送邮件给指定列表中的收件人
func (s *GoMailSender) Send(to []string, subject, body string, files []string) error {
	m := gomail.NewMessage()

	// 设置邮件头
	m.SetHeader("From", s.config.UserName)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan") // 抄送
	m.SetBody("text/plain", body)

	// 添加附件
	for _, file := range files {
		m.Attach(file)
	}

	// 发送邮件
	if err := s.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
