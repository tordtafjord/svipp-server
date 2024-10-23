package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"io"
	"log"
	"net/http"
	"svipp-server/internal/auth"
	"svipp-server/internal/database"
	"svipp-server/internal/models"
)

type AuthMiddleware struct {
	authService *auth.Service
	domain      string
}

func NewAuthMiddleware(authService *auth.Service, domain string) *AuthMiddleware {
	return &AuthMiddleware{authService: authService, domain: domain}
}

func (a *AuthMiddleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, ok := a.authService.IsAuthenticated(r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), auth.SessionContextKey, session) // Changed context key
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *AuthMiddleware) AuthOrUnauth(authHandler func(http.ResponseWriter, *http.Request),
	unauthHandler func(http.ResponseWriter, *http.Request),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, ok := a.authService.IsAuthenticated(r)
		if !ok {
			unauthHandler(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), auth.SessionContextKey, session)
		authHandler(w, r.WithContext(ctx))
	}
}

func (a *AuthMiddleware) RedirectIfAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := a.authService.IsAuthenticated(r)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		// User is authenticated, redirect to "/"
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (a *AuthMiddleware) ApiKeyAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the session id directly from the Authorization header
		apiKey := r.Header.Get("Authorization")
		if apiKey == "" {
			log.Println("No Authorization header present")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Hash the API key
		var keyHash [32]byte = sha256.Sum256([]byte(apiKey))

		session, ok := a.authService.ValidateApiKey(r.Context(), keyHash)
		if !ok {
			log.Printf("Token not valid %v", keyHash)
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
			session, ok := r.Context().Value(auth.SessionContextKey).(database.GetSessionRow)
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

func CacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set cache-control header
		w.Header().Set("Cache-Control", "public, max-age=31536000") // Cache for 1 year
		next.ServeHTTP(w, r)
	})
}
