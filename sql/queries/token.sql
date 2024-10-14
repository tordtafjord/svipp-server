-- name: InsertToken :one
INSERT INTO token (token, expires_at, user_id, role)
VALUES ($1, $2, $3, $4)
    RETURNING user_id, role, expires_at;


-- name: GetSession :one
SELECT user_id, role, expires_at FROM token WHERE token = $1 AND expires_at > CURRENT_TIMESTAMP;

-- name: DeleteSession :exec
DELETE FROM token WHERE token = $1;
