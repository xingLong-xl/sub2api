package handler

import (
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// ToolChatHandler 处理工具页聊天的服务端历史。
type ToolChatHandler struct {
	store service.ToolChatHistoryStore
}

func NewToolChatHandler(store service.ToolChatHistoryStore) *ToolChatHandler {
	return &ToolChatHandler{store: store}
}

type toolChatSessionsResponse struct {
	Sessions      []service.ToolChatSession `json:"sessions"`
	ContextLimit  int                       `json:"context_limit"`
	RetentionDays int                       `json:"retention_days"`
}

type toolChatSessionsRequest struct {
	Sessions []service.ToolChatSession `json:"sessions"`
}

type toolChatPreferencesResponse struct {
	Preferences *service.ToolChatPreferences `json:"preferences"`
}

// ListSessions 返回当前用户的服务端聊天历史。
// GET /api/v1/tools/chat/sessions
func (h *ToolChatHandler) ListSessions(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}

	sessions, err := h.store.Get(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	// 读取时也做一次保留期清理，避免只有读操作的老数据一直留在 Redis。
	sessions, changed := service.MaintainToolChatSessions(sessions, time.Now().UnixMilli())
	if changed {
		_ = h.store.Save(c.Request.Context(), subject.UserID, sessions)
	}
	if sessions == nil {
		sessions = []service.ToolChatSession{}
	}
	response.Success(c, toolChatSessionsResponse{
		Sessions:      sessions,
		ContextLimit:  service.ToolChatContextLimit,
		RetentionDays: service.ToolChatRetentionDays,
	})
}

// SaveSessions 用当前用户提交的会话列表覆盖服务端历史。
// PUT /api/v1/tools/chat/sessions
func (h *ToolChatHandler) SaveSessions(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}

	var body toolChatSessionsRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, "Invalid chat history payload")
		return
	}

	sessions, err := service.NormalizeToolChatSessions(body.Sessions)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.store.Save(c.Request.Context(), subject.UserID, sessions); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"saved": true})
}

// GetPreferences 返回当前账号的工具页连接偏好。
// GET /api/v1/tools/chat/preferences
func (h *ToolChatHandler) GetPreferences(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}
	preferences, err := h.store.GetPreferences(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, toolChatPreferencesResponse{Preferences: preferences})
}

// SavePreferences 保存当前账号的工具页连接偏好，供其他浏览器恢复。
// PUT /api/v1/tools/chat/preferences
func (h *ToolChatHandler) SavePreferences(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}
	var body service.ToolChatPreferences
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, "Invalid tool preferences payload")
		return
	}
	preferences, err := service.NormalizeToolChatPreferences(body)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.store.SavePreferences(c.Request.Context(), subject.UserID, preferences); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, toolChatPreferencesResponse{Preferences: &preferences})
}
