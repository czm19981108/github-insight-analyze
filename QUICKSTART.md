# 快速入门指南

只需 5 分钟即可启动和运行 OSS Insight 趋势通知器！

## 前置要求

- 已安装 Go 1.21+ ([下载](https://golang.org/dl/))
- SMTP 凭据（Gmail、Outlook 或任何邮件服务）
- GitHub 账户（用于自动执行）

## 步骤 1：安装 Go

如果尚未安装 Go，请从 [golang.org](https://golang.org/dl/) 下载并安装

验证安装：
```bash
go version
```

## 步骤 2：克隆和构建

```bash
# 克隆仓库
git clone https://github.com/yourusername/github-insight-analyze.git
cd github-insight-analyze

# 下载依赖
go mod download

# 构建应用程序
go build -o notifier ./cmd/notifier
```

或使用 Make：
```bash
make build
```

## 步骤 3：配置

### 选项 A：使用环境变量（推荐用于测试）

```bash
# 复制示例文件
cp .env.example .env

# 编辑 .env 文件设置您的配置
nano .env  # 或使用您喜欢的编辑器
```

必需的设置：
```bash
SMTP_HOST="smtp.gmail.com"
SMTP_PORT="587"
SMTP_USERNAME="your-email@gmail.com"
SMTP_PASSWORD="your-app-password"
EMAIL_FROM="your-email@gmail.com"
EMAIL_TO="recipient@example.com"
QUERY_LANGUAGE="go"
QUERY_PERIOD="daily"
```

### 选项 B：使用配置文件

```bash
# 复制示例配置
cp configs/config.example.yaml configs/config.yaml

# 编辑配置文件
nano configs/config.yaml
```

## 步骤 4：Gmail 设置（如果使用 Gmail）

1. 访问 [Google 账户安全](https://myaccount.google.com/security)
2. 启用 **两步验证**
3. 访问 [应用专用密码](https://myaccount.google.com/apppasswords)
4. 为"邮件"创建新的应用专用密码
5. 复制密码并在配置中使用

## 步骤 5：测试运行

```bash
# 使用环境变量运行
./notifier

# 或使用配置文件运行
./notifier -config configs/config.yaml

# 或使用 Make
make run
```

您应该会看到类似以下的输出：
```
2025/01/07 10:30:00 Loading configuration...
2025/01/07 10:30:00 Configuration loaded successfully
2025/01/07 10:30:00 - Language: go
2025/01/07 10:30:00 - Period: daily
2025/01/07 10:30:00 Creating API client...
2025/01/07 10:30:00 Fetching trending repositories...
2025/01/07 10:30:02 Successfully fetched 100 repositories
2025/01/07 10:30:02 Formatting data...
2025/01/07 10:30:02 Creating email client...
2025/01/07 10:30:02 Sending email...
2025/01/07 10:30:05 Email sent successfully!
```

## 步骤 6：设置 GitHub Actions（可选）

用于自动生成每日报告：

### 1. 推送到 GitHub

```bash
git init
git add .
git commit -m "Initial commit"
git branch -M main
git remote add origin https://github.com/yourusername/github-insight-analyze.git
git push -u origin main
```

### 2. 添加密钥

访问：**仓库 → Settings → Secrets and variables → Actions**

添加以下密钥：
- `SMTP_HOST`: smtp.gmail.com
- `SMTP_PORT`: 587
- `SMTP_USERNAME`: your-email@gmail.com
- `SMTP_PASSWORD`: your-app-password
- `EMAIL_FROM`: your-email@gmail.com
- `EMAIL_TO`: recipient@example.com

### 3. 启用 Actions

- 访问 **Actions** 标签页
- 启用工作流
- 工作流将每天上海时间 07:30 运行

### 4. 手动测试

- 访问 **Actions** → **Daily Trending Report**
- 点击 **Run workflow**
- 选择分支和参数
- 点击 **Run workflow**

## 常见问题

### 问题：邮件发送失败

**解决方案**：检查 SMTP 凭据和端口。Gmail 用户必须使用应用专用密码。

### 问题：API 超时

**解决方案**：在配置中增加超时时间：
```bash
export API_TIMEOUT=60
```

### 问题：找不到 Go 命令

**解决方案**：安装 Go 或将其添加到 PATH：
```bash
export PATH=$PATH:/usr/local/go/bin
```

## 下一步

- 在 `pkg/formatter/formatter.go` 中自定义邮件模板
- 在 `.github/workflows/daily-report.yml` 中调整计划时间
- 在配置中添加多个收件人
- 探索不同的语言和时间段

## 常用命令

```bash
# 构建
make build

# 运行
make run

# 使用配置运行
make run-config

# 测试
make test

# 清理构建产物
make clean

# 格式化代码
make fmt

# 显示所有命令
make help

# 检查版本
./notifier -version
```

## 支持

- 阅读完整的 [README.md](README.md)
- 查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解开发信息
- 在 [GitHub](https://github.com/yourusername/github-insight-analyze/issues) 上报告问题

## 示例输出

您将收到一封包含精美 HTML 报告的邮件，显示：

- 前 100 个热门仓库
- 仓库描述
- 星标数和增长情况
- Fork 数
- 编程语言
- 仓库直接链接

祝您使用愉快！🚀
