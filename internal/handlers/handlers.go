package handlers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"svipp-server/internal/auth"
	"svipp-server/internal/config"
	"svipp-server/internal/database"
	"svipp-server/pkg/maps"
	"svipp-server/pkg/sms"
)

type Handler struct {
	db                *database.Queries
	authService       *auth.Service
	smsService        *sms.TwilioClient
	mapsService       *maps.MapsService
	domain            string
	businessSubDomain string
}

var validate *validator.Validate

func NewHandler(srv *config.Services, domain string) *Handler {
	return &Handler{
		db:                srv.DB,
		authService:       srv.AuthService,
		smsService:        srv.SmsClient,
		mapsService:       srv.MapsClient,
		domain:            domain,
		businessSubDomain: "bedrift." + domain,
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
		switch {
		case err.Tag() == "required":
			errorMessages = append(errorMessages, fmt.Sprintf("%s er et krav", err.Field()))
		case err.Tag() == "email":
			errorMessages = append(errorMessages, fmt.Sprintf("'%s' er ikke en gyldig email", err.Value()))
		case err.Tag() == "min":
			errorMessages = append(errorMessages, fmt.Sprintf("%s må være minst %s tegn langt", err.Field(), err.Param()))
		case err.Tag() == "e164":
			errorMessages = append(errorMessages, fmt.Sprintf("'%s' er ikke et gyldig telefonnummer", err.Value()))
		case err.Tag() == "eqfield":
			errorMessages = append(errorMessages, "Passord er ikke identiske")
		case err.Field() == "OrgNumber":
			errorMessages = append(errorMessages, fmt.Sprintf("'%s' er ikke et gyldig organisasjonsnummer", err.Value()))
		default:
			errorMessages = append(errorMessages, fmt.Sprintf("%s tilfredstiller ikke %s krav", err.Field(), err.Tag()))
		}
	}
	return errorMessages
}
