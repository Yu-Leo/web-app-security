-- +goose Up
-- +goose StatementBegin
ALTER TABLE request_logs
  ADD COLUMN request_http_id TEXT,
  ADD COLUMN fragment TEXT,
  ADD COLUMN request_body TEXT,
  ADD COLUMN source_service TEXT,
  ADD COLUMN source_labels JSONB,
  ADD COLUMN destination_service TEXT,
  ADD COLUMN destination_labels JSONB,
  ADD COLUMN context_extensions JSONB,
  ADD COLUMN metadata_context JSONB,
  ADD COLUMN route_metadata_context JSONB,
  DROP COLUMN sni,
  DROP COLUMN tls_presented,
  DROP COLUMN tls_subject,
  DROP COLUMN tls_uri_san,
  DROP COLUMN tls_dns_san,
  DROP COLUMN tls_peer_certificate;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE request_logs
  ADD COLUMN sni TEXT,
  ADD COLUMN tls_presented BOOLEAN,
  ADD COLUMN tls_subject TEXT,
  ADD COLUMN tls_uri_san TEXT,
  ADD COLUMN tls_dns_san TEXT,
  ADD COLUMN tls_peer_certificate TEXT,
  DROP COLUMN route_metadata_context,
  DROP COLUMN metadata_context,
  DROP COLUMN context_extensions,
  DROP COLUMN destination_labels,
  DROP COLUMN destination_service,
  DROP COLUMN source_labels,
  DROP COLUMN source_service,
  DROP COLUMN request_body,
  DROP COLUMN fragment,
  DROP COLUMN request_http_id;
-- +goose StatementEnd
