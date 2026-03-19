-- name: CreateResource :one
INSERT INTO resources (name, url_pattern, security_profile_id, traffic_profile_id)
VALUES ($1, $2, $3, $4)
RETURNING id, name, url_pattern, security_profile_id, traffic_profile_id, created_at, updated_at;

-- name: GetResource :one
SELECT id, name, url_pattern, security_profile_id, traffic_profile_id, created_at, updated_at
FROM resources
WHERE id = $1;

-- name: ListResources :many
SELECT id, name, url_pattern, security_profile_id, traffic_profile_id, created_at, updated_at
FROM resources
ORDER BY id;

-- name: UpdateResource :one
UPDATE resources
SET name = $2,
    url_pattern = $3,
    security_profile_id = $4,
    traffic_profile_id = $5,
    updated_at = now()
WHERE id = $1
RETURNING id, name, url_pattern, security_profile_id, traffic_profile_id, created_at, updated_at;

-- name: DeleteResource :one
DELETE FROM resources
WHERE id = $1
RETURNING id;
