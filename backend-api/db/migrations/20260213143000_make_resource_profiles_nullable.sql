-- +goose Up
-- +goose StatementBegin
ALTER TABLE resources
  ALTER COLUMN security_profile_id DROP NOT NULL,
  ALTER COLUMN traffic_profile_id DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE resources
  ALTER COLUMN security_profile_id SET NOT NULL,
  ALTER COLUMN traffic_profile_id SET NOT NULL;
-- +goose StatementEnd
