package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"svipp-server/internal/httputil"
	"svipp-server/internal/models"
	"time"
)

type ShopifyRateRequest struct {
	Rate RateDetails `json:"rate"`
}

type RateDetails struct {
	Origin      Address `json:"origin"`
	Destination Address `json:"destination"`
	Items       []Item  `json:"items"`
	Currency    string  `json:"currency"`
	Locale      string  `json:"locale"`
}

type Address struct {
	CountryCode string  `json:"country"`
	PostalCode  string  `json:"postal_code"`
	Province    string  `json:"province"`
	City        string  `json:"city"`
	Name        *string `json:"name"`
	Address1    string  `json:"address1"`
	Address2    string  `json:"address2"`
	Address3    *string `json:"address3"`
	Phone       *string `json:"phone"`
	Fax         *string `json:"fax"`
	Email       *string `json:"email"`
	AddressType *string `json:"address_type"`
	CompanyName *string `json:"company_name"`
}

type Item struct {
	Name               string      `json:"name"`
	SKU                string      `json:"sku"`
	Quantity           int         `json:"quantity"`
	Grams              int         `json:"grams"`
	Price              int         `json:"price"`
	Vendor             string      `json:"vendor"`
	RequiresShipping   bool        `json:"requires_shipping"`
	Taxable            bool        `json:"taxable"`
	FulfillmentService string      `json:"fulfillment_service"`
	Properties         interface{} `json:"properties"`
	ProductID          int64       `json:"product_id"`
	VariantID          int64       `json:"variant_id"`
}

type ShopifyRate struct {
	ServiceName     string                `json:"service_name"`
	Description     string                `json:"description"`
	ServiceCode     models.DeliveryOption `json:"service_code"`
	Currency        string                `json:"currency"`
	TotalPrice      int64                 `json:"total_price"`
	PhoneRequired   bool                  `json:"phone_required" default:"true"`
	MinDeliveryDate *string               `json:"min_delivery_date,omitempty"`
	MaxDeliveryDate *string               `json:"max_delivery_date,omitempty"`
}

type ShopifyRates struct {
	Rates []ShopifyRate `json:"rates"`
}

func (h *Handler) ShopifyCallback(w http.ResponseWriter, r *http.Request) {
	// Read the raw request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body)) // Replace the body for further use

	// Store the raw JSON in the database
	err = h.db.CreateShopifyRequest(r.Context(), body)
	if err != nil {
		log.Printf("Failed to store Shopify request: %v", err)
		// Note: We're logging the error but not returning, to ensure we still process the request
	}

	var rateRequest ShopifyRateRequest
	err = json.NewDecoder(r.Body).Decode(&rateRequest)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	currentTime := time.Now()
	maxTime := currentTime.Add(90 * time.Minute)
	currentTimeString := currentTime.Format("2006-01-02 15:04:05 -0700")
	maxTimeString := maxTime.Format("2006-01-02 15:04:05 -0700")

	// Create a slice of ShopifyRate
	rates := []ShopifyRate{
		{
			ServiceName:     "Nå",
			Description:     "Leveres snartest mulig av vårt leveringsbud",
			ServiceCode:     models.Express,
			Currency:        "NOK",
			TotalPrice:      15000,
			PhoneRequired:   true,
			MinDeliveryDate: &currentTimeString,
			MaxDeliveryDate: &maxTimeString,
		},
		{
			ServiceName:     "I dag",
			Description:     "Leveres i løpet av dagen fra vårt leveringsbud",
			ServiceCode:     models.Today,
			Currency:        "NOK",
			TotalPrice:      12500,
			PhoneRequired:   true,
			MinDeliveryDate: &currentTimeString,
			MaxDeliveryDate: &maxTimeString,
		},
		{
			ServiceName:     "Senere",
			Description:     "Leveres innen 1-3 dager av vårt leveringsbud",
			ServiceCode:     models.Later,
			Currency:        "NOK",
			TotalPrice:      10000,
			PhoneRequired:   true,
			MinDeliveryDate: &currentTimeString,
			MaxDeliveryDate: &maxTimeString,
		},
		// Add more rates here if needed
	}
	response := &ShopifyRates{
		Rates: rates,
	}

	// TODO: Add logic to populate the response with actual data
	// Encode and send the response
	httputil.JSONResponse(w, http.StatusOK, response)
}
