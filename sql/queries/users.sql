-- name: CreateUser :one
INSERT INTO users (name, phone, email, password, device_token, temporary)
VALUES ($1, $2, $3, $4, $5, false)
RETURNING id, name, phone, email;


-- name: GetUserBasicInfoById :one
SELECT id, name, phone, email
FROM users
WHERE id = $1;


-- name: CreateTemporaryUser :one
INSERT INTO users (phone, temporary)
VALUES ($1, true)
RETURNING id;

-- name: GetOrCreateTempUser :one
WITH existing_user AS (
    SELECT * FROM users WHERE phone = $1
)
INSERT INTO users (phone, temporary)
SELECT $1, true
WHERE NOT EXISTS (SELECT 1 FROM existing_user)
RETURNING *;

SELECT * FROM existing_user
UNION ALL
SELECT * FROM users WHERE phone = $1 AND NOT EXISTS (SELECT 1 FROM existing_user);



-- name: GetUserByEmail :one
SELECT *
FROM users WHERE email = $1 AND email IS NOT NULL;



-- name: GetUserByPhone :one
SELECT id
FROM users WHERE phone = $1;


-- name: CreateDriver :exec
INSERT INTO driver (id)
VALUES ($1);



-- name: GetDriverById :one
SELECT *
FROM driver WHERE id = $1;

-- name: GetDeviceTokenByUserID :one
SELECT device_token
FROM users
WHERE id = $1;

-- name: UpdateDeviceTokenByUserID :exec
UPDATE users
SET device_token = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;
