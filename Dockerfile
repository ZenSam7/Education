# Устанавливаем имя для этапа первого сборки
FROM golang:1.22.4-alpine3.20 AS builder
LABEL author="ZenSam7"
WORKDIR /app/education

COPY . .
# Компилируем приложение
RUN go build -o main main.go

# На втором этапе устанавливаем только alpine чтобы прога занимала как можно меньше места
FROM alpine:3.20
WORKDIR /app/education
# Тут мы просто из builder копируем, а не сами файлы
COPY --from=builder /app/education/main .
# Отдельно загружаем конфигурацию и миграции
COPY .env .
COPY templates ./templates
COPY db/migration ./migration

EXPOSE 8080

# Запускаем приложение
CMD ["/app/education/main"]
