package middleware

import (
	"net/http"
	"strings"
)

type authInterceptor struct {
	validToken string
}

// NewAuthInterceptor creates a new instance of authInterceptor with the provided validToken.
// The interceptor can be used to validate authentication tokens in incoming requests.
//
// validToken: the token string that will be considered valid for authentication.
// Returns: a pointer to an authInterceptor initialized with the given token.
func NewAuthInterceptor(validToken string) *authInterceptor {
	return &authInterceptor{validToken: validToken}
}

// abortWithStatus sends an HTTP error response with the specified status code and message.
// It writes the error message to the response writer and sets the HTTP status code accordingly.
//
// Parameters:
//   - w: The http.ResponseWriter to write the response to.
//   - code: The HTTP status code to set in the response.
//   - message: The error message to include in the response body.
func (a *authInterceptor) abortWithStatus(w http.ResponseWriter, code int, message string) {
	http.Error(w, message, code)
}

// MiddlewareFunc returns an HTTP middleware that checks for a valid Bearer token in the Authorization header.
// If the header is missing or the token is invalid, it responds with HTTP 401 Unauthorized and an error message.
// Otherwise, it calls the next handler in the chain.
func (a *authInterceptor) MiddlewareFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				a.abortWithStatus(w, http.StatusUnauthorized, "Missing authorization header")
				return
			}

			token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
			if token != a.validToken {
				a.abortWithStatus(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			next.ServeHTTP(w, r)
		},
	)
}
