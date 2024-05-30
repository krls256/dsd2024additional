-- +goose Up
-- +goose StatementBegin
create table if not exists accounts
(
    created_at timestamptz default now(),
    updated_at timestamptz,


    id         uuid primary key,
    name       varchar(40),
    login      varchar,
    password   varchar
);

create table if not exists tokens
(
    created_at    timestamptz default now(),
    id            uuid        default gen_random_uuid() primary key,
    account_id    uuid references accounts (id) on delete cascade,
    access_token  text        not null,
    refresh_token text        not null,
    expires_at    timestamptz not null
);

create index if not exists idx_account_id on tokens using btree (account_id);
create index if not exists idx_access_token on tokens using hash (access_token);
create index if not exists idx_refresh_token on tokens using hash (refresh_token);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index idx_refresh_token;
drop index idx_access_token;
drop index idx_account_id;

drop table if exists tokens;
drop table if exists accounts;
-- +goose StatementEnd
