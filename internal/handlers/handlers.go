package handlers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"html/template"
	"io/fs"
	"svipp-server/assets"
	"svipp-server/internal/auth"
	"svipp-server/internal/config"
	"svipp-server/internal/database"
)

type Handler struct {
	db         *database.Queries
	jwtService *auth.JWTService
	templates  *template.Template
}

var validate *validator.Validate

func NewHandler(cfg *config.Config) *Handler {
	templates, err := parseTemplates()
	if err != nil {
		// Handle error (e.g., log it or panic)
		panic(err)
	}

	return &Handler{
		db:         cfg.DB.DBQ,
		jwtService: auth.NewJWTService(cfg),
		templates:  templates,
	}
}

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func parseTemplates() (*template.Template, error) {
	subFS, err := fs.Sub(assets.EmbeddedFiles, "templates")
	if err != nil {
		return nil, fmt.Errorf("failed to create sub-filesystem: %w", err)
	}
	tmpl := template.Must(template.ParseFS(subFS, "*.gohtml"))

	return tmpl, nil
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
