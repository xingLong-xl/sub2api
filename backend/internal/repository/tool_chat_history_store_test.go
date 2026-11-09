package repository

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestToolChatHistoryStoreRoundTripAndIsolation(t *testing.T) {
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = rdb.Close() })
	store := NewToolChatHistoryStore(rdb)
	ctx := context.Background()

	userOne := []service.ToolChatSession{{
		ID:        "s1",
		Title:     "hello",
		CreatedAt: 100,
		Messages: []service.ToolChatMessage{
			{Role: "user", Content: "hi"},
			{Role: "assistant", Content: "hello"},
		},
	}}
	require.NoError(t, store.Save(ctx, 1, userOne))

	got, err := store.Get(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, userOne, got)
	require.True(t, mr.Exists(toolChatHistoryKey(1)))
	require.Equal(t, toolChatHistoryTTL, mr.TTL(toolChatHistoryKey(1)))

	other, err := store.Get(ctx, 2)
	require.NoError(t, err)
	require.Empty(t, other)
}

func TestToolChatHistoryStoreEmptyAndOverwrite(t *testing.T) {
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = rdb.Close() })
	store := NewToolChatHistoryStore(rdb)
	ctx := context.Background()

	missing, err := store.Get(ctx, 7)
	require.NoError(t, err)
	require.Empty(t, missing)

	require.NoError(t, store.Save(ctx, 7, nil))
	empty, err := store.Get(ctx, 7)
	require.NoError(t, err)
	require.Empty(t, empty)

	require.NoError(t, store.Save(ctx, 7, []service.ToolChatSession{{ID: "new"}}))
	got, err := store.Get(ctx, 7)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, "new", got[0].ID)
}
