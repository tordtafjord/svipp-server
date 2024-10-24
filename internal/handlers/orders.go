package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
	"svipp-server/internal/auth"
	"svipp-server/internal/database"
	"svipp-server/internal/httputil"
	"svipp-server/internal/models"
	"time"
)

// Real expiration time 15.5 min, shown to customers 15min:
const quoteExpirationDuration = 15*time.Minute + 30*time.Second

type orderQuoteRequest struct {
	PickupAddress   string `json:"pickupAddress" validate:"required"`
	DeliveryAddress string `json:"deliveryAddress" validate:"required"`
}

type orderQuoteResponse struct {
	PickupAddress   string      `json:"pickupAddress"`
	DeliveryAddress string      `json:"deliveryAddress"`
	DistanceMeters  int32       `json:"distanceMeters"`
	PriceOptions    QuotePrices `json:"priceOptions"`
	ExpiresAt       time.Time   `json:"expiresAt"`
}

type newOrderRequest struct {
	PickupAddress   string `json:"pickupAddress" validate:"required"`
	DeliveryAddress string `json:"deliveryAddress" validate:"required"`
	Phone           string `json:"phone" validate:"required,e164"`
	PriceOption     string `json:"priceOption" validate:"required"`
	IsSender        bool   `json:"isSender" validate:"required"`
}

type QuotePrices struct {
	Prices map[models.DeliveryOption]int32
}

func NewQuotePrices() QuotePrices {
	return QuotePrices{
		Prices: make(map[models.DeliveryOption]int32),
	}
}

// Get delivery cost, prices locked and guaranteed for 15 minuted
func (h *Handler) GetOrderQuote(w http.ResponseWriter, r *http.Request) {
	session, err := auth.GetSessionFromCtx(r.Context())
	if err != nil {
		httputil.InternalServerErrorResponse(w, "Failed to retrieve session of user %v", err, false)
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

	// TODO: Replace with calculation service
	prices := NewQuotePrices()
	prices.Prices[models.Express] = 15000
	prices.Prices[models.Today] = 12500
	prices.Prices[models.Tomorrow] = 10000
	prices.Prices[models.Later] = 10000

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
		UserID:          session.UserID,
		PickupAddress:   params.PickupAddress,
		DeliveryAddress: params.DeliveryAddress,
		DistanceMeters:  meters,
		DrivingSeconds:  seconds,
		PriceOptions:    priceOptions,
		ExpiresAt:       expiresTimestamptz,
	})
	if err != nil {
		httputil.InternalServerErrorResponse(w, "Failed to upsert quote", err, false)
	}

	quote := orderQuoteResponse{
		PickupAddress:   params.PickupAddress,
		DeliveryAddress: params.DeliveryAddress,
		DistanceMeters:  meters,
		PriceOptions:    prices,
		ExpiresAt:       expires.Add(-30 * time.Second),
	}
	httputil.JSONResponse(w, http.StatusOK, quote)
}

// handlerNewOrder handles the creation of a new order from a quote
func (h *Handler) NewOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session, err := auth.GetSessionFromCtx(ctx)
	if err != nil {
		httputil.InternalServerErrorResponse(w, "Failed to retrieve session of user %v", err, false)
		return
	}

	var params newOrderRequest
	// Parse and validate the JSON request body
	if err = json.NewDecoder(r.Body).Decode(&params); err != nil {
		httputil.BadRequestResponse(w, err, false)
		return
	}
	if validationErrors := validateStruct(params); validationErrors != nil {
		httputil.ValidationFailedResponse(w, validationErrors, false)
		return
	}

	orderQuote, err := h.db.GetOrderQuote(ctx, database.GetOrderQuoteParams{
		UserID:          session.UserID,
		PickupAddress:   params.PickupAddress,
		DeliveryAddress: params.DeliveryAddress,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			// Order quote not found or expired
			httputil.JSONResponse(w, http.StatusGone, map[string]string{
				"error":   "Order quote has expired or not found",
				"message": "Please request a new quote",
			})
		} else {
			httputil.InternalServerErrorResponse(w, "Failed to query db for orderQuote", err, false)
		}
		return
	}

	otherUser, err := h.db.GetOrCreateTempUser(ctx, database.GetOrCreateTempUserParams{
		Phone:     &params.Phone,
		FirstName: nil,
		LastName:  nil,
		Email:     nil,
	})
	if err != nil {
		httputil.InternalServerErrorResponse(w, "Failed to GetOrCreateTempUser", err, false)
		return
	}

	// Parse stored quotePrices
	var quotePrices QuotePrices
	if err = json.Unmarshal(orderQuote.PriceOptions, &quotePrices); err != nil {
		httputil.BadRequestResponse(w, err, false)
		return
	}
	price, exists := quotePrices.Prices[models.DeliveryOption(params.PriceOption)]
	if !exists {
		msg := fmt.Sprintf("Chosen price option %v does not exits in %v", params.PriceOption, quotePrices.Prices)
		httputil.ErrorResponse(w, http.StatusBadRequest, msg, msg, false)
		return
	}

	var recipientId, senderId int64
	var status string
	if session.Role != models.RoleBusiness.String() {
		// If user is not from a webshop or merchant order
		if params.IsSender {
			recipientId = otherUser.ID
			senderId = session.UserID
		} else {
			recipientId = session.UserID
			senderId = otherUser.ID
		}
		status = models.Pending.String()
	} else {
		// Web shop or merchant order
		recipientId = otherUser.ID
		senderId = session.UserID
		status = models.Confirmed.String()
	}

	// Set price from price options
	newOrder, err := h.db.CreateOrder(ctx, database.CreateOrderParams{
		UserID:          session.UserID,
		SenderID:        senderId,
		RecipientID:     recipientId,
		PickupAddress:   orderQuote.PickupAddress,
		DeliveryAddress: orderQuote.DeliveryAddress,
		DistanceMeters:  orderQuote.DistanceMeters,
		DrivingSeconds:  orderQuote.DrivingSeconds,
		PriceCents:      price,
		Status:          status,
	})
	if err != nil {
		httputil.InternalServerErrorResponse(w, "Failed to create new order", err, false)
		return
	}

	// TODO Notify user sms with order url
	msg := fmt.Sprintf("Pakke på vei til deg. https://svipp.app/orders/%s", newOrder.PublicID)
	h.smsService.SendSMSAsync(*otherUser.Phone, msg)

	httputil.JSONResponse(w, http.StatusCreated, newOrder)
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
