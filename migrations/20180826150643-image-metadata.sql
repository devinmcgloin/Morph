-- +migrate Up
CREATE TABLE content.image_metadata (
    image_id integer NOT NULL CONSTRAINT image_metadata_pkey PRIMARY KEY CONSTRAINT image_metadata_images_id_fk REFERENCES content.images,
    aperture double precision,
    exposure_time text,
    focal_length integer,
    iso integer,
    make text,
    model text,
    lens_make text,
    lens_model text,
    pixel_xd integer NOT NULL,
    pixel_yd integer NOT NULL,
    capture_time timestamp
);

CREATE UNIQUE INDEX image_metadata_image_id_uindex ON content.image_metadata (image_id);

-- +migrate Down
DROP TABLE content.image_metadata;

