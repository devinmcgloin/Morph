package main

const createScript = `
create EXTENSION postgis;

CREATE TYPE STAT_TYPE AS ENUM ('view', 'download');
CREATE TYPE COLOR_TYPE AS ENUM ('shade', 'specific');
CREATE TYPE CONTENT_TYPE AS ENUM ('user', 'image');

--- colors
create SCHEMA colors;
create table colors.clr
(
  name text not null,
  hex text not null,
  type color_type not null
)
;

create unique index clr_name_hex_type_uindex
  on colors.clr (name, hex, type)
;

-- Content

create schema content;
create table content.colors
(
  id serial not null
    constraint color_pkey
    primary key,
  red integer not null,
  green integer not null,
  blue integer not null,
  hue integer not null,
  saturation integer not null,
  val integer not null,
  color text not null,
  shade text not null,
  cielab cube
)
;

create unique index color_id_uindex
  on content.colors (id)
;

create index colors_cielab_index
  on content.colors (cielab)
;

create table content.image_color_bridge
(
  image_id integer not null
    constraint image_color_bridge_images_id_fk
    references content.images,
  color_id integer not null
    constraint image_color_bridge_color_id_fk
    references content.colors,
  pixel_fraction double precision not null,
  score double precision not null
)
;

create table content.image_geo
(
  image_id integer not null
    constraint image_geo_pkey
    primary key,
  loc geography(Point,4326),
  dir numeric,
  description text
)
;

create unique index image_geo_image_id_uindex
  on content.image_geo (image_id)
;


create table content.image_label_bridge
(
  image_id integer not null
    constraint image_label_bridge_images_id_fk
    references content.images,
  label_id integer
    constraint image_label_bridge_image_tags_id_fk
    references content.labels,
  score double precision
)
;

create unique index image_label_bridge_image_id_tag_id_uindex
  on content.image_label_bridge (image_id, label_id)
;

create table content.image_landmark_bridge
(
  image_id integer not null
    constraint image_landmark_bridge_images_id_fk
    references content.images,
  landmark_id integer not null
    constraint image_landmark_bridge_landmarks_id_fk
    references content.landmarks (id),
  score double precision not null
)
;

create unique index image_landmark_bridge_image_id_landmark_id_uindex
  on content.image_landmark_bridge (image_id, landmark_id)
;

create table content.image_metadata
(
  image_id integer not null
    constraint image_metadata_pkey
    primary key
    constraint image_metadata_images_id_fk
    references content.images,
  aperture double precision,
  exposure_time text,
  focal_length integer,
  iso integer,
  make text,
  model text,
  lens_make text,
  lens_model text,
  pixel_xd integer not null,
  pixel_yd integer not null,
  capture_time timestamp
)
;

create unique index image_metadata_image_id_uindex
  on content.image_metadata (image_id)
;


create table content.image_stats
(
  image_id integer not null,
  date date default ('now'::text)::date not null,
  stat_type stat_type not null,
  total integer default 0 not null
)
;

create unique index image_stats_image_id_date_uindex
  on content.image_stats (image_id, date)
;

create table content.image_tag_bridge
(
  image_id integer not null
    constraint image_tag_bridge_images_id_fk
    references content.images,
  tag_id integer
    constraint image_tag_bridge_image_tags_id_fk
    references content.image_tags
)
;

create unique index image_tag_bridge_image_id_tag_id_uindex
  on content.image_tag_bridge (image_id, tag_id)
;


create table content.image_tags
(
  id serial not null
    constraint image_tags_pkey
    primary key,
  description text not null
)
;

create unique index image_tags_id_uindex
  on content.image_tags (id)
;

create unique index image_tags_description_uindex
  on content.image_tags (description)
;

create table content.images
(
  id serial not null
    constraint images_pkey
    primary key,
  publish_time timestamp with time zone default timezone('UTC'::text, now()) not null,
  last_modified timestamp with time zone default timezone('UTC'::text, now()) not null,
  user_id integer not null,
  featured boolean default false not null,
  shortcode varchar(12) not null,
  views integer default 0,
  favorites integer default 0,
  title text,
  description text
)
;

create index index_images_on_ranking
  on content.images (ranking(id, views + favorites, featured::integer + 3))
;

create table content.labels
(
  id serial not null
    constraint image_labels_pkey
    primary key,
  description text not null
)
;

create unique index image_labels_id_uindex
  on content.labels (id)
;

create unique index image_labels_description_uindex
  on content.labels (description)
;

create table content.landmarks
(
  id serial not null,
  description text,
  location geography(Point,4326) not null
)
;

create unique index landmarks_id_uindex
  on content.landmarks (id)
;

create table content.user_favorites
(
  user_id integer,
  image_id integer,
  created_at timestamp with time zone default timezone('UTC'::text, now())
)
;
create table content.user_follows
(
  user_id integer not null
    constraint user_follow_users_id_fk
    references users (id),
  followed_id integer not null
    constraint user_follows_pkey
    primary key
    constraint user_followed_users_id_fk
    references users (id),
  created_at timestamp with time zone default timezone('UTC'::text, now())
)
;

create unique index user_follows_user_id_uindex
  on content.user_follows (user_id, followed_id)
;


create table content.users
(
  id serial not null,
  username varchar(100) not null
    constraint users_pkey
    primary key,
  name text,
  email varchar(100) not null,
  bio text,
  url varchar(50),
  twitter text,
  instagram text,
  featured boolean default false not null,
  admin boolean default false not null,
  created_at timestamp with time zone default timezone('UTC'::text, now()) not null,
  last_modified timestamp with time zone default timezone('UTC'::text, now()) not null,
  location text
)
;

create unique index users_id_uindex
  on content.users (id)
;

--- permissions
create table permissions.can_delete
(
  user_id integer not null,
  o_id integer not null,
  type content_type not null
)
;

create unique index can_delete_user_id_o_id_uindex
  on permissions.can_delete (user_id, o_id)
;

create table permissions.can_edit
(
  user_id integer not null,
  o_id integer not null,
  type content_type not null
)
;

create unique index can_edit_user_id_o_id_uindex
  on permissions.can_edit (user_id, o_id)
;

create table permissions.can_view
(
  user_id integer not null,
  o_id integer not null,
  type content_type not null
)
;

create unique index can_view_user_id_o_id_type_uindex
  on permissions.can_view (user_id, o_id, type)
;




CREATE OR REPLACE FUNCTION popularity(count INTEGER, weight INTEGER DEFAULT 3)
  RETURNS INTEGER AS $$
SELECT count * weight
$$ LANGUAGE SQL IMMUTABLE;

CREATE OR REPLACE FUNCTION ranking(id INTEGER, counts INTEGER, weight INTEGER)
  RETURNS INTEGER AS $$
SELECT id + popularity(counts, weight)
$$ LANGUAGE SQL IMMUTABLE;

CREATE INDEX index_images_on_ranking
  ON content.images (ranking(id, views + favorites, featured :: INT + 3) DESC);


-- Stats Triggers

CREATE OR REPLACE FUNCTION flush_all_stats()
  RETURNS TRIGGER AS
$BODY$
BEGIN
  IF NEW.total <> OLD.total
  THEN
    UPDATE content.images
    SET views   = (SELECT sum(total)
                   FROM content.image_stats
                   WHERE stat_type = 'view' AND image_id = NEW.image_id),
      favorites = (SELECT count(*)
                   FROM content.user_favorites
                   WHERE image_id = NEW.image_id)
    WHERE images.id = OLD.image_id;
  END IF;

  RETURN NEW;
END;
$BODY$

LANGUAGE plpgsql;

CREATE TRIGGER flush_stats
AFTER UPDATE ON content.image_stats
FOR EACH ROW
EXECUTE PROCEDURE flush_all_stats();

CREATE OR REPLACE FUNCTION log_stat(id INT, t STAT_TYPE)
  RETURNS BOOLEAN
LANGUAGE plpgsql AS
$BODY$
BEGIN
  IF (SELECT count(*)
      FROM CONTENT.image_stats
      WHERE date = CURRENT_DATE
            AND stat_type = t
            AND image_id = id) > 0
  THEN
    UPDATE content.image_stats
    SET total = total + 1
    WHERE stat_type = t AND image_id = id AND date = current_date;
  ELSE
    INSERT INTO content.image_stats VALUES (id, current_date, t, 1);
  END IF;
  RETURN TRUE;
END;
$BODY$;


--- Random

CREATE OR REPLACE FUNCTION random_image(u INTEGER DEFAULT -1)
  RETURNS INT
LANGUAGE plpgsql AS
$BODY$
BEGIN
  IF u = -1
  THEN
    RETURN (SELECT id
            FROM content.images
              INNER JOIN permissions.can_view ON o_id = id AND type = 'image'
            WHERE permissions.can_view.user_id = -1
            OFFSET random() * (SELECT count(id)
                               FROM content.images
                                 INNER JOIN permissions.can_view ON o_id = id AND type = 'image')
            LIMIT 1);
  ELSE
    RETURN (SELECT id
            FROM content.images
              INNER JOIN permissions.can_view ON o_id = id AND type = 'image'
            WHERE images.user_id = u
            OFFSET random() * (SELECT count(id)
                               FROM content.images
                                 INNER JOIN permissions.can_view ON o_id = id AND type = 'image'
                               WHERE images.user_id = u)
            LIMIT 1);
  END IF;
END;
$BODY$;


-- Stream

CREATE OR REPLACE FUNCTION stream(u INTEGER)
  RETURNS TABLE(user_id INT, o_id INT, t CONTENT_TYPE, created_at TIMESTAMP WITH TIME ZONE)
LANGUAGE plpgsql AS
$BODY$
BEGIN
  SELECT
    a.user_id,
    a.o_id,
    a.t,
    a.relation,
    a.created_at
  FROM (SELECT
          follows.user_id,
          follows.followed_id    AS o_id,
          'user' :: CONTENT_TYPE AS t,
          'follows'              AS relation,
          follows.created_at
        FROM CONTENT.user_follows AS follows
        WHERE user_id IN (SELECT followed_id
                          FROM CONTENT.user_follows AS F
                          WHERE F.user_id = 10)
        UNION
        SELECT
          favs.user_id,
          favs.image_id           AS o_id,
          'image' :: CONTENT_TYPE AS t,
          'favorites'             AS relation,

          favs.created_at
        FROM CONTENT.user_favorites AS favs
        WHERE favs.user_id IN (SELECT followed_id
                               FROM CONTENT.user_follows AS F
                               WHERE F.user_id = 10)
        UNION
        SELECT
          I.user_id,
          I.id                    AS o_id,
          'image' :: CONTENT_TYPE AS t,
          'published'             AS relation,

          I.publish_time          AS created_at
        FROM CONTENT.images AS I
        WHERE I.user_id IN (SELECT followed_id
                            FROM CONTENT.user_follows AS F
                            WHERE F.user_id = 10)) AS a;
END;
$BODY$;


-- Text Search View

CREATE MATERIALIZED VIEW searches AS
  SELECT
    u.id                                                  AS searchable_id,
    'user'                                                AS searchable_type,
    setweight(to_tsvector(coalesce(u.name, '')), 'A') ||
    setweight(to_tsvector(coalesce(u.username, '')), 'A') ||
    setweight(to_tsvector(coalesce(u.bio, '')), 'B') ||
    setweight(to_tsvector(coalesce(u.location, '')), 'D') AS term
  FROM content.users AS u

  UNION

  SELECT
    t.id                                     AS searchable_id,
    'tag'                                    AS searchable_type,
    to_tsvector(coalesce(t.description, '')) AS term
  FROM content.image_tags AS t

  UNION

  SELECT
    i.id    AS searchable_id,
    'image' AS searchable_type,
    to_tsvector(coalesce(meta.make, '')) ||
    to_tsvector(coalesce(meta.model, '')) ||
    to_tsvector(coalesce(meta.lens_make, '')) ||
    to_tsvector(coalesce(meta.lens_model, '')) ||
    to_tsvector(string_agg(coalesce(tags.description, ''), ' ')) ||
    to_tsvector(string_agg(coalesce(landmarks.description, ''), ' ')) ||
    to_tsvector(string_agg(coalesce(labels.description,''), ' '))
            AS term
  FROM content.images AS i
    LEFT JOIN content.image_metadata AS meta ON i.id = meta.image_id
    LEFT JOIN content.image_tag_bridge AS tag_bridge ON i.id = tag_bridge.image_id
    LEFT JOIN content.image_tags AS tags ON tag_bridge.tag_id = tags.id
    LEFT JOIN content.image_landmark_bridge AS landmark_bridge ON i.id = landmark_bridge.image_id
    LEFT JOIN content.landmarks AS landmarks ON landmark_bridge.landmark_id = landmarks.id
    LEFT JOIN content.image_label_bridge AS label_bridge ON i.id = label_bridge.image_id
    LEFT JOIN content.labels AS labels ON label_bridge.label_id = labels.id
  GROUP BY i.id, meta.make, meta.model, meta.lens_make, meta.lens_model;

-- REFRESH MATERIALIZED VIEW CONCURRENTLY searches;
`
