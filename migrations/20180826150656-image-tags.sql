-- +migrate Up
CREATE TABLE content.tags (
    id serial NOT NULL CONSTRAINT tags_pkey PRIMARY KEY,
    description text NOT NULL
);

CREATE UNIQUE INDEX tags_id_uindex ON content.tags (id);

CREATE UNIQUE INDEX tags_description_uindex ON content.tags (description);

-- +migrate Down
DROP TABLE content.tags;

