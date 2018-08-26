-- +migrate Up
CREATE TABLE content.streams (
    id serial NOT NULL CONSTRAINT streams_pkey PRIMARY KEY,
    title varchar(50),
    description varchar(100),
    user_id integer NOT NULL CONSTRAINT user_stream_id_fk REFERENCES content.users (id),
    created_at timestamp WITH time zone DEFAULT timezone('UTC'::text, now()) NOT NULL, updated_at timestamp WITH time zone DEFAULT timezone('UTC'::text, now()) NOT NULL);

CREATE UNIQUE INDEX streams_uindex ON content.streams (id);

CREATE UNIQUE INDEX user_uindex ON content.streams (user_id);

-- +migrate Down
DROP TABLE content.streams;

