# Файл для упрощения работы с командами
# Заменяет команды на сокращённые версии
# Как воспользоваться: Открываем терминал WSL Ubuntu и пишем:
# $ cd /mnt/c/Users/samki/GoProjects/Education
# $ make <название нашей команды>

# Создаём новый контейнер
postgres:
	sudo docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:16
# Создаём новую бд
createdb:
	sudo docker exec postgres16 createdb --username=root --owner=root education
# Удаляем бд
dropdb:
	sudo docker exec postgres16 dropdb education

# Поднимаем миграции (т.е. переходим к новой версии бд)
migrateup:
	migrate -path ./db/migration/ -database "postgresql://root:root@localhost:5432/education?sslmode=disable" up
migrateup1:
	migrate -path ./db/migration/ -database "postgresql://root:root@localhost:5432/education?sslmode=disable" up 1
# Опускаем миграции (т.е. переходим к прошлой версии бд)
migratedown:
	migrate -path ./db/migration/ -database "postgresql://root:root@localhost:5432/education?sslmode=disable" down
migratedown1:
	migrate -path ./db/migration/ -database "postgresql://root:root@localhost:5432/education?sslmode=disable" down 1

# Подключаемся к бд
connect:
	sudo docker exec -it postgres16 psql -U root education
# Удаляем и создаём новую бд со всеми миграциями
refreshdb:
	sudo make dropdb && sudo make createdb && sudo make migrateup

# Создаём код для запросов через sqlc
sqlc:
	sudo sqlc generate
# Запускаем все тесты с подробным описанием и проверкой на полное покрытие тестов
test:
	make sqlc && sudo go test -cover ./db/sqlc

# Перезапускаем нахер всё
RESET:
	sudo docker restart postgres16 && sudo make refreshdb && sudo make sqlc
# Как RESET только ещё и сервер запускаем
RESTART:
	make RESET && make server

# Запускаем cервер
server:
	sudo go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown connect refreshdb sqlc test RESET server
