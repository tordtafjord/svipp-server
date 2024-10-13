// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: token.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const insertToken = `-- name: InsertToken :one
INSERT INTO token (token, expires_at, user_id, role)
VALUES ($1, $2, $3, $4)
    RETURNING user_id, role, expires_at
`

type InsertTokenParams struct {
	Token     string             `json:"token"`
	ExpiresAt pgtype.Timestamptz `json:"expiresAt"`
	UserID    int32              `json:"userId"`
	Role      string             `json:"role"`
}

type InsertTokenRow struct {
	UserID    int32              `json:"userId"`
	Role      string             `json:"role"`
	ExpiresAt pgtype.Timestamptz `json:"expiresAt"`
}

func (q *Queries) InsertToken(ctx context.Context, arg InsertTokenParams) (InsertTokenRow, error) {
	row := q.db.QueryRow(ctx, insertToken,
		arg.Token,
		arg.ExpiresAt,
		arg.UserID,
		arg.Role,
	)
	var i InsertTokenRow
	err := row.Scan(&i.UserID, &i.Role, &i.ExpiresAt)
	return i, err
}
