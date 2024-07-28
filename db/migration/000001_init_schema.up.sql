-- Между varchar и text разницы нет

-- Крч если через go передать в качестве текстового аргумента nil то он замениться на '',
-- а '' != NULL поэтому она вставиться как пустая строка, хотя в go мы передали nil,
-- поэтому надо использовать CHECK ( title <> '' )

CREATE TABLE "users" (
    id_user     serial PRIMARY KEY,
    created_at  timestamptz DEFAULT now(),
    name        varchar NOT NULL CHECK ( name <> '' ),
    description text DEFAULT NULL,
    karma       integer DEFAULT 0 NOT NULL
);

CREATE TABLE "articles" (
    id_article  serial PRIMARY KEY,
    created_at  timestamptz DEFAULT now(),
    edited_at   timestamptz,
    title       varchar NOT NULL CHECK ( title <> '' ),
    text        text NOT NULL CHECK ( text <> '' ),
    comments    integer[],
    authors     integer[] NOT NULL,
    evaluation  integer NOT NULL DEFAULT 0
);

CREATE TABLE "comments" (
    id_comment  serial PRIMARY KEY,
    created_at  timestamptz DEFAULT now(),
    edited_at   timestamptz DEFAULT NULL,
    text        text NOT NULL CHECK ( text <> '' ),
    author      integer NOT NULL,
    evaluation  integer NOT NULL DEFAULT 0,
    FOREIGN KEY (author) REFERENCES "users" (id_user)
);

CREATE INDEX user_indx ON "users" ("id_user");

CREATE INDEX article_indx ON "articles" ("id_article");

CREATE INDEX comment_indx ON "comments" ("id_comment");
