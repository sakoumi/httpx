package httptestx

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Response[T any] struct {
	StatusCode int
	Header     http.Header
	Body       T
	RawBody    []byte
}

func (r *Response[T]) Assert(t testing.TB, wantStatusCode int, wantBody T) {
	t.Helper()

	require.Equal(t, wantStatusCode, r.StatusCode)
	assert.Equal(t, wantBody, r.Body)
}
