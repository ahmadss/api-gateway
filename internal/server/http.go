package server

import "github.com/valyala/fasthttp"

func SecurityHeaders(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	ctx.Response.Header.Set("X-Content-Type-Options", "nosniff")
	ctx.Response.Header.Set("X-Frame-Options", "DENY")
	ctx.Response.Header.Set("Referrer-Policy", "strict-origin-when-cross-origin")
}
