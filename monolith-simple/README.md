# Simple Monolith

> One of the architecture implementations in this repository: a single package with one entity, one in-memory store, and the HTTP handlers that read and write it directly.

---

## Core Principle

Keep everything in one place. There is no controller/service/repository split: the HTTP handlers are methods on the same `userStore` that holds the data, and validation, ID/timestamp generation, and the email-uniqueness rule all live right next to the map they operate on.

---

## Dependency Flow

```text
main --> router --> handlers (userStore methods) --> users map
```

- **Mixed.** Handlers, validation, and persistence are not separated into layers; they are different methods on the same type.
- There is nothing to trace across packages: the whole request lifecycle (parse, validate, store, respond) lives in two files.

---

## Directory Structure

```text
monolith-simple/
├── main.go
├── store.go
├── handlers.go
├── main_test.go
├── .air.toml
├── .env.example
├── go.mod
└── go.sum
```

No `cmd/`, no `internal/`. Everything is `package main`, so the split between files here is just for readability, not for architecture. Routing (`chi`) and config loading (`godotenv`) are the same tooling used by the other examples in this repository.

---

## When to Use

Prototypes, MVPs, small internal tools, and scripts that need a couple of HTTP endpoints without any architectural ceremony.

It is suitable for:

- Prototypes and proofs of concept
- MVPs where speed of setup matters more than long-term structure
- Small internal tools and one-off scripts
- Learning the fixed API contract before comparing it against a structured architecture

---

## Pros and Cons

**Pros**

- ✅ Fastest possible setup, nothing to wire together
- ✅ The entire request lifecycle is readable top to bottom in one place
- ✅ No indirection: a handler's logic is right there, not spread across layers
- ✅ Great for prototypes and throwaway tools

**Cons**

- ⚠️ No separation between HTTP concerns and business rules
- ⚠️ Doesn't scale well as the codebase grows
- ⚠️ Shared in-memory state means tests need distinct data per test to avoid collisions
- ⚠️ Hard to swap persistence or add another interface (CLI, gRPC) without touching everything
- ⚠️ Encourages copy-paste growth instead of reuse

---

## Running

```bash
cp .env.example .env
go run .
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
