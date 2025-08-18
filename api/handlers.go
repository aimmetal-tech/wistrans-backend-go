package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"github.com/aimmetal-tech/wistrans-backend/llm"
	"github.com/aimmetal-tech/wistrans-backend/models"
	"github.com/aimmetal-tech/wistrans-backend/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

// Handlers API处理函数集合
type Handlers struct {
	Store *store.SessionStore
	LLMClient *llm.Client
}

// NewHandlers 创建新的处理函数实例
func NewHandlers(store *store.SessionStore) (*Handlers, error) {
	// 初始化大模型客户端
	llmClient, err := llm.NewClient()
	if err != nil {
		return nil, err
	}
	
	return &Handlers{
		Store: store,
		LLMClient: llmClient,
	}, nil
}

// HealthCheck 健康检查接口
func (h *Handlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "服务运行正常",
	})
}

// CreateConversation 创建新会话
func (h *Handlers) CreateConversation(c *gin.Context) {
	// 生成新的会话ID
	conversationID := uuid.New().String()
	
	// 创建会话对象
	conversation := &models.Conversation{
		ID: conversationID,
	}
	
	// 保存到数据库
	err := h.Store.CreateConversation(conversation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建会话失败: " + err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"id": conversation.ID,
	})
}

// UpdateConversation 更新会话标题
func (h *Handlers) UpdateConversation(c *gin.Context) {
	conversationID := c.Param("id")
	
	// 获取请求体中的标题
	var req struct {
		Title string `json:"title"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}
	
	// 获取现有会话
	conversation, err := h.Store.GetConversation(conversationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "会话不存在",
		})
		return
	}
	
	// 更新标题
	conversation.Title = req.Title
	
	// 保存到数据库
	err = h.Store.UpdateConversation(conversation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新会话失败: " + err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "会话更新成功",
	})
}

// StreamResponse 流式响应结构体
type StreamResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int     `json:"index"`
		Delta        Delta   `json:"delta"`
		Logprobs     *string `json:"logprobs"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
}

// Delta 增量内容结构体
type Delta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

// TranslateSegment 翻译片段结构体
type TranslateSegment struct {
	ID   string `json:"id" binding:"required"`   // 片段ID，用于标识片段以便后续返回到前端相应位置
	Text string `json:"text" binding:"required"` // 要翻译的文本
}

// TranslateRequest 翻译请求结构体
type TranslateRequest struct {
	Target    string            `json:"target" binding:"required"`     // 目标语言
	Segments  []TranslateSegment `json:"segments" binding:"required"`   // 要翻译的文本片段
	ExtraArgs interface{}       `json:"extra_args,omitempty"`          // 额外参数，例如翻译的风格
}

// TranslateResponse 翻译响应结构体
type TranslateResponse struct {
	Target   string            `json:"target"`  // 目标语言
	Segments []TranslateSegment `json:"segments"` // 翻译后的文本片段
}

// StreamConversation 流式对话接口
func (h *Handlers) StreamConversation(c *gin.Context) {
	// 从查询参数获取参数
	conversationID := c.Query("id")
	input := c.Query("input")
	modelStr := c.Query("model")
	
	// 检查必要参数
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数id不能为空",
		})
		return
	}
	
	if input == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数input不能为空",
		})
		return
	}
	
	// 检查会话是否存在
	conversation, err := h.Store.GetConversation(conversationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "会话不存在",
		})
		return
	}
	
	// 解析模型参数
	provider, model := h.LLMClient.ParseModel(modelStr)
	
	// 获取会话历史消息
	messages, err := h.Store.GetMessagesByConversationID(conversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取会话历史失败: " + err.Error(),
		})
		return
	}
	
	// 转换为聊天消息格式
	chatMessages := store.ToChatMessages(messages)
	
	// 添加用户新消息
	userMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: input,
	}
	chatMessages = append(chatMessages, userMessage)
	
	// 保存用户消息到数据库
	userMsg := models.FromChatMessage(conversationID, userMessage)
	err = h.Store.CreateMessage(userMsg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "保存用户消息失败: " + err.Error(),
		})
		return
	}
	
	// 调用大模型API并流式返回结果
	ctx := context.Background()
	stream, err := h.LLMClient.StreamChat(ctx, provider, model, chatMessages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "调用大模型API失败: " + err.Error(),
		})
		return
	}
	defer stream.Close()
	
	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	
	// 生成响应ID
	responseID := uuid.New().String()
	
	// 流式返回结果
	c.Stream(func(w io.Writer) bool {
		// 发送开始标记
		c.SSEvent("start", gin.H{})
		c.Writer.Flush()
		
		// 流式读取并发送数据
		responseContent := ""
		for {
			chunk, err := stream.Recv()
			if err != nil {
				// 流结束或出错
				break
			}
			
			if len(chunk.Choices) > 0 {
				delta := chunk.Choices[0].Delta
				if delta.Content != "" {
					responseContent += delta.Content
					
					// 构造符合DeepSeek格式的响应
					response := StreamResponse{
						ID:      responseID,
						Object:  "chat.completion.chunk",
						Created: time.Now().Unix(),
						Model:   model,
					}
					
					response.Choices = append(response.Choices, struct {
						Index        int     `json:"index"`
						Delta        Delta   `json:"delta"`
						Logprobs     *string `json:"logprobs"`
						FinishReason *string `json:"finish_reason"`
					}{
						Index: 0,
						Delta: Delta{
							Role:    delta.Role,
							Content: delta.Content,
						},
						Logprobs:     nil,
						FinishReason: nil,
					})
					
					// 发送数据
					c.SSEvent("data", response)
					c.Writer.Flush()
				}
			}
		}
		
		// 发送结束标记，包含finish_reason
		finishReason := "stop"
		response := StreamResponse{
			ID:      responseID,
			Object:  "chat.completion.chunk",
			Created: time.Now().Unix(),
			Model:   model,
		}
		
		response.Choices = append(response.Choices, struct {
			Index        int     `json:"index"`
			Delta        Delta   `json:"delta"`
			Logprobs     *string `json:"logprobs"`
			FinishReason *string `json:"finish_reason"`
		}{
			Index:        0,
			Delta:        Delta{},
			Logprobs:     nil,
			FinishReason: &finishReason,
		})
		
		c.SSEvent("data", response)
		c.Writer.Flush()
		
		// 保存助手消息到数据库
		assistantMessage := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: responseContent,
		}
		assistantMsg := models.FromChatMessage(conversationID, assistantMessage)
		err = h.Store.CreateMessage(assistantMsg)
		if err != nil {
			// 记录错误但不中断流
			fmt.Printf("保存助手消息失败: %v\n", err)
		}
		
		// 如果这是第一条消息，自动生成对话标题
		if len(messages) == 0 {
			// 异步生成并更新标题
			go h.generateAndSetConversationTitle(conversation, chatMessages, responseContent)
		}
		
		// 发送最终结束标记
		c.SSEvent("end", gin.H{})
		c.Writer.Flush()
		
		return false
	})
}

// generateAndSetConversationTitle 生成并设置对话标题
func (h *Handlers) generateAndSetConversationTitle(conversation *models.Conversation, chatMessages []openai.ChatCompletionMessage, responseContent string) {
	// 构造生成标题的提示
	titlePrompt := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("请为以下对话生成一个25字以内的简要标题:\n\n用户: %s\n助手: %s", 
				chatMessages[len(chatMessages)-1].Content, responseContent),
		},
	}
	
	// 获取默认模型用于生成标题
	provider, model := h.LLMClient.ParseModel("")
	
	// 调用大模型API生成标题
	ctx := context.Background()
	client, _, err := h.LLMClient.GetClient(provider)
	if err != nil {
		fmt.Printf("获取大模型客户端失败: %v\n", err)
		return
	}
	
	req := openai.ChatCompletionRequest{
		Model:    model,
		Messages: titlePrompt,
	}
	
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Printf("调用大模型API生成标题失败: %v\n", err)
		return
	}
	
	if len(resp.Choices) > 0 {
		title := strings.TrimSpace(resp.Choices[0].Message.Content)
		// 限制标题长度为25个字符
		if len(title) > 25 {
			// 按字符数截取，而非字节数
			runes := []rune(title)
			if len(runes) > 25 {
				title = string(runes[:25])
			}
		}
		
		// 更新会话标题
		conversation.Title = title
		conversation.UpdatedAt = time.Now()
		err = h.Store.UpdateConversation(conversation)
		if err != nil {
			fmt.Printf("更新会话标题失败: %v\n", err)
			return
		}
	}
}

// Translate 网页翻译接口
// 该函数处理网页翻译请求，接收多个文本片段并翻译为目标语言
func (h *Handlers) Translate(c *gin.Context) {
	var req TranslateRequest
	
	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}
	
	// 准备翻译消息
	// 构造翻译提示词
	prompt := fmt.Sprintf("请将以下内容翻译为%s语言:\n", req.Target)
	for _, segment := range req.Segments {
		prompt += fmt.Sprintf("片段ID %s: %s\n", segment.ID, segment.Text)
	}
	
	// 处理extra_args参数
	if req.ExtraArgs != nil {
		// 将extra_args转换为字符串
		var extraArgsStr string
		switch v := req.ExtraArgs.(type) {
		case string:
			// 如果是字符串，直接使用
			extraArgsStr = v
		case map[string]interface{}:
			// 如果是对象，将其转换为易读的字符串
			if style, ok := v["style"]; ok {
				extraArgsStr = fmt.Sprintf("风格: %v", style)
			} else {
				// 处理其他可能的字段
				for key, value := range v {
					extraArgsStr = fmt.Sprintf("%s %s: %v", extraArgsStr, key, value)
				}
			}
		default:
			// 其他类型转换为字符串
			extraArgsStr = fmt.Sprintf("%v", v)
		}
		
		if extraArgsStr != "" {
			prompt += fmt.Sprintf("翻译要求: %s\n", extraArgsStr)
		}
	}
	
	prompt += "请按照以下JSON格式返回结果，只返回JSON，不要包含其他内容:\n"
	prompt += "{\n  \"target\": \"目标语言\",\n  \"segments\": [\n    {\"id\": \"片段ID\", \"text\": \"翻译后的文本\"}\n  ]\n}"
	
	// 使用默认模型进行翻译（这里使用Qwen）
	provider, model := h.LLMClient.ParseModel("")
	
	// 构造聊天消息
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
	}
	
	// 调用大模型API
	ctx := context.Background()
	client, _, err := h.LLMClient.GetClient(provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取大模型客户端失败: " + err.Error(),
		})
		return
	}
	
	// 创建聊天完成请求
	reqBody := openai.ChatCompletionRequest{
		Model:    model,
		Messages: messages,
	}
	
	// 获取响应
	resp, err := client.CreateChatCompletion(ctx, reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "调用大模型API失败: " + err.Error(),
		})
		return
	}
	
	// 解析响应中的JSON
	var translateResp TranslateResponse
	if len(resp.Choices) > 0 {
		content := resp.Choices[0].Message.Content
		
		// 尝试解析JSON
		if err := json.Unmarshal([]byte(content), &translateResp); err != nil {
			// 如果解析失败，尝试清理内容后再解析
			// 移除可能的代码块标记
			content = strings.TrimPrefix(content, "``json")
			content = strings.TrimSuffix(content, "```")
			content = strings.TrimSpace(content)
			
			if err := json.Unmarshal([]byte(content), &translateResp); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "解析翻译结果失败: " + err.Error(),
					"raw_content": content,
				})
				return
			}
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "大模型未返回有效内容",
		})
		return
	}
	
	c.JSON(http.StatusOK, translateResp)
}

// GetConversationDetail 获取会话详情接口
func (h *Handlers) GetConversationDetail(c *gin.Context) {
	// 从查询参数获取会话ID
	conversationID := c.Query("id")
	
	// 检查必要参数
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数id不能为空",
		})
		return
	}
	
	// 获取会话详情
	conversation, err := h.Store.GetConversation(conversationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "会话不存在",
		})
		return
	}
	
	c.JSON(http.StatusOK, conversation)
}

// GetConversationHistory 获取会话历史记录接口
func (h *Handlers) GetConversationHistory(c *gin.Context) {
	// 从查询参数获取会话ID
	conversationID := c.Query("id")
	
	// 检查必要参数
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数id不能为空",
		})
		return
	}
	
	// 获取会话历史消息
	messages, err := h.Store.GetMessagesByConversationID(conversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取会话历史失败: " + err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"conversation_id": conversationID,
		"messages":        messages,
	})
}
