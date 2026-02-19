package auth

import (
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// LoginHandler serves the login page with OIDC authentication
type LoginHandler struct {
	template      *template.Template
	oidcIssuerURL string
	oidcClientID  string
	redirectURI   string
}

// NewLoginHandler creates a new login handler
func NewLoginHandler(templatesDir string) (*LoginHandler, error) {
	tmplPath := filepath.Join(templatesDir, "auth", "login.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return nil, err
	}

	oidcIssuerURL := os.Getenv("OIDC_ISSUER_URL")
	oidcClientID := os.Getenv("OIDC_CLIENT_ID")
	redirectURI := os.Getenv("OIDC_REDIRECT_URI")

	return &LoginHandler{
		template:      tmpl,
		oidcIssuerURL: oidcIssuerURL,
		oidcClientID:  oidcClientID,
		redirectURI:   redirectURI,
	}, nil
}

// ServeHTTP handles the login page request
func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Build the OIDC authorization URL
	authURL := h.buildAuthURL()

	// Render the template
	data := struct {
		AuthURL string
	}{
		AuthURL: authURL,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.template.Execute(w, data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// buildAuthURL constructs the OIDC authorization URL
func (h *LoginHandler) buildAuthURL() string {
	params := url.Values{}
	params.Set("client_id", h.oidcClientID)
	params.Set("redirect_uri", h.redirectURI+"/auth/callback")
	params.Set("response_type", "code")
	params.Set("scope", "openid")

	return h.oidcIssuerURL + "/protocol/openid-connect/auth?" + params.Encode()
}
