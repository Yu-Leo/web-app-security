-- name: CreateTrafficProfile :one
INSERT INTO traffic_profiles (name, description, is_enabled)
VALUES ($1, $2, $3)
RETURNING id, name, description, is_enabled, created_at, updated_at;

-- name: GetTrafficProfile :one
SELECT id, name, description, is_enabled, created_at, updated_at
FROM traffic_profiles
WHERE id = $1;

-- name: ListTrafficProfiles :many
SELECT id, name, description, is_enabled, created_at, updated_at
FROM traffic_profiles
ORDER BY id;

-- name: UpdateTrafficProfile :one
UPDATE traffic_profiles
SET name = $2,
    description = $3,
    is_enabled = $4,
    updated_at = now()
WHERE id = $1
RETURNING id, name, description, is_enabled, created_at, updated_at;

-- name: DeleteTrafficProfile :one
DELETE FROM traffic_profiles
WHERE id = $1
RETURNING id;
