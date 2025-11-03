package middleware

import (
	"net/http"
)

type CorsInterceptor struct {
	ClientURL string
	Local     bool
}

type corsInterceptor struct {
	clientURL string
	local     bool
}

func NewCorsInterceptor(c CorsInterceptor) *corsInterceptor {
	return &corsInterceptor{
		clientURL: c.ClientURL,
		local:     c.Local,
	}
}

func (c *corsInterceptor) CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var origin string
		if c.local {
			origin = "*"
		} else {
			origin = c.clientURL
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Vary", "Origin")

		if !c.local && c.clientURL == "" {
			http.Error(w, "Server misconfiguration: clientURL cannot be empty", http.StatusInternalServerError)
			return
		}

		if !c.local {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
