package httpx

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

type EmptyResponse struct{}

type ErrResponse struct {
	Message string `json:"message,omitempty"`
}

type ErrConverter func(err error, l *slog.Logger) (statusCode int, i any)

type Responder struct {
	log          *slog.Logger
	errConverter ErrConverter
}

type ResponderOption func(r *Responder)

func NewResponder(opts ...ResponderOption) *Responder {
	r := &Responder{
		log: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}

	for _, o := range opts {
		o(r)
	}

	if r.errConverter == nil {
		r.errConverter = func(err error, l *slog.Logger) (statusCode int, i any) {
			return http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)
		}
	}

	return r
}

func (r *Responder) JSON(w http.ResponseWriter, statusCode int, i any) {
	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(statusCode)

	if i != nil {
		if err := json.NewEncoder(w).Encode(i); err != nil {
			r.log.Error("something went wrong", "error", err)
		}
	}
}

func (r *Responder) OK(w http.ResponseWriter, i any) {
	r.JSON(w, http.StatusOK, i)
}

func (r *Responder) Created(w http.ResponseWriter, i any) {
	r.JSON(w, http.StatusCreated, i)
}

func (r *Responder) EmptyOK(w http.ResponseWriter) {
	r.JSON(w, http.StatusOK, nil)
}

func (r *Responder) NoContent(w http.ResponseWriter) {
	r.JSON(w, http.StatusNoContent, nil)
}

func (r *Responder) InternalServerError(w http.ResponseWriter) {
	r.JSON(w, http.StatusInternalServerError, ErrResponse{
		Message: http.StatusText(http.StatusInternalServerError),
	})
}

func (r *Responder) Redirect(w http.ResponseWriter, code int, url string) {
	if code != http.StatusMovedPermanently &&
		code != http.StatusFound &&
		code != http.StatusTemporaryRedirect &&
		code != http.StatusPermanentRedirect {
		code = http.StatusFound
	}

	w.Header().Set(HeaderLocation, url)
	w.WriteHeader(code)
}

func (r *Responder) Error(w http.ResponseWriter, err error) {
	statusCode, body := r.errConverter(err, r.log)

	r.JSON(w, statusCode, body)
}

func WithLogger(log *slog.Logger) ResponderOption {
	return func(r *Responder) {
		if log != nil {
			r.log = log
		}
	}
}

func WithErrConverter(errConverter ErrConverter) ResponderOption {
	return func(r *Responder) {
		r.errConverter = errConverter
	}
}
