package proxy

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync/atomic"

	"ag-account/internal/circuitbreaker"
	"ag-account/internal/config"
	"ag-account/internal/discovery"
	"ag-account/internal/logger"

	"github.com/hashicorp/consul/api"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

var lbCounter int64

type Handler struct {
	Discovery discovery.Discovery
	CB        circuitbreaker.CircuitBreaker
	Config    *config.GatewayConfig
}

// ✅ Fungsi utama API Gateway
func (h *Handler) Gateway(ctx *fasthttp.RequestCtx) {
	path := ctx.Path()
	pathStr := strings.TrimPrefix(string(path), "/")
	pathParts := strings.Split(pathStr, "/")

	if len(pathParts) < 1 || pathParts[0] == "" {
		respondWithJSON(ctx, fasthttp.StatusBadRequest, map[string]string{"error": "invalid service path"})
		return
	}

	serviceMap := map[string]string{
		"auth":      "auth-public-gateway",    // Untuk login, register, OTP
		"public":    "content-public-gateway", // Artikel & dokumentasi
		"account":   "account-gateway",        // Data user
		"spiritual": "spiritual-gateway",      // Aktivitas spiritual
		"family":    "family-gateway",         // Modul keluarga
		"consult":   "consult-gateway",        // Konsultasi terbuka
		"move":      "move-gateway",           // Ride/pergerakan
		"commerce":  "commerce-gateway",       // Belanja & produk
		"finance":   "finance-gateway",        // Transaksi & infaq
		"collab":    "collab-gateway",         // Pilar kolaborasi
		"education": "education-gateway",      // Pilar pendidikan
		"notify":    "notification-gateway",   // Email, push, inbox, real-time alert
	}

	serviceKey := pathParts[0]
	serviceName, exists := serviceMap[serviceKey]
	if !exists {
		respondWithJSON(ctx, fasthttp.StatusBadRequest, map[string]string{"error": "invalid service"})
		return
	}

	instances, err := h.Discovery.GetService(serviceName)
	if err != nil || len(instances) == 0 {
		respondWithJSON(ctx, fasthttp.StatusServiceUnavailable, map[string]string{"error": "service unavailable"})
		return
	}

	instance := loadBalanceSmart(instances)
	if instance == nil {
		respondWithJSON(ctx, fasthttp.StatusServiceUnavailable, map[string]string{"error": "no available sub API gateways"})
		return
	}

	targetBaseURL := fmt.Sprintf("http://%s:%d", instance.Service.Address, instance.Service.Port)

	// Ambil path setelah prefix service
	targetPath := "/" + strings.Join(pathParts[1:], "/")
	if ctx.QueryArgs().Len() > 0 {
		targetPath += "?" + ctx.QueryArgs().String()
	}

	err = circuitbreaker.Do("api_gateway", func() error {
		proxyRequest(ctx, targetBaseURL, targetPath)
		return nil
	}, func(e error) error {
		// fallback jika circuit breaker terbuka atau error
		respondWithJSON(ctx, fasthttp.StatusServiceUnavailable, map[string]string{
			"error": "service temporarily unavailable",
		})
		return nil
	})

	if err != nil {
		logger.Log.Error("Proxy error", zap.Error(err))
	}

}

func proxyRequest(ctx *fasthttp.RequestCtx, backendURL, newPath string) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(backendURL + newPath)
	req.Header.SetMethodBytes(ctx.Method())
	req.SetBody(ctx.Request.Body())

	ctx.Request.Header.VisitAll(func(key, value []byte) {
		req.Header.SetBytesKV(key, value)
	})

	if err := fasthttp.Do(req, resp); err != nil {
		logger.Log.Error("Proxy error", zap.String("url", backendURL+newPath), zap.Error(err))
		respondWithJSON(ctx, fasthttp.StatusBadGateway, map[string]string{"error": "failed to reach service"})
		return
	}

	ctx.SetStatusCode(resp.StatusCode())
	ctx.Response.SetBodyRaw(resp.Body()) // gunakan SetBodyRaw agar tidak copy ulang

	resp.Header.VisitAll(func(k, v []byte) {
		ctx.Response.Header.SetBytesKV(k, v)
	})
}

func loadBalanceSmart(instances []*api.ServiceEntry) *api.ServiceEntry {
	var healthy []*api.ServiceEntry
	for _, inst := range instances {
		for _, check := range inst.Checks {
			if check.Status == "passing" {
				healthy = append(healthy, inst)
				break
			}
		}
	}
	if len(healthy) == 0 {
		return nil
	}
	index := atomic.AddInt64(&lbCounter, 1) % int64(len(healthy))
	return healthy[index]
}

func respondWithJSON(ctx *fasthttp.RequestCtx, statusCode int, data map[string]string) {
	ctx.SetStatusCode(statusCode)
	ctx.Response.Header.Set("Content-Type", "application/json")

	if jsonData, err := json.Marshal(data); err == nil {
		ctx.Response.SetBodyRaw(jsonData)
	}
}
