# Устанавливаем имя для этапа первого сборки
FROM golang:1.22.4-alpine3.20 AS builder
WORKDIR /app/education
LABEL author="ZenSam7"

COPY . .
# Компилируем приложение в файл main(.exe) с точкой входа из main.go
RUN go build -o main main.go

# На втором этапе устанавливаем только alpine чтобы прога занимала как можно меньше места
FROM alpine:3.20
WORKDIR /app/education
COPY --from=builder /app/education/main .
# Отдельно загружаем конфигурацию
COPY .env .

EXPOSE 8080

# Запускаем приложение (но с возможностью проигнорировать эту команду)
CMD ["/app/education/main"]