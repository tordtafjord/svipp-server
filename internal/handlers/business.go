package handlers

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"net/http"
	"strconv"
	"svipp-server/assets/templates/components"
	"svipp-server/internal/auth"
	"svipp-server/internal/database"
	"svipp-server/internal/httputil"
	"svipp-server/internal/models"
	"svipp-server/internal/password"
	"svipp-server/pkg/util"
	"time"
)

type createBusinessAccountRequest struct {
	FirstName       *string `json:"firstName" validate:"required,min=2,max=100"`
	LastName        *string `json:"lastName" validate:"required,min=2,max=100"`
	Phone           string  `json:"phone" validate:"required,e164"`
	Email           string  `json:"email" validate:"required,email"`
	Password        string  `json:"password" validate:"required,min=8,max=64"`
	ConfirmPassword string  `json:"confirmPassword" validate:"required,min=8,max=64,eqfield=Password"`
	CompanyName     string  `json:"companyName" validate:"required,min=1,max=255"`
	OrgNumber       string  `json:"orgNumber" validate:"required,len=9,numeric"`
	BusinessAddress string  `json:"businessAddress" validate:"required,min=1,max=255"`
	City            string  `json:"city" validate:"required,min=1,max=100"`
	ZipCode         string  `json:"zipCode" validate:"required,min=4,max=10,numeric"`
	LocationName    string  `json:"locationName" validate:"required,min=1,max=100"`
}

type NewShopifyConfigForm struct {
	LocationName       string         `json:"locationName" validate:"required"`
	UseShopifyAddress  bool           `json:"useShopifyAddress"`
	Address            *string        `json:"address" validate:"omitempty"`
	ZipCode            *string        `json:"zipCode" validate:"omitempty,numeric"`
	City               *string        `json:"city" validate:"omitempty"`
	PickupInstructions *string        `json:"pickupInstructions" validate:"omitempty"`
	PickupWindows      []PickupWindow `json:"pickupWindows" validate:"dive"`
}

type PickupWindow struct {
	Start string `json:"start" validate:"omitempty,datetime=15:04"`
	End   string `json:"end" validate:"omitempty,datetime=15:04"`
}

func (h *Handler) CreateShopifyConfig(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form create shopify config %v", err)
		httputil.RedToastResponse(w, r, "Error: Bad Request")
		return
	}

	pickupWindows := []PickupWindow{
		PickupWindow{Start: r.FormValue("0Start"), End: r.FormValue("0End")},
		PickupWindow{Start: r.FormValue("1Start"), End: r.FormValue("1End")},
		PickupWindow{Start: r.FormValue("2Start"), End: r.FormValue("2End")},
		PickupWindow{Start: r.FormValue("3Start"), End: r.FormValue("3End")},
		PickupWindow{Start: r.FormValue("4Start"), End: r.FormValue("4End")},
		PickupWindow{Start: r.FormValue("5Start"), End: r.FormValue("5End")},
		PickupWindow{Start: r.FormValue("6Start"), End: r.FormValue("6End")},
	}

	params := NewShopifyConfigForm{
		LocationName:       r.FormValue("locationName"),
		Address:            util.StringToPtr(r.FormValue("address")),
		ZipCode:            util.StringToPtr(r.FormValue("zipCode")),
		City:               util.StringToPtr(r.FormValue("city")),
		PickupInstructions: util.StringToPtr(r.FormValue("pickupInstructions")),
		UseShopifyAddress:  r.FormValue("useShopifyAddress") == "true" || r.FormValue("useShopifyAddress") == "on",
		PickupWindows:      pickupWindows,
	}
	if validationErrors := validateStruct(params); validationErrors != nil {
		httputil.YellowToastResponse(w, r, validationErrors)
		return
	}

	userId, err := auth.GetUserIdFromCtx(r.Context())
	if err != nil {
		log.Printf("Failed to get user id")
		httputil.RedToastResponse(w, r, "Error: Intern Serverfeil")
		return
	}

	var addr *string
	if params.Address == nil || params.ZipCode == nil || params.City == nil {
		if !params.UseShopifyAddress {
			httputil.YellowToastResponse(w, r, []string{"Mangler gyldig addresse"})
			return
		}
	} else {
		address := fmt.Sprintf("%s, %s %s", params.Address, params.ZipCode, params.City)
		addr = &address
	}

	shopifyConfig := database.CreateShopifyApiKeyParams{
		BusinessID:         int64(userId),
		LocationName:       &params.LocationName,
		PickupAddress:      addr,
		PickupCoords:       nil,
		PickupInstructions: params.PickupInstructions,
	}

	apiKey, token, err := h.authService.CreateShopifyApiKey(r.Context(), shopifyConfig)
	if err != nil {
		log.Printf("Failed to create shopify api config")
		httputil.RedToastResponse(w, r, "Error: Intern Serverfeil")
		return
	}

	// Convert pickup windows to arrays for bulk insert
	var (
		daysOfWeek   = make([]int32, 0, 7) // 0-6 for Monday-Sunday
		openingTimes = make([]pgtype.Time, 0, 7)
		closingTimes = make([]pgtype.Time, 0, 7)
	)

	for i, window := range pickupWindows {
		opens, closes := util.TimeInputToPgTime(window.Start), util.TimeInputToPgTime(window.End)
		if opens.Valid && closes.Valid {
			daysOfWeek = append(daysOfWeek, int32(i))
			openingTimes = append(openingTimes, opens)
			closingTimes = append(closingTimes, closes)
			continue
		}

		if opens.Valid || closes.Valid {
			httputil.YellowToastResponse(w, r, []string{fmt.Sprintf("%s pickup window is not complete", time.Weekday(i))})
			return
		}
	}

	err = h.db.UpsertBusinessHours(r.Context(), database.UpsertBusinessHoursParams{
		ApiKey:       apiKey[:],
		DayOfWeek:    daysOfWeek,
		OpeningTimes: openingTimes,
		ClosingTimes: closingTimes,
	})
	if err != nil {
		log.Printf("Failed to insert business hours")
		httputil.RedToastResponse(w, r, "Error: Intern Serverfeil")
		return
	}

	err = components.ApiKeyModal(token).Render(r.Context(), w)
	if err != nil {
		log.Printf("Failed to render api key template")
		httputil.RedToastResponse(w, r, "Error: Intern Serverfeil")
		return
	}
}

func (h *Handler) CreateBusiness(w http.ResponseWriter, r *http.Request) {
	var params createBusinessAccountRequest

	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form create business %v", err)
		httputil.RedToastResponse(w, r, "Error: Bad Request")
		return
	}

	params.OrgNumber = r.FormValue("orgNumber")
	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	params.FirstName = &firstName
	params.LastName = &lastName
	params.Phone = fmt.Sprintf("+%s%s", r.FormValue("countryCode"), r.FormValue("phone"))
	params.Email = r.FormValue("email")
	params.Password = r.FormValue("password")
	params.ConfirmPassword = r.FormValue("confirmPassword")
	params.BusinessAddress = r.FormValue("businessAddress")
	params.CompanyName = r.FormValue("companyName")
	params.ZipCode = r.FormValue("zipCode")
	params.City = r.FormValue("city")
	params.LocationName = r.FormValue("locationName")

	if validationErrors := validateStruct(params); validationErrors != nil {
		httputil.YellowToastResponse(w, r, validationErrors)
		return
	}
	// orgnr valid, parse to int64:
	orgNr, err := strconv.ParseInt(params.OrgNumber, 10, 64)
	if err != nil {
		log.Printf("Feilet orgnr convert %v", err)
		httputil.RedToastResponse(w, r, "Ugyldig Orgnr")
		return
	}

	pswdHash, err := password.Hash(params.Password)
	if err != nil {
		log.Printf("Feilet pswdhash %v", err)
		httputil.RedToastResponse(w, r, "Error: Intern Serverfeil")
		return
	}

	ctx := r.Context()
	user, err := h.db.CreateUser(ctx, database.CreateUserParams{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Phone:     &params.Phone,
		Email:     &params.Email,
		Password:  &pswdHash,
		Temporary: new(bool), // Use new(bool) to create a pointer to false
		Role:      models.RoleBusiness.String(),
	})
	if err != nil {
		log.Printf("Oppretting av bedriftskonto mislyktes %v", err)
		httputil.YellowToastResponse(w, r, []string{"Oppretting av konto feilet. Har du en konto i fra før?"})
		return
	}
	// TODO: Double check this does not overwrite exising non temp users

	addr := fmt.Sprintf("%s, %s %s", params.BusinessAddress, params.ZipCode, params.City)
	err = h.db.CreateBusiness(ctx, database.CreateBusinessParams{
		ID:      user.ID,
		Name:    params.CompanyName,
		OrgID:   orgNr,
		Address: addr,
	})
	if err != nil {
		log.Printf("Oppretting av bedriftskonto mislyktes %v", err)
		httputil.YellowToastResponse(w, r, []string{"Oppretting av konto feilet. Har du en konto i fra før?"})
		return
	}

	sessionId, err := h.authService.CreateSession(ctx, user.ID, models.Role(user.Role))
	if err != nil {
		log.Printf("Error creating session %v", err)
		httputil.RedToastResponse(w, r, "Error: Intern Serverfeil")
		return
	}

	cookie := auth.CreateCookie(sessionId)
	http.SetCookie(w, &cookie)
	w.Header().Set("HX-Redirect", "/")
}
