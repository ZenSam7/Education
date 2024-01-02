CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "created_at" timestamptz DEFAULT (now()),
  "name" varchar NOT NULL,
  "description" varchar,
  "email" varchar,
  "karma" integer DEFAULT 0
);

CREATE TABLE "articles" (
  "id" bigserial PRIMARY KEY,
  "created_at" timestamptz DEFAULT (now()),
  "title" varchar NOT NULL,
  "text" text NOT NULL,
  "comments" bigserial,
  "from_user" bigserial,
  "evaluation" integer DEFAULT 0
);

CREATE TABLE "comments" (
  "id" bigserial PRIMARY KEY,
  "created_at" timestamptz DEFAULT (now()),
  "text" text NOT NULL,
  "from_user" bigserial,
  "evaluation" integer DEFAULT 0
);

CREATE INDEX name_ind ON "users" ("name");

CREATE INDEX title_ind ON "articles" ("title");

CREATE INDEX from_user_ind ON "articles" ("from_user");

ALTER TABLE "articles" ADD FOREIGN KEY ("comments") REFERENCES "comments"("id");

ALTER TABLE "articles" ADD FOREIGN KEY ("from_user") REFERENCES "users"("id");

ALTER TABLE "comments" ADD FOREIGN KEY ("from_user") REFERENCES "users"("id");
