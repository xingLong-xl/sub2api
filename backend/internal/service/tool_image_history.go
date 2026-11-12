package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

const (
	// ToolImageHistoryMaxItems 是单个用户最多保留的图片生成记录数。
	ToolImageHistoryMaxItems = 20
	// ToolImageHistoryRetentionDays 是服务端图片记录保留天数。
	ToolImageHistoryRetentionDays = 7
	// ToolImageHistoryMaxImages 是单条记录最多保留的图片数。
	ToolImageHistoryMaxImages = 10
	// ToolImageHistoryMaxBytes 是单用户图片历史 JSON 的软上限，避免 Redis 被 base64 图片撑爆。
	ToolImageHistoryMaxBytes = 48 << 20

	toolImageHistoryMaxIDBytes     = 128
	toolImageHistoryMaxPromptBytes = 64 << 10
	toolImageHistoryMaxModelBytes  = 200
	toolImageHistoryMaxImageBytes  = 24 << 20
	toolImageHistoryMaxParamBytes  = 200
)

// ToolImageHistoryImage 是图片生成历史中的一张图片。
// src 可以是 data URL 或上游 URL；服务端不主动下载它，只按用户隔离存储。
type ToolImageHistoryImage struct {
	Src       string `json:"src"`
	IsDataURL bool   `json:"isDataUrl"`
	Mime      string `json:"mime,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Bytes     int    `json:"bytes,omitempty"`
}

// ToolImageHistoryParams 保存生成参数，字段与前端 HistoryParams 对齐。
type ToolImageHistoryParams struct {
	Mode              string `json:"mode"`
	Size              string `json:"size"`
	Count             int    `json:"count"`
	Quality           string `json:"quality"`
	Style             string `json:"style"`
	ResponseFormat    string `json:"responseFormat"`
	Background        string `json:"background"`
	OutputFormat      string `json:"outputFormat"`
	Moderation        string `json:"moderation"`
	OutputCompression string `json:"outputCompression"`
	PartialImages     string `json:"partialImages"`
	InputFidelity     string `json:"inputFidelity"`
}

// ToolImageHistoryItem 是一条图片生成历史记录。
type ToolImageHistoryItem struct {
	ID           string                  `json:"id"`
	Prompt       string                  `json:"prompt"`
	Model        string                  `json:"model"`
	Mode         string                  `json:"mode"`
	Images       []ToolImageHistoryImage `json:"images"`
	Params       ToolImageHistoryParams  `json:"params"`
	DurationMS   int64                   `json:"duration_ms,omitempty"`
	CostEstimate *float64                `json:"cost_estimate,omitempty"`
	CreatedAt    int64                   `json:"createdAt"`
}

// ToolImageHistoryStore 按用户持久化图片生成历史（Redis 实现）。
type ToolImageHistoryStore interface {
	Get(ctx context.Context, userID int64) ([]ToolImageHistoryItem, error)
	Save(ctx context.Context, userID int64, items []ToolImageHistoryItem) error
}

var ErrToolImageHistoryTooLarge = infraerrors.BadRequest(
	"TOOL_IMAGE_HISTORY_TOO_LARGE",
	"image history is too large; reduce image size or history count",
)

func normalizeImageMode(mode string) string {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "edit":
		return "edit"
	default:
		return "generate"
	}
}

func trimParam(value string) string {
	return strings.TrimSpace(value)
}

func validateParamLength(name, value string) error {
	if len(value) > toolImageHistoryMaxParamBytes {
		return fmt.Errorf("%s is too long", name)
	}
	return nil
}

// NormalizeToolImageHistoryItem 校验并清理前端提交的单条图片历史。
func NormalizeToolImageHistoryItem(input ToolImageHistoryItem, now int64) (ToolImageHistoryItem, error) {
	if now <= 0 {
		now = time.Now().UnixMilli()
	}
	id := strings.TrimSpace(input.ID)
	if id == "" {
		return ToolImageHistoryItem{}, fmt.Errorf("id is required")
	}
	if len(id) > toolImageHistoryMaxIDBytes {
		return ToolImageHistoryItem{}, fmt.Errorf("id is too long")
	}
	prompt := strings.TrimSpace(input.Prompt)
	if prompt == "" {
		return ToolImageHistoryItem{}, fmt.Errorf("prompt is required")
	}
	if len(prompt) > toolImageHistoryMaxPromptBytes {
		return ToolImageHistoryItem{}, fmt.Errorf("prompt is too long")
	}
	model := strings.TrimSpace(input.Model)
	if model == "" {
		return ToolImageHistoryItem{}, fmt.Errorf("model is required")
	}
	if len(model) > toolImageHistoryMaxModelBytes {
		return ToolImageHistoryItem{}, fmt.Errorf("model is too long")
	}
	if len(input.Images) == 0 {
		return ToolImageHistoryItem{}, fmt.Errorf("images are required")
	}
	if len(input.Images) > ToolImageHistoryMaxImages {
		return ToolImageHistoryItem{}, fmt.Errorf("too many images: max %d", ToolImageHistoryMaxImages)
	}

	images := make([]ToolImageHistoryImage, 0, len(input.Images))
	for index, image := range input.Images {
		src := strings.TrimSpace(image.Src)
		if src == "" {
			return ToolImageHistoryItem{}, fmt.Errorf("image %d is empty", index)
		}
		if len(src) > toolImageHistoryMaxImageBytes {
			return ToolImageHistoryItem{}, fmt.Errorf("image %d is too large", index)
		}
		mime := strings.TrimSpace(image.Mime)
		if len(mime) > toolImageHistoryMaxParamBytes {
			return ToolImageHistoryItem{}, fmt.Errorf("image %d mime is too long", index)
		}
		images = append(images, ToolImageHistoryImage{
			Src:       src,
			IsDataURL: image.IsDataURL || strings.HasPrefix(src, "data:image/"),
			Mime:      mime,
			Width:     maxImageInt(image.Width, 0),
			Height:    maxImageInt(image.Height, 0),
			Bytes:     maxImageInt(image.Bytes, 0),
		})
	}

	createdAt := input.CreatedAt
	if createdAt <= 0 || createdAt > now {
		createdAt = now
	}

	itemMode := normalizeImageMode(input.Mode)
	params := input.Params
	rawParamsMode := strings.TrimSpace(params.Mode)
	params.Mode = normalizeImageMode(rawParamsMode)
	if rawParamsMode == "" {
		params.Mode = itemMode
	}
	params.Size = trimParam(params.Size)
	params.Quality = trimParam(params.Quality)
	params.Style = trimParam(params.Style)
	params.ResponseFormat = trimParam(params.ResponseFormat)
	params.Background = trimParam(params.Background)
	params.OutputFormat = trimParam(params.OutputFormat)
	params.Moderation = trimParam(params.Moderation)
	params.OutputCompression = trimParam(params.OutputCompression)
	params.PartialImages = trimParam(params.PartialImages)
	params.InputFidelity = trimParam(params.InputFidelity)
	params.Count = clampImageHistoryCount(params.Count)
	for name, value := range map[string]string{
		"size":               params.Size,
		"quality":            params.Quality,
		"style":              params.Style,
		"response_format":    params.ResponseFormat,
		"background":         params.Background,
		"output_format":      params.OutputFormat,
		"moderation":         params.Moderation,
		"output_compression": params.OutputCompression,
		"partial_images":     params.PartialImages,
		"input_fidelity":     params.InputFidelity,
	} {
		if err := validateParamLength(name, value); err != nil {
			return ToolImageHistoryItem{}, err
		}
	}
	if !utf8.ValidString(prompt) {
		return ToolImageHistoryItem{}, fmt.Errorf("prompt is not valid utf-8")
	}

	durationMS := input.DurationMS
	if durationMS < 0 {
		durationMS = 0
	}
	var costEstimate *float64
	if input.CostEstimate != nil && *input.CostEstimate >= 0 {
		value := *input.CostEstimate
		costEstimate = &value
	}

	item := ToolImageHistoryItem{
		ID:           id,
		Prompt:       prompt,
		Model:        model,
		Mode:         itemMode,
		Images:       images,
		Params:       params,
		DurationMS:   durationMS,
		CostEstimate: costEstimate,
		CreatedAt:    createdAt,
	}
	if encodedSize(item) > ToolImageHistoryMaxBytes {
		return ToolImageHistoryItem{}, ErrToolImageHistoryTooLarge
	}
	return item, nil
}

func clampImageHistoryCount(value int) int {
	if value <= 0 {
		return 1
	}
	if value > ToolImageHistoryMaxImages {
		return ToolImageHistoryMaxImages
	}
	return value
}

func encodedSize(value any) int {
	data, err := json.Marshal(value)
	if err != nil {
		return 0
	}
	return len(data)
}

// MaintainToolImageHistory 清理过期记录、重复记录和超出数量/容量上限的旧记录。
// 返回的记录按时间倒序排列。
func MaintainToolImageHistory(input []ToolImageHistoryItem, now int64) ([]ToolImageHistoryItem, bool) {
	if now <= 0 {
		now = time.Now().UnixMilli()
	}
	cutoff := now - int64(ToolImageHistoryRetentionDays)*24*60*60*1000
	changed := false

	seen := make(map[string]struct{}, len(input))
	out := make([]ToolImageHistoryItem, 0, len(input))
	for _, item := range input {
		if item.ID == "" || item.CreatedAt < cutoff {
			changed = true
			continue
		}
		if _, exists := seen[item.ID]; exists {
			changed = true
			continue
		}
		seen[item.ID] = struct{}{}
		if len(item.Images) > ToolImageHistoryMaxImages {
			item.Images = item.Images[len(item.Images)-ToolImageHistoryMaxImages:]
			changed = true
		}
		out = append(out, item)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].CreatedAt > out[j].CreatedAt })
	if len(out) > ToolImageHistoryMaxItems {
		out = out[:ToolImageHistoryMaxItems]
		changed = true
	}

	totalBytes := 0
	for i, item := range out {
		size := encodedSize(item)
		if totalBytes+size > ToolImageHistoryMaxBytes {
			out = out[:i]
			changed = true
			break
		}
		totalBytes += size
	}
	return out, changed
}

func maxImageInt(value, minimum int) int {
	if value < minimum {
		return minimum
	}
	return value
}
