package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"fuzzy/config"
)

// SecurityMiddleware adds security headers to responses
func SecurityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate nonce for CSP
		nonce := generateNonce()
		
		// Set security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Set Content Security Policy
		csp := fmt.Sprintf(
			"default-src 'self'; "+
				"style-src 'self' 'unsafe-inline'; "+
				"script-src 'self' 'nonce-%s'; "+
				"img-src 'self' data:; "+
				"font-src 'self'; "+
				"object-src 'none'; "+
				"base-uri 'self'; "+
				"form-action 'self'",
			nonce,
		)
		w.Header().Set("Content-Security-Policy", csp)
		
		// Set HSTS header if HTTPS is enabled
		if config.AppConfig.Security.HTTPSEnabled {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		
		// Store nonce in request context for templates if needed
		r = r.WithContext(r.Context())
		
		next.ServeHTTP(w, r)
	})
}

// generateNonce creates a cryptographically secure nonce for CSP
func generateNonce() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to a simple nonce if crypto/rand fails
		return "fallback-nonce"
	}
	return base64.StdEncoding.EncodeToString(bytes)
}

// AddSecurityHeaders wraps a handler function with security middleware
func AddSecurityHeaders(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		SecurityMiddleware(http.HandlerFunc(handlerFunc)).ServeHTTP(w, r)
	}
}