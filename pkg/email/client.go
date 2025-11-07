package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

// Client 邮件客户端
type Client struct {
	smtpHost string
	smtpPort int
	username string
	password string
	from     string
}

// NewClient 创建邮件客户端
func NewClient(smtpHost string, smtpPort int, username, password, from string) *Client {
	return &Client{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		username: username,
		password: password,
		from:     from,
	}
}

// Message 邮件消息
type Message struct {
	To      []string
	Subject string
	Body    string
	IsHTML  bool
}

// Send 发送邮件
func (c *Client) Send(msg *Message) error {
	if len(msg.To) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	// 构建邮件内容
	content := c.buildMessage(msg)

	// SMTP认证
	auth := smtp.PlainAuth("", c.username, c.password, c.smtpHost)

	// 发送邮件
	addr := fmt.Sprintf("%s:%d", c.smtpHost, c.smtpPort)
	err := smtp.SendMail(addr, auth, c.from, msg.To, []byte(content))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// buildMessage 构建邮件内容
func (c *Client) buildMessage(msg *Message) string {
	var sb strings.Builder

	// 邮件头
	sb.WriteString(fmt.Sprintf("From: %s\r\n", c.from))
	sb.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(msg.To, ", ")))
	sb.WriteString(fmt.Sprintf("Subject: %s\r\n", msg.Subject))
	sb.WriteString("MIME-Version: 1.0\r\n")

	// 内容类型
	if msg.IsHTML {
		sb.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	} else {
		sb.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	}

	sb.WriteString("\r\n")

	// 邮件正文
	sb.WriteString(msg.Body)

	return sb.String()
}

// SendText 发送纯文本邮件
func (c *Client) SendText(to []string, subject, body string) error {
	msg := &Message{
		To:      to,
		Subject: subject,
		Body:    body,
		IsHTML:  false,
	}
	return c.Send(msg)
}

// SendHTML 发送HTML邮件
func (c *Client) SendHTML(to []string, subject, body string) error {
	msg := &Message{
		To:      to,
		Subject: subject,
		Body:    body,
		IsHTML:  true,
	}
	return c.Send(msg)
}

// ValidateConfig 验证邮件配置
func ValidateConfig(smtpHost, username, password, from string, to []string) error {
	if smtpHost == "" {
		return fmt.Errorf("SMTP host is required")
	}
	if username == "" {
		return fmt.Errorf("SMTP username is required")
	}
	if password == "" {
		return fmt.Errorf("SMTP password is required")
	}
	if from == "" {
		return fmt.Errorf("from address is required")
	}
	if len(to) == 0 {
		return fmt.Errorf("at least one recipient is required")
	}
	return nil
}
