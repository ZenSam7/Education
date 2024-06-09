# Устанавливаем имя для этапа первого сборки
FROM golang:1.22.4-alpine3.20 AS builder
LABEL author="ZenSam7"
WORKDIR /app/education

COPY . .
# Компилируем приложение и устанавливаем migration
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz

# На втором этапе устанавливаем только alpine чтобы прога занимала как можно меньше места
FROM alpine:3.20
WORKDIR /app/education
# Тут мы просто из builder копируем, а не сами файлы
COPY --from=builder /app/education/main .
COPY --from=builder /app/education/migrate .
# Отдельно загружаем конфигурацию и миграции
COPY .env .
COPY run.sh .
RUN chmod +x ./run.sh
COPY db/migration ./migration

EXPOSE 8080

# Запускаем приложение (но с возможностью проигнорировать эту команду)
ENTRYPOINT ["/app/education/run.sh", "/app/education/main"]
