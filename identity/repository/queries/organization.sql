-- name: CreateOrganization :exec
INSERT INTO organizations (id, name, slug, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5);

-- name: GetOrganizationByID :one
SELECT id, name, slug, created_at, updated_at
FROM organizations
WHERE id = $1;

-- name: GetOrganizationBySlug :one
SELECT id, name, slug, created_at, updated_at
FROM organizations
WHERE slug = $1;

-- name: ListOrganizations :many
SELECT id, name, slug, created_at, updated_at
FROM organizations
ORDER BY created_at DESC;

-- name: UpdateOrganization :exec
UPDATE organizations
SET name = $1, slug = $2, updated_at = $3
WHERE id = $4;

-- name: DeleteOrganization :exec
DELETE FROM organizations
WHERE id = $1;
