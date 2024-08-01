CREATE TABLE "sessions" (
    id_session      uuid PRIMARY KEY NOT NULL,  -- id_session мы создаём при вызове loginUser
    issued_at       timestamptz DEFAULT now(),
    expired_at      timestamptz NOT NULL,
    refresh_token   VARCHAR NOT NULL,
    id_user         integer NOT NULL,
    client_ip       VARCHAR NOT NULL,
    blocked         boolean NOT NULL DEFAULT FALSE,
    FOREIGN KEY (id_user) REFERENCES "users" (id_user)
);
