package main

import (
	"context"
	"crypto/subtle"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"svipp-server/internal/auth"
	"svipp-server/internal/httputil"
	"svipp-server/internal/models"
)

type JWTAuthMiddleware struct {
	secretKey []byte
}

func NewJWTAuthMiddleware(secretKey []byte) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{secretKey: secretKey}
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
			claims := &auth.CustomClaims{} // Changed to CustomClaims

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				// Added signing method check
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return m.secretKey, nil
			})

			if err != nil {
				httputil.UnauthorizedResponse(w)
				return
			}

			if !token.Valid {
				httputil.UnauthorizedResponse(w)
				return
			}

			ctx := context.WithValue(r.Context(), auth.UserClaimsContextKey, claims) // Changed context key
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
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		claims := &auth.CustomClaims{} // Changed to CustomClaims
		tokenString := cookie.Value

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Added signing method check
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return m.secretKey, nil
		})
		if err != nil {
			log.Printf("Failed parsing jwt with claims %v", err)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if !token.Valid {
			log.Printf("Token not valid %v", token.Claims)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), auth.UserClaimsContextKey, claims) // Changed context key
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRole(allowedRoles ...models.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(auth.UserClaimsContextKey).(*auth.CustomClaims)
			if !ok {
				log.Printf("Error loading claims %v", r)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			for _, role := range allowedRoles {
				if claims.Role == role.String() {
					next.ServeHTTP(w, r)
					return
				}
			}

			log.Printf("%v not present in %v", claims.Role, allowedRoles)
			w.WriteHeader(http.StatusUnauthorized)
			return
		})
	}
}
