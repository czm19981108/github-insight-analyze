package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config 应用程序配置
type Config struct {
	API   APIConfig   `yaml:"api"`
	Email EmailConfig `yaml:"email"`
	Query QueryConfig `yaml:"query"`
}

// APIConfig OSS Insight API配置
type APIConfig struct {
	BaseURL string `yaml:"base_url"`
	Timeout int    `yaml:"timeout"` // 超时时间（秒）
}

// EmailConfig 邮件配置
type EmailConfig struct {
	SMTPHost     string   `yaml:"smtp_host"`
	SMTPPort     int      `yaml:"smtp_port"`
	Username     string   `yaml:"username"`
	Password     string   `yaml:"password"`
	From         string   `yaml:"from"`
	To           []string `yaml:"to"`
	Subject      string   `yaml:"subject"`
	UseHTML      bool     `yaml:"use_html"`
}

// QueryConfig 查询参数配置
type QueryConfig struct {
	Language  string `yaml:"language"`   // 编程语言，如 "go", "java", "all"
	Period    string `yaml:"period"`     // 时间范围，如 "daily", "weekly", "monthly"
	Limit     int    `yaml:"limit"`      // 获取数量，默认100
}

// Load 从配置文件加载配置
func Load(configPath string) (*Config, error) {
	config := &Config{
		API: APIConfig{
			BaseURL: "https://api.ossinsight.io",
			Timeout: 30,
		},
		Query: QueryConfig{
			Language: "all",
			Period:   "daily",
			Limit:    100,
		},
		Email: EmailConfig{
			SMTPPort: 587,
			Subject:  "GitHub Trending Repositories Report",
			UseHTML:  true,
		},
	}

	// 如果提供了配置文件路径，则从文件加载
	if configPath != "" {
		if err := loadFromFile(configPath, config); err != nil {
			return nil, fmt.Errorf("failed to load config from file: %w", err)
		}
	}

	// 从环境变量覆盖配置
	loadFromEnv(config)

	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// loadFromFile 从YAML文件加载配置
func loadFromFile(path string, config *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, config)
}

// loadFromEnv 从环境变量加载配置
func loadFromEnv(config *Config) {
	// API配置
	if v := os.Getenv("API_BASE_URL"); v != "" {
		config.API.BaseURL = v
	}
	if v := os.Getenv("API_TIMEOUT"); v != "" {
		if timeout, err := strconv.Atoi(v); err == nil {
			config.API.Timeout = timeout
		}
	}

	// 邮件配置
	if v := os.Getenv("SMTP_HOST"); v != "" {
		config.Email.SMTPHost = v
	}
	if v := os.Getenv("SMTP_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			config.Email.SMTPPort = port
		}
	}
	if v := os.Getenv("SMTP_USERNAME"); v != "" {
		config.Email.Username = v
	}
	if v := os.Getenv("SMTP_PASSWORD"); v != "" {
		config.Email.Password = v
	}
	if v := os.Getenv("EMAIL_FROM"); v != "" {
		config.Email.From = v
	}
	if v := os.Getenv("EMAIL_TO"); v != "" {
		config.Email.To = strings.Split(v, ",")
	}
	if v := os.Getenv("EMAIL_SUBJECT"); v != "" {
		config.Email.Subject = v
	}
	if v := os.Getenv("EMAIL_USE_HTML"); v != "" {
		config.Email.UseHTML = v == "true" || v == "1"
	}

	// 查询配置
	if v := os.Getenv("QUERY_LANGUAGE"); v != "" {
		config.Query.Language = v
	}
	if v := os.Getenv("QUERY_PERIOD"); v != "" {
		config.Query.Period = v
	}
	if v := os.Getenv("QUERY_LIMIT"); v != "" {
		if limit, err := strconv.Atoi(v); err == nil {
			config.Query.Limit = limit
		}
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	// 验证邮件配置
	if c.Email.SMTPHost == "" {
		return fmt.Errorf("SMTP host is required")
	}
	if c.Email.Username == "" {
		return fmt.Errorf("SMTP username is required")
	}
	if c.Email.Password == "" {
		return fmt.Errorf("SMTP password is required")
	}
	if c.Email.From == "" {
		return fmt.Errorf("email from address is required")
	}
	if len(c.Email.To) == 0 {
		return fmt.Errorf("at least one recipient email is required")
	}

	// 验证查询配置
	validPeriods := map[string]bool{
		"daily":   true,
		"weekly":  true,
		"monthly": true,
	}
	if !validPeriods[c.Query.Period] {
		return fmt.Errorf("invalid period: %s (must be daily, weekly, or monthly)", c.Query.Period)
	}

	if c.Query.Limit <= 0 || c.Query.Limit > 100 {
		return fmt.Errorf("limit must be between 1 and 100")
	}

	return nil
}
