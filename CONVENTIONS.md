# Project Conventions

## 1. Package Naming

- Use singular, lowercase names: `handler`, `service`, `repository`, `model`
- Avoid stutter: `service.UserService` is acceptable but keep names short.

## 2. Project Structure

```
internal/
├── handler/ # HTTP handlers (input parsing, response)
├── service/ # Business logic (validation, orchestration)
├── repository/ # Wraps sqlc for DB operations
├── model/ # Domain structs shared across layers
└── db/
├── generated/ # sqlc-generated code
└── queries/ # .sql query files
```

## 3. HTTP Routes

- Use plural nouns: `/users`, `/users/{id}`, `/users/register`
- Handlers for all routes live in `internal/handler`.

## 4. sqlc

- `.sql` files go in `internal/db/queries`
- Generated code goes in `internal/db/sqlc`, package name `sqlc`
- Repositories call sqlc; services never call sqlc directly.

## 5. Handlers

- Parse HTTP input → call service → send response
- Do not put business logic in handlers.

## 6. Services

- Contain business rules (validation, orchestration)
- Do not access database directly; only call repository functions.

## 7. Repositories

- Wrap sqlc queries to abstract persistence layer
- Keep them thin, only responsible for DB operations.

## 8. Testing

- Prefer real database tests (Postgres testcontainer or SQLite in-memory)
- Avoid interfaces unless mocking is required for unit tests
- If needed, define interfaces where they are consumed (e.g., in service tests).

## 9. Imports

- `internal/handler` → `handler.RegisterRoutes(mux, services...)`
- `internal/service` → `service.NewUserService(repo)`
- `internal/repository` → `repository.NewUserRepository(db)`
- `internal/db/generated` → `db.New()`
