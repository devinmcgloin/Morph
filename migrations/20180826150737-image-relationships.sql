-- +migrate Up
CREATE TABLE content.image_color_bridge (
  image_id integer NOT NULL CONSTRAINT image_color_bridge_images_id_fk REFERENCES content.images,
  color_id integer NOT NULL CONSTRAINT image_color_bridge_color_id_fk REFERENCES content.colors,
  pixel_fraction double precision NOT NULL,
  score double precision NOT NULL,
  created_at timestamp WITH time zone DEFAULT timezone('UTC'::text, now())
);

CREATE TABLE content.image_label_bridge (
  image_id integer NOT NULL CONSTRAINT image_label_bridge_images_id_fk REFERENCES content.images,
  label_id integer CONSTRAINT image_label_bridge_image_tags_id_fk REFERENCES content.labels,
  score double precision,
  created_at timestamp WITH time zone DEFAULT timezone('UTC'::text, now())
);

CREATE UNIQUE INDEX image_label_bridge_image_id_tag_id_uindex ON content.image_label_bridge (image_id, label_id);

CREATE TABLE content.image_landmark_bridge (
  image_id integer NOT NULL CONSTRAINT image_landmark_bridge_images_id_fk REFERENCES content.images,
  landmark_id integer NOT NULL CONSTRAINT image_landmark_bridge_landmarks_id_fk REFERENCES content.landmarks (id),
  score double precision NOT NULL,
  created_at timestamp WITH time zone DEFAULT timezone('UTC'::text, now())
);

CREATE UNIQUE INDEX image_landmark_bridge_image_id_landmark_id_uindex ON content.image_landmark_bridge (image_id, landmark_id);

CREATE TABLE content.image_tag_bridge (
  image_id integer NOT NULL CONSTRAINT image_tag_bridge_images_id_fk REFERENCES content.images,
  tag_id integer CONSTRAINT tag_bridge_image_tags_id_fk REFERENCES content.tags,
  created_at timestamp WITH time zone DEFAULT timezone('UTC'::text, now())
);

CREATE UNIQUE INDEX image_tag_bridge_image_id_tag_id_uindex ON content.image_tag_bridge (image_id, tag_id);

-- +migrate Down
DROP TABLE content.image_landmark_bridge;

DROP TABLE content.image_color_bridge;

DROP TABLE content.image_label_bridge;

DROP TABLE content.image_tag_bridge;

