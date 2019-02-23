
-- +migrate Up
ALTER TABLE anime ADD COLUMN processed BOOLEAN DEFAULT false;
-- +migrate Down
