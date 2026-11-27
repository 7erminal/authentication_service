package api

import (
	"context"
	"net/http"

	"authentication_service/internal/trace"

	"github.com/beego/beego/v2/client/httplib"
)

type Client struct {
	Request *Request
	// Ctx is the request context. It cancels the outbound call when the
	// inbound request goes away and carries the tracing IDs. Optional: nil
	// means context.Background().
	Ctx context.Context
}

func (c *Client) SendRequest() (*http.Response, error) {
	ctx := c.Ctx
	if ctx == nil {
		ctx = context.Background()
	}
	// TODO: BaseURLとPathのvalidation
	beegoRequest := httplib.NewBeegoRequestWithCtx(ctx,
		c.Request.BaseURL+c.Request.Path,
		c.Request.Method.String())
	for key, value := range c.Request.HeaderField {
		beegoRequest.Header(key, value)
	}
	if c.Request.ForwardTraceHeaders {
		// Forward the SAME X-Request-Id / X-Trace-Id / X-Session-Id we
		// received; never generate new ones for downstream calls.
		trace.InjectHeaders(ctx, beegoRequest.GetRequest().Header)
	}
	for key, value := range c.Request.Params {
		beegoRequest.Param(key, value)
	}
	res, err := beegoRequest.Response()
	if err != nil {
		return nil, err
	}
	return res, nil
}
