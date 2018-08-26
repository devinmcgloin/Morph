-- +migrate Up

CREATE SCHEMA content;

CREATE TABLE content.images (
  id serial NOT NULL CONSTRAINT images_pkey PRIMARY KEY,
  user_id integer NOT NULL,
  featured boolean DEFAULT FALSE NOT NULL,
  shortcode varchar(12) NOT NULL,
  views integer DEFAULT 0,
  favorites integer DEFAULT 0,
  title text,
  description text,
  publish_time timestamp WITH time zone DEFAULT timezone('UTC'::text, now()) NOT NULL, 
  last_modified timestamp WITH time zone DEFAULT timezone('UTC'::text, now()) NOT NULL
);

-- +migrate Down

DROP TABLE content.images;
DROP SCHEMA content;

