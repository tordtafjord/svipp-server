-- name: UpsertQuote :exec
INSERT INTO order_quote(
    user_id,
    pickup_address,
    delivery_address,
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
        $7)
ON CONFLICT (user_id, pickup_address ,delivery_address)
DO UPDATE SET
              distance_meters = excluded.distance_meters,
              driving_seconds = excluded.driving_seconds,
              price_options = excluded.price_options,
              expires_at = excluded.expires_at,
              created_at = CURRENT_TIMESTAMP;


-- name: GetOrderQuote :one
SELECT
    pickup_address,
    delivery_address,
    distance_meters,
    driving_seconds,
    price_options
FROM order_quote
WHERE
    user_id = $1 AND
    pickup_address = $2 AND
    delivery_address = $3 AND
    expires_at > NOW();
