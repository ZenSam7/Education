ALTER TABLE "users" ADD COLUMN "email" varchar UNIQUE NOT NULL CHECK ( email <> '' );
ALTER TABLE "users" ADD CONSTRAINT email_must_be_valid CHECK ( email ~ '^[a-zA-Z0-9.!#$%&''*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$' );
ALTER TABLE "users" ADD COLUMN "password_hash" varchar UNIQUE NOT NULL CHECK ( password_hash <> '' );
