include .env

POSTGRES_URL = "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${DB_HOST}:5432/education?sslmode=${DB_SSL_MODE}"

# Создаём новый контейнер с бд
db: volume
	find . -type f -name "*.sh" -exec dos2unix {} \;
	docker run --name db \
		-p 5432:5432 \
		--net education_net \
		--env-file .env \
		-e PGPASSWORD=root \
		-v db_data \
		-v ./init_db.sh:/docker-entrypoint-initdb.d/init_db.sh \
		-d postgres:16 \
		-c wal_level=replica \
		-c max_wal_senders=2 \
		-c max_replication_slots=2

replic: volume
	docker run --name db_repl  \
		--network education_net \
		--env-file .env \
		-e PGPASSWORD=root \
		-u postgres \
		-v db_repl \
		--rm -d postgres:16 \
		/bin/bash -c "\
			pg_basebackup -h db -D /var/lib/postgresql/data -U root -Fp -Xs -P && \
			echo \"primary_conninfo = 'host=db port=5432 user=root password=root'\" >> /var/lib/postgresql/data/postgresql.conf && \
			echo \"standby_mode = 'on'\" > /var/lib/postgresql/data/standby.signal && \
			chmod 0750 /var/lib/postgresql/data && \
			exec postgres -c 'hot_standby=on'"

# Создаём новую бд
createdb:
	docker exec db createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} education
# Удаляем бд
dropdb:
	docker exec db dropdb education

redis: volume
	docker run --name redis \
		-p 6379:6379 \
		--network education_net \
		-v redis \
		-d redis:7

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
refreshdb: dropdb createdb migrateup

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
RESTART: RESET server

# Команда для сборки и запуска приложения
server:
	sudo go run main.go

net:
	docker network create education_net
volume:
	docker volume create db_data
	docker volume create db_repl
	docker volume create redis

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
generate: mock sqlc proto db_doc

# Команда для запуска всего стека
all: volume net db replic redis server

# Команда для остановки и удаления всех контейнеров и volumes
clean:
	docker rm -f db db_repl redis edu || true
	docker volume rm db_data db_repl redis || true
	docker network rm education_net || true

from_scratch: mock sqlc proto clean all

.PHONY: db createdb dropdb migrateup migrateup1 migratedown migratedown1 makemigrate db_doc all clean
.PHONY: connect refreshdb sqlc test RESET RESTART server net proto redis mock volume generate from_scratch
