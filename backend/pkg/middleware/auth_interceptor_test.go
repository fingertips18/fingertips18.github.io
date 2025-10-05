package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthInterceptor_MiddlewareFunc(t *testing.T) {
	const validToken = "secret123"
	interceptor := NewAuthInterceptor(validToken)

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	tests := []struct {
		name       string
		authHeader string
		wantCode   int
		wantBody   string
		wantCalled bool
	}{
		{
			name:       "missing header",
			authHeader: "",
			wantCode:   http.StatusUnauthorized,
			wantBody:   "Missing authorization header\n",
			wantCalled: false,
		},
		{
			name:       "invalid token",
			authHeader: "Bearer wrongtoken",
			wantCode:   http.StatusUnauthorized,
			wantBody:   "Invalid token\n",
			wantCalled: false,
		},
		{
			name:       "valid token",
			authHeader: "Bearer " + validToken,
			wantCode:   http.StatusOK,
			wantBody:   "OK",
			wantCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextCalled = false
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()

			interceptor.MiddlewareFunc(next).ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wantCode, res.StatusCode)
			body := rec.Body.String()
			assert.Equal(t, tt.wantBody, body)
			assert.Equal(t, tt.wantCalled, nextCalled)
		})
	}
}
