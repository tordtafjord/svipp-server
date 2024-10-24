// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Business struct {
	ID        int64              `json:"id"`
	Name      string             `json:"name"`
	OrgID     int64              `json:"orgId"`
	Address   string             `json:"address"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt pgtype.Timestamptz `json:"updatedAt"`
	DeletedAt pgtype.Timestamptz `json:"deletedAt"`
}

type BusinessHour struct {
	ApiKey    []byte      `json:"apiKey"`
	DayOfWeek int32       `json:"dayOfWeek"`
	OpensAt   pgtype.Time `json:"opensAt"`
	ClosesAt  pgtype.Time `json:"closesAt"`
}

type Driver struct {
	ID        int64              `json:"id"`
	Status    string             `json:"status"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt pgtype.Timestamptz `json:"updatedAt"`
	DeletedAt pgtype.Timestamptz `json:"deletedAt"`
}

type Order struct {
	ID                   int64              `json:"id"`
	UserID               int64              `json:"userId"`
	SenderID             int64              `json:"senderId"`
	RecipientID          int64              `json:"recipientId"`
	DriverID             *int64             `json:"driverId"`
	ShareKey             pgtype.UUID        `json:"shareKey"`
	PickupAddress        string             `json:"pickupAddress"`
	DeliveryAddress      string             `json:"deliveryAddress"`
	PickupCoords         interface{}        `json:"pickupCoords"`
	DeliveryCoords       interface{}        `json:"deliveryCoords"`
	PickupInstructions   *string            `json:"pickupInstructions"`
	DeliveryInstructions *string            `json:"deliveryInstructions"`
	Status               string             `json:"status"`
	DistanceMeters       int32              `json:"distanceMeters"`
	DrivingSeconds       int32              `json:"drivingSeconds"`
	PriceCents           int32              `json:"priceCents"`
	DeliveryWindowStart  pgtype.Timestamptz `json:"deliveryWindowStart"`
	DeliveryWindowEnd    pgtype.Timestamptz `json:"deliveryWindowEnd"`
	CreatedAt            pgtype.Timestamptz `json:"createdAt"`
	ConfirmedAt          pgtype.Timestamptz `json:"confirmedAt"`
	AcceptedAt           pgtype.Timestamptz `json:"acceptedAt"`
	PickedUpAt           pgtype.Timestamptz `json:"pickedUpAt"`
	DeliveredAt          pgtype.Timestamptz `json:"deliveredAt"`
	UpdatedAt            pgtype.Timestamptz `json:"updatedAt"`
	CancelledAt          pgtype.Timestamptz `json:"cancelledAt"`
}

type OrderQuote struct {
	UserID          int64              `json:"userId"`
	PickupAddress   string             `json:"pickupAddress"`
	DeliveryAddress string             `json:"deliveryAddress"`
	DistanceMeters  int32              `json:"distanceMeters"`
	DrivingSeconds  int32              `json:"drivingSeconds"`
	PriceOptions    []byte             `json:"priceOptions"`
	CreatedAt       pgtype.Timestamptz `json:"createdAt"`
	ExpiresAt       pgtype.Timestamptz `json:"expiresAt"`
}

type Rating struct {
	ID        int64              `json:"id"`
	OrderID   int64              `json:"orderId"`
	RaterID   int64              `json:"raterId"`
	RatedID   int64              `json:"ratedId"`
	Rating    int16              `json:"rating"`
	Comment   *string            `json:"comment"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}

type ShopifyApiConfig struct {
	ApiKey             []byte             `json:"apiKey"`
	QuoteKey           string             `json:"quoteKey"`
	BusinessID         int64              `json:"businessId"`
	PickupAddress      *string            `json:"pickupAddress"`
	PickupCoords       interface{}        `json:"pickupCoords"`
	PickupInstructions *string            `json:"pickupInstructions"`
	CreatedAt          pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt          pgtype.Timestamptz `json:"updatedAt"`
	DeletedAt          pgtype.Timestamptz `json:"deletedAt"`
	LocationName       *string            `json:"locationName"`
}

type Token struct {
	Token     string             `json:"token"`
	ExpiresAt pgtype.Timestamptz `json:"expiresAt"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
	UserID    int64              `json:"userId"`
	Role      string             `json:"role"`
}

type User struct {
	ID          int64              `json:"id"`
	FirstName   *string            `json:"firstName"`
	LastName    *string            `json:"lastName"`
	Phone       *string            `json:"phone"`
	Email       *string            `json:"email"`
	Password    *string            `json:"password"`
	DeviceToken *string            `json:"deviceToken"`
	Temporary   *bool              `json:"temporary"`
	Role        string             `json:"role"`
	RateTotal   int32              `json:"rateTotal"`
	Rates       int32              `json:"rates"`
	CreatedAt   pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt   pgtype.Timestamptz `json:"updatedAt"`
	DeletedAt   pgtype.Timestamptz `json:"deletedAt"`
}
