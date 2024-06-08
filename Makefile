# Файл для упрощения работы с командами
# Заменяет команды на сокращённые версии
# Как воспользоваться: Открываем терминал WSL Ubuntu и пишем:
# $ cd /mnt/c/Users/samki/GoProjects/Education
# $ make <название нашей команды>

include .env

# Создаём новый контейнер
postgres:
	sudo docker run --name postgres -p 5432:5432 --net edu_net -e POSTGRES_USER=${DB_PASSWORD} -e POSTGRES_PASSWORD=${DB_USER_NAME} -d postgres:16
# Создаём новую бд
createdb:
	sudo docker exec postgres createdb --username=${DB_USER_NAME} --owner=${DB_USER_NAME} education
# Удаляем бд
dropdb:
	sudo docker exec postgres dropdb education

# Поднимаем миграции (т.е. переходим к новой версии бд)
migrateup:
	migrate -path ./db/migration/ -database "postgresql://${DB_USER_NAME}:${DB_PASSWORD}@${DB_HOST}:5432/${DB_NAME}?sslmode=${DB_SSL_MODE}" up
migrateup1:
	migrate -path ./db/migration/ -database "postgresql://${DB_USER_NAME}:${DB_PASSWORD}@${DB_HOST}:5432/${DB_NAME}?sslmode=${DB_SSL_MODE}" up 1
# Опускаем миграции (т.е. переходим к прошлой версии бд)
migratedown:
	migrate -path ./db/migration/ -database "postgresql://${DB_USER_NAME}:${DB_PASSWORD}@${DB_HOST}:5432/${DB_NAME}?sslmode=${DB_SSL_MODE}" down
migratedown1:
	migrate -path ./db/migration/ -database "postgresql://${DB_USER_NAME}:${DB_PASSWORD}@${DB_HOST}:5432/${DB_NAME}?sslmode=${DB_SSL_MODE}" down 1
# Создаём новую миграцию
makemigrate:
	migrate create -ext sql -dir migration -seq

# Подключаемся к бд
connect:
	sudo docker exec -it postgres psql -U root education
# Удаляем и создаём новую бд со всеми миграциями
refreshdb:
	sudo make dropdb && sudo make createdb && sudo make migrateup

# Создаём код для запросов через sqlc
sqlc:
	sudo sqlc generate
# Запускаем все тесты с подробным описанием и проверкой на полное покрытие тестов
test:
	make sqlc && go test -cover ./...

# Пересоздаём нахер всё
RESET:
	sudo docker restart postgres && sudo make refreshdb && sudo make sqlc
# Как RESET только ещё и сервер запускаем
RESTART:
	make RESET && make server

# Запускаем cервер
server:
	sudo go run main.go

# Собираем образ
myimage:
	sudo docker build -t education:latest .
# Меняем DB_HOST с localhost на "postgres" (имя контейнера), т.к. они подключены к одной сети edu_net
runimage:
	docker run --name edu -p 8080:8080 -e GIN_MODE=release -e DB_HOST="postgres" --net edu_net education:latest


.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 makemigrate
.PHONY: connect refreshdb sqlc test RESET RESTART server myimage runimage
