-- name: CreateSecurityProfile :one
INSERT INTO security_profiles (name, description, base_action, log_enabled, is_enabled)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, description, base_action, log_enabled, is_enabled, created_at, updated_at;

-- name: GetSecurityProfile :one
SELECT id, name, description, base_action, log_enabled, is_enabled, created_at, updated_at
FROM security_profiles
WHERE id = $1;

-- name: ListSecurityProfiles :many
SELECT id, name, description, base_action, log_enabled, is_enabled, created_at, updated_at
FROM security_profiles
ORDER BY id;

-- name: UpdateSecurityProfile :one
UPDATE security_profiles
SET name = $2,
    description = $3,
    base_action = $4,
    log_enabled = $5,
    is_enabled = $6,
    updated_at = now()
WHERE id = $1
RETURNING id, name, description, base_action, log_enabled, is_enabled, created_at, updated_at;

-- name: DeleteSecurityProfile :one
DELETE FROM security_profiles
WHERE id = $1
RETURNING id;
