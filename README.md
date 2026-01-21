# Go Architecture

![Go Version](https://img.shields.io/badge/Go-1.24.3-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Build](https://img.shields.io/badge/Build-Passing-brightgreen?style=for-the-badge)

> A comprehensive, side-by-side comparison of 9 software architecture patterns in Go, all implementing the exact same User CRUD API.

## Table of Contents

- [Why This Project Exists](#why-this-project-exists)
- [The Fixed User CRUD API](#the-fixed-user-crud-api)
- [User Domain Model](#user-domain-model)
- [Business Rules](#business-rules)
- [What's Intentionally Simple](#whats-intentionally-simple)
- [Architecture Comparison](#architecture-comparison)
- [The 9 Architectures](#the-9-architectures)
- [Directory Structure](#directory-structure)
- [How to Use This Repository](#how-to-use-this-repository)
- [Contributing](#contributing)
- [License](#license)

---

## Why This Project Exists

**The Problem:** Choosing the right architecture is hard. Every article, book, and framework advocates for different approaches, but practical comparisons are rare.

**The Solution:** This repository implements **the exact same User CRUD API** in 9 different architectural patterns. No framework noise, no complex business logic—just pure architectural comparison.

### Learning Objectives

By exploring this repository, you'll understand:

✅ **Dependency Flow** - How each architecture manages dependencies between layers
✅ **Testability** - Which patterns make testing easier or harder
✅ **Complexity Trade-offs** - Initial setup cost vs. long-term maintainability
✅ **Domain Protection** - How business rules are isolated from infrastructure
✅ **Real-world Applicability** - When to use each pattern in production

### Who This Is For

- **Go developers** learning architectural patterns
- **Tech leads** evaluating architecture for new projects
- **Teams** debating monolith vs microservices vs modular approaches
- **Students** studying software design and architecture
- **Anyone** who wants to see architecture patterns in action, not just theory

---

## The Fixed User CRUD API

Every architecture implements the **exact same HTTP API**:

```http
POST   /users          # Create a new user
GET    /users/{id}     # Get user by ID
GET    /users          # List all users
PUT    /users/{id}     # Update user
DELETE /users/{id}     # Delete user
GET    /health         # Health check
```

### Example Request/Response

**Create User:**
```bash
POST /users
Content-Type: application/json

{
  "name": "Antônio Prado",
  "email": "antonio@antonioeprado.dev"
}
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Antônio Prado",
  "email": "antonio@antonioeprado.dev",
  "createdAt": "2026-01-20T00:00:00Z"
}
```

---

## User Domain Model

The `User` entity is **immutable** and consistent across all architectures:

```go
type User struct {
    ID        string    // UUID generated at creation
    Name      string    // Required field
    Email     string    // Must be unique
    CreatedAt time.Time // Generated in domain layer
}
```

**Key Design Decisions:**

- **Immutable** - Once created, the structure doesn't change (updates create new instances)
- **UUID for ID** - Universally unique, generated in domain logic (not database)
- **Domain Timestamps** - `CreatedAt` is set by domain logic, not infrastructure
- **Simple Fields** - No password, no roles—focus stays on architecture

---

## Business Rules

All architectures enforce these **exact same rules**:

| Rule | Description | Where Enforced |
|------|-------------|----------------|
| **Email Uniqueness** | No two users can have the same email | Domain/Service Layer |
| **Name Required** | Name field cannot be empty | Domain Validation |
| **ID Generation** | UUID generated at user creation | Domain Logic |
| **CreatedAt Generation** | Timestamp set when user is created | Domain Logic |
| **Delete Strategy** | Logical OR Physical delete (consistent per architecture) | Repository Layer |

### Why These Rules?

- **Simple enough** to not distract from architecture
- **Complex enough** to show how validation flows through layers
- **Real-world** - email uniqueness requires repository interaction
- **Testable** - easy to verify rule enforcement in different layers

---

## Architecture Comparison

| Architecture | Separation | Testability | Complexity | Dependency Flow | Best For | Learning Curve |
|--------------|-----------|-------------|------------|-----------------|----------|----------------|
| [**Simple Monolith 🏗️**](./monolith-simple/) | ⭐ Low | ⭐⭐ Moderate | ⭐ Very Low | Mixed | Prototypes, MVPs, small tools | ⭐ Very Easy |
| [**Layered 🏗️**](./layered/) | ⭐⭐⭐ Good | ⭐⭐⭐ Good | ⭐⭐ Low | Top → Down | Traditional CRUD apps | ⭐⭐ Easy |
| [**MVC 🏗️**](./mvc/) | ⭐⭐⭐ Good | ⭐⭐⭐ Good | ⭐⭐ Low | Controller → Model | Web apps, admin panels | ⭐⭐ Easy |
| [**Clean Architecture 🏗️**](./clean-architecture/) | ⭐⭐⭐⭐⭐ Excellent | ⭐⭐⭐⭐⭐ Excellent | ⭐⭐⭐⭐ High | Inward | Complex domains, long-lived systems | ⭐⭐⭐⭐ Hard |
| [**Hexagonal 🏗️**](./hexagonal/) | ⭐⭐⭐⭐⭐ Excellent | ⭐⭐⭐⭐⭐ Excellent | ⭐⭐⭐⭐ High | Outside → Inside | Swappable infrastructure | ⭐⭐⭐⭐ Hard |
| [**DDD 🏗️**](./ddd/) | ⭐⭐⭐⭐⭐ Excellent | ⭐⭐⭐⭐ Very Good | ⭐⭐⭐⭐⭐ Very High | Domain-centric | Rich domains, complex logic | ⭐⭐⭐⭐⭐ Very Hard |
| [**Modular Monolith 🏗️**](./modular-monolith/) | ⭐⭐⭐⭐ Very Good | ⭐⭐⭐⭐ Very Good | ⭐⭐⭐ Moderate | Module boundaries | Scalable monoliths | ⭐⭐⭐ Moderate |
| [**Event-Driven 🏗️**](./event-driven/) | ⭐⭐⭐⭐ Very Good | ⭐⭐⭐ Good | ⭐⭐⭐⭐ High | Event flow | Async workflows | ⭐⭐⭐⭐ Hard |
| [**Microservices 🏗️**](./microservices/) | ⭐⭐⭐⭐⭐ Excellent | ⭐⭐⭐⭐ Very Good | ⭐⭐⭐⭐⭐ Very High | Service boundaries | Distributed teams | ⭐⭐⭐⭐⭐ Very Hard |

### Key Insights

- 💡 **Complexity vs Benefits:** More complex architectures provide better separation but higher initial cost
- 💡 **No Silver Bullet:** The "best" architecture depends on your team, domain, and scale
- 💡 **Start Simple:** Most projects should start simple and evolve complexity when needed
- 💡 **Testability Correlation:** Better separation usually means better testability

---

## The 9 Architectures

Click on any architecture to see its detailed implementation and explanation:

### 1. [Simple Monolith 🏗️](./monolith-simple/)
Single package, minimal structure. All code in one place.

**When to use:** Prototypes, MVPs, small internal tools, scripts with HTTP endpoints.

---

### 2. [Layered Architecture 🏗️](./layered/)
Traditional N-tier: Presentation → Business → Data layers.

**When to use:** Traditional CRUD applications, teams familiar with layered patterns.

---

### 3. [MVC (Model-View-Controller) 🏗️](./mvc/)
Separate data (Model), presentation (View), and orchestration (Controller).

**When to use:** Web applications with UI, form-heavy applications, APIs following MVC frameworks.

---

### 4. [Clean Architecture 🏗️](./clean-architecture/)
Uncle Bob's Clean Architecture. Dependency Inversion—outer layers depend on inner layers.

**When to use:** Complex, long-lived systems requiring maximum testability.

---

### 5. [Hexagonal Architecture (Ports & Adapters) 🏗️](./hexagonal/)
Application core with ports (interfaces) and adapters (implementations) on the outside.

**When to use:** Need to swap infrastructure easily, heavy testing requirements, multiple input sources.

---

### 6. [Domain-Driven Design (DDD) 🏗️](./ddd/)
Rich domain model with ubiquitous language, aggregates, entities, and value objects.

**When to use:** Complex business domains, collaboration with domain experts, long-lived evolving systems.

---

### 7. [Modular Monolith 🏗️](./modular-monolith/)
Single deployable with strong module boundaries and explicit contracts.

**When to use:** Growing teams needing autonomy, monolith that needs better structure, preparing for microservices.

---

### 8. [Event-Driven Architecture 🏗️](./event-driven/)
Components communicate via events, enabling loose coupling and async workflows.

**When to use:** Async workflows (email, notifications), need to decouple systems, event sourcing.

---

### 9. [Microservices 🏗️](./microservices/)
Multiple independent services, each owning its data and deployable separately.

**When to use:** Large distributed teams, independent scaling needs, polyglot requirements.

---

## Directory Structure

```
go-architecture/
├── 🏗️ monolith-simple/              # 1. Simple monolith
├── 🏗️ layered/                      # 2. Layered
├── 🏗️ mvc/                          # 3. MVC pattern
├── 🏗️ clean-architecture/           # 4. Clean Architecture
├── 🏗️ hexagonal/                    # 5. Hexagonal (Ports & Adapters)
├── 🏗️ ddd/                          # 6. Domain-Driven Design
├── 🏗️ modular-monolith/             # 7. Modular Monolith
├── 🏗️ event-driven/                 # 8. Event-Driven
└── 🏗️ microservices/                # 9. Microservices
```

---

## How to Use This Repository

### For Learning

1. **Start with the comparison table** - Understand the landscape
2. **Pick an architecture** - Click through to its README
3. **Study the implementation** - See how it handles the same API
4. **Compare with others** - Notice the differences in structure
5. **Run the examples** - Test them locally

### For Evaluation

1. Review the [Architecture Comparison Table](#architecture-comparison)
2. Check implementations for patterns you're considering
3. Run the same API tests against different architectures
4. Evaluate complexity vs benefits for your context

### Running Examples

```bash
# Navigate to any architecture
cd monolith-simple

# Run the service
go run main.go

# In another terminal, test it
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Antônio Prado","email":"antonio@antonioeprado.dev"}'

# Get user
curl http://localhost:8080/users/{id}

# List users
curl http://localhost:8080/users

# Health check
curl http://localhost:8080/health
```

---

## Contributing

We welcome contributions! Here's how you can help:

### Improving Existing Implementations

- Found a better way to implement a pattern? Open a PR!
- See a bug or inconsistency? Create an issue!
- Want to add tests? Contributions welcome!

### Documentation

- Improve explanations in READMEs
- Add more diagrams
- Translate to other languages

### Guidelines

- ✅ Keep the API contract identical
- ✅ Maintain the same business rules
- ✅ Focus on architecture, not features
- ✅ Include tests and documentation
- ✅ Use Go idiomatic patterns

---

## License

MIT License - see LICENSE file for details.
