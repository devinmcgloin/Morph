-- +migrate Up
CREATE TABLE content.user_follows_user (
    user_id integer NOT NULL CONSTRAINT user_follow_users_id_fk REFERENCES content.users (id),
    followed_id integer NOT NULL CONSTRAINT user_follows_pkey PRIMARY KEY CONSTRAINT user_followed_users_id_fk REFERENCES content.users (id),
    created_at timestamp WITH time zone DEFAULT timezone('UTC'::text, now()));

CREATE UNIQUE INDEX user_follows_user_id_uindex ON content.user_follows_user (user_id, followed_id);

-- +migrate Down
DROP TABLE content.user_follows_user;

