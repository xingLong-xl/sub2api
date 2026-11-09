package service

import (
	"context"
	"fmt"
	"strings"
	"time"
)

const (
	// ToolChatMaxSessions 是单个用户最多保留的会话数。
	ToolChatMaxSessions = 20
	// ToolChatMaxMessagesPerSession 是单个会话最多保留的消息数。
	ToolChatMaxMessagesPerSession = 200
	// ToolChatContextLimit 是每次聊天请求默认携带的最近上下文条数。
	ToolChatContextLimit = 20
	// ToolChatRetentionDays 是服务端聊天记录保留天数，超过后自动清理。
	ToolChatRetentionDays = 7

	toolChatMaxSessionIDBytes = 128
	toolChatMaxTitleBytes     = 200
	toolChatMaxMessageBytes   = 64 << 10
	toolChatMaxSessionBytes   = 8 << 20
)

// ToolChatMessage 是工具页聊天中的一条消息。只持久化 user/assistant 两种角色；
// system prompt 由前端在每次请求时临时拼装，不写入历史。
type ToolChatMessage struct {
	Role      string `json:"role"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"createdAt,omitempty"`
}

// ToolChatSession 是工具页聊天的一个会话。
type ToolChatSession struct {
	ID        string            `json:"id"`
	Title     string            `json:"title"`
	CreatedAt int64             `json:"createdAt"`
	Messages  []ToolChatMessage `json:"messages"`
}

// ToolChatPreferences 是账号级的工具页连接偏好。
// 只保存 API Key 的 ID，不保存密钥明文；模型名用于跨浏览器恢复选择。
type ToolChatPreferences struct {
	APIKeyID   *int64 `json:"api_key_id"`
	ChatModel  string `json:"chat_model"`
	ImageModel string `json:"image_model"`
	VideoModel string `json:"video_model"`
}

const toolChatMaxModelNameBytes = 200

// NormalizeToolChatPreferences 清理并校验连接偏好，避免超长模型名进入 Redis。
func NormalizeToolChatPreferences(input ToolChatPreferences) (ToolChatPreferences, error) {
	out := ToolChatPreferences{
		ChatModel:  strings.TrimSpace(input.ChatModel),
		ImageModel: strings.TrimSpace(input.ImageModel),
		VideoModel: strings.TrimSpace(input.VideoModel),
	}
	if input.APIKeyID != nil && *input.APIKeyID > 0 {
		id := *input.APIKeyID
		out.APIKeyID = &id
	}
	for name, value := range map[string]string{
		"chat_model":  out.ChatModel,
		"image_model": out.ImageModel,
		"video_model": out.VideoModel,
	} {
		if len(value) > toolChatMaxModelNameBytes {
			return ToolChatPreferences{}, fmt.Errorf("%s is too long", name)
		}
	}
	return out, nil
}

// ToolChatHistoryStore 持久化工具页聊天历史和账号级连接偏好。实现使用 Redis，按用户隔离。
type ToolChatHistoryStore interface {
	Get(ctx context.Context, userID int64) ([]ToolChatSession, error)
	Save(ctx context.Context, userID int64, sessions []ToolChatSession) error
	GetPreferences(ctx context.Context, userID int64) (*ToolChatPreferences, error)
	SavePreferences(ctx context.Context, userID int64, preferences ToolChatPreferences) error
}

// MaintainToolChatSessions 清理超过保留期的消息、丢弃空会话，并返回是否有改动。
// 时间戳缺失的旧数据按会话创建时间补齐；未来时间会被钳制到当前时间。
func MaintainToolChatSessions(input []ToolChatSession, now int64) ([]ToolChatSession, bool) {
	if now <= 0 {
		now = time.Now().UnixMilli()
	}
	cutoff := now - int64(ToolChatRetentionDays)*24*60*60*1000

	out := make([]ToolChatSession, 0, len(input))
	changed := false
	for _, session := range input {
		if session.CreatedAt <= 0 {
			session.CreatedAt = now
			changed = true
		}

		messages := make([]ToolChatMessage, 0, len(session.Messages))
		for _, message := range session.Messages {
			createdAt := message.CreatedAt
			if createdAt <= 0 {
				createdAt = session.CreatedAt
				changed = true
			}
			if createdAt > now {
				createdAt = now
				changed = true
			}
			message.CreatedAt = createdAt
			if createdAt < cutoff {
				changed = true
				continue
			}
			messages = append(messages, message)
		}

		// 没有任何未过期消息的会话整体丢弃，避免空壳长期占用 Redis。
		if len(messages) == 0 {
			changed = true
			continue
		}
		if len(messages) > ToolChatMaxMessagesPerSession {
			messages = messages[len(messages)-ToolChatMaxMessagesPerSession:]
			changed = true
		}
		if len(messages) != len(session.Messages) {
			changed = true
		}
		session.Messages = messages
		out = append(out, session)
	}

	if len(out) > ToolChatMaxSessions {
		out = out[:ToolChatMaxSessions]
		changed = true
	}
	return out, changed
}

// NormalizeToolChatSessions 校验并规范化前端提交的会话列表。
// 这里做服务端兜底，避免客户端把无限长度或非法角色的历史写进 Redis。
func NormalizeToolChatSessions(input []ToolChatSession) ([]ToolChatSession, error) {
	if len(input) > ToolChatMaxSessions {
		return nil, fmt.Errorf("too many sessions: max %d", ToolChatMaxSessions)
	}

	out := make([]ToolChatSession, 0, len(input))
	seenIDs := make(map[string]struct{}, len(input))
	for index := range input {
		source := input[index]
		id := strings.TrimSpace(source.ID)
		if id == "" {
			return nil, fmt.Errorf("session %d: id is required", index)
		}
		if len(id) > toolChatMaxSessionIDBytes {
			return nil, fmt.Errorf("session %d: id is too long", index)
		}
		if _, exists := seenIDs[id]; exists {
			return nil, fmt.Errorf("session %d: duplicate id", index)
		}
		seenIDs[id] = struct{}{}

		title := strings.TrimSpace(source.Title)
		if len(title) > toolChatMaxTitleBytes {
			return nil, fmt.Errorf("session %d: title is too long", index)
		}
		sourceMessages := source.Messages
		if len(sourceMessages) > ToolChatMaxMessagesPerSession {
			// 会话过长时只保留最近的消息，避免保存失败导致整段历史停止同步。
			sourceMessages = sourceMessages[len(sourceMessages)-ToolChatMaxMessagesPerSession:]
		}

		createdAt := source.CreatedAt
		if createdAt <= 0 {
			createdAt = time.Now().UnixMilli()
		}

		messages := make([]ToolChatMessage, 0, len(sourceMessages))
		sessionBytes := 0
		for messageIndex := range sourceMessages {
			message := sourceMessages[messageIndex]
			role := strings.TrimSpace(message.Role)
			if role != "user" && role != "assistant" {
				return nil, fmt.Errorf("session %d message %d: invalid role", index, messageIndex)
			}
			if len(message.Content) > toolChatMaxMessageBytes {
				return nil, fmt.Errorf("session %d message %d: content is too long", index, messageIndex)
			}
			sessionBytes += len(message.Content)
			messages = append(messages, ToolChatMessage{Role: role, Content: message.Content, CreatedAt: message.CreatedAt})
		}
		if sessionBytes > toolChatMaxSessionBytes {
			return nil, fmt.Errorf("session %d: content is too large", index)
		}

		out = append(out, ToolChatSession{
			ID:        id,
			Title:     title,
			CreatedAt: createdAt,
			Messages:  messages,
		})
	}
	maintained, _ := MaintainToolChatSessions(out, time.Now().UnixMilli())
	return maintained, nil
}
