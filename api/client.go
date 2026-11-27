package api

import (
	"authentication_service/internal/trace"
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/astaxie/beego/httplib"
)

type Client struct {
	Request *Request
	Type_   string
	// Internal must be explicitly set by service-to-service callers. External
	// providers (Google OAuth, etc.) must not receive correlation metadata.
	Internal bool
}

// SendRequest is retained for non-request callers. It never invents trace IDs.
func (c *Client) SendRequest() (*http.Response, error) {
	return c.SendRequestContext(context.Background())
}

func (c *Client) SendRequestContext(ctx context.Context) (*http.Response, error) {
	fullURL := c.Request.BaseURL + c.Request.Path
	beegoRequest := httplib.NewBeegoRequest(fullURL, c.Request.Method.String())
	for key, value := range c.Request.HeaderField {
		beegoRequest.Header(key, value)
	}
	for key, value := range c.Request.FileField {
		beegoRequest.PostFile(key, value)
	}
	if c.Type_ == "params" {
		for key, value := range c.Request.Params {
			beegoRequest.Param(key, value)
		}
	} else if c.Type_ == "body" {
		var err error
		if c.Request.InterfaceParams == nil {
			_, err = beegoRequest.JSONBody(c.Request.Params)
		} else {
			_, err = beegoRequest.JSONBody(c.Request.InterfaceParams)
		}
		if err != nil {
			return nil, err
		}
	}
	request := beegoRequest.GetRequest()
	*request = *request.WithContext(ctx)
	setTraceHeaders := func(h http.Header, internal bool) {
		h.Del(trace.SessionHeader)
		h.Del(trace.TraceHeader)
		h.Del(trace.RequestHeader)
		if internal {
			trace.Forward(ctx, h)
		}
	}
	setTraceHeaders(request.Header, c.Internal)
	origin, err := url.Parse(fullURL)
	if err != nil {
		return nil, err
	}
	beegoRequest.SetCheckRedirect(func(req *http.Request, via []*http.Request) error {
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		setTraceHeaders(req.Header, c.Internal && req.URL.Scheme == origin.Scheme && req.URL.Host == origin.Host)
		return nil
	})
	beegoRequest.Debug(false)
	response, err := beegoRequest.Response()
	if err != nil {
		// Errors often include URLs, query values, or credentials. Never log them.
		trace.Logger(ctx).Error("outbound HTTP failed", "internal", c.Internal)
		return nil, err
	}
	trace.Logger(ctx).Debug("outbound HTTP completed", "internal", c.Internal, "status", response.StatusCode)
	return response, nil
}
