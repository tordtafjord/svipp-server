package main

import (
	"context"
	"crypto/subtle"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"svipp-server/internal/httputil"
)

type JWTAuthMiddleware struct {
	secretKey []byte
}

func NewJWTAuthMiddleware(secretKey []byte) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{secretKey: secretKey}
}

// New type for context key
type contextKey string

const UserClaimsContextKey contextKey = "userClaims"

// New custom claims struct
type CustomClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (m *JWTAuthMiddleware) JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httputil.UnauthorizedResponse(w, (r.Header.Get("HX-Request") == "true")) // Changed to app method
			return
		}

		// Changed to use constant-time comparison
		const bearerPrefix = "Bearer "
		if len(authHeader) > len(bearerPrefix) && subtle.ConstantTimeCompare([]byte(authHeader[:len(bearerPrefix)]), []byte(bearerPrefix)) == 1 {
			tokenString := authHeader[len(bearerPrefix):]
			claims := &CustomClaims{} // Changed to CustomClaims

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				// Added signing method check
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return m.secretKey, nil
			})

			if err != nil {
				httputil.UnauthorizedResponse(w, (r.Header.Get("HX-Request") == "true")) // New method for detailed error handling
				return
			}

			if !token.Valid {
				httputil.UnauthorizedResponse(w, (r.Header.Get("HX-Request") == "true"))
				return
			}

			ctx := context.WithValue(r.Context(), UserClaimsContextKey, claims) // Changed context key
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			httputil.UnauthorizedResponse(w, (r.Header.Get("HX-Request") == "true"))
		}
	})
}

func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserClaimsContextKey).(*CustomClaims)
			if !ok {
				httputil.UnauthorizedResponse(w, (r.Header.Get("HX-Request") == "true"))
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
