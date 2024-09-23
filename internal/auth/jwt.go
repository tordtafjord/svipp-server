package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"svipp-server/internal/config"
	"time"
)

type JWTService struct {
	jwtSecret *[]byte
}

// New type for context key
type contextKey string

const UserClaimsContextKey contextKey = "userClaims"

type CustomClaims struct {
	UserID string `json:"user_id"`
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
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(*s.jwtSecret)
}

func GetUserIdFromContext(ctx context.Context) (int32, error) {
	claims, ok := ctx.Value(UserClaimsContextKey).(*jwt.MapClaims)
	if !ok {
		return 0, errors.New("Failed to get claims from context")
	}

	if userId, ok := (*claims)["userId"].(float64); ok {
		return int32(userId), nil
	}

	return 0, errors.New("Failed to get user id from claims")
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
