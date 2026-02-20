package design

import (
	. "goa.design/goa/v3/dsl"
)

// Auth Service — OIDC/Keycloak authentication endpoints.
// All methods use NoSecurity() since they handle the authentication flow itself.
var _ = Service("Auth", func() {
	Description("Authentication endpoints for OIDC/Keycloak login, callback, and token refresh.")

	Method("login", func() {
		Description("Returns the Keycloak OIDC authorization URL for initiating the login flow.")
		NoSecurity()
		Result(func() {
			Attribute("auth_url", String, "Keycloak OIDC authorization URL")
			Required("auth_url")
		})
		HTTP(func() {
			GET("/auth/login")
			Response(StatusOK)
		})
	})

	Method("callback", func() {
		Description("Handles the OIDC callback, exchanges authorization code for tokens.")
		NoSecurity()
		Payload(func() {
			Attribute("code", String, "Authorization code from OIDC provider")
			Required("code")
		})
		Result(func() {
			Attribute("access_token", String, "JWT access token")
			Attribute("token_type", String, "Token type (Bearer)")
			Attribute("expires_in", Int, "Token expiry in seconds")
			Required("access_token", "token_type", "expires_in")
		})
		HTTP(func() {
			GET("/auth/callback")
			Param("code")
			Response(StatusOK)
		})
	})

	Method("refresh", func() {
		Description("Exchanges a refresh token (from HttpOnly cookie) for a new access token.")
		NoSecurity()
		Result(func() {
			Attribute("access_token", String, "JWT access token")
			Attribute("token_type", String, "Token type (Bearer)")
			Attribute("expires_in", Int, "Token expiry in seconds")
			Required("access_token", "token_type", "expires_in")
		})
		HTTP(func() {
			POST("/auth/refresh")
			Response(StatusOK)
		})
	})
})
