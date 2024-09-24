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

type ErrorData struct {
	Messages []string
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
			errorMessages = append(errorMessages, fmt.Sprintf("%s er et krav", err.Field()))
		case "email":
			errorMessages = append(errorMessages, fmt.Sprintf("'%s' er ikke en gyldig email", err.Value()))
		case "min":
			errorMessages = append(errorMessages, fmt.Sprintf("%s må være minst %s tegn langt", err.Field(), err.Param()))
		case "e164":
			errorMessages = append(errorMessages, fmt.Sprintf("'%s' er ikke et gyldig telefonnummer", err.Value()))
		case "eqfield":
			errorMessages = append(errorMessages, "Passord er ikke identiske")
		default:
			errorMessages = append(errorMessages, fmt.Sprintf("%s tilfredstiller ikke %s krav", err.Field(), err.Tag()))
		}
	}
	return errorMessages
}
