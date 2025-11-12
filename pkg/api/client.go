package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	Rank            int    `json:"rank"`
	URL             string `json:"url"`
	HTMLURL         string `json:"html_url"` // GitHub API uses html_url
	Owner           string `json:"owner"`
}

// TrendingResponse API响应 (OSSInsight format)
type TrendingResponse struct {
	Data []Repository `json:"data"`
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

	// 尝试解析 GitHub Search API 响应
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

	// 尝试解析 OSSInsight 格式
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
	// 由于 OSSInsight API trending端点需要预缓存且不稳定
	// 我们使用 GitHub Search API 来模拟 trending 功能
	// 通过搜索最近pushed且star数高的仓库来获取trending repos

	// 计算日期范围 - 使用 pushed 而不是 created，获取最近活跃的项目
	var pushedAfter string
	switch period {
	case "daily", "past_day":
		pushedAfter = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	case "weekly", "past_7_days":
		pushedAfter = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	case "monthly", "past_month", "past_28_days":
		pushedAfter = time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	default:
		pushedAfter = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	}

	// 构建GitHub Search API URL
	// 搜索最近更新且stars数量较多的仓库
	endpoint := "https://api.github.com/search/repositories"

	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}

	// 构建查询字符串 - 搜索星标数大于50且最近有推送的仓库
	query := fmt.Sprintf("stars:>50 pushed:>%s", pushedAfter)
	if language != "" && language != "all" {
		query += fmt.Sprintf(" language:%s", language)
	}

	q := u.Query()
	q.Set("q", query)
	q.Set("sort", "stars")
	q.Set("order", "desc")
	if limit > 0 && limit <= 100 {
		q.Set("per_page", fmt.Sprintf("%d", limit))
	} else {
		q.Set("per_page", "30")
	}
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
