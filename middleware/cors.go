package middleware

import (
	"net/http"
	"strings"

	"github.com/sakoumi/httpx"
)

type CORSConfig struct {
	AllowOrigins     []string `env:"CORS_ALLOW_ORIGINS"`
	AllowMethods     []string `env:"CORS_ALLOW_METHODS" envDefault:"GET,POST,PUT,PATCH,DELETE,OPTIONS,HEAD"`
	AllowHeaders     []string `env:"CORS_ALLOW_HEADERS" envDefault:"Content-Type,Authorization"`
	ExposeHeaders    []string `env:"CORS_EXPOSE_HEADERS"`
	AllowCredentials bool     `env:"CORS_ALLOW_CREDENTIALS" envDefault:"false"`
	MaxAge           int      `env:"CORS_MAX_AGE" envDefault:"0"`
}

func CORS(cfg CORSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get(httpx.HeaderOrigin)

			if isOriginAllowed(origin, cfg.AllowOrigins) {
				w.Header().Set(httpx.HeaderCORSAllowOrigin, origin)

				if cfg.AllowCredentials {
					w.Header().Set(httpx.HeaderCORSAllowCredentials, "true")
				}
			}

			if r.Method == http.MethodOptions {
				w.Header().Set(httpx.HeaderCORSAllowMethods, strings.Join(cfg.AllowMethods, ", "))
				w.Header().Set(httpx.HeaderCORSAllowHeaders, strings.Join(cfg.AllowHeaders, ", "))

				if cfg.MaxAge > 0 {
					w.Header().Set(httpx.HeaderCORSMaxAge, string(rune(cfg.MaxAge)))
				}

				if len(cfg.ExposeHeaders) > 0 {
					w.Header().Set(httpx.HeaderCORSExposeHeaders, strings.Join(cfg.ExposeHeaders, ", "))
				}

				w.WriteHeader(http.StatusOK)
				return
			}

			if len(cfg.ExposeHeaders) > 0 {
				w.Header().Set(httpx.HeaderCORSExposeHeaders, strings.Join(cfg.ExposeHeaders, ", "))
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if origin == "" {
		return false
	}

	for _, allowed := range allowedOrigins {
		if allowed == "*" {
			return true
		}

		if allowed == origin {
			return true
		}

		if after, ok := strings.CutPrefix(allowed, "*."); ok {
			domain := after
			if strings.HasSuffix(origin, domain) {
				return true
			}
		}
	}

	return false
}
