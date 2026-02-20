package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	genauth "digital-contracting-service/gen/auth"

	"goa.design/clue/log"
)

// authSvc implements the generated auth.Service interface.
type authSvc struct {
	oidcIssuerURL string
	oidcClientID  string
	redirectURI   string
}

// NewAuth returns the Auth service implementation.
func NewAuth() genauth.Service {
	return &authSvc{
		oidcIssuerURL: os.Getenv("OIDC_ISSUER_URL"),
		oidcClientID:  os.Getenv("OIDC_CLIENT_ID"),
		redirectURI:   os.Getenv("OIDC_REDIRECT_URI"),
	}
}

// Login returns the Keycloak OIDC authorization URL.
func (s *authSvc) Login(ctx context.Context) (*genauth.LoginResult, error) {
	log.Printf(ctx, "auth.login")
	params := url.Values{}
	params.Set("client_id", s.oidcClientID)
	params.Set("redirect_uri", s.redirectURI+"/auth/callback")
	params.Set("response_type", "code")
	params.Set("scope", "openid")
	authURL := s.oidcIssuerURL + "/protocol/openid-connect/auth?" + params.Encode()
	return &genauth.LoginResult{AuthURL: authURL}, nil
}

// keycloakTokenResponse is the raw response from Keycloak's token endpoint.
type keycloakTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Error        string `json:"error,omitempty"`
	ErrorDesc    string `json:"error_description,omitempty"`
}

// Callback exchanges the authorization code for tokens.
// The refresh_token cookie is set by a custom response encoder wrapper in http.go.
func (s *authSvc) Callback(ctx context.Context, p *genauth.CallbackPayload) (*genauth.CallbackResult, error) {
	log.Printf(ctx, "auth.callback")

	tokenResp, err := s.exchangeCodeForToken(ctx, p.Code)
	if err != nil {
		return nil, fmt.Errorf("token exchange failed: %w", err)
	}

	// Stash the refresh token in context so the response encoder can set the cookie.
	// This is picked up by the custom encoder wrapper in http.go.
	SetRefreshTokenInContext(ctx, tokenResp.RefreshToken)

	return &genauth.CallbackResult{
		AccessToken: tokenResp.AccessToken,
		TokenType:   tokenResp.TokenType,
		ExpiresIn:   tokenResp.ExpiresIn,
	}, nil
}

// Refresh exchanges the refresh_token (from HttpOnly cookie) for a new access token.
func (s *authSvc) Refresh(ctx context.Context) (*genauth.RefreshResult, error) {
	log.Printf(ctx, "auth.refresh")

	// Extract *http.Request from context (injected by RequestContextMiddleware).
	r, ok := HTTPRequestFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing HTTP request in context")
	}

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		return nil, fmt.Errorf("missing refresh token cookie")
	}

	tokenResp, err := s.refreshAccessToken(ctx, cookie.Value)
	if err != nil {
		return nil, fmt.Errorf("token refresh failed: %w", err)
	}

	return &genauth.RefreshResult{
		AccessToken: tokenResp.AccessToken,
		TokenType:   tokenResp.TokenType,
		ExpiresIn:   tokenResp.ExpiresIn,
	}, nil
}

// exchangeCodeForToken POSTs the auth code to Keycloak's token endpoint.
func (s *authSvc) exchangeCodeForToken(ctx context.Context, code string) (*keycloakTokenResponse, error) {
	tokenEndpoint := s.oidcIssuerURL + "/protocol/openid-connect/token"
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("client_id", s.oidcClientID)
	data.Set("redirect_uri", s.redirectURI+"/auth/callback")

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

	var tokenResp keycloakTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}
	if tokenResp.Error != "" {
		return nil, fmt.Errorf("%s: %s", tokenResp.Error, tokenResp.ErrorDesc)
	}
	return &tokenResp, nil
}

// refreshAccessToken asks Keycloak for a new access token using the refresh token.
func (s *authSvc) refreshAccessToken(ctx context.Context, refreshToken string) (*keycloakTokenResponse, error) {
	tokenEndpoint := s.oidcIssuerURL + "/protocol/openid-connect/token"
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", s.oidcClientID)

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

	var tokenResp keycloakTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}
	if tokenResp.Error != "" {
		return nil, fmt.Errorf("%s: %s", tokenResp.Error, tokenResp.ErrorDesc)
	}
	return &tokenResp, nil
}
