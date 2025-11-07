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

// Client OSS Insight API 客户端
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// Repository 仓库信息
type Repository struct {
	RepoID          int64   `json:"repo_id"`
	RepoName        string  `json:"repo_name"`
	Description     string  `json:"description"`
	Language        string  `json:"language"`
	Stars           int     `json:"stars"`
	Forks           int     `json:"forks"`
	Stargazers      int     `json:"stargazers"`
	StarsDelta      int     `json:"stars_delta"`
	ForksDelta      int     `json:"forks_delta"`
	StargazersDelta int     `json:"stargazers_delta"`
	Rank            int     `json:"rank"`
	URL             string  `json:"url"`
	Owner           string  `json:"owner"`
}

// TrendingResponse API响应
type TrendingResponse struct {
	Data []Repository `json:"data"`
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

	// 解析响应
	var result TrendingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// 设置排名和URL
	for i := range result.Data {
		result.Data[i].Rank = i + 1
		if result.Data[i].URL == "" {
			result.Data[i].URL = fmt.Sprintf("https://github.com/%s", result.Data[i].RepoName)
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
	// OSS Insight API endpoint
	// 注意：这里使用的是示例endpoint，实际使用时需要根据OSS Insight的真实API调整
	endpoint := fmt.Sprintf("%s/v1/repos/trending", c.baseURL)

	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}

	// 添加查询参数
	q := u.Query()
	if language != "" && language != "all" {
		q.Set("language", language)
	}
	q.Set("period", period)
	q.Set("limit", fmt.Sprintf("%d", limit))
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
