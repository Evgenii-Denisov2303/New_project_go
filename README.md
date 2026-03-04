# Go Training Project

Небольшой учебный REST API для практики Go (стандартная библиотека, без фреймворков).

## Что внутри

- `GET /health` - проверка сервиса
- `GET /tasks` - список задач
- `POST /tasks` - создать задачу
- `PATCH /tasks/{id}/done` - отметить задачу как выполненную
- `DELETE /tasks/{id}` - удалить задачу

## Запуск

```bash
go run ./cmd/api
```

Сервер стартует на `http://localhost:8080`.

## Примеры запросов

```bash
curl http://localhost:8080/health
```

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Изучить указатели в Go"}'
```

```bash
curl http://localhost:8080/tasks
```

```bash
curl -X PATCH http://localhost:8080/tasks/1/done
```

```bash
curl -X DELETE http://localhost:8080/tasks/1
```

## Тесты

```bash
go test ./...
```
