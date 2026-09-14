# Clean Architecture

> One of the architecture implementations in this repository, organized into four concentric rings where source code dependencies only ever point inward.

---

## Core Principle

Four rings, each depending only on the ring inside it:

- **Entities** (`entities`) - the `User` type and the rules that make one valid on its own (name/email required, ID and `CreatedAt` generated). Depends on nothing.
- **Use Cases** (`usecases`) - the `UserInteractor`, plus two ports it defines itself: `UserRepository` (persistence) and `UserOutputPort` (presentation). It depends only on Entities and on interfaces it owns.
- **Interface Adapters** (`adapters`) - a Gateway (`InMemoryUserRepository`) implementing `UserRepository`, a Presenter (`HTTPUserPresenter`) implementing `UserOutputPort`, and a Controller decoding HTTP requests into use case input.
- **Frameworks & Drivers** - `cmd/api`, `config`, `middleware`, `routes`, `server`. The composition root (`server.Build`) is the only place that wires a concrete Gateway into the Interactor and a concrete Controller into the router.

The dependency inversion happens twice: the Gateway is an outer type satisfying an inner interface, and so is the Presenter. Neither the Interactor nor the Use Cases ring imports either of them.

---

## Dependency Flow

```text
main --> server --> router --> controller --> usecase (interactor) --> entities
                                    |                  |
                                    |                  +--> repository (port, implemented by gateway)
                                    +--> presenter (implements the output port)
```

- The **Controller** depends on `UserInputPort`, an interface, not on the concrete `UserInteractor`.
- The **Interactor** depends on `UserRepository` and reports through `UserOutputPort` - both interfaces it defines itself.
- The **Gateway** and the **Presenter** are outer types that implement those interfaces; nothing inward knows they exist.

---

## Directory Structure

```text
clean-architecture/
├── cmd/api/main.go
├── internal/
│   ├── entities/
│   ├── usecases/
│   ├── adapters/
│   │   ├── dto/
│   │   ├── response/
│   │   ├── presenter/
│   │   ├── controller/
│   │   └── gateway/
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

Complex, long-lived systems that need maximum testability and the ability to swap infrastructure (persistence, delivery mechanism) without touching business rules.

It is suitable for:

- Systems expected to evolve for years, across multiple teams
- Domains where business rules need to be tested in isolation, with no HTTP server or database involved
- Codebases likely to swap a storage engine, a web framework, or add a second delivery mechanism (CLI, gRPC) later

---

## Pros and Cons

**Pros**

- ✅ Business rules are testable with plain fakes, no HTTP or real persistence needed
- ✅ Swapping the storage engine only means writing a new Gateway
- ✅ Swapping the delivery mechanism only means writing a new Controller/Presenter pair
- ✅ The Dependency Rule makes accidental coupling to infrastructure hard to introduce by accident

**Cons**

- ⚠️ Significant boilerplate for a CRUD example this small: five interfaces, input/output structs, and four packages before writing a single business rule
- ⚠️ New contributors need to learn the ring vocabulary before the codebase reads naturally
- ⚠️ Wiring the composition root by hand grows tedious as more use cases are added
- ⚠️ Overkill for prototypes or short-lived tools

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
