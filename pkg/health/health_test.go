package health

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/valyala/fasthttp"
)

func TestLivenessHandler(t *testing.T) {
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	req.Header.SetHost("foobar.com")
	req.SetRequestURI("/health/liveness")
	ctx.Init(&req, nil, nil)
	LivenessHandler(&ctx)
	validateResponse(t, &ctx, `{"is_live":true}`)
}

func TestReadinessHandler(t *testing.T) {
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	req.Header.SetHost("foobar.com")
	req.SetRequestURI("/health/readiness")
	ctx.Init(&req, nil, nil)
	ReadinessHandler(&ctx)
	validateResponse(t, &ctx, `{"is_ready":true}`)
}

func validateResponse(t *testing.T, ctx *fasthttp.RequestCtx, expectBody string) {
	var resp fasthttp.Response
	s := ctx.Response.String()
	br := bufio.NewReader(bytes.NewBufferString(s))
	if err := resp.Read(br); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		t.Fatalf("unexpected statusCode: %d. Expecting %d", resp.StatusCode(), fasthttp.StatusOK)
	}
	if string(resp.Header.Peek("Connection")) != "keep-alive" {
		t.Fatalf("unexpected header value: %q. Expecting %q", resp.Header.Peek("Connection"), "keep-alive")
	}
	if string(resp.Header.Peek("Content-Type")) != "application/json" {
		t.Fatalf("unexpected header value: %q. Expecting %q", resp.Header.Peek("Content-Type"), "application/json")
	}
	body := resp.Body()
	if !bytes.Equal(body, []byte(expectBody)) {
		t.Errorf("Body Actual val not equal to expected result")
	}
}
