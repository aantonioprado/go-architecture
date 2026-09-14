# Layered (N-tier)

> One of the architecture implementations in this repository, organizing the application into Presentation, Business, and Data layers with a strict top-down dependency flow.

---

## Core Principle

Layered (N-tier) separates the application into three responsibilities, where each layer only calls the layer directly below it:

- **Presentation** - the `handler` package. Handles HTTP requests and responses, and never talks to persistence directly.
- **Business** - the `service` package. Owns validation, ID/timestamp generation, and the email-uniqueness rule.
- **Data** - the `repository` package. Owns persistence only, with no validation or business rules.

The `model` package holds only the `User` entity's data, with no behavior of its own; all logic that acts on it lives in the Business layer.

---

## Dependency Flow

```text
main --> server --> router --> handler --> service --> repository --> model
                                    |
                                    +--> dto
                                    +--> response
```

- The **Handler** only calls the **Service**; it never reaches into the **Repository** or **Model** directly.
- The **Service** validates input, generates the ID/timestamp, enforces email uniqueness, and is the only layer allowed to call the **Repository**.
- The **Repository** is responsible for persistence only and has no HTTP, JSON, or business-rule responsibilities.

---

## Directory Structure

```text
layered/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   ├── dto/
│   ├── handler/
│   ├── middleware/
│   ├── model/
│   ├── repository/
│   ├── response/
│   ├── routes/
│   ├── server/
│   └── service/
├── .air.toml
├── .env.example
├── go.mod
└── go.sum
```

---

## When to Use

Layered works well for traditional CRUD applications where teams want a clear split between request handling, business rules, and persistence, and are comfortable with a strict top-down call chain.

It is suitable for:

- Traditional CRUD applications
- Teams already familiar with layered/N-tier patterns
- Applications where business rules need a dedicated home separate from both HTTP handling and storage

---

## Pros and Cons

**Pros**

- ✅ Clear, well-known separation of concerns
- ✅ Business rules live in one dedicated place
- ✅ Easy to test each layer in isolation (handler, service, repository)
- ✅ Predictable, strictly top-down call chain
- ✅ Familiar to most enterprise teams

**Cons**

- ⚠️ More boilerplate than a simpler design for very small CRUD services
- ⚠️ Strict layering can feel like ceremony for trivial operations
- ⚠️ The service layer can become a dumping ground if not kept focused
- ⚠️ An anemic model pushes all behavior into the service, which can grow large
- ⚠️ Cross-cutting concerns (transactions spanning multiple entities, for example) don't have an obvious home

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
