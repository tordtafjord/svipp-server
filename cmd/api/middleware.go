package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"svipp-server/internal/auth"
	"svipp-server/internal/database"
	"svipp-server/internal/models"
)

type AuthMiddleware struct {
	authService *auth.Service
}

func NewAuthMiddleware(authService *auth.Service) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (a *AuthMiddleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from the cookie
		cookie, err := r.Cookie(auth.CookieName)
		if err != nil {
			log.Printf("No cookie present %v", cookie)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tokenString := cookie.Value
		session, ok := a.authService.ValidateToken(r.Context(), tokenString)
		if !ok {
			log.Printf("Token not valid %v", tokenString)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), auth.SessionContextKey, session) // Changed context key
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *AuthMiddleware) ApiKeyAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token directly from the Authorization header
		apiKey := r.Header.Get("Authorization")
		if apiKey == "" {
			log.Println("No Authorization header present")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Hash the API key
		hasher := sha256.New()
		hasher.Write([]byte(apiKey))
		hashedApiKey := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

		session, ok := a.authService.ValidateToken(r.Context(), hashedApiKey)
		if !ok {
			log.Printf("Token not valid %v", hashedApiKey)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), auth.SessionContextKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRole(allowedRoles ...models.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, ok := r.Context().Value(auth.SessionContextKey).(*database.GetSessionRow)
			if !ok {
				log.Printf("Error loading session from context")
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			for _, role := range allowedRoles {
				if session.Role == role.String() {
					next.ServeHTTP(w, r)
					return
				}
			}

			log.Printf("Role %v not present in required roles %v", session.Role, allowedRoles)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		})
	}
}

func LogRequestBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}
		// Read the body
		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)
		}

		// Restore the io.ReadCloser to its original state
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Use the body content
		log.Printf("Endpoint: %s, Request body: %s", r.URL.Path, bodyBytes)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
