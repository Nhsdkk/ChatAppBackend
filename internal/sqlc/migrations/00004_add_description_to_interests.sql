-- +goose Up
-- +goose StatementBegin
ALTER TABLE interests ADD COLUMN description text NOT NULL default '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE interests DROP COLUMN description;
-- +goose StatementEnd
