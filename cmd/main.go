package main

import (
	"ag-account/internal/circuitbreaker"
	"ag-account/internal/config"
	"ag-account/internal/discovery"
	"ag-account/internal/logger"
	"ag-account/internal/proxy"
	"fmt"
	"log"
	"os"
	"strconv"

	// "sync/atomic"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	discoveryPkg "ag-account/internal/discovery"
	// path ke log_event.pb.go
)

func main() {
	var err error

	// 1️⃣ Load konfigurasi
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("❌ Gagal memuat konfigurasi", err)
	}

	logger.Initialize(zap.InfoLevel)
	defer logger.Sync()

	// 3️⃣ Setup Service Discovery (Consul + Redis)
	consulDiscovery, err := discovery.NewConsul(cfg.ConsulAddr)
	if err != nil {
		logger.Log.Fatal("❌ Gagal inisialisasi Consul", zap.Error(err))
	}

	// 4️⃣ Konfigurasi Circuit Breaker
	cb := circuitbreaker.CircuitBreaker{
		Timeout:               cfg.HystrixTimeout,
		MaxConcurrentRequests: 100,
		ErrorPercentThreshold: 25,
	}
	cb.Configure("api_gateway")

	// conn, err := discovery.NewGRPCConnConsul("account-gateway", "consul-server:8500")
	// if err != nil {
	// 	log.Fatalf("❌ Gagal connect ke account-gateway via gRPC: %v", err)
	// }
	// accountClient := pb.NewAccountGatewayClient(conn)

	// 5️⃣ Membuat Proxy Handler
	handler := &proxy.Handler{
		Discovery: consulDiscovery,
		CB:        cb,
		Config:    cfg,
	}

	registerAPIGatewayToConsul()

	// 6️⃣ Jalankan FastHTTP Server
	log.Println("🚀 FastHTTP server mulai...")
	startFastHTTPServer(handler, cfg, consulDiscovery)
}

func registerAPIGatewayToConsul() {
	// 🔧 Ambil ENV dari .env atau Docker ENV
	consulAddr := os.Getenv("CONSUL_ADDR")
	serviceName := os.Getenv("CONSUL_SERVICE_NAME")
	serviceID := os.Getenv("CONSUL_SERVICE_ID")
	serviceHost := os.Getenv("CONSUL_SERVICE_HOST")
	servicePortStr := os.Getenv("CONSUL_SERVICE_PORT")

	// ✅ Validasi input
	if consulAddr == "" || serviceName == "" || serviceID == "" || serviceHost == "" || servicePortStr == "" {
		log.Fatal("❌ Semua variabel ENV CONSUL_ADDR, CONSUL_SERVICE_NAME, CONSUL_SERVICE_ID, CONSUL_SERVICE_HOST, dan CONSUL_SERVICE_PORT wajib diisi")
	}

	servicePort, err := strconv.Atoi(servicePortStr)
	if err != nil {
		log.Fatalf("❌ Port tidak valid: %v", err)
	}

	// 📦 Buat config dan client Consul
	config := api.DefaultConfig()
	config.Address = consulAddr

	// 🔐 Otomatis baca CONSUL_HTTP_TOKEN dari ENV
	// Atau bisa diset manual: config.Token = os.Getenv("CONSUL_HTTP_TOKEN")

	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("❌ Gagal membuat client Consul: %v", err)
	}

	// 📡 Definisi service yang akan diregistrasikan
	reg := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: serviceHost,
		Port:    servicePort,
		Check: &api.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d/health", serviceHost, servicePort),
			Interval:                       "10s",
			Timeout:                        "2s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}

	// 🚀 Register ke Consul
	if err := client.Agent().ServiceRegister(reg); err != nil {
		log.Fatalf("❌ Gagal register service ke Consul: %v", err)
	}

	log.Printf("✅ API Gateway %s berhasil register ke Consul di %s", serviceID, consulAddr)
}

func startFastHTTPServer(handler *proxy.Handler, cfg *config.GatewayConfig, discovery *discoveryPkg.ConsulDiscovery) {
	server := &fasthttp.Server{
		Handler:               fastHTTPHandler(handler, cfg, discovery),
		ReadTimeout:           1 * time.Second,
		WriteTimeout:          1 * time.Second,
		MaxConnsPerIP:         1000,
		MaxRequestsPerConn:    100,
		NoDefaultServerHeader: true,
		DisableKeepalive:      false, // aktifkan keepalive
		GetOnly:               false,
		ReduceMemoryUsage:     true,
		Concurrency:           0,           // default unlimited, sesuaikan dengan CPU core jika perlu
		MaxRequestBodySize:    1024 * 1024, // kalau tidak perlu besar, batasi!
	}

	address := fmt.Sprintf(":%d", cfg.Server.HTTPPort)
	logger.Log.Info("🚀 Starting FastHTTP server", zap.String("address", address))

	if err := server.ListenAndServe(address); err != nil {
		logger.Log.Fatal("🔥 FastHTTP server gagal berjalan!", zap.Error(err))
	}
}

func fastHTTPHandler(handler *proxy.Handler, cfg *config.GatewayConfig, discovery discoveryPkg.Discovery) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())

		if path == "/health" {
			fmt.Printf("[Gateway] Health check accessed: %s\n", path)
			ctx.Success("application/json", []byte(`{"status":"ok"}`))
			return
		}

		if path == "/auth/health" {
			fmt.Println("[Gateway] Forward ke auth-public-gateway")

			statusCode, body, err := fasthttp.Get(nil, "http://localhost:21003/health")
			if err != nil || statusCode != 200 {
				ctx.SetStatusCode(fasthttp.StatusServiceUnavailable)
				ctx.SetContentType("application/json")
				ctx.SetBodyString(`{"error":"auth-public-gateway not available"}`)
				return
			}

			ctx.SetStatusCode(fasthttp.StatusOK)
			ctx.SetContentType("application/json")
			ctx.SetBody(body)
			return
		}

		start := time.Now()

		// 🪵 Log sebelum proses
		fmt.Printf("[Gateway] Mulai request: %s\n", path)

		handler.Gateway(ctx)

		duration := time.Since(start)
		status := ctx.Response.StatusCode()

		// 🪵 Log sesudah proses
		fmt.Printf("[Gateway] Selesai request: %s | Status: %d | Durasi: %v\n", path, status, duration)

		// Structured log
		logger.Log.Info("📤 Request selesai",
			zap.String("path", path),
			zap.Int("status", status),
			zap.Duration("duration", duration),
		)
	}
}
