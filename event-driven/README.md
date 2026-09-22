# Event-Driven Architecture

> One of the architecture implementations in this repository, organized around a publisher that broadcasts what happened, and independent subscribers that react without the publisher knowing they exist.

---

## Core Principle

`user.UserService` never calls `notification` or `audit` directly. After successfully creating a user, it publishes a `UserCreated` value onto a `Bus` (`internal/shared/events`) and moves on. `notification` and `audit` each subscribed to that same event type at startup, independently of one another; neither knows the other exists, and the publisher knows neither of them.

The bus dispatches to every subscriber **asynchronously**, in its own goroutine, so a slow or failing subscriber never delays the HTTP response or takes down another subscriber. This is an in-process bus, not a real message broker: there is no persistence, no retry, no delivery guarantee beyond "best effort while the process is alive." A `Wait()` method exists purely so tests (and the smoke test below) can deterministically wait for the asynchronous handlers to finish before asserting on their side effects; production code never calls it.

---

## Dependency Flow

```text
user.UserController --> user.UserService --> shared/events.Bus
                                                    ^
                              notification.Listener |  audit.Listener
                              (both subscribe at startup, independently)
```

- `user.UserService` depends on `EventPublisher`, an interface it defines itself with a single `Publish` method - it has no idea who, if anyone, is listening.
- `notification.Listener` and `audit.Listener` depend on `shared/events.Bus` and on `user.UserCreated` (the event's shape), never on `user.UserService` or on each other.

---

## Directory Structure

```text
event-driven/
├── cmd/api/main.go
├── internal/
│   ├── user/
│   ├── notification/
│   ├── audit/
│   ├── health/
│   └── shared/
│       ├── events/
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

Workflows where one action should trigger several independent reactions (a welcome email, an audit trail, a cache invalidation) without the code that performs the action having to know about, call, or wait for all of them.

It is suitable for:

- Systems expected to grow more reactions to the same event over time, added without touching the publisher
- Workloads where a slow side effect (sending an email) should never block the primary request
- Codebases that will eventually move this same publish/subscribe shape onto a real broker, keeping the call sites unchanged

---

## Pros and Cons

**Pros**

- ✅ Adding a third subscriber to `UserCreated` touches zero lines in `user.UserService`
- ✅ A subscriber that panics or errors cannot break `POST /users` or any other subscriber
- ✅ The HTTP response returns as soon as the user is persisted; subscribers run after, not in the request's critical path
- ✅ The same `Bus` shape maps directly onto a real broker later: `Publish` becomes a produce, `Subscribe` becomes a consumer group

**Cons**

- ⚠️ No delivery guarantee: if the process crashes between `Publish` and a handler running, that handler's work is simply lost
- ⚠️ Debugging "why didn't the email get sent" means tracing an async call graph instead of reading a stack trace
- ⚠️ Only `UserCreated` has subscribers here; `UserUpdated`/`UserDeleted` were left unpublished on purpose, to avoid events with no consumer
- ⚠️ Overkill for a single reaction to a single action - a direct function call would be simpler and just as decoupled with an interface

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

**Create user** (watch the process log for `[notification]` and `[audit]` lines appearing shortly after the response, since dispatch is asynchronous):

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
