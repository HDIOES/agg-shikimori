
-- +migrate Up
ALTER TABLE ANIME ADD COLUMN poster_url VARCHAR(1024);
-- +migrate Down
