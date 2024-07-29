ALTER TABLE "users" ADD COLUMN "email_verified" bool NOT NULL DEFAULT false;

CREATE TABLE "verify_emails" (
    id_verify_email serial PRIMARY KEY,
    id_user integer UNIQUE NOT NULL,  -- Уникальный, чтобы не было >1 запроса от одного пользователя
    secret_key varchar NOT NULL,
    expired_at timestamptz NOT NULL DEFAULT (NOW() + INTERVAL '10 minutes'),
    FOREIGN KEY (id_user) REFERENCES "users" (id_user)
);
