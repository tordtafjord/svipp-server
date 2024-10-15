-- name: CreateBusiness :exec
INSERT INTO business (name, org_id, address)
VALUES ($1, $2, $3);
