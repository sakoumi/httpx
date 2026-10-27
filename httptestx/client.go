package httptestx

import (
	"bytes"
	"encoding/json"
	"io"
	"maps"
	"mime"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/sakoumi/httpx"
)

type Option func(*Client)

func WithTransport(rt http.RoundTripper) Option {
	return func(c *Client) {
		c.http.Transport = rt
	}
}

func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.http.Timeout = d }
}

func WithDefaultHeaders(h map[string]string) Option {
	return func(c *Client) {
		c.headers = maps.Clone(h)
	}
}

type Client struct {
	t       testing.TB
	baseURL string
	http    *http.Client
	headers map[string]string
}

func New(t testing.TB, baseURL string, opts ...Option) *Client {
	t.Helper()

	c := &Client{
		t:       t,
		baseURL: baseURL,
		http: &http.Client{
			Timeout:       10 * time.Second,
			CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		},
	}

	for _, o := range opts {
		o(c)
	}

	return c
}

func (c *Client) WithT(t testing.TB) *Client {
	cp := *c
	cp.t = t

	return &cp
}

func (c *Client) Get[T any](path string, opts ...RequestOption) *Response[T] {
	c.t.Helper()

	return c.Do[T](http.MethodGet, path, nil, opts...)
}

func (c *Client) Delete[T any](path string, opts ...RequestOption) *Response[T] {
	c.t.Helper()

	return c.Do[T](http.MethodDelete, path, nil, opts...)
}

func (c *Client) Post[T any](path string, in any, opts ...RequestOption) *Response[T] {
	c.t.Helper()

	return c.Do[T](http.MethodPost, path, in, opts...)
}

func (c *Client) Put[T any](path string, in any, opts ...RequestOption) *Response[T] {
	c.t.Helper()

	return c.Do[T](http.MethodPut, path, in, opts...)
}

func (c *Client) Patch[T any](path string, in any, opts ...RequestOption) *Response[T] {
	c.t.Helper()

	return c.Do[T](http.MethodPatch, path, in, opts...)
}

func (c *Client) Do[T any](method, path string, in any, opts ...RequestOption) *Response[T] {
	c.t.Helper()

	body, contentType := c.encodeBody(method, path, in)

	req, err := http.NewRequestWithContext(c.t.Context(), method, c.baseURL+path, body)
	require.NoError(c.t, err, "build request %s %s", method, path)

	if contentType != "" {
		req.Header.Set(httpx.HeaderContentType, contentType)
	}

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	for _, o := range opts {
		o(req)
	}

	resp, err := c.http.Do(req)
	require.NoError(c.t, err, "%s %s", method, path)

	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	require.NoError(c.t, err, "read response %s %s", method, path)

	out := &Response[T]{StatusCode: resp.StatusCode, Header: resp.Header, RawBody: raw}

	if len(raw) == 0 {
		return out
	}

	if _, isEmpty := any(out.Body).(httpx.EmptyResponse); isEmpty {
		c.t.Fatalf("%s %s: empty body expected, status %d, body %q", method, path, resp.StatusCode, raw)
	}

	if !isJSON(resp.Header.Get(httpx.HeaderContentType)) {
		return out
	}

	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()

	require.NoError(c.t, dec.Decode(&out.Body),
		"decode %s %s: status %d, body %q", method, path, resp.StatusCode, raw)

	return out
}

func (c *Client) encodeBody(method, path string, in any) (io.Reader, string) {
	c.t.Helper()

	switch v := in.(type) {
	case nil:
		return nil, ""
	case []byte:
		return bytes.NewReader(v), ""
	case string:
		return strings.NewReader(v), ""
	case io.Reader:
		return v, ""
	default:
		b, err := json.Marshal(in)
		require.NoError(c.t, err, "encode request %s %s", method, path)

		return bytes.NewReader(b), httpx.ContentTypeJSON
	}
}

func isJSON(contentType string) bool {
	mt, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}

	return mt == "application/json" || strings.HasSuffix(mt, "+json")
}
