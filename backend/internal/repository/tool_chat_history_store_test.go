package repository

import (
	"context"
	"testing"
	"time"

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

func TestToolChatPreferencesStoreRoundTripAndIsolation(t *testing.T) {
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = rdb.Close() })
	store := NewToolChatHistoryStore(rdb)
	ctx := context.Background()

	missing, err := store.GetPreferences(ctx, 1)
	require.NoError(t, err)
	require.Nil(t, missing)

	keyID := int64(9)
	want := service.ToolChatPreferences{
		APIKeyID:   &keyID,
		ChatModel:  "gpt-5",
		ImageModel: "gpt-image-1",
	}
	require.NoError(t, store.SavePreferences(ctx, 1, want))
	got, err := store.GetPreferences(ctx, 1)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, want, *got)
	require.Equal(t, time.Duration(0), mr.TTL(toolChatPreferencesKey(1)))

	other, err := store.GetPreferences(ctx, 2)
	require.NoError(t, err)
	require.Nil(t, other)
}
