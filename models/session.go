package models

import (
	"time"

	"github.com/sashabaranov/go-openai"
)

// Conversation 会话结构
type Conversation struct {
	ConversationID string    `json:"conversation_id"`
	Title          string    `json:"title"`
	User_id        string    `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Message 消息结构
type Message struct {
	MessageID      int       `json:"message_id"`
	ConversationID string    `json:"conversation_id"`
	Role           string    `json:"role"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}

// Session 会话结构
type Session struct {
	Conversation *Conversation
	History      []openai.ChatCompletionMessage
	Model        string
	Service      string
}

// ToChatMessage 转换为OpenAI聊天消息格式
func (m *Message) ToChatMessage() openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{
		Role:    m.Role,
		Content: m.Content,
	}
}

// FromChatMessage 从OpenAI聊天消息格式转换
func FromChatMessage(conversationID string, msg openai.ChatCompletionMessage) *Message {
	return &Message{
		ConversationID: conversationID,
		Role:           msg.Role,
		Content:        msg.Content,
		CreatedAt:      time.Now(),
	}
}
