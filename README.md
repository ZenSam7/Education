# Делаю бекенд учебного сайт

### Чему я тут научился (что я тут использовал):
- Кончно же **Golang** (интерфейсы, контексты, горутины ...)
- Сервер на **gRPC с HTTP Gateway**
- Зачем-то сделал ещё и **RESTful API на фреймворке Gin**
- Авторизация по 2м токенам (paseto или jwt) и верификация почты
- Использовал **Redis** и как брокер сообщений, и как кэш


- Развернул приложение в **Docker** (написал Dockerfile и docker-compose)
- **Mock-тесты** и обычные **unit-тесты**, с общим покрытием ~40%
- Поработал с **CI/CD** (GitHub Actions)
- Лучше освоил **linux**


- Имеется репликация данных в отдельный контейнер
- SQLC для генерации кода на Golang из SQL запроса
- Где необходимо, сделал транзакции
- Бд реализовал на **PostgreSQL**


  Мелочи:
- Сделал автогенерацию документации (в ./doc/)
- <p>>6000 строк отлаженного, работающего, лично написанного кода (+5500 сгенерированного)</p>
- Понял что Makefile — очень удобная штука
- Migrate для, собственно, миграций 
- Настроил красивый логгер


# Как запустить:
Можно запустить через Docker Compose:
```shell
docker compose up
```

Или можно выполнить команду на ubuntu:
```shell
make all
```

Убрать всё что было создано:
```shell
make clean
```

## Сам проект:
Архитектура бд:
<img src="./doc/schema.png" width=1500>
