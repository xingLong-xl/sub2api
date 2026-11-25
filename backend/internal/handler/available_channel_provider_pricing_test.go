package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func providerSignature(secret string, ts int64) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(strconv.FormatInt(ts, 10)))
	return hex.EncodeToString(mac.Sum(nil))
}

func TestProviderPricingAuthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("PROVIDER_PRICING_SECRET", "test-secret")
	h := &AvailableChannelHandler{}

	t.Run("missing", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/provider/pricing", nil)
		require.False(t, h.providerPricingAuthorized(c))
		require.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("valid", func(t *testing.T) {
		ts := time.Now().Unix()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/provider/pricing", nil)
		c.Request.Header.Set("X-Hvoy-Ts", strconv.FormatInt(ts, 10))
		c.Request.Header.Set("X-Hvoy-Sign", providerSignature("test-secret", ts))
		require.True(t, h.providerPricingAuthorized(c))
		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad_signature", func(t *testing.T) {
		ts := time.Now().Unix()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/provider/pricing", nil)
		c.Request.Header.Set("X-Hvoy-Ts", strconv.FormatInt(ts, 10))
		c.Request.Header.Set("X-Hvoy-Sign", "bad")
		require.False(t, h.providerPricingAuthorized(c))
		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
