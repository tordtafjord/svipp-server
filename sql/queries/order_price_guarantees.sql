-- name: UpsertQuote :exec
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
              expires_at = excluded.expires_at;