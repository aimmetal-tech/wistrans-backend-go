package store

import (
	"database/sql"

	"github.com/aimmetal-tech/wistrans-backend/models"
	"github.com/sashabaranov/go-openai"
)

// SessionStore 会话存储接口
type SessionStore struct {
	DB *sql.DB
}

// NewSessionStore 创建新的会话存储实例
func NewSessionStore(db *sql.DB) *SessionStore {
	return &SessionStore{DB: db}
}

// CreateConversation 创建新会话
func (s *SessionStore) CreateConversation(conversation *models.Conversation) error {
	_, err := s.DB.Exec(`
		INSERT INTO conversations (conversation_id, title, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, conversation.ConversationID, conversation.Title, conversation.User_id, conversation.CreatedAt, conversation.UpdatedAt)
	return err
}

// UpdateConversation 更新会话
func (s *SessionStore) UpdateConversation(conversation *models.Conversation) error {
	_, err := s.DB.Exec(`
		UPDATE conversations
		SET title = $1, updated_at = $2
		WHERE conversation_id = $3
	`, conversation.Title, conversation.UpdatedAt, conversation.ConversationID)
	return err
}

// GetConversation 获取会话
func (s *SessionStore) GetConversation(id string) (*models.Conversation, error) {
	conversation := &models.Conversation{}
	err := s.DB.QueryRow(`
		SELECT conversation_id, title, user_id, created_at, updated_at
		FROM conversations
		WHERE conversation_id = $1
	`, id).Scan(&conversation.ConversationID, &conversation.Title, &conversation.User_id, &conversation.CreatedAt, &conversation.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return conversation, nil
}

// CreateMessage 创建消息
func (s *SessionStore) CreateMessage(message *models.Message) error {
	_, err := s.DB.Exec(`
		INSERT INTO messages (conversation_id, role, content, created_at)
		VALUES ($1, $2, $3, $4)
	`, message.ConversationID, message.Role, message.Content, message.CreatedAt)
	return err
}

// GetMessagesByConversationID 获取会话的所有消息
func (s *SessionStore) GetMessagesByConversationID(conversationID string) ([]*models.Message, error) {
	rows, err := s.DB.Query(`
		SELECT message_id, conversation_id, role, content, created_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC
	`, conversationID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		message := &models.Message{}
		err := rows.Scan(&message.MessageID, &message.ConversationID, &message.Role, &message.Content, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// ToChatMessages 转换为OpenAI聊天消息格式数组
func ToChatMessages(messages []*models.Message) []openai.ChatCompletionMessage {
	var chatMessages []openai.ChatCompletionMessage
	for _, msg := range messages {
		chatMessages = append(chatMessages, msg.ToChatMessage())
	}
	return chatMessages
}