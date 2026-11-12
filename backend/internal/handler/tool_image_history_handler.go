package handler

import (
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// ToolImageHistoryHandler 处理图片生成工具的服务端历史。
type ToolImageHistoryHandler struct {
	store service.ToolImageHistoryStore
}

func NewToolImageHistoryHandler(store service.ToolImageHistoryStore) *ToolImageHistoryHandler {
	return &ToolImageHistoryHandler{store: store}
}

type toolImageHistoryResponse struct {
	Items         []service.ToolImageHistoryItem `json:"items"`
	MaxItems      int                            `json:"max_items"`
	MaxBytes      int                            `json:"max_bytes"`
	RetentionDays int                            `json:"retention_days"`
}

type toolImageHistorySaveResponse struct {
	Saved         bool `json:"saved"`
	MaxItems      int  `json:"max_items"`
	MaxBytes      int  `json:"max_bytes"`
	RetentionDays int  `json:"retention_days"`
}

// ListHistory 返回当前登录用户的图片生成历史。
// GET /api/v1/tools/image/history
func (h *ToolImageHistoryHandler) ListHistory(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}
	items, err := h.store.Get(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	maintained, changed := service.MaintainToolImageHistory(items, time.Now().UnixMilli())
	if changed {
		_ = h.store.Save(c.Request.Context(), subject.UserID, maintained)
	}
	response.Success(c, toolImageHistoryResponse{
		Items:         maintained,
		MaxItems:      service.ToolImageHistoryMaxItems,
		MaxBytes:      service.ToolImageHistoryMaxBytes,
		RetentionDays: service.ToolImageHistoryRetentionDays,
	})
}

// UpdateHistory 替换一条图片生成记录（用于删除记录中的单张图片）。
// PUT /api/v1/tools/image/history/:id
func (h *ToolImageHistoryHandler) UpdateHistory(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}
	id := strings.TrimSpace(c.Param("id"))
	var item service.ToolImageHistoryItem
	if err := c.ShouldBindJSON(&item); err != nil || id == "" || item.ID != id {
		response.BadRequest(c, "Invalid image history payload")
		return
	}
	normalized, err := service.NormalizeToolImageHistoryItem(item, time.Now().UnixMilli())
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	existing, err := h.store.Get(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	found := false
	for index := range existing {
		if existing[index].ID == normalized.ID {
			existing[index] = normalized
			found = true
			break
		}
	}
	if !found {
		response.NotFound(c, "Image history record not found")
		return
	}
	items, _ := service.MaintainToolImageHistory(existing, time.Now().UnixMilli())
	if err := h.store.Save(c.Request.Context(), subject.UserID, items); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, toolImageHistorySaveResponse{
		Saved:         true,
		MaxItems:      service.ToolImageHistoryMaxItems,
		MaxBytes:      service.ToolImageHistoryMaxBytes,
		RetentionDays: service.ToolImageHistoryRetentionDays,
	})
}

// DeleteHistory 删除一条图片生成记录。
// DELETE /api/v1/tools/image/history/:id
func (h *ToolImageHistoryHandler) DeleteHistory(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}
	id := strings.TrimSpace(c.Param("id"))
	existing, err := h.store.Get(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	items := make([]service.ToolImageHistoryItem, 0, len(existing))
	found := false
	for _, item := range existing {
		if item.ID == id {
			found = true
			continue
		}
		items = append(items, item)
	}
	if !found {
		response.NotFound(c, "Image history record not found")
		return
	}
	if err := h.store.Save(c.Request.Context(), subject.UserID, items); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

// SaveHistory 追加一条图片生成记录。
// POST /api/v1/tools/image/history
func (h *ToolImageHistoryHandler) SaveHistory(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}
	var item service.ToolImageHistoryItem
	if err := c.ShouldBindJSON(&item); err != nil {
		response.BadRequest(c, "Invalid image history payload")
		return
	}
	now := time.Now().UnixMilli()
	normalized, err := service.NormalizeToolImageHistoryItem(item, now)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	existing, err := h.store.Get(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	items, _ := service.MaintainToolImageHistory(append([]service.ToolImageHistoryItem{normalized}, existing...), now)
	found := false
	for _, candidate := range items {
		if candidate.ID == normalized.ID {
			found = true
			break
		}
	}
	if !found {
		response.BadRequest(c, "image history is too large; reduce image size or history count")
		return
	}
	if err := h.store.Save(c.Request.Context(), subject.UserID, items); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, toolImageHistorySaveResponse{
		Saved:         true,
		MaxItems:      service.ToolImageHistoryMaxItems,
		MaxBytes:      service.ToolImageHistoryMaxBytes,
		RetentionDays: service.ToolImageHistoryRetentionDays,
	})
}
