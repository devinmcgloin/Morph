
-- +migrate Up
CREATE UNIQUE INDEX users_searchable_id_searchable_type_uindex
  ON searches (searchable_id, searchable_type);
-- +migrate Down
DROP UNIQUE INDEX users_searchable_id_searchable_type_uindex;