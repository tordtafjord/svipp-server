// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: shopify-requests.sql

package database

import (
	"context"
)

const createShopifyRequest = `-- name: CreateShopifyRequest :exec
INSERT INTO shopify_requests (raw_json)
VALUES ($1)
`

func (q *Queries) CreateShopifyRequest(ctx context.Context, rawJson []byte) error {
	_, err := q.db.Exec(ctx, createShopifyRequest, rawJson)
	return err
}
