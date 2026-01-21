# Contributing to Go Architecture Patterns

Thank you for your interest in contributing to this project! This guide will help you get started.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Contribution Guidelines](#contribution-guidelines)
- [Submitting Changes](#submitting-changes)

---

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code. Please be respectful, inclusive, and constructive in all interactions.

---

## How Can I Contribute?

### 1. Adding a New Architecture

To add a new architectural pattern:

1. **Create a new directory** named after the architecture (e.g., `cqrs/`, `onion/`)
2. **Implement the exact same User CRUD API** as specified in the main README
3. **Follow the same business rules**:
   - Email must be unique
   - Name is required
   - ID generated at creation (UUID)
   - CreatedAt generated in domain layer
   - Delete strategy (logical or physical, but consistent)
4. **Use the same persistence approach**: In-memory map or SQLite (no heavy ORMs)
5. **Create a detailed README** in your architecture directory with:
   - Core principle explanation
   - ASCII diagram showing dependency flow
   - Package/directory structure
   - When to use this architecture
   - Pros and cons
   - Example usage

### 2. Improving Existing Implementations

- **Better implementation?** Open a PR with your improvements
- **Found a bug?** Create an issue with steps to reproduce
- **Want to add tests?** Test contributions are always welcome!
- **Documentation improvements?** Fix typos, improve clarity, add examples

### 3. Documentation

- Improve explanations in READMEs
- Add more diagrams (ASCII art preferred)
- Translate documentation to other languages
- Add code comments for clarity

---

## Development Setup

### Prerequisites

- Go 1.24.3 or higher
- Git

### Getting Started

1. **Fork the repository**

2. **Clone your fork**
   ```bash
   git clone https://github.com/aantonioprado/go-architecture.git
   cd go-architecture
   ```

3. **Create a branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

4. **Make your changes**

5. **Test your implementation**
   ```bash
   cd your-architecture/
   go run main.go

   # Test the API
   curl -X POST http://localhost:8080/users \
     -H "Content-Type: application/json" \
     -d '{"name":"Antônio Prado","email":"antonio@antonioeprado.dev"}'
   ```

---

## Contribution Guidelines

### The Fixed Contract

All architectures MUST implement the same API:

```
POST   /users          # Create a new user
GET    /users/{id}     # Get user by ID
GET    /users          # List all users
PUT    /users/{id}     # Update user
DELETE /users/{id}     # Delete user
GET    /health         # Health check
```

### User Model (Immutable)

```go
type User struct {
    ID        string    // UUID
    Name      string    // Required
    Email     string    // Unique
    CreatedAt time.Time // Generated in domain
}
```

### What NOT to Include

❌ Authentication/Authorization
❌ Password handling
❌ Advanced pagination
❌ Caching layers
❌ Real message queues
❌ Heavy ORMs
❌ API versioning
❌ Rate limiting

**Keep it simple.** Focus on architecture, not features.

### Code Style

- **Follow Go conventions**: Use `gofmt` and `golint`
- **Idiomatic Go**: Write code that feels natural to Go developers
- **Comments**: Explain architectural decisions, not obvious code
- **Package naming**: Use clear, descriptive names
- **Error handling**: Handle errors properly (don't ignore them)

## Submitting Changes

### Pull Request Process

1. **Update documentation** - If you changed behavior, update the README
2. **Test thoroughly** - Ensure all endpoints work as expected
3. **Commit with clear messages**:
   ```
   feat(hexagonal): add hexagonal architecture implementation
   fix(layered): correct email uniqueness validation
   docs(readme): improve architecture comparison table
   ```
4. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```
5. **Open a Pull Request** with:
   - Clear description of changes
   - Why this change is valuable
   - Screenshots/examples if applicable

### PR Checklist

- [ ] API contract matches specification
- [ ] Business rules are enforced
- [ ] README is complete and clear
- [ ] Code follows Go conventions
- [ ] No unnecessary dependencies
- [ ] Tested manually (all endpoints work)
- [ ] ASCII diagram included (for new architectures)

---

## Questions?

- **Open an issue** for questions about contributing
- **Start a discussion** for architectural debates
- **Check existing issues** before creating a new one

---

## Recognition

All contributors will be recognized in the project. Thank you for making this resource better for the Go community!
