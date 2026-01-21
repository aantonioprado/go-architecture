## Description
<!-- Describe your changes in detail -->

## Type of Change
- [ ] New architecture implementation
- [ ] Bug fix
- [ ] Documentation update
- [ ] Code improvement/refactoring
- [ ] Other (please describe):

## Architecture Affected
<!-- Which architecture(s) does this PR affect? -->
- [ ] monolith-simple
- [ ] layered
- [ ] mvc
- [ ] clean-architecture
- [ ] hexagonal
- [ ] ddd
- [ ] modular-monolith
- [ ] event-driven
- [ ] microservices
- [ ] General/Multiple

## Checklist

### For All PRs
- [ ] My code follows the Go conventions (`gofmt`, `golint`)
- [ ] I have updated documentation (README) where necessary
- [ ] I have tested my changes manually
- [ ] My commit messages are clear and descriptive

### For New Architecture Implementations
- [ ] API contract matches the specification (POST/GET/PUT/DELETE /users, GET /health)
- [ ] All business rules are enforced (email unique, name required, ID/CreatedAt generated)
- [ ] Uses in-memory map or SQLite (no heavy ORMs)
- [ ] README includes:
  - [ ] Core principle explanation
  - [ ] ASCII diagram of dependency flow
  - [ ] Directory/package structure
  - [ ] When to use this architecture
  - [ ] Pros and cons
  - [ ] Example usage
- [ ] No Auth, pagination, caching, or other excluded features
- [ ] Updated main README with link to new architecture

### For Bug Fixes
- [ ] I have added tests that prove my fix is effective
- [ ] Existing functionality is not broken

## Testing
<!-- Describe how you tested your changes -->

**Manual testing commands:**
```bash
# Example:
cd architecture-name/
go run main.go

# Test create user
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","email":"test@example.com"}'
```

## Screenshots (if applicable)
<!-- Add screenshots to help explain your changes -->

## Additional Notes
<!-- Any additional information that reviewers should know -->
