-- name: CreateSecurityRule :one
INSERT INTO security_rules (
  profile_id,
  name,
  description,
  priority,
  rule_type,
  action,
  conditions,
  ml_model_id,
  ml_threshold,
  dry_run,
  is_enabled
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id, profile_id, name, description, priority, rule_type, action, conditions, ml_model_id, ml_threshold, dry_run, is_enabled, created_at, updated_at;

-- name: GetSecurityRule :one
SELECT id, profile_id, name, description, priority, rule_type, action, conditions, ml_model_id, ml_threshold, dry_run, is_enabled, created_at, updated_at
FROM security_rules
WHERE id = $1;

-- name: ListSecurityRules :many
SELECT id, profile_id, name, description, priority, rule_type, action, conditions, ml_model_id, ml_threshold, dry_run, is_enabled, created_at, updated_at
FROM security_rules
ORDER BY id;

-- name: UpdateSecurityRule :one
UPDATE security_rules
SET profile_id = $2,
    name = $3,
    description = $4,
    priority = $5,
    rule_type = $6,
    action = $7,
    conditions = $8,
    ml_model_id = $9,
    ml_threshold = $10,
    dry_run = $11,
    is_enabled = $12,
    updated_at = now()
WHERE id = $1
RETURNING id, profile_id, name, description, priority, rule_type, action, conditions, ml_model_id, ml_threshold, dry_run, is_enabled, created_at, updated_at;

-- name: DeleteSecurityRule :one
DELETE FROM security_rules
WHERE id = $1
RETURNING id;
