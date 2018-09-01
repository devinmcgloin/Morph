-- +migrate Up
CREATE TABLE content.users (
    id serial NOT NULL,
    username varchar(100) NOT NULL CONSTRAINT users_pkey PRIMARY KEY,
    name text,
    email varchar(100) NOT NULL,
    bio text,
    url varchar(50),
    twitter text,
    instagram text,
    avatar_id UUID,
    featured boolean DEFAULT FALSE NOT NULL,
    admin boolean DEFAULT FALSE NOT NULL,
    created_at timestamp WITH time zone DEFAULT timezone('UTC'::text, now()) NOT NULL, last_modified timestamp WITH time zone DEFAULT timezone('UTC'::text, now()) NOT NULL, LOCATION text);

CREATE UNIQUE INDEX users_id_uindex ON content.users (id);

-- +migrate Down
DROP TABLE content.users;

