CREATE TABLE "images" (
    id_image serial PRIMARY KEY,
    name     varchar NOT NULL,
    content  BYTEA NOT NULL,
    id_user  integer NOT NULL,
    FOREIGN  KEY (id_user) REFERENCES "users" (id_user)
);

ALTER TABLE "articles" ADD COLUMN "id_images" integer[];
ALTER TABLE "users" ADD COLUMN "avatar" integer NOT NULL;
ALTER TABLE "users" ADD FOREIGN KEY (avatar) REFERENCES images(id_image);
