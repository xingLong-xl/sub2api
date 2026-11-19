package service

import (
	"errors"
	"strings"
)

// ErrToolHistoryStorageFull 表示 Redis 达到 maxmemory 后拒绝写入。
// 图片/聊天本身可以继续展示，只是服务端历史无法保存。
var ErrToolHistoryStorageFull = errors.New("tool history storage is full")

// IsRedisOOM reports the Redis OOM error returned by noeviction mode.
func IsRedisOOM(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "oom command not allowed") ||
		strings.Contains(message, "used memory > 'maxmemory'") ||
		strings.Contains(message, "used memory > maxmemory")
}
