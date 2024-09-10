// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Driver struct {
	ID        int32              `json:"id"`
	Status    string             `json:"status"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt pgtype.Timestamptz `json:"updatedAt"`
	DeletedAt pgtype.Timestamptz `json:"deletedAt"`
}

type Order struct {
	ID              int32              `json:"id"`
	UserID          int32              `json:"userId"`
	SenderID        int32              `json:"senderId"`
	RecipientID     int32              `json:"recipientId"`
	DriverID        *int32             `json:"driverId"`
	PickupAddress   string             `json:"pickupAddress"`
	DeliveryAddress string             `json:"deliveryAddress"`
	Status          string             `json:"status"`
	Distance        int32              `json:"distance"`
	DrivingMinutes  float64            `json:"drivingMinutes"`
	Price           float64            `json:"price"`
	CreatedAt       pgtype.Timestamptz `json:"createdAt"`
	ConfirmedAt     pgtype.Timestamptz `json:"confirmedAt"`
	AcceptedAt      pgtype.Timestamptz `json:"acceptedAt"`
	PickedUpAt      pgtype.Timestamptz `json:"pickedUpAt"`
	DeliveredAt     pgtype.Timestamptz `json:"deliveredAt"`
	UpdatedAt       pgtype.Timestamptz `json:"updatedAt"`
	CancelledAt     pgtype.Timestamptz `json:"cancelledAt"`
}

type Rating struct {
	ID        int32              `json:"id"`
	OrderID   int32              `json:"orderId"`
	RaterID   int32              `json:"raterId"`
	RateeID   int32              `json:"rateeId"`
	Rating    int32              `json:"rating"`
	Comment   *string            `json:"comment"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}

type ShopifyRequest struct {
	ID        int32              `json:"id"`
	RawJson   []byte             `json:"rawJson"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}

type TempOrder struct {
	ID              int32              `json:"id"`
	UserID          int32              `json:"userId"`
	PickupAddress   string             `json:"pickupAddress"`
	DeliveryAddress string             `json:"deliveryAddress"`
	Distance        *int32             `json:"distance"`
	DrivingMinutes  *float64           `json:"drivingMinutes"`
	Price           *float64           `json:"price"`
	CreatedAt       pgtype.Timestamptz `json:"createdAt"`
}

type User struct {
	ID          int32              `json:"id"`
	Name        *string            `json:"name"`
	Phone       string             `json:"phone"`
	Email       *string            `json:"email"`
	Password    *string            `json:"password"`
	DeviceToken *string            `json:"deviceToken"`
	Temporary   *bool              `json:"temporary"`
	RateTotal   int32              `json:"rateTotal"`
	Rates       int32              `json:"rates"`
	CreatedAt   pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt   pgtype.Timestamptz `json:"updatedAt"`
	DeletedAt   pgtype.Timestamptz `json:"deletedAt"`
	Role        string             `json:"role"`
}
