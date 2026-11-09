package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/redis/go-redis/v9"
)

const toolChatHistoryKeyPrefix = "tool_chat:sessions:"
const toolChatPreferencesKeyPrefix = "tool_chat:preferences:"

// 比保留期多留一天，让长时间不活跃的用户键也能自动过期。
const toolChatHistoryTTL = (service.ToolChatRetentionDays + 1) * 24 * time.Hour

type toolChatHistoryStore struct {
	rdb *redis.Client
}

// NewToolChatHistoryStore 返回 Redis 实现的工具页聊天历史存储。
// 每个用户最多 20 个会话，只保留最近 7 天消息；键本身 8 天后过期。
func NewToolChatHistoryStore(rdb *redis.Client) service.ToolChatHistoryStore {
	return &toolChatHistoryStore{rdb: rdb}
}

func toolChatHistoryKey(userID int64) string {
	return fmt.Sprintf("%s%d", toolChatHistoryKeyPrefix, userID)
}

func toolChatPreferencesKey(userID int64) string {
	return fmt.Sprintf("%s%d", toolChatPreferencesKeyPrefix, userID)
}

func (s *toolChatHistoryStore) Get(ctx context.Context, userID int64) ([]service.ToolChatSession, error) {
	data, err := s.rdb.Get(ctx, toolChatHistoryKey(userID)).Bytes()
	if errors.Is(err, redis.Nil) {
		return []service.ToolChatSession{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get tool chat history: %w", err)
	}
	if len(data) == 0 {
		return []service.ToolChatSession{}, nil
	}
	var sessions []service.ToolChatSession
	if err := json.Unmarshal(data, &sessions); err != nil {
		return nil, fmt.Errorf("decode tool chat history: %w", err)
	}
	if sessions == nil {
		sessions = []service.ToolChatSession{}
	}
	return sessions, nil
}

func (s *toolChatHistoryStore) Save(ctx context.Context, userID int64, sessions []service.ToolChatSession) error {
	if sessions == nil {
		sessions = []service.ToolChatSession{}
	}
	data, err := json.Marshal(sessions)
	if err != nil {
		return fmt.Errorf("encode tool chat history: %w", err)
	}
	if err := s.rdb.Set(ctx, toolChatHistoryKey(userID), data, toolChatHistoryTTL).Err(); err != nil {
		return fmt.Errorf("save tool chat history: %w", err)
	}
	return nil
}

// GetPreferences 返回账号级连接偏好；没有保存过时返回 nil。
func (s *toolChatHistoryStore) GetPreferences(ctx context.Context, userID int64) (*service.ToolChatPreferences, error) {
	data, err := s.rdb.Get(ctx, toolChatPreferencesKey(userID)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get tool chat preferences: %w", err)
	}
	if len(data) == 0 {
		return nil, nil
	}
	var preferences service.ToolChatPreferences
	if err := json.Unmarshal(data, &preferences); err != nil {
		return nil, fmt.Errorf("decode tool chat preferences: %w", err)
	}
	return &preferences, nil
}

// SavePreferences 保存账号级连接偏好，不设置 TTL，跨浏览器长期生效。
func (s *toolChatHistoryStore) SavePreferences(ctx context.Context, userID int64, preferences service.ToolChatPreferences) error {
	data, err := json.Marshal(preferences)
	if err != nil {
		return fmt.Errorf("encode tool chat preferences: %w", err)
	}
	if err := s.rdb.Set(ctx, toolChatPreferencesKey(userID), data, 0).Err(); err != nil {
		return fmt.Errorf("save tool chat preferences: %w", err)
	}
	return nil
}
