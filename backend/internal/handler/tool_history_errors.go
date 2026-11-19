package handler

import (
	"errors"
	"net/http"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// respondToolHistoryStorageError maps Redis maxmemory/noeviction failures to a
// stable reason so the frontend can tell the user that only history persistence
// failed, while keeping the current chat/image visible.
func respondToolHistoryStorageError(c *gin.Context, err error) bool {
	if !errors.Is(err, service.ErrToolHistoryStorageFull) {
		return false
	}
	response.ErrorWithDetails(c, http.StatusInsufficientStorage,
		"history storage is full; the current item was not saved",
		"TOOL_HISTORY_STORAGE_FULL", nil)
	return true
}
