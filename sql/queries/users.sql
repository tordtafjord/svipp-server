-- name: CreateUser :one
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
RETURNING id, first_name, last_name, phone, email, role;


-- name: GetUserBasicInfoById :one
SELECT id, first_name, last_name, phone, email
FROM users
WHERE id = $1;



-- name: GetOrCreateTempUser :one
INSERT INTO users (phone, first_name, last_name, email, temporary)
VALUES ($1, $2, $3, $4, true)
ON CONFLICT (phone) DO UPDATE SET phone = EXCLUDED.phone
RETURNING id, first_name, last_name, phone, email, role, device_token;



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
