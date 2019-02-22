
-- +migrate Up
ALTER TABLE anime ADD COLUMN score REAL;
ALTER TABLE anime ADD COLUMN duration REAL;
ALTER TABLE anime ADD COLUMN rating REAL;
ALTER TABLE anime ADD COLUMN franchase VARCHAR(255);
-- +migrate Down
