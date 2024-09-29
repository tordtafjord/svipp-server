// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: order_price_guarantees.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const upsertQuote = `-- name: UpsertQuote :exec
INSERT INTO order_price_guarantees (
    user_id,
    pickup_addr,
    delivery_addr,
    distance_meters,
    driving_seconds,
    price_options,
    expires_at
)
VALUES (
           $1,
           $2,
           $3,
           $4,
           $5,
           $6,
           $7 )
ON CONFLICT (delivery_addr, pickup_addr, user_id)
DO UPDATE SET
              distance_meters = excluded.distance_meters,
              driving_seconds = excluded.driving_seconds,
              price_options = excluded.price_options,
              expires_at = excluded.expires_at
`

type UpsertQuoteParams struct {
	UserID         int32              `json:"userId"`
	PickupAddr     string             `json:"pickupAddr"`
	DeliveryAddr   string             `json:"deliveryAddr"`
	DistanceMeters int32              `json:"distanceMeters"`
	DrivingSeconds int32              `json:"drivingSeconds"`
	PriceOptions   []byte             `json:"priceOptions"`
	ExpiresAt      pgtype.Timestamptz `json:"expiresAt"`
}

func (q *Queries) UpsertQuote(ctx context.Context, arg UpsertQuoteParams) error {
	_, err := q.db.Exec(ctx, upsertQuote,
		arg.UserID,
		arg.PickupAddr,
		arg.DeliveryAddr,
		arg.DistanceMeters,
		arg.DrivingSeconds,
		arg.PriceOptions,
		arg.ExpiresAt,
	)
	return err
}