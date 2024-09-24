package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"svipp-server/internal/auth"
	"svipp-server/internal/database"
	"svipp-server/internal/httputil"
	"svipp-server/internal/models"
	"svipp-server/internal/password"
)

type createUserRequest struct {
	Name            *string `json:"name" validate:"omitempty,min=2,max=100"`
	Phone           string  `json:"phone" validate:"required,e164"`
	Email           string  `json:"email" validate:"required,email"`
	Password        string  `json:"password" validate:"required,min=8,max=50"`
	ConfirmPassword string  `json:"confirmPassword" validate:"required,min=8,max=50,eqfield=Password"`
	DeviceToken     *string `json:"deviceToken" validate:"omitempty"`
}

func (h *Handler) CreateUser(writer http.ResponseWriter, request *http.Request) {
	var params createUserRequest
	isHtmx := request.Header.Get("HX-Request") == "true"

	// Parse the signup request body
	if !isHtmx {
		if err := json.NewDecoder(request.Body).Decode(&params); err != nil {
			httputil.BadRequestResponse(writer, err, false)
			return
		}
	} else {
		if err := request.ParseForm(); err != nil {
			httputil.BadRequestResponse(writer, err, true)
			return
		}
		name := request.FormValue("name")
		phone := fmt.Sprintf("+%s%s", request.FormValue("countryCode"), request.FormValue("phone"))
		params.Name = &name
		params.Phone = phone
		params.Email = request.FormValue("email")
		params.Password = request.FormValue("password")
		params.ConfirmPassword = request.FormValue("confirmPassword")
	}

	if validationErrors := validateStruct(params); validationErrors != nil {
		httputil.ValidationFailedResponse(writer, validationErrors, isHtmx)
		return
	}

	pswdHash, err := password.Hash(params.Password)
	if err != nil {
		httputil.InternalServerErrorResponse(writer, "Failed hashing password", err, isHtmx)
		return
	}

	user, err := h.db.CreateUser(request.Context(), database.CreateUserParams{
		Name:        params.Name,
		Phone:       params.Phone,
		Email:       &params.Email,
		Password:    &pswdHash,
		DeviceToken: params.DeviceToken,
		Temporary:   new(bool), // Use new(bool) to create a pointer to false
		Role:        models.RoleUser.String(),
	})
	if err != nil {
		httputil.ErrorResponse(writer, http.StatusConflict, fmt.Sprintf("Failed to create user: %v", err), "Mislyktes i å lage en ny bruke, har du en konto fra før av?", isHtmx)
		return
	}

	token, err := h.jwtService.GenerateJWT(user.ID, user.Role)
	if err != nil {
		httputil.InternalServerErrorResponse(writer, "Error generating token", err, isHtmx)
		return
	}

	if !isHtmx {
		httputil.JSONResponse(writer, 200, map[string]string{"token": token})
		return
	}

	cookie := h.jwtService.GenerateJwtCookie(token)
	http.SetCookie(writer, &cookie)
	writer.Header().Set("HX-Redirect", "/home")
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
