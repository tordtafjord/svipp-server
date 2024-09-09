package handlers

import (
	"net/http"
	"svipp-server/internal/httputil"
	"time"
)

type ShopifyRateResponse struct {
	ServiceName     string  `json:"service_name"`
	Description     string  `json:"description"`
	ServiceCode     string  `json:"service_code"`
	Currency        string  `json:"currency"`
	TotalPrice      int64   `json:"total_price"`
	PhoneRequired   bool    `json:"phone_required" default:"true"`
	MinDeliveryDate *string `json:"min_delivery_date,omitempty"`
	MaxDeliveryDate *string `json:"max_delivery_date,omitempty"`
}

func (h *Handler) ShopifyCallback(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()
	maxTime := currentTime.Add(90 * time.Minute)
	currentTimeString := currentTime.Format("2006-01-02 15:04:05 -0700")
	maxTimeString := maxTime.Format("2006-01-02 15:04:05 -0700")

	// Create a new ShopifyRateResponse
	response := &ShopifyRateResponse{
		ServiceName:   "Ekspress Levering",
		Description:   "Hentes snartest mulig av vårt leveringsbud",
		ServiceCode:   "express",
		Currency:      "NOK",
		TotalPrice:    15000,
		PhoneRequired: true,

		MinDeliveryDate: &currentTimeString,
		MaxDeliveryDate: &maxTimeString,
	}

	// TODO: Add logic to populate the response with actual data
	// Encode and send the response
	httputil.JSONResponse(w, http.StatusOK, response)
}
