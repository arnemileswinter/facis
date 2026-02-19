package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// RefreshHandler handles token refresh using the refresh token cookie
type RefreshHandler struct {
	oidcIssuerURL string
	oidcClientID  string
}

// NewRefreshHandler creates a new refresh handler
func NewRefreshHandler() *RefreshHandler {
	return &RefreshHandler{
		oidcIssuerURL: os.Getenv("OIDC_ISSUER_URL"),
		oidcClientID:  os.Getenv("OIDC_CLIENT_ID"),
	}
}

// RefreshTokenResponse represents the response from a token refresh
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Error       string `json:"error,omitempty"`
	ErrorDesc   string `json:"error_description,omitempty"`
}

// ServeHTTP handles the refresh token request
func (h *RefreshHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the refresh token from the cookie
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing refresh token"})
		return
	}

	// Exchange the refresh token for a new access token
	tokenResp, err := h.refreshAccessToken(r.Context(), cookie.Value)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Return the new access token as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenResp)
}

// refreshAccessToken exchanges a refresh token for a new access token
func (h *RefreshHandler) refreshAccessToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error) {
	tokenEndpoint := h.oidcIssuerURL + "/protocol/openid-connect/token"

	// Prepare the token refresh request
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", h.oidcClientID)

	// Make the token request
	req, err := http.NewRequestWithContext(ctx, "POST", tokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokenResp RefreshTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	// Check for errors in the response
	if tokenResp.Error != "" {
		return nil, fmt.Errorf("%s: %s", tokenResp.Error, tokenResp.ErrorDesc)
	}

	return &tokenResp, nil
}
