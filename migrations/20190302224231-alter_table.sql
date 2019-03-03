
-- +migrate Up
ALTER TABLE anime ADD COLUMN lastmodifytime TIMESTAMP;
-- +migrate Down
