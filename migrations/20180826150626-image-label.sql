-- +migrate Up
CREATE TABLE content.labels (
    id serial NOT NULL CONSTRAINT image_labels_pkey PRIMARY KEY,
    description text NOT NULL
);

CREATE UNIQUE INDEX image_labels_id_uindex ON content.labels (id);

CREATE UNIQUE INDEX image_labels_description_uindex ON content.labels (description);

-- +migrate Down
DROP TABLE content.labels;

