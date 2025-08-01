-- +goose Up
-- +goose StatementBegin
CREATE TYPE role_type AS ENUM (
    'USER',
    'ADMIN'
);
ALTER TABLE users ADD COLUMN role role_type not null default 'USER'::role_type;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN role;
DROP TYPE role_type;
-- +goose StatementEnd
