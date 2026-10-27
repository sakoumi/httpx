package httptestx

import (
	"context"
	"net/http"
	"net/url"

	"github.com/sakoumi/httpx"
)

type RequestOption func(*http.Request)

func WithHeader(k, v string) RequestOption {
	return func(r *http.Request) { r.Header.Set(k, v) }
}

func WithBearer(token string) RequestOption {
	return WithHeader(httpx.HeaderAuthorization, "Bearer "+token)
}

func WithQuery(q url.Values) RequestOption {
	return func(r *http.Request) {
		existing := r.URL.Query()
		for k, vs := range q {
			for _, v := range vs {
				existing.Add(k, v)
			}
		}
		r.URL.RawQuery = existing.Encode()
	}
}

func WithQueryParam(k, v string) RequestOption {
	return func(r *http.Request) {
		q := r.URL.Query()
		q.Set(k, v)
		r.URL.RawQuery = q.Encode()
	}
}

func WithContext(ctx context.Context) RequestOption {
	return func(r *http.Request) { *r = *r.WithContext(ctx) }
}
