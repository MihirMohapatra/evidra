-- name: CreateAPIKey :exec
INSERT INTO api_keys (id, organization_id, name, key_hash, key_prefix, is_active, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetAPIKeyByID :one
SELECT id, organization_id, name, key_hash, key_prefix, is_active, created_at, updated_at
FROM api_keys
WHERE id = $1;

-- name: GetAPIKeyByKeyHash :one
SELECT id, organization_id, name, key_hash, key_prefix, is_active, created_at, updated_at
FROM api_keys
WHERE key_hash = $1;

-- name: ListAPIKeysByOrganization :many
SELECT id, organization_id, name, key_hash, key_prefix, is_active, created_at, updated_at
FROM api_keys
WHERE organization_id = $1
ORDER BY created_at DESC;

-- name: UpdateAPIKey :exec
UPDATE api_keys
SET name = $1, key_hash = $2, key_prefix = $3, is_active = $4, updated_at = $5
WHERE id = $6;

-- name: DeleteAPIKey :exec
DELETE FROM api_keys
WHERE id = $1;
