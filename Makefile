include .env

POSTGRES_URL = "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${DB_HOST}:5432/education?sslmode=${DB_SSL_MODE}"

# Создаём новый контейнер с бд
db:
	docker run --name db -p 5432:5432 -v db_data:/var/lib/db_data/data \
 		-e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres:16

replic:
	docker run --name db_repl -v db_repl:/var/lib/postgresql/data \
 		-e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres:16
	# Записываем настройки реплики
	docker exec -i db sh -c \
		'echo "wal_level = logical\nmax_replication_slots = 2\nmax_wal_senders = 2" >> \
		/var/lib/postgresql/data/postgresql.conf'
	docker exec -i db sh -c \
		'echo "host    replication     репликация    192.168.1.0/24   md5" >> \
		/var/lib/postgresql/data/pg_hba.conf'
	docker exec -i db psql -c "SELECT pg_reload_conf();"
	docker exec -i db psql -c "SELECT * FROM pg_create_logical_replication_slot('edu_replic', 'pgoutput');"

# Создаём новую бд
createdb:
	docker exec db createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} education
# Удаляем бд
dropdb:
	docker exec db dropdb education

redis:
	docker run --name redis -v redis:/var/lib/redis/data -p 6379:6379 -d redis:7

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

# Экспортируем схему из контейнера db в .sql
# Дальше при помощи https://dbdiagram.io/d уже делаем что захотим
db_doc:
	docker exec -it db pg_dump -h localhost -p 5432 -d education -U root -s -F p -E UTF-8 -f /bin/abc.sql
	docker cp db:/bin/abc.sql ./doc/schema.sql
	sql2dbml doc/schema.sql --postgres -o doc/schema.dbml
	rm dbml-error.log

# Подключаемся к бд
connect:
	docker exec -it db psql -U root education
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
	mockgen -source=db/sqlc/querier.go -destination=my_mocks/db.go -package=my_mocks
	mockgen -source=redis/worker/distributor.go -destination=my_mocks/worker.go -package=my_mocks
	mockgen -source=token/maker.go -destination=my_mocks/token.go -package=my_mocks

# Пересоздаём нахер всё
RESET:
	docker restart db && make refreshdb && make sqlc && make mock && make proto
# Как RESET только ещё и сервер запускаем
RESTART:
	make RESET && make server

# Запускаем cервер
server:
	sudo go run main.go

net:
	docker network create education_net
volume:
	docker volume create redis
	docker volume create db_data
	docker volume create db_repl

# Если не работает proto, надо сделать эти 2 команды
# export GOPATH=$HOME/go
# PATH=$PATH:$GOPATH/bin
proto:
	protoc --proto_path=proto --go_out=protobuf --go-grpc_out=protobuf \
		   --openapiv2_out=doc --openapiv2_opt=allow_merge=true,merge_file_name=gRPC_API_doc \
		   --grpc-gateway_out=protobuf --grpc-gateway_opt=paths=source_relative \
		   --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --experimental_allow_proto3_optional \
		   proto/*.proto

# Генерируем всё что можно генерировать
generate:
	make mock && make sqlc && make proto && make db_doc

.PHONY: db createdb dropdb migrateup migrateup1 migratedown migratedown1 makemigrate db_doc
.PHONY: connect refreshdb sqlc test RESET RESTART server net proto redis mock volume generate
