
-- +migrate Up
CREATE TABLE new (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    body TEXT NOT NULL
);
-- +migrate Down
