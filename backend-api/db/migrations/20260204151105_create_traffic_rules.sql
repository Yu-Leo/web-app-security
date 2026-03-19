-- +goose Up
-- +goose StatementBegin
CREATE TABLE traffic_rules (
  id BIGSERIAL PRIMARY KEY,
  profile_id BIGINT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  priority INT NOT NULL,
  dry_run BOOLEAN NOT NULL DEFAULT false,
  match_all BOOLEAN NOT NULL DEFAULT false,
  requests_limit INT NOT NULL,
  period_seconds INT NOT NULL,
  conditions JSONB,
  is_enabled BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT fk_traffic_rules_profile
    FOREIGN KEY (profile_id)
    REFERENCES traffic_profiles(id)
    ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE traffic_rules;
-- +goose StatementEnd
