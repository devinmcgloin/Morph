CREATE FUNCTION popularity(count integer, weight integer default 3) RETURNS integer AS $$
  SELECT count * weight
$$ LANGUAGE SQL IMMUTABLE;

CREATE FUNCTION ranking(id integer, counts integer, weight integer) RETURNS integer AS $$
  SELECT id + popularity(counts, weight)
$$ LANGUAGE SQL IMMUTABLE;

CREATE INDEX index_images_on_ranking
  ON content.images (ranking(id, (SELECT count(*) FROM content.user_favorites
	WHERE image_id = id) + (SELECT count(*) FROM content.image_stats
	WHERE image_id = id AND type = 'view'), 3) DESC);


ALTER TABLE content.image_stats RENAME COLUMN timestamp TO date;
ALTER TABLE content.image_stats ALTER COLUMN date TYPE DATE USING date::DATE;
ALTER TABLE content.image_stats ALTER COLUMN date SET DEFAULT current_date;
CREATE UNIQUE INDEX image_stats_image_id_date_uindex ON content.image_stats (image_id, date);
ALTER TABLE content.image_stats RENAME COLUMN type TO stat_type;
ALTER TABLE content.image_stats RENAME COLUMN count TO total;

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


CREATE OR REPLACE FUNCTION log_stat(id int, t STAT_TYPE)
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
  return true;
END;
$BODY$;

  select log_stat(3, 'view');

update content.colors SET cielab = cube(ARRAY[5.4, 32.3, 3.4]) WHERE id = 1;

SELECT color, shade, cielab from content.colors ORDER BY cielab;

INSERT INTO content.colors (red, green, blue, hue, saturation, val, shade, color, cielab)
			VALUES (34, 100, 145, 54, 23, 12, 'blue', 'tet blue', '(57.918510, 0.904568, -11.776117)'::cube);

SELECT bridge.image_id, color, shade, cielab
FROM content.colors AS colors
  INNER JOIN content.image_color_bridge AS bridge ON colors.id = bridge.color_id
ORDER BY '(56.21196714265635, -2.9428863271913075, 5.295922123682706)'::cube <-> cielab;

(SELECT bridge.image_id
 FROM content.landmarks AS landmark
   JOIN content.image_landmark_bridge AS bridge ON landmark.id = bridge.landmark_id
 WHERE to_tsvector(landmark.description) @@ to_tsquery('city'))
UNION
(SELECT bridge.image_id
 FROM content.labels AS labels
   JOIN content.image_label_bridge AS bridge ON labels.id = bridge.label_id
 WHERE to_tsvector(labels.description) @@ to_tsquery('city'))
UNION
(SELECT bridge.image_id
 FROM content.image_tags AS tags
   JOIN content.image_tag_bridge AS bridge ON tags.id = bridge.tag_id
 WHERE to_tsvector(tags.description) @@ to_tsquery('city'))