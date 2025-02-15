PROJECT_NAME:= $(shell go list -m)
BASE_DIR:=$(CURDIR)
LOCAL_BIN:=$(CURDIR)/bin

# установка зависимостей
install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# получение зависимостей
get-deps:
	go mod tidy

# Сборка проекта
build:
	go build -o ./avito-shop-app ./cmd

# Запуск проекта
run: build
	./avito-shop-app


# Cборка докер контейнера
build-docker:
	docker-compose build

# Запуск докер контейнера
run-docker:
	docker-compose up -d

# Генерация документации swag
#swag-docs:
#	swag init -g ./cmd/main.go -o docs

# Удаление артефактов сборки
clean:
	rm -rf $(LOCAL_BIN)

# Помощь: выводит список доступных команд
help:
	@echo "Доступные команды:"
	@echo "  install-deps   Устанавливает зависимости проекта"
	@echo "  get-deps       Загружает зависимости проекта"
	@echo "  build          Собирает проект"
	@echo "  run           Запускает проект"
	@echo "  build-docker   Собирает докер контейнер"
	@echo "  run-docker     Запускает докер контейнер"
	@echo "  clean          Очищает сгенерированные файлы"
#	@echo "  swag-docs      Генерирует документацию swagger"

.default: help