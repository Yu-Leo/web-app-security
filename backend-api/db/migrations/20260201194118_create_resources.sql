-- +goose Up
-- +goose StatementBegin
CREATE TABLE resources (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  url_pattern TEXT NOT NULL,
  security_profile_id BIGINT NOT NULL,
  traffic_profile_id BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE resources;
-- +goose StatementEnd
