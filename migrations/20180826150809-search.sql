-- +migrate Up
CREATE MATERIALIZED VIEW searches
AS
SELECT
    u.id AS searchable_id,
    'user' AS searchable_type,
    setweight(to_tsvector(coalesce(u.name, '')), 'A') || setweight(to_tsvector(coalesce(u.username, '')), 'A') || setweight(to_tsvector(coalesce(u.bio, '')), 'B') || setweight(to_tsvector(coalesce(u.location, '')), 'D') AS term
FROM
    content.users AS u
UNION
SELECT
    t.id AS searchable_id,
    'tag' AS searchable_type,
    to_tsvector(coalesce(t.description, '')) AS term
FROM
    content.tags AS t
UNION
SELECT
    i.id AS searchable_id,
    'image' AS searchable_type,
    to_tsvector(coalesce(meta.make, '')) || to_tsvector(coalesce(meta.model, '')) || to_tsvector(coalesce(meta.lens_make, '')) || to_tsvector(coalesce(meta.lens_model, '')) || to_tsvector(string_agg(coalesce(tags.description, ''), ' ')) || to_tsvector(string_agg(coalesce(landmarks.description, ''), ' ')) || to_tsvector(string_agg(coalesce(labels.description, ''), ' ')) AS term
FROM
    content.images AS i
    LEFT JOIN content.image_metadata AS meta ON i.id = meta.image_id
    LEFT JOIN content.image_tag_bridge AS tag_bridge ON i.id = tag_bridge.image_id
    LEFT JOIN content.tags AS tags ON tag_bridge.tag_id = tags.id
    LEFT JOIN content.image_landmark_bridge AS landmark_bridge ON i.id = landmark_bridge.image_id
    LEFT JOIN content.landmarks AS landmarks ON landmark_bridge.landmark_id = landmarks.id
    LEFT JOIN content.image_label_bridge AS label_bridge ON i.id = label_bridge.image_id
    LEFT JOIN content.labels AS labels ON label_bridge.label_id = labels.id
GROUP BY
    i.id,
    meta.make,
    meta.model,
    meta.lens_make,
    meta.lens_model;

-- +migrate Down
DROP MATERIALIZED VIEW searches;

