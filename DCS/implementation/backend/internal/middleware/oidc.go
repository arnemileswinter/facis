package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

// holds OIDC provider configuration
type OIDCConfig struct {
	// Example: https://keycloak.example.com/auth/realms/dcs
	IssuerURL string
	// Example: "dcs-service". "aud" claim in JWT must match this value.
	ClientID string
}

// validate JWT tokens from OIDC providers
type OIDCValidator struct {
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
	config   OIDCConfig
}

// connects to the OIDC provider to get public keys
func NewOIDCValidator(ctx context.Context, config OIDCConfig) (*OIDCValidator, error) {
	provider, err := oidc.NewProvider(ctx, config.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to discover OIDC provider: %w", err)
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: config.ClientID})

	return &OIDCValidator{
		provider: provider,
		verifier: verifier,
		config:   config,
	}, nil
}

// Returns roles extracted from verified token
func (v *OIDCValidator) ValidateToken(ctx context.Context, token string) ([]string, error) {
	idToken, err := v.verifier.Verify(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("token verification failed: %w", err)
	}

	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse token claims: %w", err)
	}

	return extractRoles(claims), nil
}

// extractRoles extracts role information from JWT claims
// Uses Keycloak standard: realm_access.roles
func extractRoles(claims map[string]interface{}) []string {
	if ra, ok := claims["realm_access"].(map[string]interface{}); ok {
		if r, ok := ra["roles"].([]interface{}); ok {
			roles := make([]string, 0, len(r))
			for _, role := range r {
				if roleStr, ok := role.(string); ok {
					roles = append(roles, roleStr)
				}
			}
			return roles
		}
	}
	return []string{}
}

// Expected format: "Authorization: Bearer <token>"
func ExtractBearerToken(authHeader string) (string, error) {
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", fmt.Errorf("invalid authorization header format")
	}
	return strings.TrimPrefix(authHeader, bearerPrefix), nil
}

// access to roles from request context
type AuthContext struct {
	roles []string
}

// extract roles from context
func GetRoles(ctx context.Context) []string {
	if ac, ok := ctx.Value("auth").(AuthContext); ok {
		return ac.roles
	}
	return []string{}
}

// check if the context contains a specific role
func HasRole(ctx context.Context, requiredRole string) bool {
	roles := GetRoles(ctx)
	for _, role := range roles {
		if role == requiredRole {
			return true
		}
	}
	return false
}

// injects roles into the request context
func InjectAuthContext(ctx context.Context, roles []string) context.Context {
	return context.WithValue(ctx, "auth", AuthContext{roles: roles})
}
