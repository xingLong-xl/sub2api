package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type toolImageHistoryStoreStub struct {
	items []service.ToolImageHistoryItem
	saved []service.ToolImageHistoryItem
}

func (s *toolImageHistoryStoreStub) Get(_ context.Context, _ int64) ([]service.ToolImageHistoryItem, error) {
	return s.items, nil
}

func (s *toolImageHistoryStoreStub) Save(_ context.Context, _ int64, items []service.ToolImageHistoryItem) error {
	s.saved = items
	return nil
}

func newToolImageHistoryTestContext(method, path string, body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(method, path, bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set(string(middleware2.ContextKeyUser), middleware2.AuthSubject{UserID: 7})
	return c, rec
}

func TestToolImageHistoryHandlerList(t *testing.T) {
	now := time.Now().UnixMilli()
	store := &toolImageHistoryStoreStub{items: []service.ToolImageHistoryItem{{
		ID:        "img-1",
		Prompt:    "cat",
		Model:     "gpt-image-2",
		Mode:      "generate",
		Images:    []service.ToolImageHistoryImage{{Src: "data:image/png;base64,aGVsbG8="}},
		CreatedAt: now,
	}}}
	h := NewToolImageHistoryHandler(store)
	ctx, rec := newToolImageHistoryTestContext(http.MethodGet, "/api/v1/tools/image/history", nil)
	h.ListHistory(ctx)
	require.Equal(t, http.StatusOK, rec.Code)
	var body struct {
		Data struct {
			Items         []service.ToolImageHistoryItem `json:"items"`
			MaxItems      int                            `json:"max_items"`
			RetentionDays int                            `json:"retention_days"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.Len(t, body.Data.Items, 1)
	require.Equal(t, service.ToolImageHistoryMaxItems, body.Data.MaxItems)
	require.Equal(t, service.ToolImageHistoryRetentionDays, body.Data.RetentionDays)
}

func TestToolImageHistoryHandlerSave(t *testing.T) {
	now := time.Now().UnixMilli()
	store := &toolImageHistoryStoreStub{}
	h := NewToolImageHistoryHandler(store)
	payload, err := json.Marshal(service.ToolImageHistoryItem{
		ID:        "img-1",
		Prompt:    "cat",
		Model:     "gpt-image-2",
		Mode:      "generate",
		Images:    []service.ToolImageHistoryImage{{Src: "data:image/png;base64,aGVsbG8="}},
		Params:    service.ToolImageHistoryParams{Size: "1024x1024", Count: 1},
		CreatedAt: now,
	})
	require.NoError(t, err)
	ctx, rec := newToolImageHistoryTestContext(http.MethodPost, "/api/v1/tools/image/history", payload)
	h.SaveHistory(ctx)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Len(t, store.saved, 1)
	require.Equal(t, "img-1", store.saved[0].ID)
}

func TestToolImageHistoryHandlerSaveRejectsEmptyPrompt(t *testing.T) {
	h := NewToolImageHistoryHandler(&toolImageHistoryStoreStub{})
	payload, err := json.Marshal(service.ToolImageHistoryItem{
		ID:        "img-1",
		Model:     "gpt-image-2",
		Images:    []service.ToolImageHistoryImage{{Src: "data:image/png;base64,aGVsbG8="}},
		CreatedAt: time.Now().UnixMilli(),
	})
	require.NoError(t, err)
	ctx, rec := newToolImageHistoryTestContext(http.MethodPost, "/api/v1/tools/image/history", payload)
	h.SaveHistory(ctx)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}
