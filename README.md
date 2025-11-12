# GitHub Trending Notifier

A Go-based application that fetches GitHub trending repositories and sends automated email reports. Perfect for staying updated with the latest trending projects in your favorite programming languages.

## Features

- Fetch trending repositories from GitHub API
- Support for language filtering (Go, Java, Python, JavaScript, etc.)
- Multiple time periods (daily, weekly, monthly)
- Beautiful HTML email templates
- Plain text email support
- Automated daily reports via GitHub Actions
- Configurable via environment variables or YAML files
- Comprehensive error handling and logging

## Project Structure

```
.
├── cmd/
│   └── notifier/          # Main application entry point
├── pkg/
│   ├── api/               # GitHub API client
│   ├── email/             # Email sending functionality
│   └── formatter/         # Data formatting (text & HTML)
├── internal/
│   └── config/            # Configuration management
├── configs/
│   └── config.example.yaml # Example configuration file
├── .github/
│   └── workflows/         # GitHub Actions workflows
├── go.mod
├── go.sum
└── README.md
```

## Prerequisites

- Go 1.21 or later
- SMTP server credentials (e.g., Gmail, SendGrid, Mailgun)
- GitHub account (for automated execution)

## Installation

### 1. Install Go

Download and install Go from [golang.org](https://golang.org/dl/)

Verify installation:
```bash
go version
```

### 2. Clone Repository

```bash
git clone https://github.com/yourusername/github-insight-analyze.git
cd github-insight-analyze
```

### 3. Install Dependencies

```bash
go mod download
```

### 4. Build Application

```bash
go build -o notifier ./cmd/notifier
```

## Configuration

### Option 1: Configuration File

1. Copy the example configuration:
```bash
cp configs/config.example.yaml configs/config.yaml
```

2. Edit `configs/config.yaml`:
```yaml
api:
  base_url: "https://api.github.com"
  timeout: 30

email:
  smtp_host: "smtp.gmail.com"
  smtp_port: 587
  username: "your-email@gmail.com"
  password: "your-app-password"
  from: "your-email@gmail.com"
  to:
    - "recipient@example.com"
  subject: "GitHub Trending Repositories Report"
  use_html: true

query:
  language: "go"      # Options: go, java, python, javascript, all, etc.
  period: "daily"     # Options: daily, weekly, monthly
  limit: 100
```

### Option 2: Environment Variables

Set the following environment variables:

```bash
# SMTP Configuration
export SMTP_HOST="smtp.gmail.com"
export SMTP_PORT="587"
export SMTP_USERNAME="your-email@gmail.com"
export SMTP_PASSWORD="your-app-password"
export EMAIL_FROM="your-email@gmail.com"
export EMAIL_TO="recipient1@example.com,recipient2@example.com"
export EMAIL_SUBJECT="GitHub Trending Report"
export EMAIL_USE_HTML="true"

# API Configuration
export API_BASE_URL="https://api.github.com"
export API_TIMEOUT="30"

# Query Configuration
export QUERY_LANGUAGE="go"
export QUERY_PERIOD="daily"
export QUERY_LIMIT="100"
```

### Gmail Setup

If using Gmail, you need to create an App Password:

1. Go to [Google Account Security](https://myaccount.google.com/security)
2. Enable 2-Step Verification
3. Go to [App Passwords](https://myaccount.google.com/apppasswords)
4. Generate a new app password for "Mail"
5. Use this password in your configuration

## Usage

### Run Locally

Using configuration file:
```bash
./notifier -config configs/config.yaml
```

Using environment variables:
```bash
./notifier
```

Check version:
```bash
./notifier -version
```

### GitHub Actions Automated Execution

#### 1. Set Up Secrets

Go to your GitHub repository:
- Settings → Secrets and variables → Actions → New repository secret

Add the following secrets:
- `SMTP_HOST`: Your SMTP server host
- `SMTP_PORT`: SMTP port (usually 587)
- `SMTP_USERNAME`: Your email username
- `SMTP_PASSWORD`: Your email password or app password
- `EMAIL_FROM`: Sender email address
- `EMAIL_TO`: Recipient email addresses (comma-separated)

#### 2. Configure Workflow

The workflow is already configured in `.github/workflows/daily-report.yml`

**Schedule**: Runs daily at 07:30 Shanghai time (23:30 UTC previous day)

**Manual Trigger**: You can also trigger manually from GitHub Actions tab with custom parameters:
- Language (e.g., go, java, python)
- Period (daily, weekly, monthly)

#### 3. Enable Actions

1. Go to your repository on GitHub
2. Click on "Actions" tab
3. Enable GitHub Actions if not already enabled
4. The workflow will run automatically according to the schedule

### Manual Trigger via GitHub Actions

1. Go to Actions tab
2. Select "Daily Trending Report" workflow
3. Click "Run workflow"
4. Select branch and input parameters (optional)
5. Click "Run workflow" button

## API Reference

### GitHub API

This project uses the GitHub Search API to fetch trending repositories.

**Endpoint**: `https://api.github.com/search/repositories`

**Query Parameters**:
- `q`: Search query (e.g., `stars:>50 pushed:>2025-01-11 language:go`)
- `sort`: Sort by (stars)
- `order`: Sort order (desc)
- `per_page`: Number of results (1-100)

**Note**: GitHub API has a rate limit of 60 requests/hour for unauthenticated requests. For higher limits, consider adding authentication.

## Troubleshooting

### Email Not Sending

1. **Check SMTP credentials**: Ensure username and password are correct
2. **Gmail users**: Make sure you're using an App Password, not your regular password
3. **Firewall**: Ensure port 587 is not blocked
4. **TLS/SSL**: Some SMTP servers require TLS on port 587

### API Errors

1. **Timeout**: Increase `API_TIMEOUT` value
2. **Rate limiting**: GitHub API has rate limits for unauthenticated requests
3. **Network issues**: Check internet connectivity

### GitHub Actions Not Running

1. **Check secrets**: Ensure all required secrets are set
2. **Workflow enabled**: Make sure GitHub Actions are enabled for your repository
3. **Cron syntax**: Verify the cron schedule is correct
4. **Logs**: Check workflow logs for specific errors

## Development

### Running Tests

```bash
go test ./...
```

### Code Structure

- `cmd/notifier/main.go`: Application entry point
- `pkg/api/client.go`: GitHub API client implementation
- `pkg/email/client.go`: Email sending functionality
- `pkg/formatter/formatter.go`: Data formatting (text & HTML)
- `internal/config/config.go`: Configuration management

### Adding New Features

1. Fork the repository
2. Create a feature branch
3. Implement your changes
4. Add tests
5. Submit a pull request

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License

## Acknowledgments

- [GitHub](https://github.com) for providing the Search API
- GitHub for hosting and Actions automation

## Support

If you encounter any issues or have questions:

1. Check the [Troubleshooting](#troubleshooting) section
2. Search existing [GitHub Issues](https://github.com/yourusername/github-insight-analyze/issues)
3. Create a new issue with detailed information

## Roadmap

- [ ] Support for multiple notification channels (Slack, Discord, Telegram)
- [ ] Web dashboard for viewing reports
- [ ] Database storage for historical data
- [ ] Custom filtering rules
- [ ] Repository recommendations based on user interests
- [ ] Weekly/Monthly digest summaries

---

Made with ❤️ by the Open Source Community
