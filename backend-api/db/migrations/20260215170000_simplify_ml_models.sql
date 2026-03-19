-- +goose Up
-- +goose StatementBegin
ALTER TABLE ml_models
  ADD COLUMN IF NOT EXISTS model_data BYTEA NOT NULL DEFAULT '\\x'::bytea;

ALTER TABLE ml_models
  DROP COLUMN IF EXISTS description,
  DROP COLUMN IF EXISTS version,
  DROP COLUMN IF EXISTS status,
  DROP COLUMN IF EXISTS config,
  DROP COLUMN IF EXISTS artifact_url,
  DROP COLUMN IF EXISTS created_at,
  DROP COLUMN IF EXISTS updated_at;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ml_models
  ADD COLUMN IF NOT EXISTS description TEXT,
  ADD COLUMN IF NOT EXISTS version TEXT NOT NULL DEFAULT 'v1',
  ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'inactive',
  ADD COLUMN IF NOT EXISTS config JSONB,
  ADD COLUMN IF NOT EXISTS artifact_url TEXT,
  ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

ALTER TABLE ml_models
  DROP COLUMN IF EXISTS model_data;
-- +goose StatementEnd
