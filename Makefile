# Файл для упрощения работы с командами
# Заменяет команды на сокращённые версии
# Как воспользоваться: Открываем терминал WSL Ubuntu и пишем:
# $ cd /mnt/c/Users/samki/GoProjects/Education
# $ make <название нашей команды>

include .env

POSTGRES_URL = "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${DB_HOST}:5432/education?sslmode=${DB_SSL_MODE}"

# Создаём новый контейнер
postgres:
	docker run --name postgres -p 5432:5432 --net edu_net -e POSTGRES_USER=${POSTGRES_USER} \
 		-e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres:16
# Создаём новую бд
createdb:
	docker exec postgres createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} education
# Удаляем бд
dropdb:
	docker exec postgres dropdb education

# Поднимаем миграции (т.е. переходим к новой версии бд)
migrateup:
	migrate -path ./db/migration/ -database ${POSTGRES_URL} up
migrateup1:
	migrate -path ./db/migration/ -database ${POSTGRES_URL} up 1
# Опускаем миграции (т.е. переходим к прошлой версии бд)
migratedown:
	migrate -path ./db/migration/ -database ${POSTGRES_URL} down
migratedown1:
	migrate -path ./db/migration/ -database ${POSTGRES_URL} down 1
# Создаём новую миграцию
makemigrate:
	migrate create -ext sql -dir migration -seq

# Подключаемся к бд
connect:
	docker exec -it postgres psql -U root education
# Удаляем и создаём новую бд со всеми миграциями
refreshdb:
	make dropdb && make createdb && make migrateup

# Создаём код для запросов через sqlc
sqlc:
	sqlc generate
# Запускаем все тесты с подробным описанием, проверкой на полное покрытие тестов и без кеширования
test:
	sudo go test -count=1 -cover ./...

# Пересоздаём нахер всё
RESET:
	docker restart postgres && make refreshdb && make sqlc
# Как RESET только ещё и сервер запускаем
RESTART:
	make RESET && make server

# Запускаем cервер
server:
	sudo go run main.go

# Собираем образ
myimage:
	docker build -t education:latest .
# Меняем DB_HOST с localhost на "postgres" (имя контейнера), т.к. они подключены к одной сети edu_net
runimage:
	docker run --name edu -p 8080:8080 -e GIN_MODE=release -e DB_HOST="postgres" --net edu_net education:latest
net:
	docker network create edu_net

proto:
	rm pb/*.go
	protoc --proto_path=proto --go_out=pb --go-grpc_out=pb \
		--go_opt=paths=source_relative --go-grpc_opt=paths=source_relative proto/*.proto

evans:
	evans --host localhost --port 1213 -r repl

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 makemigrate
.PHONY: connect refreshdb sqlc test RESET RESTART server myimage runimage net proto evans
