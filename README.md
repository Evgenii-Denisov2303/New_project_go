# Inkflow

Inkflow is a multi-service Go backend project with:

- `auth` HTTP service
- `articles` HTTP service
- `notification` service
- `profile` gRPC service
- PostgreSQL
- Kafka
- Docker Compose

## Stack

- Go `1.26.1`
- PostgreSQL `17`
- Apache Kafka `4.0.2`
- JWT
- gRPC + protobuf
- Docker Compose

## Services

- `auth` -> `http://localhost:8080`
- `articles` -> `http://localhost:8081`
- `notification` -> `http://localhost:8082`
- `profile` gRPC -> `localhost:50051`
- `postgres` -> `localhost:5432`
- `kafka` -> `localhost:9092`

## Features

- User registration and login
- Password hashing with bcrypt
- JWT token generation and validation
- Protected HTTP endpoints
- Article creation and reading
- Profile storage in PostgreSQL
- Profile access over gRPC
- Event-driven notification flow through Kafka

## Project structure

```text
cmd/
  auth/
  articles/
  notification/
  profile/
  migrator/

internal/
  clients/
  config/
  domain/models/
  services/
  storage/
  transport/

migrations/
proto/
```

## Run with Docker

```powershell
docker compose up -d --build
docker compose ps
```

Stop:

```powershell
docker compose down
```

## Run tests

```powershell
go test ./...
```

## Local entry points

```powershell
go run ./cmd/auth
go run ./cmd/articles
go run ./cmd/notification
go run ./cmd/profile
go run ./cmd/migrator
```

## HTTP examples

Register:

```powershell
Invoke-WebRequest `
  -Uri http://localhost:8080/register `
  -Method POST `
  -ContentType "application/json" `
  -Body '{"email":"test@example.com","password":"123456"}'
```

Login:

```powershell
$response = Invoke-WebRequest `
  -Uri http://localhost:8080/login `
  -Method POST `
  -ContentType "application/json" `
  -Body '{"email":"test@example.com","password":"123456"}'

$token = ($response.Content | ConvertFrom-Json).token
```

Get profile:

```powershell
Invoke-WebRequest `
  -Uri http://localhost:8080/profile `
  -Headers @{ Authorization = "Bearer $token" }
```

Create article:

```powershell
Invoke-WebRequest `
  -Uri http://localhost:8081/articles `
  -Method POST `
  -Headers @{ Authorization = "Bearer $token" } `
  -ContentType "application/json" `
  -Body '{"title":"My article","content":"Hello from Inkflow"}'
```

## gRPC example

Install `grpcurl`:

```powershell
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
$env:PATH += ";$HOME\\go\\bin"
```

Call `GetProfile`:

```powershell
$json = '{"user_id":1}'

grpcurl -plaintext `
  -import-path proto `
  -proto profile/v1/profile.proto `
  -d $json `
  localhost:50051 `
  profile.v1.ProfileService/GetProfile
```

## Kafka flow

Flow:

1. `articles` creates an article
2. `articles` publishes `article.created`
3. `notification` consumes the event from Kafka
4. `notification` processes the event

Check logs:

```powershell
docker logs notification-service --tail 100
```

Expected log:

```text
notification: article created: title=Kafka article author=user@example.com
```

## Migrations

The migrator applies:

- `001_create_users_table.sql`
- `002_create_articles_table.sql`
- `003_create_profiles_table.sql`
