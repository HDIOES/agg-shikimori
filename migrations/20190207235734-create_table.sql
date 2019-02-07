
-- +migrate Up
CREATE TABLE ANIME (
    ID BIGINT PRIMARY KEY,
    NAME VARCHAR (255) NOT NULL
);
-- +migrate Down
