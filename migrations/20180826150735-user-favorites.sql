-- +migrate Up
CREATE TABLE content.user_favorites (
  user_id integer NOT NULL CONSTRAINT user_follow_users_id_fk REFERENCES content.users (id),
  image_id integer NOT NULL CONSTRAINT image_favorites_pkey PRIMARY KEY CONSTRAINT user_followed_users_id_fk REFERENCES content.images (id),
  created_at timestamp WITH time zone DEFAULT timezone('UTC'::text, now()));

CREATE UNIQUE INDEX user_favorites_id_uindex ON content.user_favorites (user_id, image_id);

-- +migrate Down
DROP TABLE content.user_favorites;

