# OSS Insight 热门趋势通知器

一个基于 Go 语言开发的应用程序，用于从 OSS Insight API 获取 GitHub 热门仓库并发送自动化邮件报告。帮助你及时了解你喜欢的编程语言的最新热门项目。

## 功能特性

- 从 OSS Insight API 获取热门仓库
- 支持编程语言过滤（Go、Java、Python、JavaScript 等）
- 多个时间周期（每日、每周、每月）
- 精美的 HTML 邮件模板
- 纯文本邮件支持
- 通过 GitHub Actions 自动生成每日报告
- 可通过环境变量或 YAML 文件配置
- 完善的错误处理和日志记录

## 项目结构

```
.
├── cmd/
│   └── notifier/          # 主程序入口
├── pkg/
│   ├── api/               # OSS Insight API 客户端
│   ├── email/             # 邮件发送功能
│   └── formatter/         # 数据格式化（文本和 HTML）
├── internal/
│   └── config/            # 配置管理
├── configs/
│   └── config.example.yaml # 配置文件示例
├── .github/
│   └── workflows/         # GitHub Actions 工作流
├── go.mod
├── go.sum
└── README.md
```

## 前置要求

- Go 1.21 或更高版本
- SMTP 服务器凭证（如 Gmail、SendGrid、Mailgun）
- GitHub 账号（用于自动化执行）

## 安装步骤

### 1. 安装 Go

从 [golang.org](https://golang.org/dl/) 下载并安装 Go

验证安装：
```bash
go version
```

### 2. 克隆仓库

```bash
git clone https://github.com/yourusername/ossinsight-analyze.git
cd ossinsight-analyze
```

### 3. 安装依赖

```bash
go mod download
```

### 4. 构建应用

```bash
go build -o notifier ./cmd/notifier
```

## 配置说明

### 方式 1：配置文件

1. 复制示例配置文件：
```bash
cp configs/config.example.yaml configs/config.yaml
```

2. 编辑 `configs/config.yaml`：
```yaml
api:
  base_url: "https://api.ossinsight.io"
  timeout: 30

email:
  smtp_host: "smtp.gmail.com"
  smtp_port: 587
  username: "your-email@gmail.com"
  password: "your-app-password"
  from: "your-email@gmail.com"
  to:
    - "recipient@example.com"
  subject: "GitHub 热门仓库报告"
  use_html: true

query:
  language: "go"      # 选项：go、java、python、javascript、all 等
  period: "daily"     # 选项：daily、weekly、monthly
  limit: 100
```

### 方式 2：环境变量

设置以下环境变量：

```bash
# SMTP 配置
export SMTP_HOST="smtp.gmail.com"
export SMTP_PORT="587"
export SMTP_USERNAME="your-email@gmail.com"
export SMTP_PASSWORD="your-app-password"
export EMAIL_FROM="your-email@gmail.com"
export EMAIL_TO="recipient1@example.com,recipient2@example.com"
export EMAIL_SUBJECT="GitHub 热门趋势报告"
export EMAIL_USE_HTML="true"

# API 配置
export API_BASE_URL="https://api.ossinsight.io"
export API_TIMEOUT="30"

# 查询配置
export QUERY_LANGUAGE="go"
export QUERY_PERIOD="daily"
export QUERY_LIMIT="100"
```

### Gmail 设置

如果使用 Gmail，需要创建应用专用密码：

1. 访问 [Google 账号安全](https://myaccount.google.com/security)
2. 启用两步验证
3. 访问[应用专用密码](https://myaccount.google.com/apppasswords)
4. 为"邮件"生成新的应用专用密码
5. 在配置中使用此密码

## 使用方法

### 本地运行

使用配置文件：
```bash
./notifier -config configs/config.yaml
```

使用环境变量：
```bash
./notifier
```

查看版本：
```bash
./notifier -version
```

### GitHub Actions 自动化执行

#### 1. 设置密钥

进入你的 GitHub 仓库：
- Settings → Secrets and variables → Actions → New repository secret

添加以下密钥：
- `SMTP_HOST`：你的 SMTP 服务器主机
- `SMTP_PORT`：SMTP 端口（通常是 587）
- `SMTP_USERNAME`：你的邮箱用户名
- `SMTP_PASSWORD`：你的邮箱密码或应用专用密码
- `EMAIL_FROM`：发件人邮箱地址
- `EMAIL_TO`：收件人邮箱地址（逗号分隔）

#### 2. 配置工作流

工作流已在 `.github/workflows/daily-report.yml` 中配置

**执行计划**：每天上海时间 07:30 执行（UTC 前一天 23:30）

**手动触发**：你也可以在 GitHub Actions 标签页手动触发，并自定义参数：
- 编程语言（如 go、java、python）
- 时间周期（daily、weekly、monthly）

#### 3. 启用 Actions

1. 在 GitHub 上访问你的仓库
2. 点击"Actions"标签
3. 如果尚未启用，启用 GitHub Actions
4. 工作流将按计划自动运行

### 通过 GitHub Actions 手动触发

1. 进入 Actions 标签
2. 选择"Daily Trending Report"工作流
3. 点击"Run workflow"
4. 选择分支并输入参数（可选）
5. 点击"Run workflow"按钮

## API 参考

### OSS Insight API

本项目使用 OSS Insight API 获取热门仓库。

**接口地址**：`https://api.ossinsight.io/v1/repos/trending`

**查询参数**：
- `language`：编程语言过滤器
- `period`：时间周期（daily、weekly、monthly）
- `limit`：结果数量（1-100）

## 故障排查

### 邮件发送失败

1. **检查 SMTP 凭证**：确保用户名和密码正确
2. **Gmail 用户**：确保使用的是应用专用密码，而非常规密码
3. **防火墙**：确保端口 587 未被阻止
4. **TLS/SSL**：某些 SMTP 服务器在端口 587 上需要 TLS

### API 错误

1. **超时**：增加 `API_TIMEOUT` 值
2. **速率限制**：OSS Insight API 可能有速率限制
3. **网络问题**：检查互联网连接

### GitHub Actions 未运行

1. **检查密钥**：确保所有必需的密钥已设置
2. **工作流已启用**：确保仓库已启用 GitHub Actions
3. **Cron 语法**：验证 cron 计划是否正确
4. **日志**：检查工作流日志中的具体错误

## 开发指南

### 运行测试

```bash
go test ./...
```

### 代码结构

- `cmd/notifier/main.go`：应用程序入口
- `pkg/api/client.go`：OSS Insight API 客户端实现
- `pkg/email/client.go`：邮件发送功能
- `pkg/formatter/formatter.go`：数据格式化（文本和 HTML）
- `internal/config/config.go`：配置管理

### 添加新功能

1. Fork 本仓库
2. 创建功能分支
3. 实现你的更改
4. 添加测试
5. 提交 Pull Request

## 贡献指南

欢迎贡献！请随时提交 Pull Request。

## 许可证

MIT License

## 致谢

- 感谢 [OSS Insight](https://ossinsight.io) 提供热门仓库 API
- 感谢 GitHub 提供托管和 Actions 自动化

## 支持

如果遇到任何问题或有疑问：

1. 查看[故障排查](#故障排查)部分
2. 搜索现有的 [GitHub Issues](https://github.com/yourusername/ossinsight-analyze/issues)
3. 创建新的 issue 并提供详细信息

## 路线图

- [ ] 支持多种通知渠道（Slack、Discord、Telegram）
- [ ] 用于查看报告的 Web 仪表板
- [ ] 历史数据的数据库存储
- [ ] 自定义过滤规则
- [ ] 基于用户兴趣的仓库推荐
- [ ] 每周/每月摘要汇总

---

## 将本地项目关联到 GitHub 远程仓库

如果你是第一次将本地项目推送到 GitHub，请按照以下步骤操作：

### 步骤 1：在 GitHub 创建新仓库

1. 登录 [GitHub](https://github.com)
2. 点击右上角的 `+` 按钮，选择 "New repository"
3. 填写仓库名称（如 `ossinsight-analyze`）
4. 选择仓库类型（公开或私有）
5. **不要**勾选"Initialize this repository with a README"（因为本地已有文件）
6. 点击"Create repository"

### 步骤 2：初始化本地 Git 仓库

在项目根目录下打开命令行，执行以下命令：

```bash
# 初始化 Git 仓库
git init

# 添加所有文件到暂存区
git add .

# 创建第一次提交
git commit -m "Initial commit: OSS Insight Trending Notifier"
```

### 步骤 3：关联远程仓库

将本地仓库与 GitHub 远程仓库关联（将 `yourusername` 替换为你的 GitHub 用户名）：

```bash
# 添加远程仓库
git remote add origin https://github.com/yourusername/ossinsight-analyze.git

# 验证远程仓库是否添加成功
git remote -v
```

### 步骤 4：推送代码到 GitHub

```bash
# 将本地代码推送到 GitHub（首次推送）
git push -u origin main
```

**注意**：如果你的默认分支是 `master` 而非 `main`，请将上述命令中的 `main` 替换为 `master`，或者使用以下命令重命名分支：

```bash
# 重命名当前分支为 main
git branch -M main

# 然后推送
git push -u origin main
```

### 步骤 5：验证推送成功

1. 在浏览器中访问你的 GitHub 仓库页面
2. 确认所有文件都已成功上传
3. 检查 README.md 是否正确显示

### 后续推送

完成首次推送后，后续的代码更新只需执行：

```bash
# 添加修改的文件
git add .

# 提交更改
git commit -m "描述你的更改"

# 推送到 GitHub
git push
```

### 使用 SSH 方式（推荐）

如果你配置了 SSH 密钥，可以使用 SSH 地址关联远程仓库：

```bash
# 使用 SSH 地址添加远程仓库
git remote add origin git@github.com:yourusername/ossinsight-analyze.git

# 推送代码
git push -u origin main
```

SSH 方式的优点是无需每次输入用户名和密码。

### 常见问题

**Q: 推送时提示需要身份验证**
A: GitHub 已不再支持密码验证，需要使用个人访问令牌（Personal Access Token）：
1. 访问 GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. 生成新令牌，勾选 `repo` 权限
3. 推送时使用令牌作为密码

**Q: 推送被拒绝（rejected）**
A: 可能是远程仓库有本地没有的提交，先执行：
```bash
git pull origin main --rebase
git push origin main
```

---

用 ❤️ 由开源社区制作
