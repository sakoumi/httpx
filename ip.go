package httpx

import (
	"net"
	"net/http"
	"strings"
)

func ClientIP(r *http.Request) string {
	if xff := r.Header.Get(HeaderXForwardedFor); xff != "" {
		ips := strings.Split(xff, ",")
		if ip := strings.TrimSpace(ips[0]); ip != "" {
			return ip
		}
	}

	if xri := r.Header.Get(HeaderXRealIP); xri != "" {
		return xri
	}

	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return ip
	}

	return r.RemoteAddr
}
