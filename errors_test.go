package httpx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrInvalidPathParam(t *testing.T) {
	err := NewErrInvalidPathParam("some_param")

	got := err.Error()
	assert.Equal(t, "invalid path param 'some_param'", got)
}

func TestErrInvalidQueryParam(t *testing.T) {
	err := NewErrInvalidQueryParam("some_param")

	got := err.Error()
	assert.Equal(t, "invalid query param 'some_param'", got)
}
