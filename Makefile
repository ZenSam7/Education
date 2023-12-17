# Этот файл всегда запускается когда мы собираем проект в бинарник

.PHONY: build

build:
	go build -c main.go

.DEFAULT_GOAL := build