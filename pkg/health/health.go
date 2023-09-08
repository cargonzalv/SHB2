package health

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

// Handler for processing incoming get request for liveness api
func LivenessHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.SetContentType("application/json")
	fmt.Fprintf(ctx, `{"is_live":true}`)
}

// Handler for processing incoming get request for readiness api
func ReadinessHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.SetContentType("application/json")
	fmt.Fprintf(ctx, `{"is_ready":true}`)
}
