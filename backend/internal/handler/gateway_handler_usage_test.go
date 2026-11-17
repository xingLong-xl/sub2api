package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestUsageUnrestrictedIncludesWeeklyWindowStart(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/usage", nil)

	weeklyWindowStart := time.Date(2026, time.July, 13, 0, 30, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	c.Set(string(middleware.ContextKeySubscription), &service.UserSubscription{
		WeeklyWindowStart: &weeklyWindowStart,
	})

	handler := &GatewayHandler{}
	handler.usageUnrestricted(
		c,
		context.Background(),
		&service.APIKey{Group: &service.Group{
			Name:             "Weekly plan",
			SubscriptionType: service.SubscriptionTypeSubscription,
		}},
		middleware.AuthSubject{},
		nil,
		nil,
		nil,
	)

	require.Equal(t, http.StatusOK, recorder.Code)
	var response struct {
		Subscription struct {
			WeeklyWindowStart *time.Time `json:"weekly_window_start"`
		} `json:"subscription"`
	}
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	require.NotNil(t, response.Subscription.WeeklyWindowStart)
	require.True(t, weeklyWindowStart.Equal(*response.Subscription.WeeklyWindowStart))
}

func TestConvertUsageAmountsUSD(t *testing.T) {
	payload := map[string]any{
		"unit":      "CNY",
		"remaining": 72.0,
		"quota": map[string]any{
			"unit":      "CNY",
			"limit":     72.0,
			"used":      0.0,
			"remaining": 72.0,
		},
		"daily_usage": []any{
			map[string]any{"cost": 7.2, "actual_cost": 7.2},
		},
	}
	convertUsageAmounts(payload, 1/7.2, "USD")
	require.Equal(t, "USD", payload["unit"])
	require.InDelta(t, 10.0, payload["remaining"], 0.0001)
	quota := payload["quota"].(map[string]any)
	require.Equal(t, "USD", quota["unit"])
	require.InDelta(t, 10.0, quota["remaining"], 0.0001)
	daily := payload["daily_usage"].([]any)[0].(map[string]any)
	require.InDelta(t, 1.0, daily["cost"], 0.0001)
}
