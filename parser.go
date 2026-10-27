package httpx

import (
	"fmt"
	"net/http"
	"strconv"
	"uuid"
)

func ParsePathParamAs[T any](r *http.Request, name string) (T, error) {
	var zero T

	v, err := ParseValueAs[T](r.PathValue(name))
	if err != nil {
		return zero, fmt.Errorf("%w: %s", NewErrInvalidPathParam(name), err)
	}

	return v, nil
}

func ParsePathParamOrAs[T any](r *http.Request, name string, defaultValue T) (T, error) {
	value := r.PathValue(name)
	if value == "" {
		return defaultValue, nil
	}

	var zero T
	v, err := ParseValueAs[T](value)
	if err != nil {
		return zero, fmt.Errorf("%w: %v", NewErrInvalidPathParam(name), err)
	}

	return v, nil
}

func ParseQueryParamAs[T any](r *http.Request, name string) (T, error) {
	var zero T

	v, err := ParseValueAs[T](r.URL.Query().Get(name))
	if err != nil {
		return zero, fmt.Errorf("%w: %v", NewErrInvalidQueryParam(name), err)
	}

	return v, nil
}

func ParseQueryParamOrAs[T any](r *http.Request, name string, defaultValue T) (T, error) {
	value := r.URL.Query().Get(name)
	if value == "" {
		return defaultValue, nil
	}

	var zero T
	v, err := ParseValueAs[T](value)
	if err != nil {
		return zero, fmt.Errorf("%w: %v", NewErrInvalidQueryParam(name), err)
	}

	return v, nil
}

func ParseValueAs[T any](value string) (T, error) {
	var zero T

	switch any(zero).(type) {
	case uuid.UUID:
		res, err := uuid.Parse(value)
		if err != nil {
			return zero, err
		}
		return any(res).(T), nil

	case int:
		res, err := strconv.Atoi(value)
		if err != nil {
			return zero, err
		}
		return any(res).(T), nil

	case int64:
		res, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return zero, err
		}
		return any(res).(T), nil

	case string:
		return any(value).(T), nil

	default:
		return zero, fmt.Errorf("unsupported param type %T", zero)
	}
}
