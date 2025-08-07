package server

import (
	"crypto/tls"
	"fmt"

	"github.com/valyala/fasthttp"
)

// NewTLSConfig 创建 TLS 配置
func NewTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load key pair: %w", err)
	}

	return &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
	}, nil
}

// RedirectHTTPS 处理 HTTP 到 HTTPS 的重定向
func RedirectHTTPS(ctx *fasthttp.RequestCtx, httpsPort int) {
	host := string(ctx.Host())
	path := string(ctx.Path())
	query := ctx.URI().QueryString()

	target := fmt.Sprintf("https://%s%s", host, path)
	if httpsPort != 443 {
		target = fmt.Sprintf("https://%s:%d%s", host, httpsPort, path)
	}
	if len(query) > 0 {
		target += "?" + string(query)
	}

	ctx.Response.Header.Set("Location", target)
	ctx.SetStatusCode(fasthttp.StatusPermanentRedirect)
}
