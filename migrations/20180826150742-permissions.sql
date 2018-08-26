-- +migrate Up
CREATE SCHEMA permissions;

CREATE TYPE CONTENT_TYPE AS ENUM ( 'user',
    'image'
);

CREATE TABLE permissions.can_delete (
    user_id integer NOT NULL,
    o_id integer NOT NULL,
    TYPE content_type NOT NULL
);

CREATE UNIQUE INDEX can_delete_user_id_o_id_uindex ON permissions.can_delete (user_id, o_id);

CREATE TABLE permissions.can_edit (
    user_id integer NOT NULL,
    o_id integer NOT NULL,
    TYPE content_type NOT NULL
);

CREATE UNIQUE INDEX can_edit_user_id_o_id_uindex ON permissions.can_edit (user_id, o_id);

CREATE TABLE permissions.can_view (
    user_id integer NOT NULL,
    o_id integer NOT NULL,
    TYPE content_type NOT NULL
);

CREATE UNIQUE INDEX can_view_user_id_o_id_type_uindex ON permissions.can_view (user_id, o_id, TYPE);

-- +migrate Down
DROP TABLE permissions.can_view;

DROP TABLE permissions.can_delete;

DROP TABLE permissions.can_edit;

DROP SCHEMA permissions;

DROP TYPE CONTENT_TYPE;

