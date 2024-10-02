// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: orders.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const confirmOrderById = `-- name: ConfirmOrderById :one
UPDATE orders
SET
    status = 'CONFIRMED',
    confirmed_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, user_id, sender_id, recipient_id, driver_id, pickup_address, delivery_address, status, distance, driving_seconds, price_cents, created_at, confirmed_at, accepted_at, picked_up_at, delivered_at, updated_at, cancelled_at, public_id, pickup_coords, delivery_coords, delivery_window_start, delivery_window_end
`

func (q *Queries) ConfirmOrderById(ctx context.Context, id int32) (Order, error) {
	row := q.db.QueryRow(ctx, confirmOrderById, id)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SenderID,
		&i.RecipientID,
		&i.DriverID,
		&i.PickupAddress,
		&i.DeliveryAddress,
		&i.Status,
		&i.Distance,
		&i.DrivingSeconds,
		&i.PriceCents,
		&i.CreatedAt,
		&i.ConfirmedAt,
		&i.AcceptedAt,
		&i.PickedUpAt,
		&i.DeliveredAt,
		&i.UpdatedAt,
		&i.CancelledAt,
		&i.PublicID,
		&i.PickupCoords,
		&i.DeliveryCoords,
		&i.DeliveryWindowStart,
		&i.DeliveryWindowEnd,
	)
	return i, err
}

const createOrder = `-- name: CreateOrder :one
INSERT INTO orders (
    user_id,
    sender_id,
    recipient_id,
    pickup_address,
    delivery_address,
    distance,
    driving_seconds,
    price_cents,
    status
)
VALUES (
           $1,
           $2,
           $3,
           $4,
           $5,
           $6,
           $7,
           $8,
           $9
       )
RETURNING pickup_address, delivery_address, distance, price_cents, status, public_id::text
`

type CreateOrderParams struct {
	UserID          int32  `json:"userId"`
	SenderID        int32  `json:"senderId"`
	RecipientID     int32  `json:"recipientId"`
	PickupAddress   string `json:"pickupAddress"`
	DeliveryAddress string `json:"deliveryAddress"`
	Distance        int32  `json:"distance"`
	DrivingSeconds  int32  `json:"drivingSeconds"`
	PriceCents      int32  `json:"priceCents"`
	Status          string `json:"status"`
}

type CreateOrderRow struct {
	PickupAddress   string `json:"pickupAddress"`
	DeliveryAddress string `json:"deliveryAddress"`
	Distance        int32  `json:"distance"`
	PriceCents      int32  `json:"priceCents"`
	Status          string `json:"status"`
	PublicID        string `json:"publicId"`
}

func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (CreateOrderRow, error) {
	row := q.db.QueryRow(ctx, createOrder,
		arg.UserID,
		arg.SenderID,
		arg.RecipientID,
		arg.PickupAddress,
		arg.DeliveryAddress,
		arg.Distance,
		arg.DrivingSeconds,
		arg.PriceCents,
		arg.Status,
	)
	var i CreateOrderRow
	err := row.Scan(
		&i.PickupAddress,
		&i.DeliveryAddress,
		&i.Distance,
		&i.PriceCents,
		&i.Status,
		&i.PublicID,
	)
	return i, err
}

const getOrderDriverIdByOrderId = `-- name: GetOrderDriverIdByOrderId :many
SELECT driver_id
FROM orders
WHERE id = $1
`

func (q *Queries) GetOrderDriverIdByOrderId(ctx context.Context, id int32) ([]*int32, error) {
	rows, err := q.db.Query(ctx, getOrderDriverIdByOrderId, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*int32
	for rows.Next() {
		var driver_id *int32
		if err := rows.Scan(&driver_id); err != nil {
			return nil, err
		}
		items = append(items, driver_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getOrderInfoByPublicId = `-- name: GetOrderInfoByPublicId :one
SELECT
    o.pickup_address,
    o.pickup_coords,
    o.delivery_address,
    o.delivery_coords,
    o.status,
    o.distance,
    o.driving_seconds,
    o.created_at,
    o.confirmed_at,
    o.accepted_at,
    o.picked_up_at,
    o.delivered_at,
    o.cancelled_at,
    sender.name AS sender_name,
    -- Add more sender columns as needed
    driver.name AS driver_name,
    driver.rates AS driver_rates,
    driver.rate_total AS driver_rate_total
FROM
    orders o
        LEFT JOIN
    users sender ON o.sender_id = sender.id
        LEFT JOIN
    users driver ON o.driver_id = driver.id
WHERE o.public_id = $1::uuid
`

type GetOrderInfoByPublicIdRow struct {
	PickupAddress   string             `json:"pickupAddress"`
	PickupCoords    interface{}        `json:"pickupCoords"`
	DeliveryAddress string             `json:"deliveryAddress"`
	DeliveryCoords  interface{}        `json:"deliveryCoords"`
	Status          string             `json:"status"`
	Distance        int32              `json:"distance"`
	DrivingSeconds  int32              `json:"drivingSeconds"`
	CreatedAt       pgtype.Timestamptz `json:"createdAt"`
	ConfirmedAt     pgtype.Timestamptz `json:"confirmedAt"`
	AcceptedAt      pgtype.Timestamptz `json:"acceptedAt"`
	PickedUpAt      pgtype.Timestamptz `json:"pickedUpAt"`
	DeliveredAt     pgtype.Timestamptz `json:"deliveredAt"`
	CancelledAt     pgtype.Timestamptz `json:"cancelledAt"`
	SenderName      *string            `json:"senderName"`
	DriverName      *string            `json:"driverName"`
	DriverRates     *int32             `json:"driverRates"`
	DriverRateTotal *int32             `json:"driverRateTotal"`
}

// Add more receiver columns as needed
func (q *Queries) GetOrderInfoByPublicId(ctx context.Context, dollar_1 pgtype.UUID) (GetOrderInfoByPublicIdRow, error) {
	row := q.db.QueryRow(ctx, getOrderInfoByPublicId, dollar_1)
	var i GetOrderInfoByPublicIdRow
	err := row.Scan(
		&i.PickupAddress,
		&i.PickupCoords,
		&i.DeliveryAddress,
		&i.DeliveryCoords,
		&i.Status,
		&i.Distance,
		&i.DrivingSeconds,
		&i.CreatedAt,
		&i.ConfirmedAt,
		&i.AcceptedAt,
		&i.PickedUpAt,
		&i.DeliveredAt,
		&i.CancelledAt,
		&i.SenderName,
		&i.DriverName,
		&i.DriverRates,
		&i.DriverRateTotal,
	)
	return i, err
}

const getOrdersByDriverId = `-- name: GetOrdersByDriverId :many
SELECT id, user_id, sender_id, recipient_id, driver_id, pickup_address, delivery_address, status, distance, driving_seconds, price_cents, created_at, confirmed_at, accepted_at, picked_up_at, delivered_at, updated_at, cancelled_at, public_id, pickup_coords, delivery_coords, delivery_window_start, delivery_window_end
FROM orders
WHERE driver_id = $1
ORDER BY created_at DESC
`

func (q *Queries) GetOrdersByDriverId(ctx context.Context, driverID *int32) ([]Order, error) {
	rows, err := q.db.Query(ctx, getOrdersByDriverId, driverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Order
	for rows.Next() {
		var i Order
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.SenderID,
			&i.RecipientID,
			&i.DriverID,
			&i.PickupAddress,
			&i.DeliveryAddress,
			&i.Status,
			&i.Distance,
			&i.DrivingSeconds,
			&i.PriceCents,
			&i.CreatedAt,
			&i.ConfirmedAt,
			&i.AcceptedAt,
			&i.PickedUpAt,
			&i.DeliveredAt,
			&i.UpdatedAt,
			&i.CancelledAt,
			&i.PublicID,
			&i.PickupCoords,
			&i.DeliveryCoords,
			&i.DeliveryWindowStart,
			&i.DeliveryWindowEnd,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getOrdersByUserId = `-- name: GetOrdersByUserId :many
SELECT id, user_id, sender_id, recipient_id, driver_id, pickup_address, delivery_address, status, distance, driving_seconds, price_cents, created_at, confirmed_at, accepted_at, picked_up_at, delivered_at, updated_at, cancelled_at, public_id, pickup_coords, delivery_coords, delivery_window_start, delivery_window_end
FROM orders
WHERE user_id = $1
ORDER BY created_at DESC
`

func (q *Queries) GetOrdersByUserId(ctx context.Context, userID int32) ([]Order, error) {
	rows, err := q.db.Query(ctx, getOrdersByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Order
	for rows.Next() {
		var i Order
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.SenderID,
			&i.RecipientID,
			&i.DriverID,
			&i.PickupAddress,
			&i.DeliveryAddress,
			&i.Status,
			&i.Distance,
			&i.DrivingSeconds,
			&i.PriceCents,
			&i.CreatedAt,
			&i.ConfirmedAt,
			&i.AcceptedAt,
			&i.PickedUpAt,
			&i.DeliveredAt,
			&i.UpdatedAt,
			&i.CancelledAt,
			&i.PublicID,
			&i.PickupCoords,
			&i.DeliveryCoords,
			&i.DeliveryWindowStart,
			&i.DeliveryWindowEnd,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const setDriverIdByOrderId = `-- name: SetDriverIdByOrderId :one
UPDATE orders
SET
    driver_id = $1,  -- New driver_id to be set
    status = 'ACCEPTED',
    accepted_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $2      -- Order ID for which the driver_id needs to be updated
  AND driver_id IS NULL  -- Only update if driver_id is not already set
RETURNING id, user_id, sender_id, recipient_id, driver_id, pickup_address, delivery_address, status, distance, driving_seconds, price_cents, created_at, confirmed_at, accepted_at, picked_up_at, delivered_at, updated_at, cancelled_at, public_id, pickup_coords, delivery_coords, delivery_window_start, delivery_window_end
`

type SetDriverIdByOrderIdParams struct {
	DriverID *int32 `json:"driverId"`
	ID       int32  `json:"id"`
}

func (q *Queries) SetDriverIdByOrderId(ctx context.Context, arg SetDriverIdByOrderIdParams) (Order, error) {
	row := q.db.QueryRow(ctx, setDriverIdByOrderId, arg.DriverID, arg.ID)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SenderID,
		&i.RecipientID,
		&i.DriverID,
		&i.PickupAddress,
		&i.DeliveryAddress,
		&i.Status,
		&i.Distance,
		&i.DrivingSeconds,
		&i.PriceCents,
		&i.CreatedAt,
		&i.ConfirmedAt,
		&i.AcceptedAt,
		&i.PickedUpAt,
		&i.DeliveredAt,
		&i.UpdatedAt,
		&i.CancelledAt,
		&i.PublicID,
		&i.PickupCoords,
		&i.DeliveryCoords,
		&i.DeliveryWindowStart,
		&i.DeliveryWindowEnd,
	)
	return i, err
}
