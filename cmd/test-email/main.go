package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/github-insight-analyze/trending-notifier/utils"
	"gopkg.in/gomail.v2"
)

func main() {
	// 使用 Windows 命名互斥体防止重复运行
	// 这是系统级别的锁，比文件锁更可靠
	mutex, err := utils.CreateNamedMutex("Global\\OSSInsightTestEmailSender")
	if err != nil {
		log.Printf("⚠️  %v", err)
		log.Println("程序已退出（防止重复发送邮件）")
		return
	}
	defer mutex.Release()

	runID := time.Now().UnixNano()
	// 生成唯一的运行 ID
	log.Printf("========== 运行 ID: %d ==========", runID)
	// 加载 .env 文件
	// 如果不在这里手动加载.env文件会导致程序去读取系统环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("警告: 未找到 .env 文件，将使用系统环境变量")
	}

	log.Println("========== 程序开始执行 ==========")

	// 从环境变量读取配置
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	// 加一步格式转换
	smtpPort, err2 := utils.StringToInt(smtpPortStr)
	if err2 != nil {
		log.Fatalf("SMTP端口格式错误: %v", err2)
	}

	fmt.Println("当前时间(标准格式):", utils.FormatTime("yyyy-MM-dd HH:mm:ss"))
	fmt.Println("当前时间(紧凑格式):", utils.FormatTime("yyyyMMddHHmmss"))
	fmt.Println("当前时间(仅日期):", utils.FormatTime("yyyyMMdd"))

	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("EMAIL_FROM")
	to := os.Getenv("EMAIL_TO") // 你的测试邮箱

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "SMTP配置测试")

	// 添加时间戳以区分不同的发送
	timestamp := utils.FormatTime("yyyy-MM-dd HH:mm:ss")
	m.SetBody("text/html",
		fmt.Sprintf("<h1>测试成功！</h1><p>您的SMTP配置正确。</p><p>发送时间: %s</p><p>运行ID: %d</p>", timestamp, runID))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	log.Println("========== 准备发送邮件 ==========")
	if err3 := d.DialAndSend(m); err3 != nil {
		log.Fatal("发送失败: ", err3)
	}

	log.Println("========== 测试邮件发送成功！ ==========")
	log.Println("========== 程序执行完毕 ==========")
}
