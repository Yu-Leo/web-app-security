-- name: CreateRequestLog :one
INSERT INTO request_logs (
  resource_id,
  occurred_at,
  client_ip,
  method,
  path,
  status_code,
  action,
  rule_id,
  profile_id,
  user_agent,
  country,
  latency_ms,
  request_id,
  metadata,
  host,
  scheme,
  protocol,
  authority,
  query,
  source_port,
  destination_ip,
  destination_port,
  source_principal,
  source_service,
  source_labels,
  destination_service,
  destination_labels,
  request_http_id,
  fragment,
  request_headers,
  request_body_size,
  request_body,
  context_extensions,
  metadata_context,
  route_metadata_context
)
VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
  $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
  $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35
)
RETURNING id, resource_id, occurred_at, client_ip, method, path, status_code, action, rule_id, profile_id, user_agent, country, latency_ms, request_id, metadata, host, scheme, protocol, authority, query, source_port, destination_ip, destination_port, source_principal, source_service, source_labels, destination_service, destination_labels, request_http_id, fragment, request_headers, request_body_size, request_body, context_extensions, metadata_context, route_metadata_context;

-- name: GetRequestLog :one
SELECT id, resource_id, occurred_at, client_ip, method, path, status_code, action, rule_id, profile_id, user_agent, country, latency_ms, request_id, metadata, host, scheme, protocol, authority, query, source_port, destination_ip, destination_port, source_principal, source_service, source_labels, destination_service, destination_labels, request_http_id, fragment, request_headers, request_body_size, request_body, context_extensions, metadata_context, route_metadata_context
FROM request_logs
WHERE id = $1;

-- name: ListRequestLogs :many
SELECT id, resource_id, occurred_at, client_ip, method, path, status_code, action, rule_id, profile_id, user_agent, country, latency_ms, request_id, metadata, host, scheme, protocol, authority, query, source_port, destination_ip, destination_port, source_principal, source_service, source_labels, destination_service, destination_labels, request_http_id, fragment, request_headers, request_body_size, request_body, context_extensions, metadata_context, route_metadata_context
FROM request_logs
ORDER BY id;
