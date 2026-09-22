# Modular Monolith

> One of the architecture implementations in this repository, organized as a single deployable where business logic lives inside a module with a clear boundary, instead of being spread across generic technical layers shared by everything.

---

## Core Principle

A single binary, a single process, but the code is not one big shared tree: business logic lives inside `internal/user`, and nothing outside that package reaches into how it validates, stores, or serves a `User`. `internal/health` is trivial and sits outside the module because it has no business rule of its own. `internal/shared` is what the application needs to exist at all (config, logging, routing, the composition root) - infrastructure every module would depend on, not a module itself.

Go's own `internal/` visibility already keeps this whole tree unreachable from outside the repository's Go module. With a single business module, there is no sibling to be blocked from reaching into `internal/user` - the point being illustrated here is the shape a module takes, not a live demonstration of the compiler stopping a violation. A second module would get its own nested `internal/` inside its own folder, so that even a neighbor inside the same binary could only reach it through whatever it chooses to export, the same discipline that already governs `internal/user` from the rest of the codebase.

---

## Dependency Flow

```text
shared/server (composition root) --> user.NewUserController --> user.NewUserService --> user.NewInMemoryUserRepository
                                          |
                                          +--> shared/routes --> user.UserController (HTTP entry point)
```

- `shared/routes` and `shared/server` depend on the `user` package's exported types (`UserController`, `NewUserService`, `NewInMemoryUserRepository`) - never on anything unexported inside it.
- Nothing inside `internal/user` imports `shared` or `health` - the module doesn't know the transport it's served over is HTTP, or that a composition root exists.

---

## Directory Structure

```text
modular-monolith/
├── cmd/api/main.go
├── internal/
│   ├── user/
│   ├── health/
│   └── shared/
│       ├── config/
│       ├── middleware/
│       ├── response/
│       ├── routes/
│       └── server/
├── .air.toml
├── .env.example
├── go.mod
└── go.sum
```

---

## When to Use

Growing teams that need to own a slice of the codebase without stepping on each other, monoliths whose folders have started blurring into each other, or systems being deliberately prepared for an eventual split into services.

It is suitable for:

- A single team today, several teams tomorrow, each expected to own a module
- Codebases where "just import it, it's the same binary" has already caused unwanted coupling once
- A deployment that should stay a monolith for now, but whose modules are cut along the same seams a future microservice extraction would use

---

## Pros and Cons

**Pros**

- ✅ A module's internals are unreachable by construction, not by convention someone has to remember
- ✅ A team can rewrite everything inside `internal/user` without any other package noticing, as long as the exported surface stays the same
- ✅ Splitting a module out into its own service later is mostly a transport change, not a rewrite, because the boundary was already there
- ✅ Still one binary, one deploy, no network calls between modules, no distributed transactions

**Cons**

- ⚠️ With a single module, as here, there is no neighbor for the boundary to actually be tested against - the discipline is real, but nothing is currently being kept out
- ⚠️ Small systems pay for a structure that only earns its cost once there is more than one module fighting for space
- ⚠️ `internal/` boundaries are a Go-specific trick; the same discipline in another language needs a build tool or a lint rule to enforce it instead

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
