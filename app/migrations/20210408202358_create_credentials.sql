-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table credentials (
  id bigserial primary key,
  app varchar(255) not null,
  username varchar(225) not null,
  password varchar(225) not null,
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz
);

create index credentials_app_idx ON credentials(LOWER(app));

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop index if exists credentials_app_idx;

drop table if exists credentials;
