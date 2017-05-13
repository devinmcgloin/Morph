CREATE TABLE image_label_bridge
(
    image_id INTEGER NOT NULL,
    label_id INTEGER,
    CONSTRAINT image_label_bridge_images_id_fk FOREIGN KEY (image_id) REFERENCES images (id),
    CONSTRAINT image_label_bridge_image_tags_id_fk FOREIGN KEY (label_id) REFERENCES image_labels (id)
);
CREATE UNIQUE INDEX image_label_bridge_image_id_tag_id_uindex ON image_label_bridge (image_id, label_id);
CREATE TABLE image_labels
(
    id INTEGER DEFAULT nextval('content.image_labels_id_seq'::regclass) PRIMARY KEY NOT NULL,
    description TEXT NOT NULL
);
CREATE UNIQUE INDEX image_labels_id_uindex ON image_labels (id);
CREATE UNIQUE INDEX image_labels_description_uindex ON image_labels (description);
CREATE TABLE image_metadata
(
    image_id INTEGER PRIMARY KEY NOT NULL,
    aperture TEXT,
    exposure_time TEXT,
    focal_length TEXT,
    iso INTEGER,
    make TEXT,
    model TEXT,
    lens_make TEXT,
    lens_model TEXT,
    pixel_xd INTEGER NOT NULL,
    pixel_yd INTEGER NOT NULL,
    capture_time TIMESTAMP,
    CONSTRAINT image_metadata_images_id_fk FOREIGN KEY (image_id) REFERENCES images (id)
);
CREATE UNIQUE INDEX image_metadata_image_id_uindex ON image_metadata (image_id);
CREATE TABLE image_tag_bridge
(
    image_id INTEGER NOT NULL,
    tag_id INTEGER,
    CONSTRAINT image_tag_bridge_images_id_fk FOREIGN KEY (image_id) REFERENCES images (id),
    CONSTRAINT image_tag_bridge_image_tags_id_fk FOREIGN KEY (tag_id) REFERENCES image_tags (id)
);
CREATE UNIQUE INDEX image_tag_bridge_image_id_tag_id_uindex ON image_tag_bridge (image_id, tag_id);
CREATE TABLE image_tags
(
    id INTEGER DEFAULT nextval('content.image_tags_id_seq'::regclass) PRIMARY KEY NOT NULL,
    description TEXT NOT NULL
);
CREATE UNIQUE INDEX image_tags_id_uindex ON image_tags (id);
CREATE UNIQUE INDEX image_tags_description_uindex ON image_tags (description);
CREATE TABLE images
(
    id INTEGER DEFAULT nextval('content.images_id_seq'::regclass) PRIMARY KEY NOT NULL,
    publish_time TIMESTAMP WITH TIME ZONE DEFAULT timezone('UTC'::text, now()) NOT NULL,
    last_modified TIMESTAMP WITH TIME ZONE DEFAULT timezone('UTC'::text, now()) NOT NULL,
    owner_id INTEGER NOT NULL,
    featured BOOLEAN DEFAULT false NOT NULL,
    downloads INTEGER DEFAULT 0 NOT NULL,
    views INTEGER DEFAULT 0 NOT NULL,
    shortcode VARCHAR(12) NOT NULL
);
CREATE TABLE user_favorites
(
    user_id INTEGER NOT NULL,
    image_id INTEGER PRIMARY KEY NOT NULL,
    CONSTRAINT user_favorites_users_id_fk FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT user_favorites_images_id_fk FOREIGN KEY (image_id) REFERENCES images (id)
);
CREATE UNIQUE INDEX user_favorites_user_id_image_id_uindex ON user_favorites (user_id, image_id);
CREATE TABLE user_follows
(
    user_id INTEGER NOT NULL,
    followed_id INTEGER PRIMARY KEY NOT NULL,
    CONSTRAINT user_follow_users_id_fk FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT user_followed_users_id_fk FOREIGN KEY (followed_id) REFERENCES users (id)
);
CREATE UNIQUE INDEX user_follows_user_id_uindex ON user_follows (user_id, followed_id);
CREATE TABLE users
(
    id INTEGER DEFAULT nextval('content.users_id_seq'::regclass) NOT NULL,
    username VARCHAR(100) PRIMARY KEY NOT NULL,
    name TEXT,
    email VARCHAR(100) NOT NULL,
    bio TEXT,
    url VARCHAR(50),
    password TEXT NOT NULL,
    salt TEXT NOT NULL,
    featured BOOLEAN DEFAULT false NOT NULL,
    admin BOOLEAN DEFAULT false NOT NULL,
    views INTEGER DEFAULT 0 NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('UTC'::text, now()) NOT NULL,
    last_modified TIMESTAMP WITH TIME ZONE DEFAULT timezone('UTC'::text, now()) NOT NULL
);
CREATE UNIQUE INDEX users_id_uindex ON users (id);
CREATE TABLE can_delete
(
    user_id INTEGER NOT NULL,
    o_id INTEGER NOT NULL,
    type CONTENT_TYPE NOT NULL
);
CREATE TABLE can_edit
(
    user_id INTEGER NOT NULL,
    o_id INTEGER NOT NULL,
    type CONTENT_TYPE NOT NULL
);
CREATE TABLE can_view
(
    user_id INTEGER NOT NULL,
    o_id INTEGER NOT NULL,
    type CONTENT_TYPE NOT NULL
);
