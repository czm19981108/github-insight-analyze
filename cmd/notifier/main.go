package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ossinsight-analyze/trending-notifier/internal/config"
	"github.com/ossinsight-analyze/trending-notifier/pkg/api"
	"github.com/ossinsight-analyze/trending-notifier/pkg/email"
	"github.com/ossinsight-analyze/trending-notifier/pkg/formatter"
)

var (
	configPath = flag.String("config", "", "Path to configuration file")
	version    = flag.Bool("version", false, "Show version information")
)

const appVersion = "1.0.0"

func main() {
	flag.Parse()

	// 显示版本信息
	if *version {
		fmt.Printf("OSS Insight Trending Notifier v%s\n", appVersion)
		os.Exit(0)
	}

	// 加载配置
	log.Println("Loading configuration...")
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Configuration loaded successfully")
	log.Printf("- Language: %s", cfg.Query.Language)
	log.Printf("- Period: %s", cfg.Query.Period)
	log.Printf("- Limit: %d", cfg.Query.Limit)
	log.Printf("- Recipients: %v", cfg.Email.To)

	// 运行主逻辑
	if err := run(cfg); err != nil {
		log.Fatalf("Application error: %v", err)
	}

	log.Println("Email sent successfully!")
}

func run(cfg *config.Config) error {
	ctx := context.Background()

	// 创建API客户端
	log.Println("Creating API client...")
	apiClient := api.NewClient(
		cfg.API.BaseURL,
		time.Duration(cfg.API.Timeout)*time.Second,
	)

	// 获取trending repositories
	log.Printf("Fetching trending repositories (language: %s, period: %s)...",
		cfg.Query.Language, cfg.Query.Period)

	repos, err := apiClient.GetTrendingRepos(
		ctx,
		cfg.Query.Language,
		cfg.Query.Period,
		cfg.Query.Limit,
	)
	if err != nil {
		return fmt.Errorf("failed to fetch trending repositories: %w", err)
	}

	if len(repos) == 0 {
		log.Println("Warning: No repositories returned from API")
		return fmt.Errorf("no repositories found")
	}

	log.Printf("Successfully fetched %d repositories", len(repos))

	// 格式化数据
	log.Println("Formatting data...")
	var formattedContent string

	if cfg.Email.UseHTML {
		htmlFormatter := formatter.NewHTMLFormatter()
		formattedContent, err = htmlFormatter.Format(repos, cfg.Query.Language, cfg.Query.Period)
		if err != nil {
			return fmt.Errorf("failed to format data as HTML: %w", err)
		}
	} else {
		textFormatter := formatter.NewTextFormatter()
		formattedContent, err = textFormatter.Format(repos, cfg.Query.Language, cfg.Query.Period)
		if err != nil {
			return fmt.Errorf("failed to format data as text: %w", err)
		}
	}

	log.Println("Data formatted successfully")

	// 创建邮件客户端
	log.Println("Creating email client...")
	emailClient := email.NewClient(
		cfg.Email.SMTPHost,
		cfg.Email.SMTPPort,
		cfg.Email.Username,
		cfg.Email.Password,
		cfg.Email.From,
	)

	// 发送邮件
	log.Printf("Sending email to %d recipients...", len(cfg.Email.To))

	msg := &email.Message{
		To:      cfg.Email.To,
		Subject: cfg.Email.Subject,
		Body:    formattedContent,
		IsHTML:  cfg.Email.UseHTML,
	}

	if err := emailClient.Send(msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
