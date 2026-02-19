package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// CallbackHandler handles the OIDC callback and exchanges authorization code for tokens
type CallbackHandler struct {
	template      *template.Template
	oidcIssuerURL string
	oidcClientID  string
}

// NewCallbackHandler creates a new callback handler
func NewCallbackHandler(templatesDir string) (*CallbackHandler, error) {
	tmplPath := filepath.Join(templatesDir, "auth", "callback.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return nil, err
	}

	return &CallbackHandler{
		template:      tmpl,
		oidcIssuerURL: os.Getenv("OIDC_ISSUER_URL"),
		oidcClientID:  os.Getenv("OIDC_CLIENT_ID"),
	}, nil
}

// TokenResponse represents the token endpoint response from the OIDC provider
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Error        string `json:"error,omitempty"`
	ErrorDesc    string `json:"error_description,omitempty"`
}

// ServeHTTP handles the callback request
func (h *CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the authorization code from query parameters
	code := r.URL.Query().Get("code")
	if code == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "<h1>Error</h1><p>missing authorization code</p>")
		return
	}

	// Exchange the code for tokens
	tokenResp, err := h.exchangeCodeForToken(r.Context(), code)
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<h1>Error</h1><p>%s</p>", html.EscapeString(err.Error()))
		return
	}

	// Set refresh token as secure HttpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenResp.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSite(http.SameSiteLaxMode),
		Path:     "/auth/refresh",
		MaxAge:   7 * 24 * 60 * 60, // 7 days
	})

	// Render the template with token data (excluding refresh token)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.template.Execute(w, tokenResp); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// exchangeCodeForToken exchanges the authorization code for tokens
func (h *CallbackHandler) exchangeCodeForToken(ctx context.Context, code string) (*TokenResponse, error) {
	tokenEndpoint := h.oidcIssuerURL + "/protocol/openid-connect/token"

	// Prepare the token request
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("client_id", h.oidcClientID)
	// The redirect_uri must match exactly what was sent in the authorization request
	redirectURI := os.Getenv("OIDC_REDIRECT_URI") + "/auth/callback"
	data.Set("redirect_uri", redirectURI)

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

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	// Check for errors in the response
	if tokenResp.Error != "" {
		return nil, fmt.Errorf("%s: %s", tokenResp.Error, tokenResp.ErrorDesc)
	}

	return &tokenResp, nil
}
