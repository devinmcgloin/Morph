CREATE SEQUENCE images_id_seq;

CREATE SEQUENCE users_id_seq;

CREATE SEQUENCE image_tags_id_seq;

CREATE SEQUENCE image_labels_id_seq;

CREATE TABLE content.images
  (
     id            SERIAL NOT NULL CONSTRAINT images_pkey PRIMARY KEY,
     publish_time  TIMESTAMP DEFAULT Timezone('UTC'::text, Now()) NOT NULL,
     last_modified TIMESTAMP DEFAULT Timezone('UTC'::text, Now()) NOT NULL,
     owner         INTEGER NOT NULL,
     featured      BOOLEAN DEFAULT FALSE NOT NULL,
     downloads     INTEGER DEFAULT 0 NOT NULL,
     VIEWS         INTEGER DEFAULT 0 NOT NULL,
     location      POINT,
     shortcode     VARCHAR(12) NOT NULL
  );

CREATE TABLE content.users
  (
     id            SERIAL NOT NULL,
     username      VARCHAR(100) NOT NULL CONSTRAINT users_pkey PRIMARY KEY,
     name          TEXT,
     email         VARCHAR(100) NOT NULL,
     bio           TEXT,
     url           VARCHAR(50),
     password      TEXT NOT NULL,
     salt          TEXT NOT NULL,
     featured      BOOLEAN DEFAULT FALSE NOT NULL,
     ADMIN         BOOLEAN DEFAULT FALSE NOT NULL,
     VIEWS         INTEGER DEFAULT 0 NOT NULL,
     created_at    TIMESTAMP DEFAULT Timezone('UTC'::text, Now()) NOT NULL,
     last_modified TIMESTAMP DEFAULT Timezone('UTC'::text, Now()) NOT NULL
  );

CREATE TABLE content.image_metadata
  (
     image_id      INTEGER NOT NULL CONSTRAINT image_metadata_pkey PRIMARY KEY
     CONSTRAINT
     image_metadata_images_id_fk REFERENCES content.images,
     aperature     TEXT,
     exposure_time TEXT,
     focal_length  TEXT,
     iso           INTEGER,
     make          TEXT,
     model         TEXT,
     lens_make     INTEGER,
     lens_model    INTEGER,
     pixel_xd      INTEGER NOT NULL,
     pixel_yd      INTEGER NOT NULL,
     capture_time  TIMESTAMP,
     direction     DOUBLE PRECISION
  );

CREATE UNIQUE INDEX image_metadata_image_id_uindex
  ON image_metadata (image_id);

CREATE TABLE content.image_tags
  (
     id          SERIAL NOT NULL CONSTRAINT image_tags_pkey PRIMARY KEY,
     description TEXT NOT NULL
  );

CREATE UNIQUE INDEX image_tags_id_uindex
  ON image_tags (id);

CREATE UNIQUE INDEX image_tags_description_uindex
  ON image_tags (description);

CREATE TABLE content.image_labels
  (
     id          SERIAL NOT NULL CONSTRAINT image_labels_pkey PRIMARY KEY,
     description TEXT NOT NULL
  );

CREATE UNIQUE INDEX image_labels_id_uindex
  ON image_labels (id);

CREATE UNIQUE INDEX image_labels_description_uindex
  ON image_labels (description);

CREATE TABLE content.image_tag_bridge
  (
     image_id INTEGER NOT NULL CONSTRAINT image_tag_bridge_images_id_fk
     REFERENCES
     content.images,
     tag_id   INTEGER CONSTRAINT image_tag_bridge_image_tags_id_fk REFERENCES
     content.image_tags
  );

CREATE UNIQUE INDEX image_tag_bridge_image_id_tag_id_uindex
  ON image_tag_bridge (image_id, tag_id);

CREATE TABLE content.image_label_bridge
  (
     image_id INTEGER NOT NULL CONSTRAINT image_label_bridge_images_id_fk
     REFERENCES
     content.images,
     label_id INTEGER CONSTRAINT image_label_bridge_image_tags_id_fk REFERENCES
     content.image_labels
  );

CREATE UNIQUE INDEX image_label_bridge_image_id_tag_id_uindex
  ON image_label_bridge (image_id, label_id);

CREATE TABLE permissions.can_edit
  (
     user_id INTEGER NOT NULL,
     o_id    INTEGER NOT NULL,
     TYPE    CONTENT_TYPE NOT NULL
  );

CREATE TABLE permissions.can_delete
  (
     user_id INTEGER NOT NULL,
     o_id    INTEGER NOT NULL,
     TYPE    CONTENT_TYPE NOT NULL
  );

CREATE TABLE permissions.can_view
  (
     user_id INTEGER NOT NULL,
     o_id    INTEGER NOT NULL,
     TYPE    CONTENT_TYPE NOT NULL
  );
