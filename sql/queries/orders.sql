-- name: CreateOrder :one
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
RETURNING pickup_address, delivery_address, distance, price_cents, status, public_id;


-- name: GetOrderInfoByPublicId :one
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
-- Add more receiver columns as needed
FROM
    orders o
        LEFT JOIN
    users sender ON o.sender_id = sender.id
        LEFT JOIN
    users driver ON o.driver_id = driver.id
WHERE o.public_id = $1::uuid;


-- name: GetOrderDriverIdByOrderId :many
SELECT driver_id
FROM orders
WHERE id = $1;

-- name: GetOrdersByUserId :many
SELECT *
FROM orders
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetOrdersByDriverId :many
SELECT *
FROM orders
WHERE driver_id = $1
ORDER BY created_at DESC;

-- name: ConfirmOrderById :one
UPDATE orders
SET
    status = 'CONFIRMED',
    confirmed_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: SetDriverIdByOrderId :one
UPDATE orders
SET
    driver_id = $1,  -- New driver_id to be set
    status = 'ACCEPTED',
    accepted_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $2      -- Order ID for which the driver_id needs to be updated
  AND driver_id IS NULL  -- Only update if driver_id is not already set
RETURNING *;


