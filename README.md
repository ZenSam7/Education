# Это просто сайт для освоения бекенда

### Чему я тут научился (что я тут использовал):
- Кончно же **Golang** (интерфейсы, контексты, горутины ...)
- Сервер на **gRPC с HTTP Gateway**
- Зачем-то сделал ещё и **RESTful API на фреймворке Gin**
- Авторизация по токенам и верификация почты
- Использовал Redis как брокер сообщений


- Развернул приложение в **Docker** (написал Dockerfile и docker-compose)
- **Mock-тесты** и обычные **unit-тесты**, с общим покрытием ~50%
- Поработал с **CI/CD** (GitHub Actions)
- Лучше освоил **linux**


- Бд реализовал на **PostgreSQL**
- SQLC для генерации кода на Golang из SQL запроса 
- Где необходимо, сделал транзакции


  Мелочи:
- <p>>5000 строк отлаженного, работающего, лично написанного кода (+5500 сгенерированного)</p>
- Понял что Makefile — очень удобная штука
- Migrate для, собственно, миграций 
- Настроил красивый логгер
- Генерация документации — приятная вещь


В планах:
- Сделать фронтенд
- Использовать Redis ещё и как кэш

# Как запустить:
Перейти в директорию проекта и выполнить команду на ubuntu (если возникла ошибка с подключением в бд
попробуйте вызвать команду ещё раз):

```shell
make redis
make postgres
make createdb
make server
```

Или можно запустить через Docker Compose:
```shell
docker compose up
```
