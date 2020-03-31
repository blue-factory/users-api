-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  id uuid primary key default gen_random_uuid() UNIQUE,
  name varchar(255) NOT NULL,
  email varchar(255)  NOT NULL UNIQUE,
  password varchar(255) NOT NULL,
  verified BOOLEAN NOT NULL DEFAULT FALSE,
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

create trigger update_users_update_at
before update on users for each row execute procedure update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
