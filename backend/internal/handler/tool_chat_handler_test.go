package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	servermiddleware "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type toolChatHistoryStoreStub struct {
	sessions []service.ToolChatSession
	saved    []service.ToolChatSession
	err      error
}

func (s *toolChatHistoryStoreStub) Get(_ context.Context, _ int64) ([]service.ToolChatSession, error) {
	return s.sessions, s.err
}

func (s *toolChatHistoryStoreStub) Save(_ context.Context, _ int64, sessions []service.ToolChatSession) error {
	s.saved = sessions
	return s.err
}

func newToolChatTestContext(method, path, body string) (*gin.Context, *httptest.ResponseRecorder) {
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(method, path, strings.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Set(string(servermiddleware.ContextKeyUser), servermiddleware.AuthSubject{UserID: 42})
	return ctx, rec
}

func TestToolChatHandlerListSessions(t *testing.T) {
	now := time.Now().UnixMilli()
	store := &toolChatHistoryStoreStub{sessions: []service.ToolChatSession{{
		ID:        "s1",
		CreatedAt: now,
		Messages:  []service.ToolChatMessage{{Role: "user", Content: "hi", CreatedAt: now}},
	}}}
	h := NewToolChatHandler(store)
	ctx, rec := newToolChatTestContext(http.MethodGet, "/api/v1/tools/chat/sessions", "")

	h.ListSessions(ctx)

	require.Equal(t, http.StatusOK, rec.Code)
	var body struct {
		Code int `json:"code"`
		Data struct {
			Sessions      []service.ToolChatSession `json:"sessions"`
			ContextLimit  int                       `json:"context_limit"`
			RetentionDays int                       `json:"retention_days"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.Equal(t, 0, body.Code)
	require.Len(t, body.Data.Sessions, 1)
	require.Equal(t, service.ToolChatContextLimit, body.Data.ContextLimit)
	require.Equal(t, service.ToolChatRetentionDays, body.Data.RetentionDays)
}

func TestToolChatHandlerSaveSessions(t *testing.T) {
	store := &toolChatHistoryStoreStub{}
	h := NewToolChatHandler(store)
	now := time.Now().UnixMilli()
	payload := fmt.Sprintf(`{"sessions":[{"id":"s1","title":"hello","createdAt":%d,"messages":[{"role":"user","content":"hi","createdAt":%d}]}]}`, now, now)
	ctx, rec := newToolChatTestContext(http.MethodPut, "/api/v1/tools/chat/sessions", payload)

	h.SaveSessions(ctx)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Len(t, store.saved, 1)
	require.Equal(t, "s1", store.saved[0].ID)
}

func TestToolChatHandlerRejectsInvalidRole(t *testing.T) {
	store := &toolChatHistoryStoreStub{}
	h := NewToolChatHandler(store)
	ctx, rec := newToolChatTestContext(http.MethodPut, "/api/v1/tools/chat/sessions", `{"sessions":[{"id":"s1","messages":[{"role":"system","content":"x"}]}]}`)

	h.SaveSessions(ctx)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Empty(t, store.saved)
}
