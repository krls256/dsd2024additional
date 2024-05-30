-- +goose Up
-- +goose StatementBegin
create table if not exists profiles
(
    created_at timestamptz default now(),
    updated_at timestamptz,


    id         uuid primary key,
    nickname   varchar,
    about_me   text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
