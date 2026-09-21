# Domain-Driven Design (DDD)

> One of the architecture implementations in this repository, organized around a rich domain model instead of a passive one waiting to be told what to do.

---

## Core Principle

Everything revolves around the `User` **Aggregate Root**, not a data structure some outer layer manipulates:

- **Value Object** (`domain/user.Email`) - a type with no identity of its own, defined by its value, self-validating at construction (`NewEmail`).
- **Aggregate Root** (`domain/user.User`) - owns its state (private fields, only reachable through getters) and its own invariants. `Register` and `ChangeDetails` are behavior on the aggregate, not steps an outer service performs on a struct's public fields.
- **Domain Events** (`domain/user.Event` and friends) - `Register`/`ChangeDetails` record `UserRegistered`/`UserDetailsChanged` on the aggregate itself. No publisher or message broker: the Application Service reads `Events()` after persisting and clears them. It is the tactical pattern, not the infrastructure around it.
- **Repository** (`domain/user.Repository`) - one per aggregate, defined by the domain, persisting the aggregate as a whole (`Save`, not separate `Create`/`Update`).
- **Application Service** (`application/user.Service`) - thin orchestration: loads the aggregate, calls its behavior, persists it, publishes its events. Anything that needs to reach across aggregates (the email-uniqueness check) lives here, never inside `User` itself.

---

## Dependency Flow

```text
interfaces/http --> application/user (Service) --> domain/user (Aggregate Root)
                                                          ^
                                    infrastructure/persistence/memory (Repository)
```

- `application/user.Service` depends on `domain/user.Repository`, an interface the domain defines, not on the concrete `memory.UserRepository`.
- `interfaces/http.UserHandler` depends on `application/user.Service`, an interface, not on the concrete `applicationService`.
- Every dependency arrow points at `domain/user`; nothing in `domain/user` imports outward.

Unlike `hexagonal` in this repository, where the core is a symmetric pair of ports around an otherwise plain service, here the aggregate itself carries behavior and records its own events. The Application Service is not a stand-in for business logic, it is a coordinator.

---

## Directory Structure

```text
ddd/
├── cmd/api/main.go
├── internal/
│   ├── domain/user/
│   ├── application/user/
│   ├── infrastructure/persistence/memory/
│   ├── interfaces/http/
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

Complex business domains where the rules are worth modeling explicitly, and where collaboration with domain experts benefits from code that reads in the same language they use.

It is suitable for:

- Domains where an entity enforces real invariants, not just required-field checks
- Systems that will eventually need to react to what happened (a registration, a change) without necessarily needing a message broker on day one
- Long-lived, evolving systems where a rich model pays for itself over time

---

## Pros and Cons

**Pros**

- ✅ The aggregate owns its invariants; there is no way to construct or mutate a `User` into an invalid state from outside `domain/user`
- ✅ Domain Events make "what happened" explicit and inspectable, instead of being implied by whichever fields changed
- ✅ The Application Service stays thin: reading it tells you what happens, reading the aggregate tells you why
- ✅ A Value Object like `Email` can grow real invariants later without touching anything outside `domain/user`

**Cons**

- ⚠️ A single aggregate this small does not need its own bounded context folder; the ceremony is easiest to justify once there is more than one aggregate
- ⚠️ Domain Events with no subscriber are a pattern demonstration, not a working notification system
- ⚠️ Cross-aggregate rules (like email uniqueness) still need somewhere to live outside the aggregate, and that boundary takes explanation
- ⚠️ Overkill for a domain whose rules really are just "two required fields"

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
