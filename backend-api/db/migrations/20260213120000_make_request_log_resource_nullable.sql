-- +goose Up
-- +goose StatementBegin
ALTER TABLE request_logs
  ALTER COLUMN resource_id DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE request_logs
  ALTER COLUMN resource_id SET NOT NULL;
-- +goose StatementEnd
