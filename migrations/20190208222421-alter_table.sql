-- +migrate Up
CREATE SEQUENCE anime_id_seq;
ALTER TABLE anime ALTER COLUMN id SET DEFAULT nextval('anime_id_seq');
ALTER SEQUENCE anime_id_seq OWNED BY anime.id;
ALTER TABLE anime ADD COLUMN external_id VARCHAR (255);
-- +migrate Down