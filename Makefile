# Файл для упрощения работы с командами
# Заменяет команды на сокращённые версии
# Как воспользоваться: Открываем терминал WSL Ubuntu и пишем:
# $ cd /mnt/c/Users/samki/GoProjects/Education
# $ make <сокращённое название команды>

postgres:
	sudo docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:16
createdb:
	sudo docker exec -it postgres16 createdb --username=root --owner=root habr2
dropdb:
	sudo docker exec -it postgres16 dropdb habr2

.PHONY: postgres createdb dropdb
