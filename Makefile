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
	sudo docker exec -it postgres16 createdb --username=root --owner=root education
# Удаляем бд
dropdb:
	sudo docker exec -it postgres16 dropdb education
# Поднимаем миграции (т.е. переходим к текущей версии бд)
migrateup:
	migrate -path /mnt/c/Users/samki/GoProjects/Education/db/migration/ -database "postgresql://root:root@localhost:5432/education?sslmode=disable" up
# Опускаем миграции (т.е. переходим к прошлой версии бд)
migratedown:
	migrate -path /mnt/c/Users/samki/GoProjects/Education/db/migration/ -database "postgresql://root:root@localhost:5432/education?sslmode=disable" down
# Подключаемся к бд
connect:
	sudo docker exec -it postgres16 psql -U root education
# Удаляем и создаём новую бд со всеми миграциями
refreshdb:
	sudo make dropdb && sudo make createdb && sudo make migrateup

# Создаём запросы через sqlc
sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb migrateup migratedown connect refreshdb
