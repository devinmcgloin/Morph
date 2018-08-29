-- +migrate Up
DROP TYPE IF EXISTS STAT_TYPE;

CREATE TABLE content.image_stats (
    image_id integer NOT NULL,
    date date DEFAULT ('now'::text) ::date NOT NULL,
    stat_type SMALLINT NOT NULL,
    total integer DEFAULT 0 NOT NULL
);

CREATE UNIQUE INDEX image_stats_image_id_date_uindex ON content.image_stats (image_id, date);

-- +migrate Down
DROP TABLE content.image_stats;

