package handlers

import (
	"encoding/json"
	"net/http"
	"svipp-server/internal/auth"
	"svipp-server/internal/database"
	"svipp-server/internal/httputil"
)

type createUserRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=2,max=100"`
	Phone       string  `json:"phone" validate:"required,e164"`
	Email       string  `json:"email" validate:"required,email"`
	Password    string  `json:"password" validate:"required,min=8,max=50"`
	DeviceToken *string `json:"device_token" validate:"omitempty"`
}

func (h *Handler) CreateUser(writer http.ResponseWriter, request *http.Request) {
	var params createUserRequest

	// Parse the JSON request body
	if err := json.NewDecoder(request.Body).Decode(&params); err != nil {
		httputil.BadRequestResponse(writer, err)
		return
	}

	if validationErrors := validateStruct(params); validationErrors != nil {
		httputil.UnvalidResponse(writer, validationErrors)
		return
	}

	user, err := h.db.CreateUser(request.Context(), database.CreateUserParams{
		Name:        params.Name,
		Phone:       params.Phone,
		Email:       &params.Email,
		Password:    &params.Password,
		DeviceToken: params.DeviceToken,
	})
	if err != nil {
		httputil.ErrorResponse(writer, http.StatusConflict, err.Error(), "User with email or phone already exists")
		return
	}

	httputil.JSONResponse(writer, http.StatusCreated, user)
}

func (h *Handler) GetMyAccount(writer http.ResponseWriter, request *http.Request) {
	userId, err := auth.GetUserIdFromContext(request.Context())
	if err != nil {
		httputil.InternalServerErrorResponse(writer, "Failed to get user id from context", err)
		return
	}
	user, err := h.db.GetUserBasicInfoById(request.Context(), userId)
	if err != nil {
		httputil.ErrorResponse(writer, http.StatusInternalServerError, err.Error(), "Failed to get user")
		return
	}
	httputil.JSONResponse(writer, http.StatusOK, user)
}
