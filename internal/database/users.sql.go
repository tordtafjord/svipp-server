// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"
)

const createDriver = `-- name: CreateDriver :exec
INSERT INTO driver (id)
VALUES ($1)
`

func (q *Queries) CreateDriver(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, createDriver, id)
	return err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (first_name, last_name, phone, email, password, device_token, temporary, role)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (phone) DO UPDATE
    SET first_name = EXCLUDED.first_name,
        last_name = EXCLUDED.last_name,
        phone = COALESCE(EXCLUDED.phone, users.phone),
        email = COALESCE(EXCLUDED.email, users.email),
        password = EXCLUDED.password,
        device_token = EXCLUDED.device_token,
        temporary = EXCLUDED.temporary,
        role = EXCLUDED.role,
        created_at = CURRENT_TIMESTAMP,
        updated_at = CURRENT_TIMESTAMP
WHERE users.temporary = true
RETURNING id, first_name, last_name, phone, email, role
`

type CreateUserParams struct {
	FirstName   *string `json:"firstName"`
	LastName    *string `json:"lastName"`
	Phone       *string `json:"phone"`
	Email       *string `json:"email"`
	Password    *string `json:"password"`
	DeviceToken *string `json:"deviceToken"`
	Temporary   *bool   `json:"temporary"`
	Role        string  `json:"role"`
}

type CreateUserRow struct {
	ID        int64   `json:"id"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Phone     *string `json:"phone"`
	Email     *string `json:"email"`
	Role      string  `json:"role"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.FirstName,
		arg.LastName,
		arg.Phone,
		arg.Email,
		arg.Password,
		arg.DeviceToken,
		arg.Temporary,
		arg.Role,
	)
	var i CreateUserRow
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Phone,
		&i.Email,
		&i.Role,
	)
	return i, err
}

const getDeviceTokenByUserID = `-- name: GetDeviceTokenByUserID :one
SELECT device_token
FROM users
WHERE id = $1
`

func (q *Queries) GetDeviceTokenByUserID(ctx context.Context, id int64) (*string, error) {
	row := q.db.QueryRow(ctx, getDeviceTokenByUserID, id)
	var device_token *string
	err := row.Scan(&device_token)
	return device_token, err
}

const getDriverById = `-- name: GetDriverById :one
SELECT id, status, created_at, updated_at, deleted_at
FROM driver WHERE id = $1
`

func (q *Queries) GetDriverById(ctx context.Context, id int64) (Driver, error) {
	row := q.db.QueryRow(ctx, getDriverById, id)
	var i Driver
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getOrCreateTempUser = `-- name: GetOrCreateTempUser :one
INSERT INTO users (phone, first_name, last_name, email, temporary)
VALUES ($1, $2, $3, $4, true)
ON CONFLICT (phone) DO UPDATE SET phone = EXCLUDED.phone
RETURNING id, first_name, last_name, phone, email, role, device_token
`

type GetOrCreateTempUserParams struct {
	Phone     *string `json:"phone"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Email     *string `json:"email"`
}

type GetOrCreateTempUserRow struct {
	ID          int64   `json:"id"`
	FirstName   *string `json:"firstName"`
	LastName    *string `json:"lastName"`
	Phone       *string `json:"phone"`
	Email       *string `json:"email"`
	Role        string  `json:"role"`
	DeviceToken *string `json:"deviceToken"`
}

func (q *Queries) GetOrCreateTempUser(ctx context.Context, arg GetOrCreateTempUserParams) (GetOrCreateTempUserRow, error) {
	row := q.db.QueryRow(ctx, getOrCreateTempUser,
		arg.Phone,
		arg.FirstName,
		arg.LastName,
		arg.Email,
	)
	var i GetOrCreateTempUserRow
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Phone,
		&i.Email,
		&i.Role,
		&i.DeviceToken,
	)
	return i, err
}

const getUserBasicInfoById = `-- name: GetUserBasicInfoById :one
SELECT id, first_name, last_name, phone, email
FROM users
WHERE id = $1
`

type GetUserBasicInfoByIdRow struct {
	ID        int64   `json:"id"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Phone     *string `json:"phone"`
	Email     *string `json:"email"`
}

func (q *Queries) GetUserBasicInfoById(ctx context.Context, id int64) (GetUserBasicInfoByIdRow, error) {
	row := q.db.QueryRow(ctx, getUserBasicInfoById, id)
	var i GetUserBasicInfoByIdRow
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Phone,
		&i.Email,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, first_name, last_name, phone, email, password, device_token, temporary, role, rate_total, rates, created_at, updated_at, deleted_at
FROM users WHERE email = $1 AND email IS NOT NULL
`

func (q *Queries) GetUserByEmail(ctx context.Context, email *string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Phone,
		&i.Email,
		&i.Password,
		&i.DeviceToken,
		&i.Temporary,
		&i.Role,
		&i.RateTotal,
		&i.Rates,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getUserByPhone = `-- name: GetUserByPhone :one
SELECT id
FROM users WHERE phone = $1
`

func (q *Queries) GetUserByPhone(ctx context.Context, phone *string) (int64, error) {
	row := q.db.QueryRow(ctx, getUserByPhone, phone)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const updateDeviceTokenByUserID = `-- name: UpdateDeviceTokenByUserID :exec
UPDATE users
SET device_token = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
`

type UpdateDeviceTokenByUserIDParams struct {
	ID          int64   `json:"id"`
	DeviceToken *string `json:"deviceToken"`
}

func (q *Queries) UpdateDeviceTokenByUserID(ctx context.Context, arg UpdateDeviceTokenByUserIDParams) error {
	_, err := q.db.Exec(ctx, updateDeviceTokenByUserID, arg.ID, arg.DeviceToken)
	return err
}
