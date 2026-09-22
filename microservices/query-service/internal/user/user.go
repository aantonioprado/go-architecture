package user

import "time"

// User here is a plain read replica: query-service trusts command-service to
// have already validated and generated these fields, it never builds a User
// from raw input itself.
type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}
