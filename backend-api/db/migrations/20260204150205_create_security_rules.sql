-- +goose Up
-- +goose StatementBegin
CREATE TABLE security_rules (
  id BIGSERIAL PRIMARY KEY,
  profile_id BIGINT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  priority INT NOT NULL,
  rule_type TEXT NOT NULL,
  action TEXT NOT NULL,
  conditions JSONB,
  dry_run BOOLEAN NOT NULL DEFAULT false,
  is_enabled BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT fk_security_rules_profile
    FOREIGN KEY (profile_id)
    REFERENCES security_profiles(id)
    ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE security_rules;
-- +goose StatementEnd
