-- +migrate Up
CREATE TABLE content.landmarks (
    id serial NOT NULL,
    description text,
    LOCATION geography (Point,
        4326) NOT NULL
);

CREATE UNIQUE INDEX landmarks_id_uindex ON content.landmarks (id);

-- +migrate Down
DROP TABLE content.landmarks;

