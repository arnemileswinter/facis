package auth

import (
	"net/http"

	"digital-contracting-service/internal/middleware"
)

// Goa HTTP middleware that validates OIDC tokens
// extract the Bearer token from the Authorization header,
// validate it with the OIDC provider, and inject auth context
func OIDCMiddleware(validator *middleware.OIDCValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			// Extract Bearer token
			token, err := middleware.ExtractBearerToken(authHeader)
			if err != nil {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			// Validate token with OIDC provider
			roles, err := validator.ValidateToken(r.Context(), token)
			if err != nil {
				http.Error(w, "token validation failed: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// Inject auth context into request
			ctx := middleware.InjectAuthContext(r.Context(), roles)

			// Continue with the next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AuthorizationMiddleware creates a middleware that checks if a user has required roles
// Use this to wrap individual endpoint handlers that require specific roles
func AuthorizationMiddleware(requiredRoles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Check if user has one of the required roles
			hasRequiredRole := false
			for _, role := range requiredRoles {
				if middleware.HasRole(ctx, role) {
					hasRequiredRole = true
					break
				}
			}

			if !hasRequiredRole {
				http.Error(w, "insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
