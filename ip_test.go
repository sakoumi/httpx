package httpx

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientIP(t *testing.T) {
	testCases := []struct {
		name       string
		remoteAddr string
		xForwarded string
		xRealIP    string
		want       string
	}{
		{
			name:       "X-Forwarded-For with single IP",
			remoteAddr: "127.0.0.1:8080",
			xForwarded: "192.168.1.100",
			want:       "192.168.1.100",
		},
		{
			name:       "X-Forwarded-For with multiple IPs",
			remoteAddr: "127.0.0.1:8080",
			xForwarded: "192.168.1.100, 10.0.0.1, 172.16.0.1",
			want:       "192.168.1.100",
		},
		{
			name:       "X-Forwarded-For with spaces",
			remoteAddr: "127.0.0.1:8080",
			xForwarded: "  192.168.1.100  ,  10.0.0.1  ",
			want:       "192.168.1.100",
		},
		{
			name:       "X-Real-IP when X-Forwarded-For empty",
			remoteAddr: "127.0.0.1:8080",
			xRealIP:    "203.0.113.42",
			want:       "203.0.113.42",
		},
		{
			name:       "RemoteAddr with port",
			remoteAddr: "192.168.0.1:54321",
			want:       "192.168.0.1",
		},
		{
			name:       "RemoteAddr IPv6 with port",
			remoteAddr: "[::1]:8080",
			want:       "::1",
		},
		{
			name:       "RemoteAddr without port",
			remoteAddr: "192.168.0.1",
			want:       "192.168.0.1",
		},
		{
			name:       "X-Forwarded-For has priority",
			remoteAddr: "127.0.0.1:8080",
			xForwarded: "203.0.113.1",
			xRealIP:    "203.0.113.2",
			want:       "203.0.113.1",
		},
		{
			name:       "X-Real-IP has priority over RemoteAddr",
			remoteAddr: "127.0.0.1:8080",
			xRealIP:    "203.0.113.42",
			want:       "203.0.113.42",
		},
		{
			name:       "X-Forwarded-For empty string",
			remoteAddr: "127.0.0.1:8080",
			xForwarded: "",
			xRealIP:    "203.0.113.42",
			want:       "203.0.113.42",
		},
		{
			name:       "X-Forwarded-For only spaces",
			remoteAddr: "127.0.0.1:8080",
			xForwarded: "   ,   ,   ",
			xRealIP:    "203.0.113.42",
			want:       "203.0.113.42",
		},
		{
			name:       "X-Forwarded-For comma only",
			remoteAddr: "127.0.0.1:8080",
			xForwarded: ",",
			xRealIP:    "203.0.113.42",
			want:       "203.0.113.42",
		},
		{
			name:       "All IPs IPv4",
			remoteAddr: "192.168.0.1:8080",
			xForwarded: "203.0.113.1",
			xRealIP:    "198.51.100.1",
			want:       "203.0.113.1",
		},
		{
			name:       "All IPs IPv6",
			remoteAddr: "[2001:db8::1]:8080",
			xForwarded: "2001:db8::2",
			want:       "2001:db8::2",
		},
		{
			name:       "Mixed IPv4 and IPv6 in X-Forwarded-For",
			remoteAddr: "127.0.0.1:8080",
			xForwarded: "203.0.113.1, 2001:db8::1",
			want:       "203.0.113.1",
		},
		{
			name:       "only RemoteAddr with port",
			remoteAddr: "192.168.1.1:3000",
			want:       "192.168.1.1",
		},
		{
			name:       "only RemoteAddr without port",
			remoteAddr: "10.0.0.1",
			want:       "10.0.0.1",
		},
		{
			name:       "empty RemoteAddr",
			remoteAddr: "",
			want:       "",
		},
		{
			name:       "Request from direct client",
			remoteAddr: "192.168.1.100:54321",
			want:       "192.168.1.100",
		},
		{
			name:       "Request through Nginx proxy",
			remoteAddr: "127.0.0.1:8080",
			xRealIP:    "203.0.113.50",
			want:       "203.0.113.50",
		},
		{
			name:       "Request through CDN (Cloudflare)",
			remoteAddr: "104.16.1.1:12345",
			xForwarded: "203.0.113.75, 104.16.0.1",
			want:       "203.0.113.75",
		},
		{
			name:       "Request through multiple proxies",
			remoteAddr: "127.0.0.1:8080",
			xForwarded: "203.0.113.99, 10.0.1.1, 10.0.2.1",
			want:       "203.0.113.99",
		},
		{
			name:       "Request from localhost",
			remoteAddr: "127.0.0.1:8080",
			want:       "127.0.0.1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.RemoteAddr = tc.remoteAddr

			if tc.xForwarded != "" {
				req.Header.Set(HeaderXForwardedFor, tc.xForwarded)
			}

			if tc.xRealIP != "" {
				req.Header.Set(HeaderXRealIP, tc.xRealIP)
			}

			got := ClientIP(req)

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestClientIP_Priority(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "127.0.0.1:8080"
	req.Header.Set(HeaderXForwardedFor, "203.0.113.1")
	req.Header.Set(HeaderXRealIP, "198.51.100.1")

	got := ClientIP(req)

	assert.Equal(t, "203.0.113.1", got)

	req.Header.Del(HeaderXForwardedFor)
	got = ClientIP(req)

	assert.Equal(t, "198.51.100.1", got)

	req.Header.Del(HeaderXRealIP)
	got = ClientIP(req)

	assert.Equal(t, "127.0.0.1", got)
}

func TestClientIP_IPv6(t *testing.T) {
	testCases := []struct {
		name       string
		remoteAddr string
		want       string
	}{
		{
			name:       "IPv6 loopback with port",
			remoteAddr: "[::1]:8080",
			want:       "::1",
		},
		{
			name:       "IPv6 full address with port",
			remoteAddr: "[2001:db8:85a3::8a2e:370:7334]:8080",
			want:       "2001:db8:85a3::8a2e:370:7334",
		},
		{
			name:       "IPv6 compressed with port",
			remoteAddr: "[::ffff:192.0.2.1]:8080",
			want:       "::ffff:192.0.2.1",
		},
		{
			name:       "IPv6 without port",
			remoteAddr: "2001:db8::1",
			want:       "2001:db8::1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.RemoteAddr = tc.remoteAddr

			got := ClientIP(req)

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestClientIP_Malformed(t *testing.T) {
	tesCases := []struct {
		name       string
		remoteAddr string
		want       string
	}{
		{
			name:       "RemoteAddr with multiple colons (malformed)",
			remoteAddr: "192.168.1.1:8080:9090",
			want:       "192.168.1.1:8080:9090",
		},
		{
			name:       "Very long RemoteAddr",
			remoteAddr: "192.168.1.1:99999",
			want:       "192.168.1.1",
		},
	}

	for _, tc := range tesCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.RemoteAddr = tc.remoteAddr

			got := ClientIP(req)

			assert.Equal(t, tc.want, got)
		})
	}
}
