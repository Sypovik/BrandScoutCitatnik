# Структура проекта "Цитатник"
```
.
├── cmd
│   └── server
│       └── main.go
└── internal
    ├── handlers
    │   ├── handlers.go
    │   └── handlers_test.go
    ├── models
    │   └── quote.go
    ├── repository
    │   ├── memory.go
    │   └── repository.go
    └── router
        └── router.go
```

```bash
tree . -I "*.md"
```