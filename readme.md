
# Проект к тестовому [заданию](https://github.com/avito-tech/tech-internship/blob/main/Tech%20Internships/Backend/Backend-trainee-assignment-winter-2025/Backend-trainee-assignment-winter-2025.md) Avito Shop

## Быстрый звапуск через docker
для запуска подготовлены файлы [Dockerfile](Dockerfile) и [docker-compose.yml](docker-compose.yaml)
просто выполните команду 
Для сборки
```bash
docker-compose build
````
Для запуска
```bash
docker-compose up -d
```

## environment параметры которые нужны для запуска 
### параметры с # необязательны
```env
#jwt secret key
SECRET_KEY = "testkey"
#JWT_EXPIRATION_TIME=12 int param value set in hours

#application mode
APP_MODE=debug

#database connection params
DB_USER=postgres
DB_PASSWORD=password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=avito_shop
DB_SSLMODE=disable
#DB_SSLROOTCERT=/path/to/cert
#DB_MAXOPENCONNS=5
#DB_MSXIDLECONNS=5
#DB_MAXLIFETIME int param value set in seconds
#DB_MAXLIFETIME=5
```

## Структура проекта
```
cmd
├── main.go
internal
├── app
│   ├── handlers
│   │   ├── models.go  -- models for handlers
│   │   └── handlers.go -- gin handlers methods
│   ├── mw 
│   │   └── middleware.go -- middleware auth methods
│   └── app.go -- main app methods
├── service
│   ├── service.go -- service init methods
│   ├── user_service.go -- user service methods
│   ├── wallet_service.go -- wallet service methods
│   └── models.go -- models for service
├── storage
│   ├── employees.go -- employees storage methods
│   ├── merch.go -- merch storage methods
│   ├── wallet.go -- wallet storage methods
│   └── db.go -- db init methods
── utils
│   └── jwt_utils.go -- jwt methods
├── migrations
│    └── init.sql -- db migrations
└── tests
    ├── e2e
    │   ├── purchase_test.go -- auth, buy merch, check balance
    │   └── transfer_test.go -- auth, transfer coin to user, check balance  
    └── helpers.go -- help methods for test

```

## для удобства имеется [makefile](makefile)
```bash
make help
```
выводит доступные команды
```bash
Доступные команды:
  install-deps   Устанавливает зависимости проекта
  get-deps       Загружает зависимости проекта
  build          Собирает проект
  run           Запускает проект
  build-docker   Собирает докер контейнер
  run-docker     Запускает докер контейнер
  clean          Очищает сгенерированные файлы
```

## прилагаю файл линтера
[.golangci.yml](.golangci.yml)

### Вопросы возникшие по тз
1. Допускалось ли использование redis или memcached?

***Не использовал ни то, ни то так как в тз не было упоминания об это!***

2. Для авторизации доступов должен использоваться JWT. Пользовательский токен доступа к API выдается после авторизации/регистрации пользователя. При первой авторизации пользователь должен создаваться автоматически.

***
Здесь непонятно если пользователя нет то вводя username и password создается новый пользователь? 
Если да то нужно ли было реализовать какой то метод подтверждения регистрации?
***
