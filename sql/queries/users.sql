-- name: CreateUser :one
INSERT INTO users (name, phone, email, password, device_token, temporary, role)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, name, phone, email, role;


-- name: GetUserBasicInfoById :one
SELECT id, name, phone, email
FROM users
WHERE id = $1;



-- name: GetOrCreateTempUser :one
INSERT INTO users (phone, name, email, temporary)
VALUES ($1, $2, $3, true)
ON CONFLICT (phone) DO UPDATE
SET name = CASE 
        WHEN users.temporary THEN $2
        ELSE users.name
    END,
    email = CASE 
        WHEN users.temporary THEN $3
        ELSE users.email
    END
RETURNING id, name, phone, email, role, device_token;



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
