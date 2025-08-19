package models

import "time"

// FetchRequest Fetch请求结构体
type FetchRequest struct {
	URL           string   `json:"url" binding:"required"`   // 要抓取的网址
	ContentType   string   `json:"content_type,omitempty"`   // 内容类型 (news, article, blog等)
	ExtractFields []string `json:"extract_fields,omitempty"` // 要提取的字段
	Language      string   `json:"language,omitempty"`       // 内容语言
	MaxLength     int      `json:"max_length,omitempty"`     // 最大内容长度
}

// FetchResponse Fetch响应结构体
type FetchResponse struct {
	URL           string                 `json:"url"`                      // 原始URL
	Title         string                 `json:"title"`                    // 页面标题
	Content       string                 `json:"content"`                  // 主要内容
	Summary       string                 `json:"summary"`                  // 内容摘要
	ExtractedData map[string]interface{} `json:"extracted_data,omitempty"` // 提取的结构化数据
	Language      string                 `json:"language"`                 // 检测的语言
	FetchTime     time.Time              `json:"fetch_time"`               // 抓取时间
	Status        string                 `json:"status"`                   // 状态 (success, error)
	Error         string                 `json:"error,omitempty"`          // 错误信息
}

// NewsItem 新闻条目结构体
type NewsItem struct {
	Title       string    `json:"title"`                  // 新闻标题
	Content     string    `json:"content"`                // 新闻内容
	Summary     string    `json:"summary"`                // 新闻摘要
	Author      string    `json:"author,omitempty"`       // 作者
	PublishDate string    `json:"publish_date,omitempty"` // 发布日期
	Category    string    `json:"category,omitempty"`     // 分类
	Tags        []string  `json:"tags,omitempty"`         // 标签
	URL         string    `json:"url"`                    // 原始URL
	FetchTime   time.Time `json:"fetch_time"`             // 抓取时间
}

// GetDefaultExtractFields 获取默认提取字段
func GetDefaultExtractFields(contentType string) []string {
	switch contentType {
	case "news":
		return []string{"title", "content", "summary", "author", "publish_date", "category"}
	case "article":
		return []string{"title", "content", "summary", "author", "publish_date"}
	case "blog":
		return []string{"title", "content", "author", "publish_date", "tags"}
	default:
		return []string{"title", "content", "summary"}
	}
}
