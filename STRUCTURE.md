# Структура проекта "Цитатник"
```
.
├── cmd
│   └── server
│       └── main.go
├── go.mod
├── internal
│   ├── handlers
│   │   ├── handlers.go
│   │   └── handlers_test.go
│   ├── models
│   │   └── quote.go
│   ├── repository
│   │   ├── errors.go
│   │   ├── memory.go
│   │   ├── memory_test.go
│   │   └── repository.go
│   └── router
│       ├── router.go
│       └── router_test.go
└── makefile
```

Команда для вывода структуры
```bash
tree . --gitignore -I "*.md"
```

`cmd/server/main.go` — точка входа в приложение, запускает HTTP-сервер.

`internal`— внутренние пакеты, которые не предназначены для использования вне проекта.

`handlers` — обработка входящих HTTP-запросов.

`middleware` — промежуточные обработчики, которые выполняются между сервером и обработчиками.

`models` — описание структур данных и бизнес-логики.

`repository` — слой доступа к данным, абстракция над хранилищем.

`router` — настройка маршрутов для HTTP-сервера.

`makefile` — скрипты для упрощения запуска тестов, сборки и других команд.