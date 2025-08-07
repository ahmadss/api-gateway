package proxy

import (
	"ag-account/internal/config"
	"ag-account/internal/server"
	"log"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/valyala/fasthttp"

	discoveryPkg "ag-account/internal/discovery"
)

// Router menggunakan fasthttp
func NewRouter(handler *Handler, cfg *config.GatewayConfig, jwks *keyfunc.JWKS, discovery discoveryPkg.Discovery) fasthttp.RequestHandler {
	// Middleware chain
	return func(ctx *fasthttp.RequestCtx) {
		// Middleware Security Headers
		server.SecurityHeaders(ctx)

		// Logging Middleware
		// logging.DetailedLoggingMiddlewareFastHTTP(ctx, logger.Log)

		// Health Check Endpoint
		if string(ctx.Path()) == "/health" {
			ctx.SetStatusCode(fasthttp.StatusOK)
			ctx.SetBody([]byte(`{"status": "ok"}`))
			return
		}

		if string(ctx.Path()) == "/debug/request" {
			log.Println("🔍 Header:", ctx.Request.Header.String())
			log.Println("🔍 Body:", string(ctx.Request.Body()))
			ctx.SetStatusCode(200)
			ctx.SetBody([]byte(`{"status":"debug ok"}`))
			return
		}

		// Tangkap semua layanan berbasis nama layanan di awal path
		handler.Gateway(ctx)
	}
}
