-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table requests (
  id bigserial primary key,
  app varchar(255) not null,
  multi boolean not null,
  number varchar(225) not null,
  message text not null,
  times integer not null,
  created_at timestamptz not null default clock_timestamp()
);

create index requests_app_idx ON requests(LOWER(app));

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop index if exists requests_app_idx;

drop table if exists requests;
