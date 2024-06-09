#!/usr/bin/env sh

# Если где-то ошибка, то скрипт завершится
set -e

echo "Ждём postgres"
sleep 2

echo "Запускаем миграции"
/app/education/migrate -path /app/education/migration/ -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${DB_HOST}:5432/education?sslmode=${DB_SSL_MODE}" -verbose up

echo "Приложение запущено"
exec "$@"  # Выполняем все переданные команды (они передаются в CMD Dockerfile)
