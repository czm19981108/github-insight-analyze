package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Client GitHub API 客户端
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// Repository 仓库信息
type Repository struct {
	RepoID          int64  `json:"repo_id"`
	RepoName        string `json:"repo_name"`
	FullName        string `json:"full_name"` // GitHub API uses full_name
	Description     string `json:"description"`
	Language        string `json:"language"`
	Stars           int    `json:"stars"`
	StargazersCount int    `json:"stargazers_count"` // GitHub API field
	Forks           int    `json:"forks"`
	ForksCount      int    `json:"forks_count"` // GitHub API field
	Stargazers      int    `json:"stargazers"`
	StarsDelta      int    `json:"stars_delta"`
	ForksDelta      int    `json:"forks_delta"`
	StargazersDelta int    `json:"stargazers_delta"`
	Pushes          int    `json:"pushes"`         // 最近时间段内的 push 数量
	PullRequests    int    `json:"pull_requests"`  // 最近时间段内的 PR 数量
	Rank            int    `json:"rank"`
	URL             string `json:"url"`
	HTMLURL         string `json:"html_url"` // GitHub API uses html_url
	Owner           string `json:"owner"`
}

// TrendingResponse API响应 (OSSInsight format)
type TrendingResponse struct {
	Data []Repository `json:"data"`
}

// OSSInsightSQLResponse OSSInsight SQL endpoint 响应格式
type OSSInsightSQLResponse struct {
	Type string `json:"type"`
	Data struct {
		Columns []struct {
			Col      string `json:"col"`
			DataType string `json:"data_type"`
			Nullable bool   `json:"nullable"`
		} `json:"columns"`
		Rows []OSSInsightRow `json:"rows"`
	} `json:"data"`
}

// OSSInsightRow OSSInsight trending repos 数据行
type OSSInsightRow struct {
	RepoID            string `json:"repo_id"`
	RepoName          string `json:"repo_name"`
	PrimaryLanguage   string `json:"primary_language"`
	Description       string `json:"description"`
	Stars             string `json:"stars"`
	Forks             string `json:"forks"`
	PullRequests      string `json:"pull_requests"`
	Pushes            string `json:"pushes"`
	TotalScore        string `json:"total_score"`
	ContributorLogins string `json:"contributor_logins"`
	CollectionNames   string `json:"collection_names"`
}

// GitHubSearchResponse GitHub Search API响应
type GitHubSearchResponse struct {
	TotalCount        int          `json:"total_count"`
	IncompleteResults bool         `json:"incomplete_results"`
	Items             []GitHubRepo `json:"items"`
}

// GitHubRepo GitHub仓库信息
type GitHubRepo struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	FullName        string `json:"full_name"`
	Description     string `json:"description"`
	Language        string `json:"language"`
	StargazersCount int    `json:"stargazers_count"`
	ForksCount      int    `json:"forks_count"`
	HTMLURL         string `json:"html_url"`
	Owner           struct {
		Login string `json:"login"`
	} `json:"owner"`
}

// NewClient 创建新的API客户端
func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// GetTrendingRepos 获取trending repositories
// language: 编程语言，如 "go", "java", "all"
// period: 时间范围，如 "daily", "weekly", "monthly"
// limit: 获取数量
func (c *Client) GetTrendingRepos(ctx context.Context, language string, period string, limit int) ([]Repository, error) {
	// 构建API URL
	apiURL, err := c.buildTrendingURL(language, period, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "OSS-Insight-Trending-Notifier/1.0")

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// 优先尝试解析 OSSInsight SQL endpoint 响应
	var ossResult OSSInsightSQLResponse
	if err := json.Unmarshal(body, &ossResult); err == nil && ossResult.Type == "sql_endpoint" && len(ossResult.Data.Rows) > 0 {
		// 转换 OSSInsight SQL 格式到统一格式
		repos := make([]Repository, 0, len(ossResult.Data.Rows))
		for i, row := range ossResult.Data.Rows {
			// 解析字符串到整数
			repoID, _ := strconv.ParseInt(row.RepoID, 10, 64)
			stars, _ := strconv.Atoi(row.Stars)
			forks, _ := strconv.Atoi(row.Forks)
			pushes, _ := strconv.Atoi(row.Pushes)
			pullRequests, _ := strconv.Atoi(row.PullRequests)

			// 提取 owner (repo_name 格式为 "owner/repo")
			owner := ""
			if idx := strings.Index(row.RepoName, "/"); idx > 0 {
				owner = row.RepoName[:idx]
			}

			repo := Repository{
				RepoID:          repoID,
				RepoName:        row.RepoName,
				FullName:        row.RepoName,
				Description:     row.Description,
				Language:        row.PrimaryLanguage,
				Stars:           stars,
				StargazersCount: stars,
				Forks:           forks,
				ForksCount:      forks,
				Pushes:          pushes,
				PullRequests:    pullRequests,
				Rank:            i + 1,
				URL:             fmt.Sprintf("https://github.com/%s", row.RepoName),
				HTMLURL:         fmt.Sprintf("https://github.com/%s", row.RepoName),
				Owner:           owner,
			}
			repos = append(repos, repo)

			// 如果达到 limit，停止添加
			if limit > 0 && len(repos) >= limit {
				break
			}
		}
		return repos, nil
	}

	// 尝试解析 GitHub Search API 响应（备用）
	var ghResult GitHubSearchResponse
	if err := json.Unmarshal(body, &ghResult); err == nil && len(ghResult.Items) > 0 {
		// 转换 GitHub 格式到统一格式
		repos := make([]Repository, 0, len(ghResult.Items))
		for i, item := range ghResult.Items {
			repo := Repository{
				RepoID:          item.ID,
				RepoName:        item.FullName,
				FullName:        item.FullName,
				Description:     item.Description,
				Language:        item.Language,
				Stars:           item.StargazersCount,
				StargazersCount: item.StargazersCount,
				Forks:           item.ForksCount,
				ForksCount:      item.ForksCount,
				Rank:            i + 1,
				URL:             item.HTMLURL,
				HTMLURL:         item.HTMLURL,
				Owner:           item.Owner.Login,
			}
			repos = append(repos, repo)
		}
		return repos, nil
	}

	// 尝试解析旧版 OSSInsight 格式（备用）
	var result TrendingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// 设置排名和URL
	for i := range result.Data {
		result.Data[i].Rank = i + 1
		if result.Data[i].URL == "" && result.Data[i].HTMLURL != "" {
			result.Data[i].URL = result.Data[i].HTMLURL
		}
		if result.Data[i].URL == "" {
			result.Data[i].URL = fmt.Sprintf("https://github.com/%s", result.Data[i].RepoName)
		}
		// 标准化 Stars 和 Forks 字段
		if result.Data[i].Stars == 0 && result.Data[i].StargazersCount > 0 {
			result.Data[i].Stars = result.Data[i].StargazersCount
		}
		if result.Data[i].Forks == 0 && result.Data[i].ForksCount > 0 {
			result.Data[i].Forks = result.Data[i].ForksCount
		}
		// 提取owner
		if result.Data[i].Owner == "" && result.Data[i].RepoName != "" {
			if idx := len(result.Data[i].RepoName); idx > 0 {
				for j := 0; j < len(result.Data[i].RepoName); j++ {
					if result.Data[i].RepoName[j] == '/' {
						result.Data[i].Owner = result.Data[i].RepoName[:j]
						break
					}
				}
			}
		}
	}

	return result.Data, nil
}

// buildTrendingURL 构建trending API URL
func (c *Client) buildTrendingURL(language string, period string, limit int) (string, error) {
	// 使用 OSSInsight Trending API 获取真正的 trending 仓库
	// 该 API 返回指定时间段内 star 增长最快的项目
	endpoint := "https://api.ossinsight.io/v1/trends/repos/"

	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}

	// 映射 period 参数到 OSSInsight API 格式
	var ossinsightPeriod string
	switch period {
	case "daily", "past_day", "past_24_hours":
		ossinsightPeriod = "past_24_hours"
	case "weekly", "past_7_days", "past_week":
		ossinsightPeriod = "past_week"
	case "monthly", "past_month", "past_28_days":
		ossinsightPeriod = "past_month"
	case "past_3_months":
		ossinsightPeriod = "past_3_months"
	default:
		ossinsightPeriod = "past_week" // 默认一周
	}

	// 映射 language 参数 - OSSInsight 使用首字母大写
	var ossinsightLanguage string
	if language == "" || language == "all" {
		ossinsightLanguage = "All"
	} else {
		// 将首字母大写（如 "go" -> "Go", "javascript" -> "JavaScript"）
		ossinsightLanguage = strings.ToUpper(language[:1]) + strings.ToLower(language[1:])
		// 特殊处理常见语言名称
		switch strings.ToLower(language) {
		case "javascript":
			ossinsightLanguage = "JavaScript"
		case "typescript":
			ossinsightLanguage = "TypeScript"
		case "c++":
			ossinsightLanguage = "C++"
		case "c#":
			ossinsightLanguage = "C#"
		case "php":
			ossinsightLanguage = "PHP"
		case "html":
			ossinsightLanguage = "HTML"
		case "css":
			ossinsightLanguage = "CSS"
		case "plpgsql":
			ossinsightLanguage = "PLpgSQL"
		case "tsql":
			ossinsightLanguage = "TSQL"
		case "hcl":
			ossinsightLanguage = "HCL"
		case "cmake":
			ossinsightLanguage = "CMake"
		case "powershell":
			ossinsightLanguage = "PowerShell"
		case "matlab":
			ossinsightLanguage = "MATLAB"
		case "objective-c":
			ossinsightLanguage = "Objective-C"
		}
	}

	q := u.Query()
	q.Set("period", ossinsightPeriod)
	q.Set("language", ossinsightLanguage)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// GetCollectionRepos 获取特定collection的repositories
// 这是一个备用方法，如果trending API不可用
func (c *Client) GetCollectionRepos(ctx context.Context, collection string, limit int) ([]Repository, error) {
	endpoint := fmt.Sprintf("%s/collections/%s/repos", c.baseURL, collection)

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("limit", fmt.Sprintf("%d", limit))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "OSS-Insight-Trending-Notifier/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result TrendingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Data, nil
}
