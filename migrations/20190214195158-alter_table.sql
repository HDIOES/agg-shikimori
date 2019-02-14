
-- +migrate Up
ALTER TABLE anime ADD CONSTRAINT external_id_unique UNIQUE (external_id);
-- +migrate Down
