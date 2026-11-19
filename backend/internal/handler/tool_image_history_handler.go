package handler

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// ToolImageHistoryHandler 处理图片生成工具的服务端历史。
type ToolImageHistoryHandler struct {
	store        service.ToolImageHistoryStore
	imageStorage *service.ImageStorageSettingService
}

func NewToolImageHistoryHandler(
	store service.ToolImageHistoryStore,
	imageStorage *service.ImageStorageSettingService,
) *ToolImageHistoryHandler {
	return &ToolImageHistoryHandler{store: store, imageStorage: imageStorage}
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

func toolImageObjectPrefix(userID int64, itemID string) string {
	sum := sha256.Sum256([]byte(itemID))
	return fmt.Sprintf("tool-images/%d/%x", userID, sum[:8])
}

func findHistoryItem(items []service.ToolImageHistoryItem, id string) *service.ToolImageHistoryItem {
	for i := range items {
		if items[i].ID == id {
			return &items[i]
		}
	}
	return nil
}

func copyPreviousObjectKeys(item *service.ToolImageHistoryItem, previous *service.ToolImageHistoryItem) {
	if item == nil || previous == nil {
		return
	}
	for i := range item.Images {
		if item.Images[i].ObjectKey != "" {
			continue
		}
		for _, oldImage := range previous.Images {
			if oldImage.ObjectKey != "" && oldImage.Src == item.Images[i].Src {
				item.Images[i].ObjectKey = oldImage.ObjectKey
				break
			}
		}
	}
}

func (h *ToolImageHistoryHandler) prepareHistoryImages(
	ctx context.Context,
	userID int64,
	item *service.ToolImageHistoryItem,
	previous *service.ToolImageHistoryItem,
) error {
	if h == nil || h.imageStorage == nil || item == nil {
		return nil
	}
	for i := range item.Images {
		image := &item.Images[i]
		if image.ObjectKey != "" || strings.TrimSpace(image.Src) == "" {
			continue
		}
		url, key, err := h.imageStorage.UploadImage(
			ctx,
			toolImageObjectPrefix(userID, item.ID),
			i,
			image.Src,
		)
		if errors.Is(err, service.ErrImageStorageDisabled) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("upload image history: %w", err)
		}
		image.Src = url
		image.IsDataURL = false
		image.ObjectKey = key
	}
	return nil
}

func imageObjectKeys(items []service.ToolImageHistoryItem) map[string]struct{} {
	keys := make(map[string]struct{})
	for _, item := range items {
		for _, image := range item.Images {
			if image.ObjectKey != "" {
				keys[image.ObjectKey] = struct{}{}
			}
		}
	}
	return keys
}

func (h *ToolImageHistoryHandler) deleteRemovedImageObjects(
	ctx context.Context,
	before, after []service.ToolImageHistoryItem,
) {
	if h == nil || h.imageStorage == nil {
		return
	}
	remaining := imageObjectKeys(after)
	for _, item := range before {
		for _, image := range item.Images {
			if image.ObjectKey == "" {
				continue
			}
			if _, ok := remaining[image.ObjectKey]; ok {
				continue
			}
			if err := h.imageStorage.DeleteImage(ctx, image.ObjectKey); err != nil {
				log.Printf("delete tool image object %s: %v", image.ObjectKey, err)
			}
		}
	}
}

func respondToolImageStorageError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	response.ErrorWithDetails(c, http.StatusServiceUnavailable,
		"image storage upload failed; history was not saved",
		"TOOL_IMAGE_STORAGE_FAILED", nil)
	return true
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
		// 读取历史不应因清理写入失败而失败；Redis 满时仍返回当前可读数据。
		if err := h.store.Save(c.Request.Context(), subject.UserID, maintained); err == nil {
			h.deleteRemovedImageObjects(c.Request.Context(), items, maintained)
		}
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
	before := append([]service.ToolImageHistoryItem(nil), existing...)
	previous := findHistoryItem(before, normalized.ID)
	if previous == nil {
		response.NotFound(c, "Image history record not found")
		return
	}
	copyPreviousObjectKeys(&normalized, previous)
	if err := h.prepareHistoryImages(c.Request.Context(), subject.UserID, &normalized, previous); err != nil {
		respondToolImageStorageError(c, err)
		return
	}
	for index := range existing {
		if existing[index].ID == normalized.ID {
			existing[index] = normalized
			break
		}
	}
	items, _ := service.MaintainToolImageHistory(existing, time.Now().UnixMilli())
	if err := h.store.Save(c.Request.Context(), subject.UserID, items); err != nil {
		// 回滚本次新上传的对象，避免 Redis 保存失败后留下孤儿图片。
		h.deleteRemovedImageObjects(c.Request.Context(), items, before)
		if !respondToolHistoryStorageError(c, err) {
			response.ErrorFrom(c, err)
		}
		return
	}
	h.deleteRemovedImageObjects(c.Request.Context(), before, items)
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
		if !respondToolHistoryStorageError(c, err) {
			response.ErrorFrom(c, err)
		}
		return
	}
	h.deleteRemovedImageObjects(c.Request.Context(), existing, items)
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
	previous := findHistoryItem(existing, normalized.ID)
	copyPreviousObjectKeys(&normalized, previous)
	if err := h.prepareHistoryImages(c.Request.Context(), subject.UserID, &normalized, previous); err != nil {
		respondToolImageStorageError(c, err)
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
		h.deleteRemovedImageObjects(c.Request.Context(), items, existing)
		response.BadRequest(c, "image history is too large; reduce image size or history count")
		return
	}
	if err := h.store.Save(c.Request.Context(), subject.UserID, items); err != nil {
		// 回滚本次新上传的对象，避免 Redis 保存失败后留下孤儿图片。
		h.deleteRemovedImageObjects(c.Request.Context(), items, existing)
		if !respondToolHistoryStorageError(c, err) {
			response.ErrorFrom(c, err)
		}
		return
	}
	h.deleteRemovedImageObjects(c.Request.Context(), existing, items)
	response.Success(c, toolImageHistorySaveResponse{
		Saved:         true,
		MaxItems:      service.ToolImageHistoryMaxItems,
		MaxBytes:      service.ToolImageHistoryMaxBytes,
		RetentionDays: service.ToolImageHistoryRetentionDays,
	})
}
