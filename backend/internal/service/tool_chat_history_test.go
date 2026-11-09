package service

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNormalizeToolChatSessions(t *testing.T) {
	now := time.Now().UnixMilli()
	sessions, err := NormalizeToolChatSessions([]ToolChatSession{
		{
			ID:        " s1 ",
			Title:     " hello ",
			CreatedAt: now,
			Messages: []ToolChatMessage{
				{Role: " user ", Content: "hi"},
				{Role: "assistant", Content: "hello"},
			},
		},
	})
	require.NoError(t, err)
	require.Len(t, sessions, 1)
	require.Equal(t, "s1", sessions[0].ID)
	require.Equal(t, "hello", sessions[0].Title)
	require.Equal(t, now, sessions[0].CreatedAt)
	require.Equal(t, "user", sessions[0].Messages[0].Role)
	require.Equal(t, "hi", sessions[0].Messages[0].Content)
	require.GreaterOrEqual(t, sessions[0].Messages[0].CreatedAt, now)
}

func TestNormalizeToolChatSessionsFillsCreatedAt(t *testing.T) {
	before := time.Now().UnixMilli()
	sessions, err := NormalizeToolChatSessions([]ToolChatSession{{
		ID:       "s1",
		Messages: []ToolChatMessage{{Role: "user", Content: "hi"}},
	}})
	require.NoError(t, err)
	require.Len(t, sessions, 1)
	require.GreaterOrEqual(t, sessions[0].CreatedAt, before)
	require.Len(t, sessions[0].Messages, 1)
	require.GreaterOrEqual(t, sessions[0].Messages[0].CreatedAt, before)
}

func TestMaintainToolChatSessionsDropsExpiredMessagesAndSessions(t *testing.T) {
	now := int64(100 * 24 * 60 * 60 * 1000)
	old := now - int64(ToolChatRetentionDays+1)*24*60*60*1000
	recent := now - 24*60*60*1000
	got, changed := MaintainToolChatSessions([]ToolChatSession{
		{
			ID:        "mixed",
			CreatedAt: old,
			Messages: []ToolChatMessage{
				{Role: "user", Content: "old", CreatedAt: old},
				{Role: "assistant", Content: "new", CreatedAt: recent},
			},
		},
		{
			ID:        "expired",
			CreatedAt: old,
			Messages:  []ToolChatMessage{{Role: "user", Content: "gone", CreatedAt: old}},
		},
		{ID: "empty", CreatedAt: now},
	}, now)

	require.True(t, changed)
	require.Len(t, got, 1)
	require.Equal(t, "mixed", got[0].ID)
	require.Len(t, got[0].Messages, 1)
	require.Equal(t, "new", got[0].Messages[0].Content)
}

func TestMaintainToolChatSessionsFillsAndClampsCreatedAt(t *testing.T) {
	now := int64(100 * 24 * 60 * 60 * 1000)
	got, changed := MaintainToolChatSessions([]ToolChatSession{{
		ID:        "s1",
		CreatedAt: now,
		Messages: []ToolChatMessage{
			{Role: "user", Content: "missing"},
			{Role: "assistant", Content: "future", CreatedAt: now + 1000},
		},
	}}, now)
	require.True(t, changed)
	require.Len(t, got, 1)
	require.Equal(t, now, got[0].Messages[0].CreatedAt)
	require.Equal(t, now, got[0].Messages[1].CreatedAt)
}

func TestNormalizeToolChatSessionsTruncatesLongHistory(t *testing.T) {
	messages := make([]ToolChatMessage, ToolChatMaxMessagesPerSession+1)
	for i := range messages {
		messages[i] = ToolChatMessage{Role: "user", Content: fmt.Sprintf("m%d", i)}
	}
	sessions, err := NormalizeToolChatSessions([]ToolChatSession{{ID: "s1", Messages: messages}})
	require.NoError(t, err)
	require.Len(t, sessions[0].Messages, ToolChatMaxMessagesPerSession)
	require.Equal(t, "m1", sessions[0].Messages[0].Content)
	require.Equal(t, "m200", sessions[0].Messages[ToolChatMaxMessagesPerSession-1].Content)
}

func TestNormalizeToolChatSessionsRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name     string
		sessions []ToolChatSession
		want     string
	}{
		{
			name:     "empty id",
			sessions: []ToolChatSession{{ID: " "}},
			want:     "id is required",
		},
		{
			name: "duplicate id",
			sessions: []ToolChatSession{
				{ID: "s1"},
				{ID: "s1"},
			},
			want: "duplicate id",
		},
		{
			name: "invalid role",
			sessions: []ToolChatSession{{
				ID:       "s1",
				Messages: []ToolChatMessage{{Role: "system", Content: "x"}},
			}},
			want: "invalid role",
		},
		{
			name: "too many sessions",
			sessions: func() []ToolChatSession {
				out := make([]ToolChatSession, ToolChatMaxSessions+1)
				for i := range out {
					out[i].ID = "s" + strings.Repeat("x", i+1)
				}
				return out
			}(),
			want: "too many sessions",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NormalizeToolChatSessions(test.sessions)
			require.Error(t, err)
			require.Contains(t, err.Error(), test.want)
		})
	}
}
