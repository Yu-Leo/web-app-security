-- name: CreateTrafficRule :one
INSERT INTO traffic_rules (
  profile_id,
  name,
  description,
  priority,
  dry_run,
  match_all,
  requests_limit,
  period_seconds,
  conditions,
  is_enabled
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id, profile_id, name, description, priority, dry_run, match_all, requests_limit, period_seconds, conditions, is_enabled, created_at, updated_at;

-- name: GetTrafficRule :one
SELECT id, profile_id, name, description, priority, dry_run, match_all, requests_limit, period_seconds, conditions, is_enabled, created_at, updated_at
FROM traffic_rules
WHERE id = $1;

-- name: ListTrafficRules :many
SELECT id, profile_id, name, description, priority, dry_run, match_all, requests_limit, period_seconds, conditions, is_enabled, created_at, updated_at
FROM traffic_rules
ORDER BY id;

-- name: UpdateTrafficRule :one
UPDATE traffic_rules
SET profile_id = $2,
    name = $3,
    description = $4,
    priority = $5,
    dry_run = $6,
    match_all = $7,
    requests_limit = $8,
    period_seconds = $9,
    conditions = $10,
    is_enabled = $11,
    updated_at = now()
WHERE id = $1
RETURNING id, profile_id, name, description, priority, dry_run, match_all, requests_limit, period_seconds, conditions, is_enabled, created_at, updated_at;

-- name: DeleteTrafficRule :one
DELETE FROM traffic_rules
WHERE id = $1
RETURNING id;
