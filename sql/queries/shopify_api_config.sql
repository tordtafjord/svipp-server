-- name: CreateShopifyApiKey :one
INSERT INTO shopify_api_config (api_key, quote_key, business_id,
                                pickup_address, pickup_coords, pickup_instructions,
                                pickup_window_start, pickup_window_end)
VALUES ($1, $2, $3, $4,
        $5, $6, $7, $8)
RETURNING business_id, pickup_address, pickup_coords, pickup_instructions, pickup_window_start, pickup_window_end;


-- name: GetQuoteKeyInfo :one
SELECT business_id, pickup_address, pickup_window_start, pickup_window_end
FROM shopify_api_config WHERE quote_key = $1 AND deleted_at IS NULL;

-- name: GetApiKeyInfo :one
SELECT business_id, pickup_address, pickup_coords, pickup_instructions, pickup_window_start, pickup_window_end
FROM shopify_api_config WHERE api_key = $1 AND deleted_at IS NULL;

