package llm

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/joho/godotenv"
)

// ModelProvider 大模型提供商枚举
type ModelProvider string

const (
	Qwen     ModelProvider = "qwen"
	DeepSeek ModelProvider = "deepseek"
	OpenAI   ModelProvider = "openai"
	Kimi     ModelProvider = "kimi"
)

// APIConfig API配置
type APIConfig struct {
	DeepSeekAPIKey string
	OpenAIAPIKey   string
	KimiAPIKey     string
	QwenAPIKey     string
}

// Client 大模型客户端
type Client struct {
	config *APIConfig
}

// NewClient 创建新的大模型客户端
func NewClient() (*Client, error) {
	// 先从系统环境变量获取
	config := &APIConfig{
		DeepSeekAPIKey: os.Getenv("DEEPSEEK_API_KEY"),
		OpenAIAPIKey:   os.Getenv("OPENAI_API_KEY"),
		KimiAPIKey:     os.Getenv("KIMI_API_KEY"),
		QwenAPIKey:     os.Getenv("QWEN_API_KEY"),
	}

	// 再从.env文件获取（如果环境变量中没有）
	if config.DeepSeekAPIKey == "" || config.OpenAIAPIKey == "" || config.KimiAPIKey == "" || config.QwenAPIKey == "" {
		// 加载.env文件
		err := godotenv.Load()
		if err != nil {
			fmt.Printf("未找到.env文件: %v\n", err)
		}

		if config.DeepSeekAPIKey == "" {
			config.DeepSeekAPIKey = os.Getenv("DEEPSEEK_API_KEY")
		}
		if config.OpenAIAPIKey == "" {
			config.OpenAIAPIKey = os.Getenv("OPENAI_API_KEY")
		}
		if config.KimiAPIKey == "" {
			config.KimiAPIKey = os.Getenv("KIMI_API_KEY")
		}
		if config.QwenAPIKey == "" {
			config.QwenAPIKey = os.Getenv("QWEN_API_KEY")
		}
	}

	// 检查是否至少填写了一个API
	if config.DeepSeekAPIKey == "" && config.OpenAIAPIKey == "" && config.KimiAPIKey == "" && config.QwenAPIKey == "" {
		return nil, fmt.Errorf("未填写API，请在系统环境变量或.env文件中填写至少一个大模型API")
	}

	return &Client{config: config}, nil
}

// ParseModel 解析模型字符串，返回提供商和模型名称
func (c *Client) ParseModel(model string) (ModelProvider, string) {
	if model == "" {
		// 默认使用Qwen
		return Qwen, "qwen-turbo-latest"
	}

	// 检查是否是 "提供商/模型" 格式
	if strings.Contains(model, "/") {
		parts := strings.Split(model, "/")
		if len(parts) >= 2 {
			provider := strings.ToLower(parts[0])
			modelName := parts[1]
			
			switch provider {
			case "qwen", "通义千问", "通义":
				return Qwen, modelName
			case "deepseek", "深度求索":
				return DeepSeek, modelName
			case "openai", "open ai", "gpt":
				return OpenAI, modelName
			case "kimi", "moonshot", "月之暗面", "月之":
				return Kimi, modelName
			default:
				// 如果提供商不匹配，尝试根据模型名匹配提供商
				return c.matchProviderByModel(modelName), modelName
			}
		}
	}
	
	// 直接指定模型的情况
	return c.matchProviderByModel(model), model
}

// matchProviderByModel 根据模型名称匹配提供商
func (c *Client) matchProviderByModel(model string) ModelProvider {
	modelLower := strings.ToLower(model)
	
	// Qwen模型匹配
	if strings.Contains(modelLower, "qwen") || strings.Contains(modelLower, "通义") {
		return Qwen
	}
	
	// DeepSeek模型匹配
	if strings.Contains(modelLower, "deepseek") {
		return DeepSeek
	}
	
	// OpenAI模型匹配
	if strings.Contains(modelLower, "gpt") {
		return OpenAI
	}
	
	// Kimi模型匹配
	if strings.Contains(modelLower, "kimi") {
		return Kimi
	}
	
	// 默认返回Qwen
	return Qwen
}

// GetClient 获取指定提供商的OpenAI客户端和模型名称
func (c *Client) GetClient(provider ModelProvider) (*openai.Client, string, error) {
	switch provider {
	case Qwen:
		if c.config.QwenAPIKey == "" {
			return nil, "", fmt.Errorf("未配置Qwen API密钥")
		}
		config := openai.DefaultConfig(c.config.QwenAPIKey)
		config.BaseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
		return openai.NewClientWithConfig(config), "qwen-turbo-latest", nil

	case DeepSeek:
		if c.config.DeepSeekAPIKey == "" {
			return nil, "", fmt.Errorf("未配置DeepSeek API密钥")
		}
		config := openai.DefaultConfig(c.config.DeepSeekAPIKey)
		config.BaseURL = "https://api.deepseek.com/v1"
		return openai.NewClientWithConfig(config), "deepseek-chat", nil

	case OpenAI:
		if c.config.OpenAIAPIKey == "" {
			return nil, "", fmt.Errorf("未配置OpenAI API密钥")
		}
		// OpenAI使用默认配置
		config := openai.DefaultConfig(c.config.OpenAIAPIKey)
		return openai.NewClientWithConfig(config), "gpt-4o", nil

	case Kimi:
		if c.config.KimiAPIKey == "" {
			return nil, "", fmt.Errorf("未配置Kimi API密钥")
		}
		config := openai.DefaultConfig(c.config.KimiAPIKey)
		config.BaseURL = "https://api.moonshot.cn/v1"
		return openai.NewClientWithConfig(config), "kimi-k2-0711-preview", nil

	default:
		return nil, "", fmt.Errorf("不支持的服务商: %s", provider)
	}
}

// StreamChat 流式对话
func (c *Client) StreamChat(ctx context.Context, provider ModelProvider, model string, messages []openai.ChatCompletionMessage) (*openai.ChatCompletionStream, error) {
	client, defaultModel, err := c.GetClient(provider)
	if err != nil {
		return nil, err
	}

	// 如果没有指定模型，则使用默认模型
	if model == "" {
		model = defaultModel
	}

	// 创建流式请求
	req := openai.ChatCompletionRequest{
		Model:    model,
		Messages: messages,
		Stream:   true,
	}

	return client.CreateChatCompletionStream(ctx, req)
}