-- name: CreateShopifyApiKey :one
INSERT INTO shopify_api_config (api_key, quote_key, business_id, location_name,
                                pickup_address, pickup_coords, pickup_instructions)
VALUES ($1, $2, $3, $4,
        $5, $6, $7)
RETURNING business_id, location_name, pickup_address, pickup_coords, pickup_instructions;

-- name: GetApiKeyInfo :one
SELECT business_id, location_name, pickup_address, pickup_coords, pickup_instructions
FROM shopify_api_config WHERE api_key = $1 AND deleted_at IS NULL;


-- name: GetQuoteKeyInfo :one
SELECT business_id, pickup_address
FROM shopify_api_config WHERE quote_key = $1 AND deleted_at IS NULL;



-- name: GetShopifyConfigs :many
SELECT * FROM shopify_api_config WHERE business_id = $1 AND deleted_at IS NULL;

-- name: InsertBusinessHours :many
INSERT INTO business_hours (api_key, day_of_week, opens_at, closes_at)
SELECT $1, unnest(@day_of_week::int[]) AS day_of_week, unnest(@opening_times::time[]), unnest(@closing_times::time[])
RETURNING day_of_week, opens_at, closes_at;

-- name: GetBusinessHours :many
SELECT opens_at, closes_at FROM business_hours WHERE api_key = $1 ORDER BY day_of_week ASC;