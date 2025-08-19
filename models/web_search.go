package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// WebSearchRequest 联网搜索请求结构体
type WebSearchRequest struct {
	Query       string                 `json:"query" binding:"required"` // 搜索查询
	MaxResults  int                    `json:"max_results,omitempty"`    // 最大结果数量
	Language    string                 `json:"language,omitempty"`       // 搜索语言
	Region      string                 `json:"region,omitempty"`         // 搜索地区
	TimeRange   string                 `json:"time_range,omitempty"`     // 时间范围
	ExtraParams map[string]interface{} `json:"extra_params,omitempty"`   // 额外参数
}

// WebSearchResult 搜索结果结构体
type WebSearchResult struct {
	Title       string `json:"title"`        // 标题
	URL         string `json:"url"`          // 链接
	Snippet     string `json:"snippet"`      // 摘要
	Source      string `json:"source"`       // 来源
	PublishedAt string `json:"published_at"` // 发布时间
	Language    string `json:"language"`     // 语言
}

// WebSearchResponse 联网搜索响应结构体
type WebSearchResponse struct {
	Query      string            `json:"query"`           // 搜索查询
	Results    []WebSearchResult `json:"results"`         // 搜索结果
	TotalCount int               `json:"total_count"`     // 总结果数
	SearchTime time.Time         `json:"search_time"`     // 搜索时间
	Status     string            `json:"status"`          // 状态
	Error      string            `json:"error,omitempty"` // 错误信息
}

// WebSearchClient 联网搜索客户端
type WebSearchClient struct {
	APIKey     string
	Endpoint   string
	HTTPClient *http.Client
}

// NewWebSearchClient 创建新的联网搜索客户端
func NewWebSearchClient() (*WebSearchClient, error) {
	// 从环境变量获取API密钥
	apiKey := os.Getenv("QWEN_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("未配置QWEN_API_KEY环境变量")
	}

	return &WebSearchClient{
		APIKey:   apiKey,
		Endpoint: "https://dashscope.aliyuncs.com/api/v1/mcps/WebSearch/sse",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// Search 执行联网搜索
func (c *WebSearchClient) Search(req WebSearchRequest) (*WebSearchResponse, error) {
	// 构造MCP请求
	mcpRequest := map[string]interface{}{
		"method": "search",
		"params": map[string]interface{}{
			"query":       req.Query,
			"max_results": req.MaxResults,
			"language":    req.Language,
			"region":      req.Region,
			"time_range":  req.TimeRange,
		},
	}

	// 添加额外参数
	if req.ExtraParams != nil {
		for key, value := range req.ExtraParams {
			mcpRequest["params"].(map[string]interface{})[key] = value
		}
	}

	// 序列化请求
	requestBody, err := json.Marshal(mcpRequest)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequest("POST", c.Endpoint, strings.NewReader(string(requestBody)))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	// 发送请求
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d", resp.StatusCode)
	}

	// 解析响应
	var mcpResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&mcpResponse); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 构造返回结果
	response := &WebSearchResponse{
		Query:      req.Query,
		SearchTime: time.Now(),
		Status:     "success",
	}

	// 解析搜索结果
	if result, ok := mcpResponse["result"]; ok {
		if resultMap, ok := result.(map[string]interface{}); ok {
			// 解析总结果数
			if totalCount, ok := resultMap["total_count"].(float64); ok {
				response.TotalCount = int(totalCount)
			}

			// 解析搜索结果列表
			if results, ok := resultMap["results"].([]interface{}); ok {
				for _, item := range results {
					if itemMap, ok := item.(map[string]interface{}); ok {
						result := WebSearchResult{}

						if title, ok := itemMap["title"].(string); ok {
							result.Title = title
						}
						if url, ok := itemMap["url"].(string); ok {
							result.URL = url
						}
						if snippet, ok := itemMap["snippet"].(string); ok {
							result.Snippet = snippet
						}
						if source, ok := itemMap["source"].(string); ok {
							result.Source = source
						}
						if publishedAt, ok := itemMap["published_at"].(string); ok {
							result.PublishedAt = publishedAt
						}
						if language, ok := itemMap["language"].(string); ok {
							result.Language = language
						}

						response.Results = append(response.Results, result)
					}
				}
			}
		}
	}

	return response, nil
}

// GetDefaultWebSearchParams 获取默认的搜索参数
func GetDefaultWebSearchParams() WebSearchRequest {
	return WebSearchRequest{
		MaxResults: 10,
		Language:   "zh",
		Region:     "CN",
		TimeRange:  "1y", // 一年内
	}
}
