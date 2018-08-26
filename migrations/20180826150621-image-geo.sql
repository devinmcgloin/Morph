-- +migrate Up
CREATE TABLE content.image_geo (
    image_id integer NOT NULL CONSTRAINT image_geo_pkey PRIMARY KEY,
    loc geography (Point,
        4326),
    dir numeric,
    description text
);

CREATE UNIQUE INDEX image_geo_image_id_uindex ON content.image_geo (image_id);

-- +migrate Down
DROP TABLE content.image_geo;

