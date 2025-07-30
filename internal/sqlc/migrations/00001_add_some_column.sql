-- +goose Up
-- +goose StatementBegin
CREATE TYPE gender AS ENUM (
    'MALE',
    'FEMALE'
    );

CREATE TYPE chat_type AS ENUM (
    'PRIVATE_CHAT',
    'GROUP_CHAT'
    );

CREATE TABLE users
(
    id               uuid primary key      default gen_random_uuid(),
    full_name        varchar(255) not null unique,
    birthday         timestamptz  not null,
    gender           gender       not null,
    email            varchar(255) not null unique,
    password         bytea        not null,
    avatar_file_name varchar(255) not null,
    online           bool         not null,
    email_verified   bool         not null default false,
    last_seen        timestamptz  not null default now(),
    created_at       timestamptz  not null default now(),
    updated_at       timestamptz  not null default now()
);

CREATE TABLE verification_codes
(
    id         uuid primary key default gen_random_uuid(),
    user_id    uuid        not null REFERENCES users (id) on delete cascade,
    code       integer     not null,
    expires_at timestamptz not null
);

CREATE TABLE interests
(
    id             uuid primary key default gen_random_uuid(),
    title          varchar(255) not null,
    icon_file_name varchar(255) not null,
    created_at     timestamptz  not null,
    updated_at     timestamptz  not null
);

CREATE TABLE user_interests
(
    user_id     uuid not null references users (id) on delete cascade,
    interest_id uuid not null references interests (id) on delete cascade,
    primary key (user_id, interest_id)
);

CREATE TABLE chats
(
    id         uuid primary key     default gen_random_uuid(),
    title      varchar(255),
    c_type     chat_type   not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

CREATE TABLE messages
(
    id                   uuid primary key     default gen_random_uuid(),
    chat_id              uuid        not null references chats (id) on delete cascade,
    sender_id            uuid        not null references users (id) on delete cascade,
    raw_text             varchar(2048),
    edited               bool        not null,
    message_reference_id uuid        references messages (id) on delete set null,
    created_at           timestamptz not null default now(),
    updated_at           timestamptz not null default now()
);

CREATE TABLE read_status
(
    user_id    uuid        not null references users (id) on delete cascade,
    message_id uuid        not null references messages (id) on delete cascade,
    read_at    timestamptz not null default now(),

    primary key (user_id, message_id)
);

CREATE TABLE attachments
(
    id         uuid primary key      default gen_random_uuid(),
    message_id uuid         not null references messages (id) on delete cascade,
    filename   varchar(255) not null,
    created_at timestamptz  not null default now(),
    updated_at timestamptz  not null default now()
);

CREATE TABLE user_chats
(
    user_id            uuid not null references users (id) on delete cascade,
    chat_id            uuid not null references chats (id) on delete cascade,
    reveal_information bool not null default false,
    blocked            bool not null default false
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_interests;
DROP TABLE user_chats;
DROP TABLE verification_codes;
DROP TABLE interests;
DROP TABLE attachments;
DROP TABLE read_status;
DROP TABLE messages;
DROP TABLE chats;
DROP TABLE users;

DROP TYPE gender;
DROP TYPE chat_type;
-- +goose StatementEnd
