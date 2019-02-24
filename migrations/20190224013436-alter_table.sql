
-- +migrate Up
ALTER TABLE anime ALTER COLUMN rating TYPE VARCHAR(255);
-- +migrate Down
