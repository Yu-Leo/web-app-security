-- +goose Up
-- +goose StatementBegin
CREATE TABLE event_logs (
  id BIGSERIAL PRIMARY KEY,
  resource_id BIGINT NOT NULL,
  occurred_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  event_type TEXT NOT NULL,
  severity TEXT NOT NULL,
  message TEXT NOT NULL,
  rule_id BIGINT,
  profile_id BIGINT,
  metadata JSONB,
  CONSTRAINT fk_event_logs_resource
    FOREIGN KEY (resource_id)
    REFERENCES resources(id)
    ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE event_logs;
-- +goose StatementEnd
