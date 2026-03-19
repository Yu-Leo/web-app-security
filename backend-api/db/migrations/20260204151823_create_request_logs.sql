-- +goose Up
-- +goose StatementBegin
CREATE TABLE request_logs (
  id BIGSERIAL PRIMARY KEY,
  resource_id BIGINT NOT NULL,
  occurred_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  client_ip TEXT NOT NULL,
  method TEXT NOT NULL,
  path TEXT NOT NULL,
  status_code INT NOT NULL,
  action TEXT NOT NULL,
  rule_id BIGINT,
  profile_id BIGINT,
  user_agent TEXT,
  country TEXT,
  latency_ms INT,
  request_id TEXT,
  metadata JSONB,
  CONSTRAINT fk_request_logs_resource
    FOREIGN KEY (resource_id)
    REFERENCES resources(id)
    ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE request_logs;
-- +goose StatementEnd
