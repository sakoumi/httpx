package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

const defaultMaxBodySize = 1 << 20

type Binder struct {
	validator   *validator.Validate
	maxBodySize int64
}

type BinderOption func(b *Binder)

func NewBinder(opts ...BinderOption) *Binder {
	b := &Binder{
		validator:   validator.New(),
		maxBodySize: defaultMaxBodySize,
	}

	for _, o := range opts {
		o(b)
	}

	return b
}

func (b *Binder) JSON(r *http.Request, i any) error {
	body := io.LimitReader(r.Body, b.maxBodySize)
	defer r.Body.Close()

	err := json.NewDecoder(body).Decode(i)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return fmt.Errorf("%w: %v", ErrBodyTooLarge, err)
		}

		return fmt.Errorf("%w: %v", ErrBodyInvalid, err)
	}

	if err := b.validator.Struct(i); err != nil {
		return fmt.Errorf("%w: %v", ErrBodyInvalid, err)
	}

	return nil
}

func WithBodyMaxSize(maxBodySize int64) BinderOption {
	return func(b *Binder) {
		if maxBodySize > 0 {
			b.maxBodySize = maxBodySize
		}
	}
}
