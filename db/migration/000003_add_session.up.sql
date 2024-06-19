CREATE TABLE "sessions" (
    id_session  serial PRIMARY KEY,
    issued_at   timestamptz DEFAULT now(),
    expired_at  timestamptz NOT NULL,
    id_user     integer NOT NULL,
    client_ip   VARCHAR NOT NULL,
    is_blocked  boolean NOT NULL DEFAULT FALSE
);

ALTER TABLE "sessions" ADD CONSTRAINT fk_user_session FOREIGN KEY (id_user) REFERENCES "users" (id_user);
