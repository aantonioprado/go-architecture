# MVC (Model-View-Controller)

> One of the architecture implementations in this repository, using the traditional MVC structure for a User CRUD API - built as a JSON REST API instead of server-rendered HTML.

---

## Core Principle

MVC separates request handling, data representation, and application data into three responsibilities:

- **Model** - represents the application's entities and data.
- **View** - defines how data is presented to the client. In this API, the `dto` and `response` packages handle JSON representation instead of HTML templates.
- **Controller** - handles HTTP requests, coordinates operations, and returns responses.

This implementation also uses a **Repository** layer for persistence. It is kept separate from the Model so that entities do not depend on how data is stored.

---

## Dependency Flow

```text
main --> server --> router --> controller --> repository --> model
                                    |
                                    +--> dto
                                    +--> response
```

- The **Controller** coordinates the request flow and communicates with the other components.
- The **Repository** is responsible for persistence and has no HTTP or JSON responsibilities.
- The **Model** contains the entities and their data without depending on the storage layer.

---

## Directory Structure

```text
mvc/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   ├── controller/
│   ├── dto/
│   ├── middleware/
│   ├── model/
│   ├── repository/
│   ├── response/
│   ├── routes/
│   └── server/
├── .air.toml
├── .env.example
├── go.mod
└── go.sum
```

---

## When to Use

MVC works well for HTTP APIs and simple CRUD applications that need to separate Controllers, data, and persistence without introducing a dedicated service or use-case layer.

It is suitable for:

- Administrative panels
- Small web applications
- APIs with simple business rules
- Applications with a straightforward request flow

---

## Pros and Cons

**Pros**

- ✅ Familiar and easy-to-understand structure
- ✅ Simple request flow
- ✅ Clear separation between Controller, Model, and Repository
- ✅ Low initial complexity
- ✅ Little configuration and ceremony
- ✅ Constructor-based dependency injection makes testing easier

**Cons**

- ⚠️ Controllers can become large as the application grows
- ⚠️ No explicit service or use-case layer
- ⚠️ Business rules involving multiple entities can become difficult to organize
- ⚠️ Less suitable for complex domains with extensive business logic
- ⚠️ Separation of responsibilities can become unclear as the application grows

---

## Running

```bash
cp .env.example .env
go run ./cmd/api
```

To run with hot reload:

```bash
air
```

---

## Example Usage

**Create user:**

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Antônio Prado","email":"antonio@antonioeprado.dev"}'
```

**List users:**

```bash
curl http://localhost:8080/users
```

**Get user by id:**

```bash
curl http://localhost:8080/users/{id}
```

**Update user:**

```bash
curl -X PUT http://localhost:8080/users/{id} \
  -H "Content-Type: application/json" \
  -d '{"name":"Antônio Elias Prado","email":"antonio@antonioeprado.dev"}'
```

**Delete user:**

```bash
curl -X DELETE http://localhost:8080/users/{id}
```

**Health check:**

```bash
curl http://localhost:8080/health
```

Response shape for a user:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Antônio Prado",
  "email": "antonio@antonioeprado.dev",
  "createdAt": "2026-01-20T00:00:00Z"
}
```
