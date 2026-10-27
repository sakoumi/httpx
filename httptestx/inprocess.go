package httptestx

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	defaultRemoteAddr = "192.168.1.1:4321"
	inProcessBaseURL  = "http://in-process"
)

type handlerTransport struct {
	h          http.Handler
	remoteAddr string
}

func (t handlerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.RemoteAddr = t.remoteAddr
	req.RequestURI = req.URL.RequestURI()

	if req.Body == nil {
		req.Body = http.NoBody
	}

	rec := httptest.NewRecorder()
	t.h.ServeHTTP(rec, req)

	resp := rec.Result()
	resp.Request = req

	return resp, nil
}

func WithHandler(h http.Handler, remoteAddr string) Option {
	return func(c *Client) {
		c.http.Transport = handlerTransport{h: h, remoteAddr: remoteAddr}
		c.http.Timeout = 0
	}
}

// NewHandler returns a client that calls h in-process (no network).
func NewHandler(t testing.TB, h http.Handler, opts ...Option) *Client {
	t.Helper()

	opts = append([]Option{WithHandler(h, defaultRemoteAddr)}, opts...)

	return New(t, inProcessBaseURL, opts...)
}
