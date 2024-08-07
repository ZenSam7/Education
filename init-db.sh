#!/bin/bash

# Просто создаём бд в docker compose

set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE education;
    GRANT ALL PRIVILEGES ON DATABASE education TO $POSTGRES_USER;
EOSQL
