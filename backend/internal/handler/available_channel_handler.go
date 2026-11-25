package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"math"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// AvailableChannelHandler 处理用户侧「可用渠道」查询。
//
// 用户侧接口委托 ChannelService.ListAvailable，并在返回前做四层过滤：
//  1. 行过滤：只保留状态为 Active 且与当前用户可访问分组有交集的渠道；
//  2. 分组过滤：渠道的 Groups 只保留用户可访问的那些；
//  3. 平台过滤：普通分组只保留自身平台模型；Composite 分组按渠道已配置的具体模型平台
//     展开。这样既防止普通分组跨平台泄漏，也让 Composite 正确展示其多平台能力；
//  4. 字段白名单：仅返回用户需要的字段（省略 BillingModelSource / RestrictModels
//     / 内部 ID / Status 等管理字段）。
type AvailableChannelHandler struct {
	channelService *service.ChannelService
	apiKeyService  *service.APIKeyService
	settingService *service.SettingService
}

// NewAvailableChannelHandler 创建用户侧可用渠道 handler。
func NewAvailableChannelHandler(
	channelService *service.ChannelService,
	apiKeyService *service.APIKeyService,
	settingService *service.SettingService,
) *AvailableChannelHandler {
	return &AvailableChannelHandler{
		channelService: channelService,
		apiKeyService:  apiKeyService,
		settingService: settingService,
	}
}

// providerPricingModel is one model+group row in the hvoy pricing protocol.
type providerPricingModel struct {
	ModelName        string            `json:"model_name"`
	GroupName        string            `json:"group_name"`
	PriceUnit        string            `json:"price_unit,omitempty"`
	InputPrice       *float64          `json:"input_price,omitempty"`
	OutputPrice      *float64          `json:"output_price,omitempty"`
	CacheInputPrice  *float64          `json:"cache_input_price,omitempty"`
	CacheCreatePrice *float64          `json:"cache_create_price,omitempty"`
	CacheCreate1h    *float64          `json:"cache_create_price_1h,omitempty"`
	UnitPrice        *float64          `json:"unit_price,omitempty"`
	Enabled          bool              `json:"enabled"`
	Note             string            `json:"note,omitempty"`
	BillingBasis     map[string]string `json:"billing_basis,omitempty"`
	Variants         []map[string]any  `json:"variants,omitempty"`
}

type providerPricingData struct {
	Currency   string                 `json:"currency"`
	PriceUnit  string                 `json:"price_unit"`
	SiteName   string                 `json:"site_name,omitempty"`
	SiteDomain string                 `json:"site_domain,omitempty"`
	UpdatedAt  string                 `json:"updated_at"`
	Models     []providerPricingModel `json:"models"`
}

type providerPricingResponse struct {
	SchemaVersion string              `json:"schema_version"`
	Success       bool                `json:"success"`
	Message       string              `json:"message"`
	Data          providerPricingData `json:"data"`
}

// ProviderPricing exposes the public hvoy Provider Pricing API.
// GET /api/provider/pricing
func (h *AvailableChannelHandler) providerPricingAuthorized(c *gin.Context) bool {
	secret := strings.TrimSpace(os.Getenv("PROVIDER_PRICING_SECRET"))
	if secret == "" {
		return true
	}
	tsRaw := strings.TrimSpace(c.GetHeader("X-Hvoy-Ts"))
	signRaw := strings.TrimSpace(c.GetHeader("X-Hvoy-Sign"))
	ts, err := strconv.ParseInt(tsRaw, 10, 64)
	if err != nil || tsRaw == "" || signRaw == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"schema_version": "1.1", "success": false, "message": "signature required"})
		c.Abort()
		return false
	}
	now := time.Now().Unix()
	if ts < now-60 || ts > now+60 {
		c.JSON(http.StatusUnauthorized, gin.H{"schema_version": "1.1", "success": false, "message": "timestamp expired"})
		c.Abort()
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(tsRaw))
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(strings.ToLower(signRaw)), []byte(expected)) {
		c.JSON(http.StatusUnauthorized, gin.H{"schema_version": "1.1", "success": false, "message": "invalid signature"})
		c.Abort()
		return false
	}
	return true
}

func (h *AvailableChannelHandler) ProviderPricing(c *gin.Context) {
	if !h.providerPricingAuthorized(c) {
		return
	}
	channels, err := h.channelService.ListAvailable(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	// Internal channel/group prices are already CNY. Do not apply USD/CNY
	// conversion here; USDCNYRate is only for usage display conversion.
	rate := 1.0
	siteName := "TokenPilot"
	siteDomain := ""
	if h.settingService != nil {
		if settings, err := h.settingService.GetPublicSettings(c.Request.Context()); err == nil {
			if strings.TrimSpace(settings.SiteName) != "" {
				siteName = strings.TrimSpace(settings.SiteName)
			}
			if u, err := url.Parse(strings.TrimSpace(settings.APIBaseURL)); err == nil {
				siteDomain = u.Hostname()
			}
		}
	}
	rows := make(map[string]providerPricingModel)
	for _, ch := range channels {
		for _, group := range ch.Groups {
			multiplier := group.RateMultiplier
			if multiplier <= 0 {
				multiplier = 1
			}
			for _, model := range ch.SupportedModels {
				if model.Platform != "" && group.Platform != "" && group.Platform != "composite" && model.Platform != group.Platform {
					continue
				}
				if override := matchProviderGroupPricing(group.ModelPricing, model.Name); override != nil {
					model.Pricing = override
				}
				if row, ok := buildProviderPricingModel(model, group.Name, multiplier, rate); ok {
					key := group.Name + "\x00" + model.Name
					if _, exists := rows[key]; !exists {
						rows[key] = row
					}
				}
			}
		}
	}
	out := make([]providerPricingModel, 0, len(rows))
	for _, row := range rows {
		out = append(out, row)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].GroupName == out[j].GroupName {
			return out[i].ModelName < out[j].ModelName
		}
		return out[i].GroupName < out[j].GroupName
	})
	c.JSON(http.StatusOK, providerPricingResponse{
		SchemaVersion: "1.1",
		Success:       true,
		Message:       "",
		Data: providerPricingData{
			Currency:   "CNY",
			PriceUnit:  "per_1m_tokens",
			SiteName:   siteName,
			SiteDomain: siteDomain,
			UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
			Models:     out,
		},
	})
}

func matchProviderGroupPricing(entries []service.ChannelModelPricing, model string) *service.ChannelModelPricing {
	modelLower := strings.ToLower(strings.TrimSpace(model))
	var wildcard *service.ChannelModelPricing
	for i := range entries {
		entry := &entries[i]
		for _, pattern := range entry.Models {
			p := strings.ToLower(strings.TrimSpace(pattern))
			if p == modelLower {
				return entry
			}
			if strings.HasSuffix(p, "*") && strings.HasPrefix(modelLower, strings.TrimSuffix(p, "*")) {
				wildcard = entry
			}
		}
	}
	return wildcard
}

func buildProviderPricingModel(model service.SupportedModel, groupName string, multiplier, rate float64) (providerPricingModel, bool) {
	p := model.Pricing
	if p == nil {
		return providerPricingModel{}, false
	}
	round := func(v float64) float64 {
		return math.Round(v*rate*multiplier*1e6) / 1e6
	}
	mode := p.BillingMode
	switch mode {
	case service.BillingModePerRequest, service.BillingModeImage, service.BillingModeVideo:
		if p.PerRequestPrice == nil || *p.PerRequestPrice <= 0 {
			return providerPricingModel{}, false
		}
		unit := round(*p.PerRequestPrice)
		priceUnit := "per_call"
		if mode == service.BillingModeVideo {
			priceUnit = "per_call"
		}
		return providerPricingModel{
			ModelName: model.Name,
			GroupName: groupName,
			PriceUnit: priceUnit,
			UnitPrice: &unit,
			Enabled:   true,
		}, true
	default:
		if p.InputPrice == nil || *p.InputPrice <= 0 {
			return providerPricingModel{}, false
		}
		in := round(*p.InputPrice * 1e6)
		row := providerPricingModel{
			ModelName:  model.Name,
			GroupName:  groupName,
			PriceUnit:  "per_1m_tokens",
			InputPrice: &in,
			Enabled:    true,
		}
		if p.OutputPrice != nil {
			v := round(*p.OutputPrice * 1e6)
			row.OutputPrice = &v
		}
		if p.CacheReadPrice != nil {
			v := round(*p.CacheReadPrice * 1e6)
			row.CacheInputPrice = &v
		}
		if p.CacheWritePrice != nil {
			v := round(*p.CacheWritePrice * 1e6)
			row.CacheCreatePrice = &v
		}
		if p.CacheWrite1hPrice != nil {
			v := round(*p.CacheWrite1hPrice * 1e6)
			row.CacheCreate1h = &v
		}
		return row, true
	}
}

// featureEnabled 返回 available-channels 开关是否启用。默认关闭（opt-in）。
func (h *AvailableChannelHandler) featureEnabled(c *gin.Context) bool {
	if h.settingService == nil {
		return false
	}
	return h.settingService.GetAvailableChannelsRuntime(c.Request.Context()).Enabled
}

// userAvailableGroup 用户可见的分组概要（白名单字段）。
//
// 前端据此区分专属 vs 公开分组（IsExclusive）、订阅 vs 标准分组（SubscriptionType，
// 订阅视觉加深），并展示默认倍率与高峰倍率规则；用户专属倍率前端走
// /groups/rates，和 API 密钥页面保持一致。
type userAvailableGroup struct {
	ID                 int64   `json:"id"`
	Name               string  `json:"name"`
	Platform           string  `json:"platform"`
	SubscriptionType   string  `json:"subscription_type"`
	RateMultiplier     float64 `json:"rate_multiplier"`
	PeakRateEnabled    bool    `json:"peak_rate_enabled"`
	PeakStart          string  `json:"peak_start"`
	PeakEnd            string  `json:"peak_end"`
	PeakRateMultiplier float64 `json:"peak_rate_multiplier"`
	IsExclusive        bool    `json:"is_exclusive"`

	// 不对外输出；List 按请求语言把 Name 换成译文后丢弃。
	i18n service.GroupI18n
}

// userSupportedModelPricing 用户可见的定价字段白名单。
type userSupportedModelPricing struct {
	BillingMode                string                   `json:"billing_mode"`
	InputPrice                 *float64                 `json:"input_price"`
	OutputPrice                *float64                 `json:"output_price"`
	CacheWritePrice            *float64                 `json:"cache_write_price"`
	CacheWrite1hPrice          *float64                 `json:"cache_write_1h_price"`
	CacheReadPrice             *float64                 `json:"cache_read_price"`
	ReasoningEffortMultipliers map[string]float64       `json:"reasoning_effort_multipliers,omitempty"`
	ImageInputPrice            *float64                 `json:"image_input_price"`
	ImageOutputPrice           *float64                 `json:"image_output_price"`
	PerRequestPrice            *float64                 `json:"per_request_price"`
	Intervals                  []userPricingIntervalDTO `json:"intervals"`
}

// userPricingIntervalDTO 定价区间白名单（去掉内部 ID、SortOrder 等前端不渲染的字段）。
type userPricingIntervalDTO struct {
	MinTokens            int      `json:"min_tokens"`
	MaxTokens            *int     `json:"max_tokens"`
	TierLabel            string   `json:"tier_label,omitempty"`
	InputPrice           *float64 `json:"input_price"`
	OutputPrice          *float64 `json:"output_price"`
	CacheWritePrice      *float64 `json:"cache_write_price"`
	CacheWrite1hPrice    *float64 `json:"cache_write_1h_price"`
	CacheReadPrice       *float64 `json:"cache_read_price"`
	InputMultiplier      *float64 `json:"input_multiplier"`
	OutputMultiplier     *float64 `json:"output_multiplier"`
	CacheWriteMultiplier *float64 `json:"cache_write_multiplier"`
	CacheReadMultiplier  *float64 `json:"cache_read_multiplier"`
	PerRequestPrice      *float64 `json:"per_request_price"`
}

// userSupportedModel 用户可见的支持模型条目。
type userSupportedModel struct {
	Name     string                     `json:"name"`
	Platform string                     `json:"platform"`
	Pricing  *userSupportedModelPricing `json:"pricing"`
}

// userChannelPlatformSection 单渠道内某个平台的子视图：用户可见的分组 + 该平台
// 支持的模型。按 platform 聚合后让前端可以把渠道名作为 row-group 一次渲染，
// 后面的平台行按 sections 顺序铺开。
type userChannelPlatformSection struct {
	Platform        string               `json:"platform"`
	Groups          []userAvailableGroup `json:"groups"`
	SupportedModels []userSupportedModel `json:"supported_models"`
}

// userAvailableChannel 用户可见的渠道条目（白名单字段）。
//
// 每个渠道聚合为一条记录，内嵌 platforms 子数组：每个 section 对应一个平台，
// 包含该平台的 groups 和 supported_models。
type userAvailableChannel struct {
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Platforms   []userChannelPlatformSection `json:"platforms"`
}

// List 列出当前用户可见的「可用渠道」。
// GET /api/v1/channels/available
func (h *AvailableChannelHandler) List(c *gin.Context) {
	subject, ok := middleware.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Feature 未启用时返回空数组（不暴露渠道信息）。检查放在认证之后，
	// 保持与未开关前的 401 行为一致：未登录先 401，登录后再按开关决定。
	if !h.featureEnabled(c) {
		response.Success(c, []userAvailableChannel{})
		return
	}

	userGroups, err := h.apiKeyService.GetAvailableGroups(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	allowedGroupIDs := make(map[int64]struct{}, len(userGroups))
	for i := range userGroups {
		allowedGroupIDs[userGroups[i].ID] = struct{}{}
	}

	channels, err := h.channelService.ListAvailable(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	acceptLanguage := c.GetHeader("Accept-Language")
	out := make([]userAvailableChannel, 0, len(channels))
	for _, ch := range channels {
		if ch.Status != service.StatusActive {
			continue
		}
		visibleGroups := filterUserVisibleGroups(ch.Groups, allowedGroupIDs)
		if len(visibleGroups) == 0 {
			continue
		}
		localizeUserAvailableGroups(visibleGroups, acceptLanguage)
		sections := buildPlatformSections(ch, visibleGroups)
		if len(sections) == 0 {
			continue
		}
		out = append(out, userAvailableChannel{
			Name:        ch.Name,
			Description: ch.Description,
			Platforms:   sections,
		})
	}

	response.Success(c, out)
}

// localizeUserAvailableGroups 把分组名称就地换成请求首选语言的版本。
func localizeUserAvailableGroups(groups []userAvailableGroup, acceptLanguage string) {
	for i := range groups {
		groups[i].Name = groups[i].i18n.LocalizeName(acceptLanguage, groups[i].Name)
	}
}

// buildPlatformSections 把一个渠道按 visibleGroups 的平台集合拆成有序的 section 列表：
// 每个 section 对应一个具体平台，只包含该平台的 groups 和 supported_models。
//
// Composite 分组可访问渠道中所有已配置的具体平台，因此会被展开到每个有支持模型的
// 平台 section。普通分组仍严格留在自身平台，避免跨平台模型信息泄漏。Composite 渠道
// 尚未配置任何模型时保留 composite section，以便前端继续展示该分组和“未配置模型”状态。
// 输出按 platform 字母序稳定排序，便于前端等效比较与回归测试。
func buildPlatformSections(
	ch service.AvailableChannel,
	visibleGroups []userAvailableGroup,
) []userChannelPlatformSection {
	groupsByPlatform := make(map[string][]userAvailableGroup, 4)
	compositeGroups := make([]userAvailableGroup, 0, 1)
	for _, g := range visibleGroups {
		if g.Platform == "" {
			continue
		}
		if g.Platform == service.PlatformComposite {
			compositeGroups = append(compositeGroups, g)
			continue
		}
		groupsByPlatform[g.Platform] = append(groupsByPlatform[g.Platform], g)
	}

	if len(compositeGroups) > 0 {
		modelPlatforms := make(map[string]struct{}, len(ch.SupportedModels))
		for i := range ch.SupportedModels {
			if platform := ch.SupportedModels[i].Platform; platform != "" {
				modelPlatforms[platform] = struct{}{}
			}
		}
		if len(modelPlatforms) == 0 {
			groupsByPlatform[service.PlatformComposite] = append(
				groupsByPlatform[service.PlatformComposite],
				compositeGroups...,
			)
		} else {
			for platform := range modelPlatforms {
				groupsByPlatform[platform] = append(groupsByPlatform[platform], compositeGroups...)
			}
		}
	}
	if len(groupsByPlatform) == 0 {
		return nil
	}

	platforms := make([]string, 0, len(groupsByPlatform))
	for p := range groupsByPlatform {
		platforms = append(platforms, p)
	}
	sort.Strings(platforms)

	sections := make([]userChannelPlatformSection, 0, len(platforms))
	for _, platform := range platforms {
		platformSet := map[string]struct{}{platform: {}}
		sections = append(sections, userChannelPlatformSection{
			Platform:        platform,
			Groups:          groupsByPlatform[platform],
			SupportedModels: toUserSupportedModels(ch.SupportedModels, platformSet),
		})
	}
	return sections
}

// filterUserVisibleGroups 仅保留用户可访问的分组。
func filterUserVisibleGroups(
	groups []service.AvailableGroupRef,
	allowed map[int64]struct{},
) []userAvailableGroup {
	visible := make([]userAvailableGroup, 0, len(groups))
	for _, g := range groups {
		if _, ok := allowed[g.ID]; !ok {
			continue
		}
		visible = append(visible, userAvailableGroup{
			ID:                 g.ID,
			Name:               g.Name,
			i18n:               g.I18n,
			Platform:           g.Platform,
			SubscriptionType:   g.SubscriptionType,
			RateMultiplier:     g.RateMultiplier,
			PeakRateEnabled:    g.PeakRateEnabled,
			PeakStart:          g.PeakStart,
			PeakEnd:            g.PeakEnd,
			PeakRateMultiplier: g.PeakRateMultiplier,
			IsExclusive:        g.IsExclusive,
		})
	}
	return visible
}

// toUserSupportedModels 将 service 层支持模型转换为用户 DTO（字段白名单）。
// 仅保留平台在 allowedPlatforms 中的条目，防止跨平台模型信息泄漏。
// allowedPlatforms 为 nil 时不做平台过滤（保留全部，供测试或明确无过滤场景使用）。
func toUserSupportedModels(
	src []service.SupportedModel,
	allowedPlatforms map[string]struct{},
) []userSupportedModel {
	out := make([]userSupportedModel, 0, len(src))
	for i := range src {
		m := src[i]
		if allowedPlatforms != nil {
			if _, ok := allowedPlatforms[m.Platform]; !ok {
				continue
			}
		}
		out = append(out, userSupportedModel{
			Name:     m.Name,
			Platform: m.Platform,
			Pricing:  toUserPricing(m.Pricing),
		})
	}
	return out
}

// toUserPricingIntervals 将定价区间转换为用户 DTO 白名单形态；nil 入参返回 nil（JSON omitempty 可省略）。
func toUserPricingIntervals(src []service.PricingInterval) []userPricingIntervalDTO {
	if src == nil {
		return nil
	}
	intervals := make([]userPricingIntervalDTO, 0, len(src))
	for _, iv := range src {
		intervals = append(intervals, userPricingIntervalDTO{
			MinTokens:            iv.MinTokens,
			MaxTokens:            iv.MaxTokens,
			TierLabel:            iv.TierLabel,
			InputPrice:           iv.InputPrice,
			OutputPrice:          iv.OutputPrice,
			CacheWritePrice:      iv.CacheWritePrice,
			CacheWrite1hPrice:    iv.CacheWrite1hPrice,
			CacheReadPrice:       iv.CacheReadPrice,
			InputMultiplier:      iv.InputMultiplier,
			OutputMultiplier:     iv.OutputMultiplier,
			CacheWriteMultiplier: iv.CacheWriteMultiplier,
			CacheReadMultiplier:  iv.CacheReadMultiplier,
			PerRequestPrice:      iv.PerRequestPrice,
		})
	}
	return intervals
}

// toUserPricing 将 service 层定价转换为用户 DTO；入参为 nil 时返回 nil。
func toUserPricing(p *service.ChannelModelPricing) *userSupportedModelPricing {
	if p == nil {
		return nil
	}
	intervals := toUserPricingIntervals(p.Intervals)
	if intervals == nil {
		// 用户侧定价的 intervals 固定输出数组（空配置为 []），保持既有契约。
		intervals = []userPricingIntervalDTO{}
	}
	billingMode := string(p.BillingMode)
	if billingMode == "" {
		billingMode = string(service.BillingModeToken)
	}
	return &userSupportedModelPricing{
		BillingMode:                billingMode,
		InputPrice:                 p.InputPrice,
		OutputPrice:                p.OutputPrice,
		CacheWritePrice:            p.CacheWritePrice,
		CacheWrite1hPrice:          p.CacheWrite1hPrice,
		CacheReadPrice:             p.CacheReadPrice,
		ReasoningEffortMultipliers: p.ReasoningEffortMultipliers,
		ImageInputPrice:            p.ImageInputPrice,
		ImageOutputPrice:           p.ImageOutputPrice,
		PerRequestPrice:            p.PerRequestPrice,
		Intervals:                  intervals,
	}
}
