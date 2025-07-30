-- +goose Up
-- +goose StatementBegin
ALTER TABLE interests ALTER COLUMN created_at SET DEFAULT now();
ALTER TABLE interests ALTER COLUMN updated_at SET DEFAULT now();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE interests ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE interests ALTER COLUMN updated_at DROP DEFAULT;
-- +goose StatementEnd
