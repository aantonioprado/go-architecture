# Microservices

> One of the architecture implementations in this repository, the only one made of genuinely independent processes talking to each other over the network instead of function calls inside one binary.

---

## Core Principle

The fixed API only has one resource, `User`, so "multiple independent services, each owning its data" is demonstrated here through a **CQRS split** (Command Query Responsibility Segregation) instead of inventing a second, unrelated business domain:

- **`command-service`** - owns the write side. It is the source of truth: it validates, generates the ID and `CreatedAt`, enforces email uniqueness against its own store, and persists first.
- **`query-service`** - owns the read side. It has its own, separate in-memory store, and never validates anything itself: it only stores whatever `command-service` tells it to.
- **`gateway`** - the only service a client ever talks to. It exposes the exact same fixed contract as every other architecture in this repository, and reverse-proxies each request to whichever backend owns that operation.

Each is its own Go module (its own `go.mod`), with no package shared between them. `command-service` and `query-service` do not import each other or a common library; the only thing that connects them is a real HTTP call `command-service` makes to `query-service` after every write, to keep the read side in sync. That call, not a shared database and not a message queue, is the entire coupling between the two.

---

## Dependency Flow

```text
client --> gateway --(POST/PUT/DELETE /users)--> command-service --(HTTP replication)--> query-service
                  \-> gateway --(GET /users, GET /users/{id})----------------------------> query-service
```

- `gateway` depends on two `httputil.ReverseProxy` instances, one per backend base URL - it never imports either service's code, it only knows two URLs.
- `command-service` depends on a `Replicator` interface it defines itself; `internal/replication` is the only concrete implementation, making an HTTP call.
- `query-service` never calls out to anyone; it is only ever called.

---

## Directory Structure

```text
microservices/
├── command-service/     # its own go.mod, cmd/api, internal/{user,replication,health,...}
├── query-service/       # its own go.mod, cmd/api, internal/{user,health,...}
├── gateway/              # its own go.mod, cmd/api, internal/{proxy,health,...}
├── docker-compose.yml
└── README.md
```

---

## When to Use

Large, distributed teams that need to own, deploy, and scale a piece of the system independently of everyone else, or workloads whose read and write traffic have very different shapes (far more reads than writes, or the other way around) and would benefit from scaling those independently.

It is suitable for:

- Organizations where a single team cannot review or safely deploy every change across the whole system anymore
- Read-heavy workloads that want a read model optimized differently from how writes are validated and stored
- Systems already comfortable with eventual consistency being a real, visible property, not an implementation detail

---

## Pros and Cons

**Pros**

- ✅ `command-service` and `query-service` can be deployed, scaled, and even rewritten independently, as long as the HTTP contract between them stays the same
- ✅ A bug in `query-service` cannot corrupt the source of truth; `command-service` never reads from it
- ✅ `gateway` gives clients one stable contract while the backends evolve independently behind it
- ✅ Nothing about this shape requires a shared database or a shared library, the two biggest sources of accidental coupling between services

**Cons**

- ⚠️ Eventual consistency is real here, not simulated: a `GET` immediately after a `POST` can, in principle, race the replication call and miss the new user
- ⚠️ Three services means three deployables, three sets of logs, three things that can be down independently - `docker-compose.yml` exists because coordinating this by hand across terminals stops being practical fast
- ⚠️ No schema registry or generated client between the services, just JSON over HTTP and a shared understanding of the wire format kept by convention, not by the compiler
- ⚠️ Overkill for a single resource this small; the CQRS split exists here to demonstrate the pattern, not because `User` actually needs two different storage strategies

---

## Running

With Docker Compose (recommended, since only this configuration wires the real network boundaries between the three services):

```bash
cp .env.example .env
docker compose build
docker compose up
```

There is a single `.env.example` at the root of `microservices/`, not one per service: `docker-compose.yml` reads it directly (`${PORT}`, `${COMMAND_SERVICE_URL}`, `${QUERY_SERVICE_URL}`) and injects each value into the right container. `PORT` is the same for all three on purpose, since each runs in its own container and never competes for a host port; only `gateway` publishes one to the host (`8080`), `command-service` and `query-service` are reachable only from inside the compose network.

Without Docker, in three separate terminals, since all three would otherwise fight over the same host ports:

```bash
cd query-service && PORT=8081 go run ./cmd/api
cd command-service && PORT=8082 QUERY_SERVICE_URL=http://localhost:8081 go run ./cmd/api
cd gateway && PORT=8080 COMMAND_SERVICE_URL=http://localhost:8082 QUERY_SERVICE_URL=http://localhost:8081 go run ./cmd/api
```

---

## Example Usage

All requests go through the gateway on port 8080, exactly like every other architecture in this repository:

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
