package handlers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"svipp-server/internal/auth"
	"svipp-server/internal/config"
	"svipp-server/internal/database"
)

type Handler struct {
	db         *database.Queries
	jwtService *auth.JWTService
}

var validate *validator.Validate

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		db:         cfg.DB.DBQ,
		jwtService: auth.NewJWTService(cfg),
	}
}

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func validateStruct(s interface{}) []string {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var errorMessages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			errorMessages = append(errorMessages, fmt.Sprintf("%s is required", err.Field()))
		case "email":
			errorMessages = append(errorMessages, fmt.Sprintf("%s must be a valid email address", err.Field()))
		case "min":
			errorMessages = append(errorMessages, fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param()))
		default:
			errorMessages = append(errorMessages, fmt.Sprintf("%s failed on the %s tag", err.Field(), err.Tag()))
		}
	}
	return errorMessages
}
