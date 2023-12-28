CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "created_at" timestapmtz DEFAULT (now()),
  "name" varchar NOT NULL,
  "description" text,
  "email" varchar,
  "karma" integer DEFAULT 0
);

CREATE TABLE "articles" (
  "id" bigserial PRIMARY KEY,
  "created_at" timestapmtz DEFAULT (now()),
  "title" varchar NOT NULL,
  "text" text NOT NULL,
  "comments" bigserial,
  "from_user" bigserial,
  "evaluation" integer DEFAULT 0
);

CREATE TABLE "comments" (
  "id" bigserial PRIMARY KEY,
  "created_at" timestapmtz DEFAULT (now()),
  "text" text NOT NULL,
  "from_user" bigserial,
  "evaluation" integer DEFAULT 0
);

CREATE INDEX ON "users" ("name");

CREATE INDEX ON "articles" ("title");

CREATE INDEX  ON "articles" ("from_user");

ALTER TABLE "articles" ADD FOREIGN KEY ("comments") REFERENCES "comments" ("id");

ALTER TABLE "articles" ADD FOREIGN KEY ("from_user") REFERENCES "users" ("id");

ALTER TABLE "comments" ADD FOREIGN KEY ("from_user") REFERENCES "users" ("id");
