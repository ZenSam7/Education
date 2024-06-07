ALTER TABLE "users" DROP CONSTRAINT email_must_be_valid;
ALTER TABLE "users" DROP COLUMN "email";
ALTER TABLE "users" DROP COLUMN "password_hash";
