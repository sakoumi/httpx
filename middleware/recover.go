package middleware

import (
	"log/slog"
	"net/http"
	"runtime"

	"github.com/sakoumi/httpx"
)

func Recover(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					buf := make([]byte, 4096)
					n := runtime.Stack(buf, false)
					stackTrace := string(buf[:n])

					log.Error(
						"request panic",
						httpx.AttrErr, err,
						httpx.AttrHTTPMethod, r.Method,
						httpx.AttrPath, r.URL.Path,
						httpx.AttrStackTrace, stackTrace,
					)

					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}

			}()

			next.ServeHTTP(w, r)
		})
	}
}
