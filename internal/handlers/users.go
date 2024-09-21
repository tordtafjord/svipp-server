package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"svipp-server/internal/auth"
	"svipp-server/internal/database"
	"svipp-server/internal/httputil"
	"svipp-server/internal/models"
)

type createUserRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=2,max=100"`
	Phone       string  `json:"phone" validate:"required,e164"`
	Email       string  `json:"email" validate:"required,email"`
	Password    string  `json:"password" validate:"required,min=8,max=50"`
	DeviceToken *string `json:"deviceToken" validate:"omitempty"`
}

func (h *Handler) CreateUser(writer http.ResponseWriter, request *http.Request) {
	var params createUserRequest

	// Parse the JSON request body
	if err := json.NewDecoder(request.Body).Decode(&params); err != nil {
		httputil.BadRequestResponse(writer, err, false)
		return
	}

	if validationErrors := validateStruct(params); validationErrors != nil {
		httputil.ValidationFailedResponse(writer, validationErrors, false)
		return
	}

	user, err := h.db.CreateUser(request.Context(), database.CreateUserParams{
		Name:        params.Name,
		Phone:       params.Phone,
		Email:       &params.Email,
		Password:    &params.Password,
		DeviceToken: params.DeviceToken,
		Temporary:   new(bool), // Use new(bool) to create a pointer to false
		Role:        models.RoleUser.String(),
	})
	if err != nil {
		httputil.ErrorResponse(writer, http.StatusConflict, fmt.Sprintf("Failed to create user: %v", err), "Mislyktes i å lage en ny bruke, har du en konto fra før av?", false)
		return
	}

	httputil.JSONResponse(writer, http.StatusCreated, user)
}

func (h *Handler) GetMyAccount(writer http.ResponseWriter, request *http.Request) {
	userId, err := auth.GetUserIdFromContext(request.Context())
	if err != nil {
		httputil.InternalServerErrorResponse(writer, "Failed to get user id from context", err, false)
		return
	}
	user, err := h.db.GetUserBasicInfoById(request.Context(), userId)
	if err != nil {
		httputil.InternalServerErrorResponse(writer, "Failed to get user", err, false)
		return
	}
	httputil.JSONResponse(writer, http.StatusOK, user)
}
