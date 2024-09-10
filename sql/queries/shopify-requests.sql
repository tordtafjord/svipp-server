-- name: CreateShopifyRequest :exec
INSERT INTO shopify_requests (raw_json)
VALUES ($1);
