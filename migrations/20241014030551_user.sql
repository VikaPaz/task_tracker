-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE table if not exists users (
    id uuid default uuid_generate_v4() primary key,
    login varchar(100),
    password varchar(100)
); 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
drop extension "uuid-ossp";
-- +goose StatementEnd
