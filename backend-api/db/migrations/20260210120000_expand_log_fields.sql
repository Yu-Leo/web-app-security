-- +goose Up
-- +goose StatementBegin
ALTER TABLE request_logs
  ADD COLUMN host TEXT,
  ADD COLUMN scheme TEXT,
  ADD COLUMN protocol TEXT,
  ADD COLUMN authority TEXT,
  ADD COLUMN query TEXT,
  ADD COLUMN source_port INT,
  ADD COLUMN destination_ip TEXT,
  ADD COLUMN destination_port INT,
  ADD COLUMN sni TEXT,
  ADD COLUMN source_principal TEXT,
  ADD COLUMN request_headers JSONB,
  ADD COLUMN request_body_size INT,
  ADD COLUMN tls_presented BOOLEAN,
  ADD COLUMN tls_subject TEXT,
  ADD COLUMN tls_uri_san TEXT,
  ADD COLUMN tls_dns_san TEXT,
  ADD COLUMN tls_peer_certificate TEXT;

ALTER TABLE event_logs
  ADD COLUMN request_id TEXT,
  ADD COLUMN client_ip TEXT,
  ADD COLUMN method TEXT,
  ADD COLUMN path TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE event_logs
  DROP COLUMN path,
  DROP COLUMN method,
  DROP COLUMN client_ip,
  DROP COLUMN request_id;

ALTER TABLE request_logs
  DROP COLUMN tls_peer_certificate,
  DROP COLUMN tls_dns_san,
  DROP COLUMN tls_uri_san,
  DROP COLUMN tls_subject,
  DROP COLUMN tls_presented,
  DROP COLUMN request_body_size,
  DROP COLUMN request_headers,
  DROP COLUMN source_principal,
  DROP COLUMN sni,
  DROP COLUMN destination_port,
  DROP COLUMN destination_ip,
  DROP COLUMN source_port,
  DROP COLUMN query,
  DROP COLUMN authority,
  DROP COLUMN protocol,
  DROP COLUMN scheme,
  DROP COLUMN host;
-- +goose StatementEnd
