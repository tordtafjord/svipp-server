-- name: CreateBusiness :exec
INSERT INTO business (id, name, org_id, address)
VALUES ($1, $2, $3, $4);
