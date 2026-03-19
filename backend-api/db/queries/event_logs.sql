-- name: CreateEventLog :one
INSERT INTO event_logs (
  resource_id,
  occurred_at,
  event_type,
  severity,
  message,
  rule_id,
  profile_id,
  metadata,
  request_id,
  client_ip,
  method,
  path
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING id, resource_id, occurred_at, event_type, severity, message, rule_id, profile_id, metadata, request_id, client_ip, method, path;

-- name: GetEventLog :one
SELECT id, resource_id, occurred_at, event_type, severity, message, rule_id, profile_id, metadata, request_id, client_ip, method, path
FROM event_logs
WHERE id = $1;

-- name: ListEventLogs :many
SELECT id, resource_id, occurred_at, event_type, severity, message, rule_id, profile_id, metadata, request_id, client_ip, method, path
FROM event_logs
ORDER BY id;
