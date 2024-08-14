#!/bin/bash
set -e

# Создаем базу данных и пользователя
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE DATABASE education;
    GRANT ALL PRIVILEGES ON DATABASE education TO $POSTGRES_USER;
EOSQL

# Создаем публикацию для репликации
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname education <<-EOSQL
    CREATE PUBLICATION my_pub FOR ALL TABLES;
EOSQL

# Добавляем запись для репликации в pg_hba.conf
echo "host    replication     root       0.0.0.0/0   md5" >> /var/lib/postgresql/data/pg_hba.conf
