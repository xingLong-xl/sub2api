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

const toolImageHistoryKeyPrefix = "tool_image:history:"

// 比保留期多留一天，让长时间不活跃用户的图片历史自动过期。
const toolImageHistoryTTL = (service.ToolImageHistoryRetentionDays + 1) * 24 * time.Hour

type toolImageHistoryStore struct {
	rdb *redis.Client
}

// NewToolImageHistoryStore 返回 Redis 实现的图片生成历史存储。
func NewToolImageHistoryStore(rdb *redis.Client) service.ToolImageHistoryStore {
	return &toolImageHistoryStore{rdb: rdb}
}

func toolImageHistoryKey(userID int64) string {
	return fmt.Sprintf("%s%d", toolImageHistoryKeyPrefix, userID)
}

func (s *toolImageHistoryStore) Get(ctx context.Context, userID int64) ([]service.ToolImageHistoryItem, error) {
	data, err := s.rdb.Get(ctx, toolImageHistoryKey(userID)).Bytes()
	if errors.Is(err, redis.Nil) {
		return []service.ToolImageHistoryItem{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get tool image history: %w", err)
	}
	if len(data) == 0 {
		return []service.ToolImageHistoryItem{}, nil
	}
	var items []service.ToolImageHistoryItem
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("decode tool image history: %w", err)
	}
	if items == nil {
		items = []service.ToolImageHistoryItem{}
	}
	return items, nil
}

func (s *toolImageHistoryStore) Save(ctx context.Context, userID int64, items []service.ToolImageHistoryItem) error {
	if items == nil {
		items = []service.ToolImageHistoryItem{}
	}
	data, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("encode tool image history: %w", err)
	}
	if err := s.rdb.Set(ctx, toolImageHistoryKey(userID), data, toolImageHistoryTTL).Err(); err != nil {
		return fmt.Errorf("save tool image history: %w", err)
	}
	return nil
}
