# Файл для упрощения работы с командами
# Заменяет команды на сокращённые версии
# Как воспользоваться: Открываем терминал WSL Ubuntu и пишем:
# $ cd /mnt/c/Users/samki/GoProjects/Education
# $ make <название нашей команды>

include .env

POSTGRES_URL = "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${DB_HOST}:5432/education?sslmode=${DB_SSL_MODE}"

# Создаём новый контейнер с бд
postgres:
	docker run --name postgres -p 5432:5432 --net edu_net -v db_data:/var/lib/postgresql/data \
 		-e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres:16
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
# Создаём пустую миграцию
makemigrate:
	migrate create -ext sql -dir ./db/migration -seq $(name)

# Экспортируем схему из контейнера postgres в .sql
# Дальше при помощи https://dbdiagram.io/d уже делаем что захотим
db_schema:
	docker exec -it postgres pg_dump -h localhost -p 5432 -d education -U root -s -F p -E UTF-8 -f /bin/abc.sql
	docker cp postgres:/bin/abc.sql ./documentation/schema.sql

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
	sudo go test -count=1 -short -cover ./...
mock:
	mockgen -source=db/sqlc/querier.go -destination=db/mockdb/querier.go -package=mockdb

# Пересоздаём нахер всё
RESET:
	docker restart postgres && make refreshdb && make sqlc && make mock && make proto
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

# Если не работает proto, надо сделать эти 2 команды
# export GOPATH=$HOME/go
# PATH=$PATH:$GOPATH/bin
proto:
	protoc --proto_path=proto --go_out=protobuf --go-grpc_out=protobuf \
		   --openapiv2_out=documentation --openapiv2_opt=allow_merge=true,merge_file_name=gRPC_API_doc \
		   --grpc-gateway_out=protobuf --grpc-gateway_opt=paths=source_relative \
		   --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --experimental_allow_proto3_optional \
		   proto/*.proto

redis:
	docker run --name redis -p 6379:6379 -d redis:7

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 makemigrate db_schema
.PHONY: connect refreshdb sqlc test RESET RESTART server myimage runimage net proto redis mock
