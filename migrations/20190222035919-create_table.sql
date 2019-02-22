
-- +migrate Up
CREATE TABLE studio (
    id SERIAL PRIMARY KEY,
    external_id BIGINT,
    studio_name VARCHAR(255),
    filtered_studio_name VARCHAR(255),
    is_real BOOLEAN,
    image_url VARCHAR(1024)
);
CREATE TABLE anime_studio (
    anime_id BIGINT REFERENCES anime(id),
    studio_id BIGINT REFERENCES studio(id),
    PRIMARY KEY(anime_id, studio_id)
);
CREATE TABLE genre (
    id SERIAL PRIMARY KEY,
    external_id BIGINT,
    genre_name VARCHAR(255),
    russian VARCHAR(255),
    kind VARCHAR(255)
);
CREATE TABLE anime_genre (
    anime_id BIGINT REFERENCES anime(id),
    genre_id BIGINT REFERENCES genre(id),
    PRIMARY KEY(anime_id, genre_id)
);
-- +migrate Down
