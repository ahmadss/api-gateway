package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	EnableHTTPS  bool
	RedirectHTTP bool
	HTTPPort     int
	HTTPSPort    int
	TLSCert      string
	TLSKey       string
}

type GatewayConfig struct {
	Server         ServerConfig
	ConsulAddr     string
	HystrixTimeout int
}

type ConsulConfig struct {
}

func Load() (*GatewayConfig, error) {
	fmt.Println("📢 Memulai pemuatan konfigurasi...")

	// Muat .env file
	_ = godotenv.Load(".env")

	// Fungsi bantu untuk parsing bool dan int
	parseBool := func(envKey string, defaultVal bool) bool {
		val := os.Getenv(envKey)
		if val == "" {
			return defaultVal
		}
		parsed, err := strconv.ParseBool(val)
		if err != nil {
			fmt.Printf("⚠️  Kesalahan parsing bool %s: %v. Gunakan default: %v\n", envKey, err, defaultVal)
			return defaultVal
		}
		return parsed
	}

	parseInt := func(envKey string, defaultVal int) int {
		val := os.Getenv(envKey)
		if val == "" {
			return defaultVal
		}
		parsed, err := strconv.Atoi(val)
		if err != nil {
			fmt.Printf("⚠️  Kesalahan parsing int %s: %v. Gunakan default: %d\n", envKey, err, defaultVal)
			return defaultVal
		}
		return parsed
	}

	cfg := &GatewayConfig{
		Server: ServerConfig{
			EnableHTTPS:  parseBool("ENABLE_HTTPS", false),
			RedirectHTTP: parseBool("REDIRECT_HTTP", false),
			HTTPPort:     parseInt("HTTP_PORT", 8080),
			HTTPSPort:    parseInt("HTTPS_PORT", 8445),
			TLSCert:      os.Getenv("TLS_CERT"),
			TLSKey:       os.Getenv("TLS_KEY"),
		},
		ConsulAddr:     os.Getenv("CONSUL_ADDR"),
		HystrixTimeout: parseInt("HYSTRIX_TIMEOUT", 5000),
	}

	if cfg.ConsulAddr == "" {
		cfg.ConsulAddr = "consul-agent-gateway:8500"
	}

	fmt.Printf("✅ Konfigurasi berhasil dimuat: %+v\n", cfg)
	return cfg, nil
}
