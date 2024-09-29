package main

import (
	"context"
	"crypto/subtle"
	"log"
	"net/http"
	"svipp-server/internal/auth"
	"svipp-server/internal/httputil"
	"svipp-server/internal/models"
)

type JWTAuthMiddleware struct {
	jwtService *auth.JWTService
}

func NewJWTAuthMiddleware(jwtService *auth.JWTService) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{jwtService: jwtService}
}

func (m *JWTAuthMiddleware) CombinedAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") == "application/json" {
			m.JwtAuthMiddleware(next).ServeHTTP(w, r)
		} else {
			m.JwtCookieAuthMiddleware(next).ServeHTTP(w, r)
		}
	})
}

func (m *JWTAuthMiddleware) JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httputil.UnauthorizedResponse(w) // Changed to app method
			return
		}

		// Changed to use constant-time comparison
		const bearerPrefix = "Bearer "
		if len(authHeader) > len(bearerPrefix) && subtle.ConstantTimeCompare([]byte(authHeader[:len(bearerPrefix)]), []byte(bearerPrefix)) == 1 {
			tokenString := authHeader[len(bearerPrefix):]

			token, ok := m.jwtService.ValidateToken(tokenString)
			if !ok {
				httputil.UnauthorizedResponse(w)
				return
			}

			ctx := context.WithValue(r.Context(), auth.UserClaimsContextKey, token.Claims)
			ctx = context.WithValue(ctx, auth.IsJsonContextKey, true)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			httputil.UnauthorizedResponse(w)
		}
	})
}

func (m *JWTAuthMiddleware) JwtCookieAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from the cookie
		cookie, err := r.Cookie("jwt")
		if err != nil {
			log.Printf("No cookie present %v", cookie)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tokenString := cookie.Value
		token, ok := m.jwtService.ValidateToken(tokenString)
		if !ok {
			log.Printf("Token not valid %v", token.Claims)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), auth.UserClaimsContextKey, token.Claims) // Changed context key
		ctx = context.WithValue(ctx, auth.IsJsonContextKey, false)                     // Added isJson context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRole(allowedRoles ...models.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(auth.UserClaimsContextKey).(*auth.CustomClaims)
			if !ok {
				log.Printf("Error loading claims %v", r)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			for _, role := range allowedRoles {
				if claims.Role == role.String() {
					next.ServeHTTP(w, r)
					return
				}
			}

			log.Printf("Role %v not present in required roles %v", claims.Role, allowedRoles)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		})
	}
}
