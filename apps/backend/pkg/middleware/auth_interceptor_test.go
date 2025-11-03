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
		{
			name:       "Bearer with no token",
			authHeader: "Bearer ",
			wantCode:   http.StatusUnauthorized,
			wantBody:   "Invalid token\n",
			wantCalled: false,
		},
		{
			name:       "extra spaces in token",
			authHeader: "Bearer  " + validToken,
			wantCode:   http.StatusUnauthorized,
			wantBody:   "Invalid token\n",
			wantCalled: false,
		},
		{
			name:       "malformed header prefix",
			authHeader: "Token " + validToken,
			wantCode:   http.StatusUnauthorized,
			wantBody:   "Invalid token\n",
			wantCalled: false,
		},
		{
			name:       "just token, no Bearer",
			authHeader: validToken,
			wantCode:   http.StatusUnauthorized,
			wantBody:   "Invalid token\n",
			wantCalled: false,
		},
		{
			name:       "Bearer with multiple spaces, empty token",
			authHeader: "Bearer     ",
			wantCode:   http.StatusUnauthorized,
			wantBody:   "Invalid token\n",
			wantCalled: false,
		},
		{
			name:       "mixed case prefix",
			authHeader: "bearer " + validToken, // Lowercase 'bearer'
			wantCode:   http.StatusUnauthorized,
			wantBody:   "Invalid token\n",
			wantCalled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rec := httptest.NewRecorder()
			interceptor.MiddlewareFunc(next).ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wantCode, res.StatusCode)
			assert.Equal(t, tt.wantBody, rec.Body.String())
			assert.Equal(t, tt.wantCalled, nextCalled)
		})
	}
}
