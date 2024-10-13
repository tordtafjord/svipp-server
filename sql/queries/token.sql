-- name: InsertToken :one
INSERT INTO token (token, expires_at, user_id, role)
VALUES ($1, $2, $3, $4)
    RETURNING user_id, role, expires_at;