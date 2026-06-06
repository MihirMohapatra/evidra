-- name: CreateUser :exec
INSERT INTO users (id, organization_id, email, password_hash, role, is_active, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetUserByID :one
SELECT id, organization_id, email, password_hash, role, is_active, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, organization_id, email, password_hash, role, is_active, created_at, updated_at
FROM users
WHERE email = $1;

-- name: ListUsersByOrganization :many
SELECT id, organization_id, email, password_hash, role, is_active, created_at, updated_at
FROM users
WHERE organization_id = $1
ORDER BY created_at DESC;

-- name: UpdateUser :exec
UPDATE users
SET email = $1, password_hash = $2, role = $3, is_active = $4, updated_at = $5
WHERE id = $6;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
