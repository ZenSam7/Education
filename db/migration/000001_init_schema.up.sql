-- Крч между varchar и text разницы нет

CREATE TABLE "users" (
  "id"          serial PRIMARY KEY,
  "created_at"  timestamptz DEFAULT (now()),
  "name"        varchar NOT NULL,
  "description" varchar,
  "email"       varchar,
  "karma"       integer DEFAULT 0
);

CREATE TABLE "articles" (
  "id"          serial PRIMARY KEY,
  "created_at"  timestamptz DEFAULT (now()),
  "title"       varchar NOT NULL,
  "text"        varchar NOT NULL,
  "comments"    integer[],
  "from_users"  varchar[],
  "evaluation"  integer DEFAULT 0
);
COMMENT ON COLUMN articles.from_users IS 'from_user nil = пользователь удалён';

CREATE TABLE "comments" (
  "id"          serial PRIMARY KEY,
  "created_at"  timestamptz DEFAULT (now()),
  "text"        varchar NOT NULL,
  "from_user"   varchar,
  "evaluation"  integer DEFAULT 0,
  "edited"      boolean DEFAULT false
);
COMMENT ON COLUMN comments.from_user IS 'from_user nil = пользователь удалён';


CREATE INDEX name_ind ON "users" ("name");

CREATE INDEX title_ind ON "articles" ("title");

CREATE INDEX from_user_ind ON "articles" ("from_users");
