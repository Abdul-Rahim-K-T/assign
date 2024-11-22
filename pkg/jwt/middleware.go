package jwt

import (
	"context"
	"net/http"
	"strings"
)

// Middleware that checks the Authorization header for a valid JWT token
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Split the header into Bearer and token parts
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims, err := ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Store the claims in the request context
		ctx := r.Context()
		r = r.WithContext(ContextWithClaims(ctx, claims))

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

func ContextWithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, ClaimsKey, claims)
}

// Extract claims from the request context
func GetClaims(r *http.Request) *Claims {
	if claims, ok := r.Context().Value(ClaimsKey).(*Claims); ok {
		return claims
	}
	return nil
}

// Context key to store JWT claims
type contextKey string

const ClaimsKey contextKey = "claims"

// ApplicantMiddleware checks if the user is an applicant
func ApplicantMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := GetClaims(r)
		if claims == nil || claims.UserType != "applicant" {
			http.Error(w, "Access forbidden: insufficient rights", http.StatusForbidden)
			return
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := GetClaims(r)
		if claims == nil || claims.UserType != "admin" {
			http.Error(w, "Access forbidden: insufficient rights", http.StatusForbidden)
			return
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
