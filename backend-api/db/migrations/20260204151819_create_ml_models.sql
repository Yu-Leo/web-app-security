-- +goose Up
-- +goose StatementBegin
CREATE TABLE ml_models (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  version TEXT NOT NULL,
  status TEXT NOT NULL,
  config JSONB,
  artifact_url TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ml_models;
-- +goose StatementEnd
