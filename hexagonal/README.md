# Hexagonal Architecture (Ports & Adapters)

> One of the architecture implementations in this repository, organized around a core that only talks to the outside world through ports it defines itself.

---

## Core Principle

The application core sits in the middle and never imports an adapter:

- **Core** (`core`) - the `domain` package holds the `User` entity and the rules that make one valid on its own (name/email required, ID and `CreatedAt` generated). `core/service` implements the primary port using only the secondary port - it has no idea an HTTP server or an in-memory map exist.
- **Ports** (`core/ports`) - two interfaces, both owned by the core:
  - `UserService` (**primary/driving port**) - what the core exposes to whoever drives it.
  - `UserRepository` (**secondary/driven port**) - what the core needs from persistence.
- **Primary/Driving Adapters** (`adapters/primary/http`) - a `UserHandler` that decodes HTTP requests, calls `UserService`, and turns the returned `(domain.User, error)` into a response. It drives the core.
- **Secondary/Driven Adapters** (`adapters/secondary/memory`) - an `InMemoryUserRepository` implementing `UserRepository`. It is driven by the core.

Both sides are symmetric: a primary adapter calls in through a primary port, a secondary adapter is called out to through a secondary port. Swapping either adapter never touches `core`.

---

## Dependency Flow

```text
adapters/primary/http --> core/ports (UserService) --> core/service --> core/ports (UserRepository) --> adapters/secondary/memory
        (driving adapter)      (driving port)            (the hexagon)        (driven port)                (driven adapter)
```

- The **UserHandler** depends on `ports.UserService`, an interface, not on the concrete `userService`.
- The **userService** depends on `ports.UserRepository`, an interface it defines itself, not on `InMemoryUserRepository`.
- Both adapters depend inward on `core`; `core` depends on neither adapter.

Unlike a callback/presenter-style output port, `UserService` here returns `(domain.User, error)` directly - any adapter (HTTP, CLI, a test) can call it and decide for itself what to do with the result.

---

## Directory Structure

```text
hexagonal/
├── cmd/api/main.go
├── internal/
│   ├── core/
│   │   ├── domain/
│   │   ├── ports/
│   │   └── service/
│   ├── adapters/
│   │   ├── primary/http/
│   │   └── secondary/memory/
│   ├── config/
│   ├── middleware/
│   ├── routes/
│   └── server/
├── .air.toml
├── .env.example
├── go.mod
└── go.sum
```

---

## When to Use

Systems that need to swap infrastructure - persistence, delivery mechanism, or both - without touching business rules, and where testing the core in isolation matters more than a rigid layering vocabulary.

It is suitable for:

- Codebases expected to gain a second delivery mechanism (CLI, gRPC, message consumer) alongside HTTP
- Domains where persistence is likely to move from an in-memory store to SQL, a document store, or an external API
- Teams that want to unit-test business rules against fakes, with zero HTTP or real database involved

---

## Pros and Cons

**Pros**

- ✅ The core never imports an adapter - dependencies only point inward
- ✅ Adding a second primary adapter (a CLI, say) needs nothing new from the core: it just calls `UserService`
- ✅ Ports return plain values, so tests call the core the same way any adapter would - no fakes for a presenter needed
- ✅ Swapping storage only means writing a new secondary adapter

**Cons**

- ⚠️ Two ports and two interfaces for a CRUD example this small is more indirection than the domain needs
- ⚠️ "Primary" vs "secondary", "driving" vs "driven" is vocabulary contributors have to learn before the layout reads naturally
- ⚠️ The composition root still wires everything by hand; that only gets more tedious as ports multiply
- ⚠️ Overkill for a prototype or a script that will never gain a second adapter on either side

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
