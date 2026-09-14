package usecases

import "time"

type CreateUserInput struct {
	Name  string
	Email string
}

type UserOutput struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}

type GetUserInput struct {
	ID string
}

type ListUsersOutput struct {
	Users []UserOutput
}

type UpdateUserInput struct {
	ID    string
	Name  string
	Email string
}
