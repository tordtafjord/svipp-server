// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: shopify_api_config.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createShopifyApiKey = `-- name: CreateShopifyApiKey :one
INSERT INTO shopify_api_config (api_key, quote_key, business_id, location_name,
                                pickup_address, pickup_coords, pickup_instructions)
VALUES ($1, $2, $3, $4,
        $5, $6, $7)
RETURNING business_id, location_name, pickup_address, pickup_coords, pickup_instructions
`

type CreateShopifyApiKeyParams struct {
	ApiKey             []byte      `json:"apiKey"`
	QuoteKey           string      `json:"quoteKey"`
	BusinessID         int64       `json:"businessId"`
	LocationName       *string     `json:"locationName"`
	PickupAddress      *string     `json:"pickupAddress"`
	PickupCoords       interface{} `json:"pickupCoords"`
	PickupInstructions *string     `json:"pickupInstructions"`
}

type CreateShopifyApiKeyRow struct {
	BusinessID         int64       `json:"businessId"`
	LocationName       *string     `json:"locationName"`
	PickupAddress      *string     `json:"pickupAddress"`
	PickupCoords       interface{} `json:"pickupCoords"`
	PickupInstructions *string     `json:"pickupInstructions"`
}

func (q *Queries) CreateShopifyApiKey(ctx context.Context, arg CreateShopifyApiKeyParams) (CreateShopifyApiKeyRow, error) {
	row := q.db.QueryRow(ctx, createShopifyApiKey,
		arg.ApiKey,
		arg.QuoteKey,
		arg.BusinessID,
		arg.LocationName,
		arg.PickupAddress,
		arg.PickupCoords,
		arg.PickupInstructions,
	)
	var i CreateShopifyApiKeyRow
	err := row.Scan(
		&i.BusinessID,
		&i.LocationName,
		&i.PickupAddress,
		&i.PickupCoords,
		&i.PickupInstructions,
	)
	return i, err
}

const getApiKeyInfo = `-- name: GetApiKeyInfo :one
SELECT business_id, location_name, pickup_address, pickup_coords, pickup_instructions
FROM shopify_api_config WHERE api_key = $1 AND deleted_at IS NULL
`

type GetApiKeyInfoRow struct {
	BusinessID         int64       `json:"businessId"`
	LocationName       *string     `json:"locationName"`
	PickupAddress      *string     `json:"pickupAddress"`
	PickupCoords       interface{} `json:"pickupCoords"`
	PickupInstructions *string     `json:"pickupInstructions"`
}

func (q *Queries) GetApiKeyInfo(ctx context.Context, apiKey []byte) (GetApiKeyInfoRow, error) {
	row := q.db.QueryRow(ctx, getApiKeyInfo, apiKey)
	var i GetApiKeyInfoRow
	err := row.Scan(
		&i.BusinessID,
		&i.LocationName,
		&i.PickupAddress,
		&i.PickupCoords,
		&i.PickupInstructions,
	)
	return i, err
}

const getBusinessHours = `-- name: GetBusinessHours :many
SELECT opens_at, closes_at FROM business_hours WHERE api_key = $1 ORDER BY day_of_week ASC
`

type GetBusinessHoursRow struct {
	OpensAt  pgtype.Time `json:"opensAt"`
	ClosesAt pgtype.Time `json:"closesAt"`
}

func (q *Queries) GetBusinessHours(ctx context.Context, apiKey []byte) ([]GetBusinessHoursRow, error) {
	rows, err := q.db.Query(ctx, getBusinessHours, apiKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetBusinessHoursRow
	for rows.Next() {
		var i GetBusinessHoursRow
		if err := rows.Scan(&i.OpensAt, &i.ClosesAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getQuoteKeyInfo = `-- name: GetQuoteKeyInfo :one
SELECT business_id, pickup_address
FROM shopify_api_config WHERE quote_key = $1 AND deleted_at IS NULL
`

type GetQuoteKeyInfoRow struct {
	BusinessID    int64   `json:"businessId"`
	PickupAddress *string `json:"pickupAddress"`
}

func (q *Queries) GetQuoteKeyInfo(ctx context.Context, quoteKey string) (GetQuoteKeyInfoRow, error) {
	row := q.db.QueryRow(ctx, getQuoteKeyInfo, quoteKey)
	var i GetQuoteKeyInfoRow
	err := row.Scan(&i.BusinessID, &i.PickupAddress)
	return i, err
}

const getShopifyConfigs = `-- name: GetShopifyConfigs :many
SELECT api_key, quote_key, business_id, pickup_address, pickup_coords, pickup_instructions, created_at, updated_at, deleted_at, location_name FROM shopify_api_config WHERE business_id = $1 AND deleted_at IS NULL
`

func (q *Queries) GetShopifyConfigs(ctx context.Context, businessID int64) ([]ShopifyApiConfig, error) {
	rows, err := q.db.Query(ctx, getShopifyConfigs, businessID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ShopifyApiConfig
	for rows.Next() {
		var i ShopifyApiConfig
		if err := rows.Scan(
			&i.ApiKey,
			&i.QuoteKey,
			&i.BusinessID,
			&i.PickupAddress,
			&i.PickupCoords,
			&i.PickupInstructions,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.LocationName,
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

const insertBusinessHours = `-- name: InsertBusinessHours :many
INSERT INTO business_hours (api_key, day_of_week, opens_at, closes_at)
SELECT $1, unnest($2::int[]) AS day_of_week, unnest($3::time[]), unnest($4::time[])
RETURNING day_of_week, opens_at, closes_at
`

type InsertBusinessHoursParams struct {
	ApiKey       []byte        `json:"apiKey"`
	DayOfWeek    []int32       `json:"dayOfWeek"`
	OpeningTimes []pgtype.Time `json:"openingTimes"`
	ClosingTimes []pgtype.Time `json:"closingTimes"`
}

type InsertBusinessHoursRow struct {
	DayOfWeek int32       `json:"dayOfWeek"`
	OpensAt   pgtype.Time `json:"opensAt"`
	ClosesAt  pgtype.Time `json:"closesAt"`
}

func (q *Queries) InsertBusinessHours(ctx context.Context, arg InsertBusinessHoursParams) ([]InsertBusinessHoursRow, error) {
	rows, err := q.db.Query(ctx, insertBusinessHours,
		arg.ApiKey,
		arg.DayOfWeek,
		arg.OpeningTimes,
		arg.ClosingTimes,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []InsertBusinessHoursRow
	for rows.Next() {
		var i InsertBusinessHoursRow
		if err := rows.Scan(&i.DayOfWeek, &i.OpensAt, &i.ClosesAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
