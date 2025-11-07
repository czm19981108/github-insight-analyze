# Quick Start Guide

Get up and running with OSS Insight Trending Notifier in 5 minutes!

## Prerequisites

- Go 1.21+ installed ([Download](https://golang.org/dl/))
- SMTP credentials (Gmail, Outlook, or any email service)
- GitHub account (for automated execution)

## Step 1: Install Go

If Go is not installed, download and install it from [golang.org](https://golang.org/dl/)

Verify installation:
```bash
go version
```

## Step 2: Clone and Build

```bash
# Clone the repository
git clone https://github.com/yourusername/ossinsight-analyze.git
cd ossinsight-analyze

# Download dependencies
go mod download

# Build the application
go build -o notifier ./cmd/notifier
```

Or use Make:
```bash
make build
```

## Step 3: Configure

### Option A: Use Environment Variables (Recommended for testing)

```bash
# Copy the example file
cp .env.example .env

# Edit .env with your settings
nano .env  # or use your preferred editor
```

Required settings:
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

### Option B: Use Configuration File

```bash
# Copy the example config
cp configs/config.example.yaml configs/config.yaml

# Edit the config file
nano configs/config.yaml
```

## Step 4: Gmail Setup (if using Gmail)

1. Go to [Google Account Security](https://myaccount.google.com/security)
2. Enable **2-Step Verification**
3. Go to [App Passwords](https://myaccount.google.com/apppasswords)
4. Create a new app password for "Mail"
5. Copy the password and use it in your configuration

## Step 5: Test Run

```bash
# Run with environment variables
./notifier

# Or run with config file
./notifier -config configs/config.yaml

# Or use Make
make run
```

You should see output like:
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

## Step 6: Set Up GitHub Actions (Optional)

For automated daily reports:

### 1. Push to GitHub

```bash
git init
git add .
git commit -m "Initial commit"
git branch -M main
git remote add origin https://github.com/yourusername/ossinsight-analyze.git
git push -u origin main
```

### 2. Add Secrets

Go to: **Repository → Settings → Secrets and variables → Actions**

Add these secrets:
- `SMTP_HOST`: smtp.gmail.com
- `SMTP_PORT`: 587
- `SMTP_USERNAME`: your-email@gmail.com
- `SMTP_PASSWORD`: your-app-password
- `EMAIL_FROM`: your-email@gmail.com
- `EMAIL_TO`: recipient@example.com

### 3. Enable Actions

- Go to **Actions** tab
- Enable workflows
- The workflow will run daily at 07:30 Shanghai time

### 4. Test Manually

- Go to **Actions** → **Daily Trending Report**
- Click **Run workflow**
- Select branch and parameters
- Click **Run workflow**

## Common Issues

### Issue: Email not sending

**Solution**: Check SMTP credentials and port. Gmail users must use App Password.

### Issue: API timeout

**Solution**: Increase timeout in configuration:
```bash
export API_TIMEOUT=60
```

### Issue: Go command not found

**Solution**: Install Go or add it to PATH:
```bash
export PATH=$PATH:/usr/local/go/bin
```

## Next Steps

- Customize email templates in `pkg/formatter/formatter.go`
- Adjust schedule in `.github/workflows/daily-report.yml`
- Add multiple recipients in configuration
- Explore different languages and time periods

## Useful Commands

```bash
# Build
make build

# Run
make run

# Run with config
make run-config

# Test
make test

# Clean build artifacts
make clean

# Format code
make fmt

# Show all commands
make help

# Check version
./notifier -version
```

## Support

- Read the full [README.md](README.md)
- Check [CONTRIBUTING.md](CONTRIBUTING.md) for development
- Report issues on [GitHub](https://github.com/yourusername/ossinsight-analyze/issues)

## Example Output

You'll receive an email with a beautiful HTML report showing:

- Top 100 trending repositories
- Repository descriptions
- Star counts and growth
- Fork counts
- Programming languages
- Direct links to repositories

Happy trending! 🚀
