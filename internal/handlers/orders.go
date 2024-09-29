package handlers

import (
	"encoding/json"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"net/http"
	"svipp-server/internal/auth"
	"svipp-server/internal/database"
	"svipp-server/internal/httputil"
	"time"
)

type orderQuoteRequest struct {
	PickupAddress   string `json:"pickupAddress" validate:"required"`
	DeliveryAddress string `json:"deliveryAddress" validate:"required"`
}

type orderQuoteResponse struct {
	PickupAddress   string         `json:"pickupAddress"`
	DeliveryAddress string         `json:"deliveryAddress"`
	DistanceMeters  int32          `json:"distanceMeters"`
	PriceOptions    map[string]int `json:"priceOptions"`
	ExpiresAt       time.Time      `json:"expiresAt"`
}

var optionPrices = map[string]int{
	"express":  15000,
	"today":    12500,
	"tomorrow": 10000,
	"later":    10000,
}

const quoteExpirationDuration = 15 * time.Minute

// Get delivery cost, prices locked and guaranteed for 15 minuted
func (h *Handler) GetOrderQuote(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserIdFromContext(r.Context())
	if err != nil {
		httputil.InternalServerErrorResponse(w, "Failed to retrieve userId of user from claims %v", err, false)
		return
	}

	var params orderQuoteRequest

	// Parse and validate the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		httputil.BadRequestResponse(w, err, false)
		return
	}
	if validationErrors := validateStruct(params); validationErrors != nil {
		httputil.ValidationFailedResponse(w, validationErrors, false)
		return
	}

	meters, seconds, err := h.mapsService.GetDistanceAndDuration(r.Context(), params.PickupAddress, params.DeliveryAddress)
	if err != nil {
		httputil.InternalServerErrorResponse(w, "Failed to get data from gmaps api %v", err, false)
		return
	}

	log.Printf("distance:%v, duration:%v, err:%v", meters, seconds, err)
	// TODO: Replace with calculation service
	prices := optionPrices

	priceOptions, err := json.Marshal(prices)
	if err != nil {
		httputil.InternalServerErrorResponse(w, "Failed to marshal price options", err, false)
		return
	}

	expires := time.Now().Add(quoteExpirationDuration)
	expiresTimestamptz := pgtype.Timestamptz{
		Time:  expires,
		Valid: true,
	}
	err = h.db.UpsertQuote(r.Context(), database.UpsertQuoteParams{
		UserID:         userId,
		PickupAddr:     params.PickupAddress,
		DeliveryAddr:   params.DeliveryAddress,
		DistanceMeters: meters,
		DrivingSeconds: seconds,
		PriceOptions:   priceOptions,
		ExpiresAt:      expiresTimestamptz,
	})
	if err != nil {
		httputil.InternalServerErrorResponse(w, "Failed to upsert quote", err, false)
	}

	quote := orderQuoteResponse{
		PickupAddress:   params.PickupAddress,
		DeliveryAddress: params.DeliveryAddress,
		DistanceMeters:  meters,
		PriceOptions:    prices,
		ExpiresAt:       expires,
	}
	httputil.JSONResponse(w, http.StatusOK, quote)
}

// handlerNewOrder handles the creation of a new order
func (h *Handler) NewOrder(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement order creation logic
	// 1. Parse request body
	// 2. Validate input
	// 3. Create order in database
	// 4. Return response

	httputil.JSONResponse(w, http.StatusCreated, map[string]string{"message": "Order created successfully"})
}

// handlerGetMyOrders retrieves orders for the authenticated user
func (h *Handler) GetMyOrders(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement logic to get user's orders
	// 1. Get user ID from context (assuming authentication middleware)
	// 2. Fetch orders from database
	// 3. Return orders in response

	// Placeholder response
	httputil.JSONResponse(w, http.StatusOK, map[string]string{"message": "User orders retrieved"})
}

// handlerConfirmOrder handles the confirmation of an order
func (h *Handler) ConfirmOrder(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement order confirmation logic
	// 1. Parse request body (likely containing order ID)
	// 2. Validate input
	// 3. Update order status in database
	// 4. Return response

	httputil.JSONResponse(w, http.StatusOK, map[string]string{"message": "Order confirmed successfully"})
}
