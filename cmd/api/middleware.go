package main

import (
	"context"
	"crypto/subtle"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"svipp-server/internal/auth"
	"svipp-server/internal/httputil"
)

type JWTAuthMiddleware struct {
	secretKey []byte
}

func NewJWTAuthMiddleware(secretKey []byte) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{secretKey: secretKey}
}

// New custom claims struct

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
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if !token.Valid {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), auth.UserClaimsContextKey, claims) // Changed context key
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(auth.UserClaimsContextKey).(*auth.CustomClaims)
			if !ok {
				httputil.UnauthorizedResponse(w)
				return
			}

			for _, role := range allowedRoles {
				if claims.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			httputil.ForbiddenResponse(w, (r.Header.Get("HX-Request") == "true"))
		})
	}
}
