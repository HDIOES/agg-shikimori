
-- +migrate Up
ALTER TABLE anime ADD COLUMN russian VARCHAR (255);
ALTER TABLE anime ADD COLUMN amine_url VARCHAR (255);
ALTER TABLE anime ADD COLUMN kind VARCHAR (255);
ALTER TABLE anime ADD COLUMN anime_status VARCHAR (255);
ALTER TABLE anime ADD COLUMN epizodes INT;
ALTER TABLE anime ADD COLUMN epizodes_aired INT;
ALTER TABLE anime ADD COLUMN aired_on DATE;
ALTER TABLE anime ADD COLUMN released_on DATE;
-- +migrate Down
