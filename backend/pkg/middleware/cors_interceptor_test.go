package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorsInterceptor_CorsMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		config         CorsInterceptor
		method         string
		reqOrigin      string
		wantOrigin     string
		wantCreds      string
		wantCode       int
		wantNextCalled bool
	}{
		{
			name: "local GET request",
			config: CorsInterceptor{
				ClientURL: "http://prod-client.com",
				Local:     true,
			},
			method:         http.MethodGet,
			wantOrigin:     "*",
			wantCreds:      "",
			wantCode:       http.StatusOK,
			wantNextCalled: true,
		},
		{
			name: "non-local GET request",
			config: CorsInterceptor{
				ClientURL: "http://prod-client.com",
				Local:     false,
			},
			method:         http.MethodGet,
			wantOrigin:     "http://prod-client.com",
			wantCreds:      "true",
			wantCode:       http.StatusOK,
			wantNextCalled: true,
		},
		{
			name: "non-local OPTIONS preflight",
			config: CorsInterceptor{
				ClientURL: "http://prod-client.com",
				Local:     false,
			},
			method:         http.MethodOptions,
			wantOrigin:     "http://prod-client.com",
			wantCreds:      "true",
			wantCode:       http.StatusOK,
			wantNextCalled: false,
		},
		{
			name: "local OPTIONS request",
			config: CorsInterceptor{
				ClientURL: "http://prod-client.com",
				Local:     true,
			},
			method:         http.MethodOptions,
			wantOrigin:     "*",
			wantCreds:      "",
			wantCode:       http.StatusOK,
			wantNextCalled: false,
		},
		{
			name: "non-local POST request",
			config: CorsInterceptor{
				ClientURL: "http://prod-client.com",
				Local:     false,
			},
			method:         http.MethodPost,
			wantOrigin:     "http://prod-client.com",
			wantCreds:      "true",
			wantCode:       http.StatusOK,
			wantNextCalled: true,
		},
		{
			name: "non-local PUT request",
			config: CorsInterceptor{
				ClientURL: "http://prod-client.com",
				Local:     false,
			},
			method:         http.MethodPut,
			wantOrigin:     "http://prod-client.com",
			wantCreds:      "true",
			wantCode:       http.StatusOK,
			wantNextCalled: true,
		},
		{
			name: "non-local DELETE request",
			config: CorsInterceptor{
				ClientURL: "http://prod-client.com",
				Local:     false,
			},
			method:         http.MethodDelete,
			wantOrigin:     "http://prod-client.com",
			wantCreds:      "true",
			wantCode:       http.StatusOK,
			wantNextCalled: true,
		},
		{
			name: "non-local request with clientURL",
			config: CorsInterceptor{
				ClientURL: "http://prod-client.com",
				Local:     false,
			},
			method:         http.MethodGet,
			wantOrigin:     "http://prod-client.com",
			wantCreds:      "true",
			wantCode:       http.StatusOK,
			wantNextCalled: true,
		},
		{
			name: "non-local request with custom Origin header",
			config: CorsInterceptor{
				ClientURL: "http://prod-client.com",
				Local:     false,
			},
			method:         http.MethodGet,
			reqOrigin:      "http://custom-origin.com",
			wantOrigin:     "http://prod-client.com",
			wantCreds:      "true",
			wantCode:       http.StatusOK,
			wantNextCalled: true,
		}, {
			name: "empty ClientURL non-local request returns 500",
			config: CorsInterceptor{
				ClientURL: "",
				Local:     false,
			},
			method:         http.MethodGet,
			wantOrigin:     "",
			wantCreds:      "",
			wantCode:       http.StatusInternalServerError,
			wantNextCalled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			})

			handler := NewCorsInterceptor(tt.config).CorsMiddleware(next)

			req := httptest.NewRequest(tt.method, "/", nil)
			if tt.reqOrigin != "" {
				req.Header.Set("Origin", tt.reqOrigin)
			} else if !tt.config.Local {
				// Simulate a default non-local Origin
				req.Header.Set("Origin", "http://some-origin.com")
			} else {
				// Optional: simulate local origin
				req.Header.Set("Origin", "http://localhost")
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wantCode, res.StatusCode)
			assert.Equal(t, tt.wantOrigin, res.Header.Get("Access-Control-Allow-Origin"))
			assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", res.Header.Get("Access-Control-Allow-Methods"))
			assert.Equal(t, "Content-Type, Authorization", res.Header.Get("Access-Control-Allow-Headers"))
			assert.Equal(t, "Origin", res.Header.Get("Vary"))

			if tt.wantCreds != "" {
				assert.Equal(t, tt.wantCreds, res.Header.Get("Access-Control-Allow-Credentials"))
			} else {
				// For local requests or empty ClientURL, credentials should not be set
				_, exists := res.Header["Access-Control-Allow-Credentials"]
				assert.False(t, exists, "credentials header should not be set when wantCreds is empty")
			}

			assert.Equal(t, tt.wantNextCalled, nextCalled)
		})
	}
}
