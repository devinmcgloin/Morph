-- +migrate Up
CREATE SCHEMA colors;

CREATE TYPE COLOR_TYPE AS ENUM ( 'shade',
    'specific'
);

CREATE TABLE colors.clr (
    name text NOT NULL,
    hex text NOT NULL,
    TYPE color_type NOT NULL
);

CREATE UNIQUE INDEX clr_name_hex_type_uindex ON colors.clr (name, hex, TYPE);

-- +migrate Down
DROP TABLE colors.clr;

DROP TYPE COLOR_TYPE;

DROP SCHEMA colors;

