-- +migrate Up
CREATE TABLE content.image_stream_bridge (
    image_id integer NOT NULL CONSTRAINT image_stream_bridge_images_id_fk REFERENCES content.images,
    stream_id integer NOT NULL CONSTRAINT image_stream_bridge_image_tags_id_fk REFERENCES content.streams,
    created_at timestamp WITH time zone DEFAULT timezone('UTC'::text, now())
);

CREATE UNIQUE INDEX image_stream_bridge_image_id_stream_id_uindex ON content.image_stream_bridge (image_id, stream_id);

-- +migrate Down
DROP TABLE content.image_stream_bridge;

