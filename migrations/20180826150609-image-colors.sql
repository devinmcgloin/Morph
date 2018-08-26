-- +migrate Up
CREATE TABLE content.colors (
    id serial NOT NULL CONSTRAINT color_pkey PRIMARY KEY,
    red integer NOT NULL,
    green integer NOT NULL,
    blue integer NOT NULL,
    hue integer NOT NULL,
    saturation integer NOT NULL,
    val integer NOT NULL,
    color text NOT NULL,
    shade text NOT NULL,
    cielab CUBE
);

CREATE UNIQUE INDEX color_id_uindex ON content.colors (id);

CREATE INDEX colors_cielab_index ON content.colors (cielab);

-- +migrate Down
DROP TABLE content.colors;

