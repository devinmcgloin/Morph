-- +migrate Up
CREATE SCHEMA permissions;

CREATE TABLE permissions.can_delete (
    user_id integer NOT NULL,
    o_id integer NOT NULL,
    class SMALLINT NOT NULL
);

CREATE UNIQUE INDEX can_delete_user_id_o_id_uindex ON permissions.can_delete (user_id, o_id, class);

CREATE TABLE permissions.can_edit (
    user_id integer NOT NULL,
    o_id integer NOT NULL,
    class SMALLINT NOT NULL
);

CREATE UNIQUE INDEX can_edit_user_id_o_id_uindex ON permissions.can_edit (user_id, o_id, class);

CREATE TABLE permissions.can_view (
    user_id integer NOT NULL,
    o_id integer NOT NULL,
    class SMALLINT NOT NULL
);

CREATE UNIQUE INDEX can_view_user_id_o_id_type_uindex ON permissions.can_view (user_id, o_id, class);

-- +migrate Down
DROP TABLE permissions.can_view;

DROP TABLE permissions.can_delete;

DROP TABLE permissions.can_edit;

DROP SCHEMA permissions;
