CREATE TABLE "sessions" (
    id_session      uuid PRIMARY KEY NOT NULL,  -- Его мы создаём при вызове loginUser
    issued_at       timestamptz DEFAULT now(),
    expired_at      timestamptz NOT NULL,
    refresh_token   VARCHAR NOT NULL,
    id_user         integer NOT NULL,
    client_ip       VARCHAR NOT NULL,
    blocked         boolean NOT NULL DEFAULT FALSE
);
