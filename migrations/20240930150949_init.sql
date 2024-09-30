-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


create table if not exists tasks
(
    id         uuid default uuid_generate_v4() primary key,
    title    text,
    description text,
    created_at timestamp,
    updated_at timestamp,
    status     varchar(100),
    owner_id uuid 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table tasks;
drop extension "uuid-ossp";
-- +goose StatementEnd
