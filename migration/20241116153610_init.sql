-- +goose Up
-- +goose StatementBegin
CREATE TABLE users(
    id serial PRIMARY KEY,
    login varchar(255) not null UNIQUE,
    password varchar(512) not null,
    salt varchar(255) not null
);
CREATE TABLE data(
    id serial PRIMARY KEY,
    user_id int not null REFERENCES users,
    data bytea not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE data;
DROP TABLE users;
-- +goose StatementEnd
