-- Крч между varchar и text разницы нет

CREATE TABLE "users" (
  id_user       serial PRIMARY KEY,
  created_at    timestamptz DEFAULT now(),
  name          varchar NOT NULL,
  description   text NOT NULL,
  email         text NOT NULL,
  karma         integer DEFAULT 0 NOT NULL
);

CREATE TABLE "articles" (
  id_article    serial PRIMARY KEY,
  created_at    timestamptz DEFAULT now(),
  edited_at     timestamptz DEFAULT NULL,
  title         varchar NOT NULL,
  text          text NOT NULL,
  comments      integer[],
  authors       integer[] NOT NULL,
  evaluation    integer NOT NULL DEFAULT 0
);

CREATE TABLE "comments" (
  id_comment    serial PRIMARY KEY,
  created_at    timestamptz DEFAULT now(),
  edited_at     timestamptz DEFAULT NULL,
  text          text NOT NULL,
  from_user     integer NOT NULL,
  evaluation    integer NOT NULL DEFAULT 0
);


CREATE INDEX user_indx ON "users" ("id_user");

CREATE INDEX article_indx ON "articles" ("id_article");

CREATE INDEX comment_indx ON "comments" ("id_comment");
