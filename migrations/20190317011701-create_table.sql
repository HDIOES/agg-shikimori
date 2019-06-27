
-- +migrate Up
CREATE TABLE ngramm (
    ngramm_value varchar(3),
    anime_id BIGINT REFERENCES anime(id),
    PRIMARY KEY(ngramm_value, anime_id)
);
CREATE INDEX idx_ngramm_value ON ngramm(ngramm_value);
-- +migrate Down
