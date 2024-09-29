package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"svipp-server/internal/config"
	"svipp-server/internal/models"
	"time"
)

type JWTService struct {
	jwtSecret *[]byte
}

// New type for context key
type contextKey string

const cookieExpiration = 24 * 3600
const UserClaimsContextKey contextKey = "userClaims"

type CustomClaims struct {
	UserID int32  `json:"userId"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTService(cfg *config.Config) *JWTService {
	return &JWTService{jwtSecret: &cfg.JWT.SecretKey}
}

func (s *JWTService) GenerateJWT(userId int32, role string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"role":   role,
		"exp":    time.Now().Add(time.Second * cookieExpiration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(*s.jwtSecret)
}

func (s *JWTService) GenerateJwtCookie(token string) http.Cookie {
	return http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Set to true if using HTTPS
		SameSite: http.SameSiteStrictMode,
		MaxAge:   cookieExpiration, // Set the expiration time in seconds (e.g., 1 hour)
	}
}

func GetUserIdFromContext(ctx context.Context) (int32, error) {
	claims, ok := ctx.Value(UserClaimsContextKey).(*CustomClaims)
	if !ok {
		return 0, errors.New("Failed to get claims from context")
	}

	return claims.UserID, nil
}

func getRoleFromContext(ctx context.Context) (string, error) {
	claims, ok := ctx.Value(UserClaimsContextKey).(*jwt.MapClaims)
	if !ok {
		return "", errors.New("Failed to get claims from context")
	}

	if role, ok := (*claims)["role"].(string); ok {
		return role, nil
	}

	return "", errors.New("Failed to get role from claims")
}

func IsUserRole(ctx context.Context) bool {
	role, err := getRoleFromContext(ctx)
	if err != nil {
		log.Println(err)
		return false
	}
	return role == models.RoleUser.String()
}

func (s *JWTService) IsAuthenticated(r http.Request) bool {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return false
	}
	claims := &CustomClaims{} // Changed to CustomClaims
	tokenString := cookie.Value

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Added signing method check
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return *s.jwtSecret, nil
	})
	if err != nil {
		return false
	}

	if !token.Valid {
		return false
	}
	return true
}
