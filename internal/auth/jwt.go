package auth

import (
	"context"
	"errors"
	"svipp-server/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTService struct {
	jwtSecret *[]byte
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
	claims, ok := ctx.Value("userClaims").(*jwt.MapClaims)
	if !ok {
		return 0, errors.New("Failed to get claims from context")
	}

	if userId, ok := (*claims)["userId"].(float64); ok {
		return int32(userId), nil
	}

	return 0, errors.New("Failed to get user id from claims")

}
