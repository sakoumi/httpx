package httpx

import (
	"errors"
	"fmt"
)

const (
	AttrAddr       = "addr"
	AttrErr        = "error"
	AttrStackTrace = "stack_trace"
	AttrHTTPMethod = "http_method"
	AttrPath       = "path"
	AttrMaxSize    = "max_size"
)

const (
	pathParam  = "path"
	queryParam = "query"
)

var (
	ErrBodyTooLarge = errors.New("request body too large")
	ErrBodyInvalid  = errors.New("invalid json")
)

type ErrInvalidParam struct {
	ParamType string
	ParamName string
}

func (e *ErrInvalidParam) Error() string {
	return fmt.Sprintf("invalid %s param '%s'", e.ParamType, e.ParamName)
}

func NewErrInvalidPathParam(paramName string) error {
	return &ErrInvalidParam{
		ParamType: pathParam,
		ParamName: paramName,
	}
}

func NewErrInvalidQueryParam(paramName string) error {
	return &ErrInvalidParam{
		ParamType: queryParam,
		ParamName: paramName,
	}
}
