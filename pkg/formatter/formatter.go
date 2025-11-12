package formatter

import (
	"fmt"
	"strings"
	"time"

	"github.com/github-insight-analyze/trending-notifier/pkg/api"
)

// Formatter 数据格式化器
type Formatter interface {
	Format(repos []api.Repository, language string, period string) (string, error)
}

// TextFormatter 纯文本格式化器
type TextFormatter struct{}

// NewTextFormatter 创建纯文本格式化器
func NewTextFormatter() *TextFormatter {
	return &TextFormatter{}
}

// Format 格式化为纯文本
func (f *TextFormatter) Format(repos []api.Repository, language string, period string) (string, error) {
	var sb strings.Builder

	// 标题
	sb.WriteString("======================================\n")
	sb.WriteString("GitHub Trending Repositories Report\n")
	sb.WriteString("======================================\n\n")

	// 查询参数
	sb.WriteString(fmt.Sprintf("Language: %s\n", formatLanguage(language)))
	sb.WriteString(fmt.Sprintf("Period: %s\n", formatPeriod(period)))
	sb.WriteString(fmt.Sprintf("Generated: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("Total Repositories: %d\n\n", len(repos)))

	// 分隔线
	sb.WriteString("--------------------------------------\n\n")

	// 仓库列表
	for i, repo := range repos {
		sb.WriteString(fmt.Sprintf("#%d  %s\n", i+1, repo.RepoName))
		sb.WriteString(fmt.Sprintf("    URL: %s\n", repo.URL))

		if repo.Description != "" {
			sb.WriteString(fmt.Sprintf("    Description: %s\n", repo.Description))
		}

		if repo.Language != "" {
			sb.WriteString(fmt.Sprintf("    Language: %s\n", repo.Language))
		}

		sb.WriteString(fmt.Sprintf("    Stars: %d", repo.Stars))
		if repo.StarsDelta > 0 {
			sb.WriteString(fmt.Sprintf(" (+%d)", repo.StarsDelta))
		}
		sb.WriteString("\n")

		sb.WriteString(fmt.Sprintf("    Forks: %d", repo.Forks))
		if repo.ForksDelta > 0 {
			sb.WriteString(fmt.Sprintf(" (+%d)", repo.ForksDelta))
		}
		sb.WriteString("\n\n")
	}

	// 页脚
	sb.WriteString("--------------------------------------\n")
	sb.WriteString("Powered by GitHub API\n")

	return sb.String(), nil
}

// HTMLFormatter HTML格式化器
type HTMLFormatter struct{}

// NewHTMLFormatter 创建HTML格式化器
func NewHTMLFormatter() *HTMLFormatter {
	return &HTMLFormatter{}
}

// Format 格式化为HTML
func (f *HTMLFormatter) Format(repos []api.Repository, language string, period string) (string, error) {
	var sb strings.Builder

	// HTML头部
	sb.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GitHub Trending Repositories Report</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Helvetica, Arial, sans-serif;
            line-height: 1.6;
            color: #24292e;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f6f8fa;
        }
        .container {
            background-color: white;
            border-radius: 6px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.12);
            padding: 24px;
        }
        h1 {
            color: #0366d6;
            border-bottom: 2px solid #0366d6;
            padding-bottom: 10px;
            margin-top: 0;
        }
        .meta {
            background-color: #f6f8fa;
            padding: 12px;
            border-radius: 6px;
            margin: 20px 0;
        }
        .meta-item {
            display: inline-block;
            margin-right: 20px;
        }
        .meta-label {
            font-weight: 600;
            color: #586069;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 20px;
        }
        th {
            background-color: #f6f8fa;
            padding: 12px;
            text-align: left;
            font-weight: 600;
            color: #24292e;
            border-bottom: 2px solid #d1d5da;
        }
        td {
            padding: 12px;
            border-bottom: 1px solid #e1e4e8;
        }
        tr:hover {
            background-color: #f6f8fa;
        }
        .rank {
            font-weight: 600;
            color: #0366d6;
            width: 50px;
        }
        .repo-name {
            font-weight: 600;
        }
        .repo-name a {
            color: #0366d6;
            text-decoration: none;
        }
        .repo-name a:hover {
            text-decoration: underline;
        }
        .description {
            color: #586069;
            font-size: 14px;
            margin-top: 4px;
        }
        .stats {
            display: flex;
            gap: 15px;
            font-size: 14px;
        }
        .stat-item {
            color: #586069;
        }
        .stat-delta {
            color: #28a745;
            font-weight: 600;
        }
        .language {
            display: inline-block;
            padding: 2px 8px;
            border-radius: 3px;
            background-color: #f1f8ff;
            color: #0366d6;
            font-size: 12px;
            font-weight: 600;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #e1e4e8;
            color: #586069;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 GitHub Trending Repositories Report</h1>

        <div class="meta">
            <div class="meta-item">
                <span class="meta-label">Language:</span> `)
	sb.WriteString(formatLanguage(language))
	sb.WriteString(`</div>
            <div class="meta-item">
                <span class="meta-label">Period:</span> `)
	sb.WriteString(formatPeriod(period))
	sb.WriteString(`</div>
            <div class="meta-item">
                <span class="meta-label">Generated:</span> `)
	sb.WriteString(time.Now().Format("2006-01-02 15:04:05"))
	sb.WriteString(`</div>
            <div class="meta-item">
                <span class="meta-label">Total:</span> `)
	sb.WriteString(fmt.Sprintf("%d repositories", len(repos)))
	sb.WriteString(`</div>
        </div>

        <table>
            <thead>
                <tr>
                    <th class="rank">#</th>
                    <th>Repository</th>
                    <th>Language</th>
                    <th>Stars</th>
                    <th>Forks</th>
                </tr>
            </thead>
            <tbody>
`)

	// 仓库列表
	for i, repo := range repos {
		sb.WriteString("                <tr>\n")
		sb.WriteString(fmt.Sprintf("                    <td class=\"rank\">%d</td>\n", i+1))

		// 仓库名称和描述
		sb.WriteString("                    <td>\n")
		sb.WriteString(fmt.Sprintf("                        <div class=\"repo-name\"><a href=\"%s\" target=\"_blank\">%s</a></div>\n",
			repo.URL, escapeHTML(repo.RepoName)))
		if repo.Description != "" {
			sb.WriteString(fmt.Sprintf("                        <div class=\"description\">%s</div>\n",
				escapeHTML(repo.Description)))
		}
		sb.WriteString("                    </td>\n")

		// 语言
		sb.WriteString("                    <td>")
		if repo.Language != "" {
			sb.WriteString(fmt.Sprintf("<span class=\"language\">%s</span>", escapeHTML(repo.Language)))
		} else {
			sb.WriteString("-")
		}
		sb.WriteString("</td>\n")

		// Stars
		sb.WriteString(fmt.Sprintf("                    <td>%s", formatNumber(repo.Stars)))
		if repo.StarsDelta > 0 {
			sb.WriteString(fmt.Sprintf(" <span class=\"stat-delta\">(+%s)</span>", formatNumber(repo.StarsDelta)))
		}
		sb.WriteString("</td>\n")

		// Forks
		sb.WriteString(fmt.Sprintf("                    <td>%s", formatNumber(repo.Forks)))
		if repo.ForksDelta > 0 {
			sb.WriteString(fmt.Sprintf(" <span class=\"stat-delta\">(+%s)</span>", formatNumber(repo.ForksDelta)))
		}
		sb.WriteString("</td>\n")

		sb.WriteString("                </tr>\n")
	}

	// HTML尾部
	sb.WriteString(`            </tbody>
        </table>

        <div class="footer">
            <p>Powered by <a href="https://api.github.com" target="_blank">GitHub API</a></p>
        </div>
    </div>
</body>
</html>`)

	return sb.String(), nil
}

// formatLanguage 格式化语言名称
func formatLanguage(language string) string {
	if language == "" || language == "all" {
		return "All Languages"
	}
	return strings.Title(language)
}

// formatPeriod 格式化时间周期
func formatPeriod(period string) string {
	switch period {
	case "daily":
		return "Past 24 Hours"
	case "weekly":
		return "Past Week"
	case "monthly":
		return "Past Month"
	default:
		return period
	}
}

// formatNumber 格式化数字（添加千位分隔符）
func formatNumber(n int) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}

	var result []byte
	for i, c := range []byte(s) {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, c)
	}
	return string(result)
}

// escapeHTML 转义HTML特殊字符
func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}
