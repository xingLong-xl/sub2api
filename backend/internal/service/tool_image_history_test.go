package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func testImageHistoryItem(id string, createdAt int64) ToolImageHistoryItem {
	return ToolImageHistoryItem{
		ID:        id,
		Prompt:    "draw a cat",
		Model:     "gpt-image-2",
		Mode:      "generate",
		Images:    []ToolImageHistoryImage{{Src: "data:image/png;base64,aGVsbG8=", IsDataURL: true, Mime: "image/png"}},
		Params:    ToolImageHistoryParams{Mode: "generate", Size: "1024x1024", Count: 1},
		CreatedAt: createdAt,
	}
}

func TestNormalizeToolImageHistoryItem(t *testing.T) {
	now := time.Now().UnixMilli()
	item, err := NormalizeToolImageHistoryItem(ToolImageHistoryItem{
		ID:        " item-1 ",
		Prompt:    "  draw  ",
		Model:     " gpt-image-2 ",
		Mode:      "edit",
		Images:    []ToolImageHistoryImage{{Src: " data:image/png;base64,aGVsbG8= ", Mime: "image/png"}},
		Params:    ToolImageHistoryParams{Size: " 2048x2048 ", Count: 99},
		CreatedAt: now + 1000,
	}, now)
	require.NoError(t, err)
	require.Equal(t, "item-1", item.ID)
	require.Equal(t, "draw", item.Prompt)
	require.Equal(t, "gpt-image-2", item.Model)
	require.Equal(t, "edit", item.Mode)
	require.Equal(t, "edit", item.Params.Mode)
	require.Equal(t, 10, item.Params.Count)
	require.True(t, item.Images[0].IsDataURL)
	require.Equal(t, now, item.CreatedAt)
}

func TestNormalizeToolImageHistoryItemRejectsInvalidInput(t *testing.T) {
	now := time.Now().UnixMilli()
	base := testImageHistoryItem("item", now)
	base.Prompt = ""
	_, err := NormalizeToolImageHistoryItem(base, now)
	require.Error(t, err)

	base = testImageHistoryItem("item", now)
	base.Images = nil
	_, err = NormalizeToolImageHistoryItem(base, now)
	require.Error(t, err)

	base = testImageHistoryItem("item", now)
	base.Images = make([]ToolImageHistoryImage, ToolImageHistoryMaxImages+1)
	for i := range base.Images {
		base.Images[i] = ToolImageHistoryImage{Src: "data:image/png;base64,aGVsbG8="}
	}
	_, err = NormalizeToolImageHistoryItem(base, now)
	require.Error(t, err)
}

func TestMaintainToolImageHistoryLimitsAndRetention(t *testing.T) {
	now := time.Now().UnixMilli()
	items := make([]ToolImageHistoryItem, 0, ToolImageHistoryMaxItems+2)
	for i := 0; i < ToolImageHistoryMaxItems+2; i++ {
		items = append(items, testImageHistoryItem(fmt.Sprintf("item-%d", i), now-int64(i)))
	}
	// 一条过期记录和一条重复记录都必须被清理。
	items = append(items, testImageHistoryItem("expired", now-int64((ToolImageHistoryRetentionDays+1)*24*60*60*1000)))
	items = append(items, testImageHistoryItem("item-0", now))

	got, changed := MaintainToolImageHistory(items, now)
	require.True(t, changed)
	require.Len(t, got, ToolImageHistoryMaxItems)
	require.Equal(t, "item-0", got[0].ID)
	require.NotContains(t, idsOfToolImageHistory(got), "expired")
	require.Equal(t, 1, countToolImageHistoryID(got, "item-0"))
}

func idsOfToolImageHistory(items []ToolImageHistoryItem) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		out = append(out, item.ID)
	}
	return out
}

func countToolImageHistoryID(items []ToolImageHistoryItem, id string) int {
	count := 0
	for _, item := range items {
		if item.ID == id {
			count++
		}
	}
	return count
}
